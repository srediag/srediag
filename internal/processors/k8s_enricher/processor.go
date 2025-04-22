package k8s_enricher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Config holds the configuration for the Kubernetes enricher processor
type Config struct {
	// IncludeSecrets determines if secret data should be included
	IncludeSecrets bool `mapstructure:"include_secrets"`

	// IncludeConfigMaps determines if ConfigMap data should be included
	IncludeConfigMaps bool `mapstructure:"include_configmaps"`

	// LabelSelector filters resources by labels
	LabelSelector string `mapstructure:"label_selector"`

	// Namespaces limits the namespaces to watch (empty means all)
	Namespaces []string `mapstructure:"namespaces"`
}

// Processor implements Kubernetes metadata enrichment
type Processor struct {
	logger   *zap.Logger
	config   Config
	client   kubernetes.Interface
	tracer   trace.Tracer
	mu       sync.RWMutex
	cache    map[string]interface{}
	stopChan chan struct{}
}

// NewProcessor creates a new Kubernetes enricher processor
func NewProcessor(config Config, logger *zap.Logger, tracer trace.Tracer) (*Processor, error) {
	// Create Kubernetes client
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &Processor{
		logger:   logger,
		config:   config,
		client:   clientset,
		tracer:   tracer,
		cache:    make(map[string]interface{}),
		stopChan: make(chan struct{}),
	}, nil
}

// Start begins watching Kubernetes resources
func (p *Processor) Start(ctx context.Context) error {
	// Initial cache population
	if err := p.refreshCache(ctx); err != nil {
		return fmt.Errorf("initial cache population failed: %w", err)
	}

	// Start watching for changes
	go p.watch(ctx)

	return nil
}

// Stop stops watching Kubernetes resources
func (p *Processor) Stop(ctx context.Context) error {
	close(p.stopChan)
	return nil
}

// OnStart is called when a span starts
func (p *Processor) OnStart(_ context.Context, _ sdktrace.ReadWriteSpan) {
	// No-op: enrichment happens on end
}

// OnEnd is called when a span ends
func (p *Processor) OnEnd(s sdktrace.ReadOnlySpan) {
	ctx := context.Background()
	_, span := p.tracer.Start(ctx, "k8s.enrich")
	defer span.End()

	// Extract pod name from attributes
	var podName, namespace string
	for _, attr := range s.Attributes() {
		switch attr.Key {
		case "k8s.pod.name":
			podName = attr.Value.AsString()
		case "k8s.namespace.name":
			namespace = attr.Value.AsString()
		}
	}

	if podName == "" || namespace == "" {
		return
	}

	// Get pod metadata
	pod, err := p.getPodMetadata(ctx, namespace, podName)
	if err != nil {
		p.logger.Error("failed to get pod metadata",
			zap.String("pod", podName),
			zap.String("namespace", namespace),
			zap.Error(err))
		return
	}

	// Add Kubernetes metadata
	span.SetAttributes(
		attribute.String("k8s.node.name", pod.Spec.NodeName),
		attribute.String("k8s.pod.uid", string(pod.UID)),
		attribute.String("k8s.pod.start_time", pod.Status.StartTime.String()),
	)

	// Add labels as attributes
	for k, v := range pod.Labels {
		span.SetAttributes(
			attribute.String(fmt.Sprintf("k8s.pod.label.%s", k), v),
		)
	}

	// Add ConfigMap versions if enabled
	if p.config.IncludeConfigMaps {
		for _, vol := range pod.Spec.Volumes {
			if vol.ConfigMap != nil {
				if cm, err := p.getConfigMapVersion(ctx, namespace, vol.ConfigMap.Name); err == nil {
					span.SetAttributes(
						attribute.String(
							fmt.Sprintf("k8s.configmap.%s.version", vol.ConfigMap.Name),
							cm.ResourceVersion,
						),
					)
				}
			}
		}
	}
}

// getPodMetadata gets pod metadata from cache or Kubernetes API
func (p *Processor) getPodMetadata(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	p.mu.RLock()
	if pod, ok := p.cache[fmt.Sprintf("%s/%s", namespace, name)].(*corev1.Pod); ok {
		p.mu.RUnlock()
		return pod, nil
	}
	p.mu.RUnlock()

	// Not in cache, get from API
	pod, err := p.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Update cache
	p.mu.Lock()
	p.cache[fmt.Sprintf("%s/%s", namespace, name)] = pod
	p.mu.Unlock()

	return pod, nil
}

// getConfigMapVersion gets ConfigMap version from cache or Kubernetes API
func (p *Processor) getConfigMapVersion(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	p.mu.RLock()
	if cm, ok := p.cache[fmt.Sprintf("cm/%s/%s", namespace, name)].(*corev1.ConfigMap); ok {
		p.mu.RUnlock()
		return cm, nil
	}
	p.mu.RUnlock()

	// Not in cache, get from API
	cm, err := p.client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Update cache
	p.mu.Lock()
	p.cache[fmt.Sprintf("cm/%s/%s", namespace, name)] = cm
	p.mu.Unlock()

	return cm, nil
}

// refreshCache refreshes the entire cache
func (p *Processor) refreshCache(ctx context.Context) error {
	// List pods
	pods, err := p.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		LabelSelector: p.config.LabelSelector,
	})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	p.mu.Lock()
	for _, pod := range pods.Items {
		p.cache[fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)] = &pod
	}
	p.mu.Unlock()

	// List ConfigMaps if enabled
	if p.config.IncludeConfigMaps {
		cms, err := p.client.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list configmaps: %w", err)
		}

		p.mu.Lock()
		for _, cm := range cms.Items {
			p.cache[fmt.Sprintf("cm/%s/%s", cm.Namespace, cm.Name)] = &cm
		}
		p.mu.Unlock()
	}

	return nil
}

// watch watches for Kubernetes resource changes
func (p *Processor) watch(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			if err := p.refreshCache(ctx); err != nil {
				p.logger.Error("cache refresh failed", zap.Error(err))
			}
		}
	}
}
