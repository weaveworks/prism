package grpcclient

import (
	"flag"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/cortexproject/cortex/pkg/util/tls"
)

// Config for a gRPC client.
type Config struct {
	MaxRecvMsgSize     int     `yaml:"max_recv_msg_size"`
	MaxSendMsgSize     int     `yaml:"max_send_msg_size"`
	UseGzipCompression bool    `yaml:"use_gzip_compression"`
	Compression        string  `yaml:"compression"`
	RateLimit          float64 `yaml:"rate_limit"`
	RateLimitBurst     int     `yaml:"rate_limit_burst"`

	BackoffOnRatelimits bool               `yaml:"backoff_on_ratelimits"`
	BackoffConfig       util.BackoffConfig `yaml:"backoff_config"`

	prefix string `yaml:"-"`
}

// RegisterFlags registers flags.
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.RegisterFlagsWithPrefix("", f)
}

// RegisterFlagsWithPrefix registers flags with prefix.
func (cfg *Config) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	cfg.prefix = prefix
	f.IntVar(&cfg.MaxRecvMsgSize, prefix+".grpc-max-recv-msg-size", 100<<20, "gRPC client max receive message size (bytes).")
	f.IntVar(&cfg.MaxSendMsgSize, prefix+".grpc-max-send-msg-size", 16<<20, "gRPC client max send message size (bytes).")
	f.BoolVar(&cfg.UseGzipCompression, prefix+".grpc-use-gzip-compression", false, "Deprecated: Use gzip compression when sending messages.")
	f.StringVar(&cfg.Compression, prefix+".grpc-compression", "", "Use compression when sending messages. Supported values are: 'gzip', 'snappy' and '' (disable compression)")
	f.Float64Var(&cfg.RateLimit, prefix+".grpc-client-rate-limit", 0., "Rate limit for gRPC client; 0 means disabled.")
	f.IntVar(&cfg.RateLimitBurst, prefix+".grpc-client-rate-limit-burst", 0, "Rate limit burst for gRPC client.")
	f.BoolVar(&cfg.BackoffOnRatelimits, prefix+".backoff-on-ratelimits", false, "Enable backoff and retry when we hit ratelimits.")

	cfg.BackoffConfig.RegisterFlags(prefix, f)
}

func (cfg *Config) Validate(log log.Logger) error {
	if cfg.UseGzipCompression {
		flagext.DeprecatedFlagsUsed.Inc()
		level.Warn(log).Log("msg", fmt.Sprintf("running with DEPRECATED flag -%s.grpc-use-gzip-compression, use -%s.grpc-use-compression instead.", cfg.prefix, cfg.prefix))
	}
	switch cfg.Compression {
	case "gzip", "snappy", "":
		// valid
	default:
		return errors.Errorf("unsupported compression type: %s", cfg.Compression)
	}
	return nil
}

// CallOptions returns the config in terms of CallOptions.
func (cfg *Config) CallOptions() []grpc.CallOption {
	var opts []grpc.CallOption
	opts = append(opts, grpc.MaxCallRecvMsgSize(cfg.MaxRecvMsgSize))
	opts = append(opts, grpc.MaxCallSendMsgSize(cfg.MaxSendMsgSize))
	compression := cfg.Compression
	if cfg.UseGzipCompression {
		compression = "gzip"
	}
	if compression != "" {
		opts = append(opts, grpc.UseCompressor(compression))
	}
	return opts
}

// DialOption returns the config as a grpc.DialOptions.
func (cfg *Config) DialOption(unaryClientInterceptors []grpc.UnaryClientInterceptor, streamClientInterceptors []grpc.StreamClientInterceptor) []grpc.DialOption {
	if cfg.BackoffOnRatelimits {
		unaryClientInterceptors = append([]grpc.UnaryClientInterceptor{NewBackoffRetry(cfg.BackoffConfig)}, unaryClientInterceptors...)
	}

	if cfg.RateLimit > 0 {
		unaryClientInterceptors = append([]grpc.UnaryClientInterceptor{NewRateLimiter(cfg)}, unaryClientInterceptors...)
	}

	return []grpc.DialOption{
		grpc.WithDefaultCallOptions(cfg.CallOptions()...),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryClientInterceptors...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamClientInterceptors...)),
	}
}

// ConfigWithTLS is the config for a grpc client with tls
type ConfigWithTLS struct {
	GRPC Config           `yaml:",inline"`
	TLS  tls.ClientConfig `yaml:",inline"`
}

// RegisterFlagsWithPrefix registers flags with prefix.
func (cfg *ConfigWithTLS) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	cfg.GRPC.RegisterFlagsWithPrefix(prefix, f)
	cfg.TLS.RegisterFlagsWithPrefix(prefix, f)
}

func (cfg *ConfigWithTLS) Validate(log log.Logger) error {
	return cfg.GRPC.Validate(log)
}

// DialOption returns the config as a grpc.DialOptions
func (cfg *ConfigWithTLS) DialOption(unaryClientInterceptors []grpc.UnaryClientInterceptor, streamClientInterceptors []grpc.StreamClientInterceptor) ([]grpc.DialOption, error) {
	opts, err := cfg.TLS.GetGRPCDialOptions()
	if err != nil {
		return nil, err
	}

	return append(opts, cfg.GRPC.DialOption(unaryClientInterceptors, streamClientInterceptors)...), nil
}
