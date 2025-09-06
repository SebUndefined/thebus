package thebus

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	DefaultTopicQueueSize = 1024
)

// ##############################################################################
// ##################################   ENUM   ##################################
// ##############################################################################

// SubscriptionStrategy define if the payload must be
// shared by the subscribers (SubscriptionStrategyPayloadShared, the default)
// or copied (SubscriptionStrategyPayloadClonedPerSubscriber)
type SubscriptionStrategy string

const (
	SubscriptionStrategyUnknown                    SubscriptionStrategy = "UNKNOWN"
	SubscriptionStrategyPayloadShared              SubscriptionStrategy = "PAYLOAD_SHARED"
	SubscriptionStrategyPayloadClonedPerSubscriber SubscriptionStrategy = "PAYLOAD_CLONED_PER_SUBSCRIBER"
)

func (enum SubscriptionStrategy) String() string {
	if len(strings.TrimSpace(string(enum))) == 0 {
		return string(SubscriptionStrategyUnknown)
	}
	return string(enum)
}

func SubscriptionStrategyValues() []SubscriptionStrategy {
	return []SubscriptionStrategy{
		SubscriptionStrategyPayloadShared,
		SubscriptionStrategyPayloadClonedPerSubscriber,
	}
}

func (enum SubscriptionStrategy) IsValid() bool {
	if slices.Contains(SubscriptionStrategyValues(), enum) {
		return true
	}
	return false
}

func (enum SubscriptionStrategy) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, enum)), nil
}

func (enum *SubscriptionStrategy) UnmarshalJSON(data []byte) error {
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	fs := SubscriptionStrategy(tmp)
	if !fs.IsValid() {
		fs = SubscriptionStrategyUnknown
	}
	*enum = fs
	return nil
}

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

	// Default for subscribers (peuvent être override par Subscribe options)
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

func DefaultConfig() *Config {
	return &Config{
		TopicQueueSize:        DefaultTopicQueueSize,
		AutoDeleteEmptyTopics: true,
		TopicIdleTTL:          0 * time.Second,
		JanitorInterval:       0 * time.Second,
		DefaultSubBufferSize:  128,
		DefaultSendTimeout:    200 * time.Millisecond,
		DefaultDropIfFull:     true,
		DefaultStrategy:       SubscriptionStrategyPayloadShared,
		IDGenerator:           DefaultIDGenerator,
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

func WithAutoDeleteEmptyTopics() Option {
	return func(cfg *Config) {
		cfg.AutoDeleteEmptyTopics = true
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

func (cfg *Config) Validate() error {
	if cfg.IDGenerator == nil {
		return ErrIDGeneratorNotSet
	}
	return nil
}

func BuildConfig(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
