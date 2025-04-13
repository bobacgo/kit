package otel

type Config struct {
	Tracer TraceConfig `mapstructure:"tracer" yaml:"tracer"`
}
