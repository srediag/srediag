package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/srediag/srediag/internal/app"
	"github.com/srediag/srediag/internal/config"
)

const (
	defaultConfigPath = "/etc/srediag/config/srediag.yaml"
	shutdownTimeout   = 30 * time.Second
)

// Exit codes
const (
	ExitSuccess       = 0
	ExitGeneralError  = 1
	ExitConfigError   = 2
	ExitPermissionErr = 3
	ExitNotFoundError = 4
	ExitTimeoutError  = 5
)

var (
	// Global flags
	configPath   string
	outputFormat string
	verbose      bool
	quiet        bool
	outputFile   string
	version      bool

	// Version information (set during build)
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"

	rootCmd = &cobra.Command{
		Use:   "srediag",
		Short: "SREDIAG - SRE Diagnostics Tool",
		Long: `SREDIAG is a tool for SRE diagnostics and monitoring.
It helps identify and diagnose issues in your infrastructure.

Usage:
  srediag <category> <command> [options]

Categories:
  diagnose    Run diagnostic checks
  analyze     Analyze system resources
  monitor     Monitor system in real-time
  profile     Profile system performance
  scan        Security scanning
  check       Compliance checking
  audit       Configuration auditing`,
		RunE: run,
	}

	// Diagnostic commands
	diagnoseCmd = &cobra.Command{
		Use:   "diagnose",
		Short: "Run diagnostics",
		Long: `Run various diagnostic checks on the system, Kubernetes, or cloud resources.
		
Examples:
  # Basic system diagnostics
  srediag diagnose system
  srediag diagnose system --resource cpu
  srediag diagnose system --resource memory
  srediag diagnose system --resource disk`,
	}

	diagnoseSystemCmd = &cobra.Command{
		Use:   "system [--resource <resource>]",
		Short: "Run system diagnostics",
		Long:  "Run diagnostics on the local system resources",
		RunE:  diagnoseSystem,
	}

	diagnoseKubernetesCmd = &cobra.Command{
		Use:   "kubernetes [--cluster <cluster>]",
		Short: "Run Kubernetes diagnostics",
		Long:  "Run diagnostics on Kubernetes clusters",
		RunE:  diagnoseKubernetes,
	}

	// Analysis commands
	analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze resources",
		Long: `Analyze various resources and provide insights.
		
Examples:
  # Process analysis
  srediag analyze process --pid 1234
  
  # Memory analysis
  srediag analyze memory --threshold 90
  
  # Bottleneck detection
  srediag analyze bottlenecks --service my-service`,
	}

	analyzeProcessCmd = &cobra.Command{
		Use:   "process [--pid <pid>]",
		Short: "Analyze process",
		Long:  "Analyze a specific process and its resource usage",
		RunE:  analyzeProcess,
	}

	analyzeMemoryCmd = &cobra.Command{
		Use:   "memory [--threshold <percent>]",
		Short: "Analyze memory usage",
		Long:  "Analyze system memory usage and identify issues",
		RunE:  analyzeMemory,
	}

	// Monitor commands
	monitorCmd = &cobra.Command{
		Use:   "monitor",
		Short: "Monitor resources",
		Long: `Monitor various resources in real-time.
		
Examples:
  # Real-time system monitoring
  srediag monitor system --interval 5s`,
	}

	monitorSystemCmd = &cobra.Command{
		Use:   "system [--interval <duration>]",
		Short: "Monitor system",
		Long:  "Monitor system resources in real-time",
		RunE:  monitorSystem,
	}

	// Security commands
	securityCmd = &cobra.Command{
		Use:   "security",
		Short: "Security operations",
		Long: `Perform security-related operations.
		
Examples:
  # Vulnerability scanning
  srediag scan vulnerabilities --severity high
  
  # Compliance checking
  srediag check compliance --standard pci-dss`,
	}

	scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Security scanning",
		Long:  "Perform security scanning operations",
	}

	scanVulnerabilitiesCmd = &cobra.Command{
		Use:   "vulnerabilities [--severity <level>]",
		Short: "Scan vulnerabilities",
		Long:  "Scan for security vulnerabilities",
		RunE:  scanVulnerabilities,
	}

	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Compliance checking",
		Long:  "Perform compliance checking operations",
	}

	checkComplianceCmd = &cobra.Command{
		Use:   "compliance [--standard <standard>]",
		Short: "Check compliance",
		Long:  "Check compliance against security standards",
		RunE:  checkCompliance,
	}
)

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "table", "output format (json/yaml/table)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-error output")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "", "output file path")
	rootCmd.PersistentFlags().BoolVar(&version, "version", false, "print version information")

	// Add diagnostic commands
	diagnoseSystemCmd.Flags().String("resource", "", "resource to diagnose (cpu/memory/disk)")
	diagnoseKubernetesCmd.Flags().String("cluster", "", "target Kubernetes cluster")
	diagnoseCmd.AddCommand(diagnoseSystemCmd, diagnoseKubernetesCmd)
	rootCmd.AddCommand(diagnoseCmd)

	// Add analysis commands
	analyzeProcessCmd.Flags().Int("pid", 0, "process ID to analyze")
	analyzeMemoryCmd.Flags().Float64("threshold", 90.0, "memory usage threshold")
	analyzeCmd.AddCommand(analyzeProcessCmd, analyzeMemoryCmd)
	rootCmd.AddCommand(analyzeCmd)

	// Add monitor commands
	monitorSystemCmd.Flags().Duration("interval", 5*time.Second, "monitoring interval")
	monitorCmd.AddCommand(monitorSystemCmd)
	rootCmd.AddCommand(monitorCmd)

	// Add security commands
	scanVulnerabilitiesCmd.Flags().String("severity", "high", "vulnerability severity level")
	scanCmd.AddCommand(scanVulnerabilitiesCmd)
	securityCmd.AddCommand(scanCmd)

	checkComplianceCmd.Flags().String("standard", "", "compliance standard to check")
	checkCmd.AddCommand(checkComplianceCmd)
	securityCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(securityCmd)

	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Load environment variables
	viper.SetEnvPrefix("SREDIAG")
	viper.AutomaticEnv()

	// Bind environment variables
	bindEnvs := []struct {
		key      string
		envVar   string
		required bool
	}{
		{"config", "SREDIAG_CONFIG", false},
		{"format", "SREDIAG_OUTPUT_FORMAT", false},
		{"log_level", "SREDIAG_LOG_LEVEL", false},
		{"api_key", "SREDIAG_API_KEY", false},
	}

	for _, env := range bindEnvs {
		if err := viper.BindEnv(env.key, env.envVar); err != nil {
			if env.required {
				fmt.Fprintf(os.Stderr, "Error binding environment variable %s: %v\n", env.envVar, err)
				os.Exit(ExitConfigError)
			}
			fmt.Fprintf(os.Stderr, "Warning: failed to bind environment variable %s: %v\n", env.envVar, err)
		}
	}
}

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Search for config in default locations
		viper.AddConfigPath("/etc/srediag/config")
		viper.AddConfigPath("$HOME/.srediag")
		viper.AddConfigPath(".")
		viper.SetConfigName("srediag")
		viper.SetConfigType("yaml")
	}

	// Set default values
	cfg := config.DefaultConfig()
	if err := viper.MergeConfigMap(map[string]interface{}{
		"core":       cfg.Core,
		"service":    cfg.Service,
		"telemetry":  cfg.Telemetry,
		"collector":  cfg.Collector,
		"diagnostic": cfg.Diagnostic,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting default config: %v\n", err)
		os.Exit(ExitConfigError)
	}

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(ExitConfigError)
		}
	}
}

// initLogger initializes the logger with the given configuration
func initLogger(cfg *config.Root) (*zap.Logger, error) {
	var config zap.Config

	// Configure logging format
	if cfg.Core.LogFormat == "console" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	// Set log level
	level := cfg.Core.LogLevel
	if verbose {
		level = "debug"
	} else if quiet {
		level = "error"
	}

	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}
	config.Level = zap.NewAtomicLevelAt(parsedLevel)

	// Create logger
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger.Named("srediag"), nil
}

func run(cmd *cobra.Command, args []string) error {
	// Print version information if requested
	if version {
		fmt.Printf("SREDIAG %s (commit: %s, built: %s)\n", Version, GitCommit, BuildDate)
		return nil
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Update version information
	cfg.Service.Version = Version

	// Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
		}
	}()

	// Create and initialize SREDIAG
	srediag, err := app.NewSREDiag(logger, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize SREDIAG: %w", err)
	}

	// Setup signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start the application
	if err := srediag.Start(ctx); err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	// Wait for interrupt signal
	<-ctx.Done()

	// Stop the application
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := srediag.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("error during shutdown: %w", err)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitGeneralError)
	}
}

// Command implementations
func diagnoseSystem(cmd *cobra.Command, args []string) error {
	resource, _ := cmd.Flags().GetString("resource")
	logger := getLogger()
	logger.Info("running system diagnostics", zap.String("resource", resource))
	return fmt.Errorf("not implemented")
}

func diagnoseKubernetes(cmd *cobra.Command, args []string) error {
	cluster, _ := cmd.Flags().GetString("cluster")
	logger := getLogger()
	logger.Info("running kubernetes diagnostics", zap.String("cluster", cluster))
	return fmt.Errorf("not implemented")
}

func analyzeProcess(cmd *cobra.Command, args []string) error {
	pid, _ := cmd.Flags().GetInt("pid")
	logger := getLogger()
	logger.Info("analyzing process", zap.Int("pid", pid))
	return fmt.Errorf("not implemented")
}

func analyzeMemory(cmd *cobra.Command, args []string) error {
	threshold, _ := cmd.Flags().GetFloat64("threshold")
	logger := getLogger()
	logger.Info("analyzing memory", zap.Float64("threshold", threshold))
	return fmt.Errorf("not implemented")
}

func monitorSystem(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	logger := getLogger()
	logger.Info("monitoring system", zap.Duration("interval", interval))
	return fmt.Errorf("not implemented")
}

func scanVulnerabilities(cmd *cobra.Command, args []string) error {
	severity, _ := cmd.Flags().GetString("severity")
	logger := getLogger()
	logger.Info("scanning vulnerabilities", zap.String("severity", severity))
	return fmt.Errorf("not implemented")
}

func checkCompliance(cmd *cobra.Command, args []string) error {
	standard, _ := cmd.Flags().GetString("standard")
	logger := getLogger()
	logger.Info("checking compliance", zap.String("standard", standard))
	return fmt.Errorf("not implemented")
}

// Helper functions
func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(ExitGeneralError)
	}
	return logger
}
