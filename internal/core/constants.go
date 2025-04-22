package core

const (
	// DefaultMetricInterval é o intervalo padrão para coleta de métricas
	DefaultMetricInterval = 10 // segundos

	// DefaultSamplingRate é a taxa de amostragem padrão para traces
	DefaultSamplingRate = 1.0

	// SamplingTypeAlwaysOn indica que todos os traces devem ser coletados
	SamplingTypeAlwaysOn = "always_on"

	// SamplingTypeAlwaysOff indica que nenhum trace deve ser coletado
	SamplingTypeAlwaysOff = "always_off"

	// SamplingTypeProbabilistic indica que os traces devem ser coletados probabilisticamente
	SamplingTypeProbabilistic = "probabilistic"

	// DefaultPluginPattern é o padrão para arquivos de plugin
	DefaultPluginPattern = "*.so"

	// DefaultPluginSymbol é o símbolo que deve ser exportado pelos plugins
	DefaultPluginSymbol = "Plugin"
)
