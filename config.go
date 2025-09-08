package thebus

import (
	"time"
)

const (
	DefaultTopicQueueSize = 1024
	DefaultSubBufferSize  = 128
	DefaultSendTimeout    = 200 * time.Millisecond
)

// ##############################################################################
// ################################   CONFIG   ##################################
// ##############################################################################

// Config is the main configuration for thebus
type Config struct {
	// Topics / queues
	TopicQueueSize        int           // default: 1024
	AutoDeleteEmptyTopics bool          // default: true
	TopicIdleTTL          time.Duration // if Janitor enabled (0 = off)
	JanitorInterval       time.Duration // 0 = off
	IDGenerator           IDGenerator   // default to DefaultIDGenerator
	CopyOnPublish         bool          // default false

	// Default for subscribers (Can be overridden by sub)
	DefaultSubBufferSize int                  // default: 128
	DefaultSendTimeout   time.Duration        // default: 200ms
	DefaultDropIfFull    bool                 // default: true
	DefaultStrategy      SubscriptionStrategy // default: SubscriptionStrategyPayloadShared

	// Limits (0 = unlimited)
	MaxTopics              int // default : 0 (unlimited)
	MaxSubscribersPerTopic int // default : 0 (unlimited)

	// Observability
	Logger       Logger                    // default: nopLogger
	Metrics      MetricsHooks              // default: nopMetrics
	PanicHandler func(topic string, v any) // default: nil (no recover)
}

func (cfg *Config) Normalize() *Config {
	if cfg.TopicQueueSize <= 0 {
		cfg.TopicQueueSize = DefaultTopicQueueSize
	}
	if cfg.TopicIdleTTL < 0 {
		cfg.TopicIdleTTL = 0
	}
	if cfg.JanitorInterval < 0 {
		cfg.JanitorInterval = 0
	}
	if cfg.IDGenerator == nil {
		cfg.IDGenerator = DefaultIDGenerator
	}
	if cfg.DefaultSubBufferSize <= 0 {
		cfg.DefaultSubBufferSize = DefaultSubBufferSize
	}
	if cfg.DefaultSendTimeout <= 0 {
		cfg.DefaultSendTimeout = DefaultSendTimeout
	}
	if !cfg.DefaultStrategy.IsValid() {
		cfg.DefaultStrategy = SubscriptionStrategyPayloadShared
	}
	if cfg.MaxTopics < 0 {
		cfg.MaxTopics = 0
	}
	if cfg.MaxSubscribersPerTopic < 0 {
		cfg.MaxSubscribersPerTopic = 0
	}
	if cfg.Logger == nil {
		cfg.Logger = &noopLogger{}
	}
	return cfg
}

func DefaultConfig() *Config {
	return &Config{
		TopicQueueSize:        DefaultTopicQueueSize,
		AutoDeleteEmptyTopics: true,
		TopicIdleTTL:          0 * time.Second,
		JanitorInterval:       0 * time.Second,
		DefaultSubBufferSize:  DefaultSubBufferSize,
		DefaultSendTimeout:    DefaultSendTimeout,
		DefaultDropIfFull:     true,
		DefaultStrategy:       SubscriptionStrategyPayloadShared,
		IDGenerator:           DefaultIDGenerator,
		Logger:                NoopLogger(),
	}
}

// ##############################################################################
// ###############################   OPTIONS   ##################################
// ##############################################################################

type Option func(*Config)

func WithTopicQueueSize(size int) Option {
	return func(cfg *Config) {
		cfg.TopicQueueSize = size
	}
}

func WithAutoDeleteEmptyTopics(b bool) Option {
	return func(cfg *Config) {
		cfg.AutoDeleteEmptyTopics = b
	}
}

func WithTopicIdleTTL(ttl time.Duration) Option {
	return func(cfg *Config) {
		cfg.TopicIdleTTL = ttl
	}
}

func WithJanitorInterval(interval time.Duration) Option {
	return func(cfg *Config) {
		cfg.JanitorInterval = interval
	}
}

func WithDefaultSubBufferSize(size int) Option {
	return func(cfg *Config) {
		cfg.DefaultSubBufferSize = size
	}
}

func WithDefaultSendTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.DefaultSendTimeout = timeout
	}
}

func WithDefaultDropIfFull() Option {
	return func(cfg *Config) {
		cfg.DefaultDropIfFull = true
	}
}

func WithDefaultStrategy(strategy SubscriptionStrategy) Option {
	return func(cfg *Config) {
		cfg.DefaultStrategy = strategy
	}
}

func WithMaxTopics(max int) Option {
	return func(cfg *Config) {
		cfg.MaxTopics = max
	}
}

func WithMaxSubscribersPerTopic(max int) Option {
	return func(cfg *Config) {
		cfg.MaxSubscribersPerTopic = max
	}
}

func WithLogger(logger Logger) Option {
	return func(cfg *Config) {
		cfg.Logger = logger
	}
}

func WithMetrics(metrics MetricsHooks) Option {
	return func(cfg *Config) {
		cfg.Metrics = metrics
	}
}

func WithPanicHandler(handler func(topic string, v any)) Option {
	return func(cfg *Config) {
		cfg.PanicHandler = handler
	}
}

func WithIDGenerator(IDGenerator IDGenerator) Option {
	return func(cfg *Config) {
		cfg.IDGenerator = IDGenerator
	}
}

func WithCopyOnPublish(b bool) Option {
	return func(cfg *Config) {
		cfg.CopyOnPublish = b
	}
}

func BuildConfig(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
