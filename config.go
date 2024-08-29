package posthog

import (
	"net/http"
	"time"
)

// Instances of this type carry the different configuration options that may
// be set when instantiating a client.
//
// Each field's zero-value is either meaningful or interpreted as using the
// default value defined by the library.
type Config struct {

	// The endpoint to which the client connect and send their messages, set to
	// `DefaultEndpoint` by default.
	Endpoint string

	// Specifying a Personal API key will make feature flag evaluation more performant,
	// but it's not required for feature flags.  If you don't have a personal API key,
	// you can leave this field empty, and all of the relevant feature flag evaluation
	// methods will still work.
	// Information on how to get a personal API key: https://posthog.com/docs/api/overview
	PersonalApiKey string

	// The flushing interval of the client. Messages will be sent when they've
	// been queued up to the maximum batch size or when the flushing interval
	// timer triggers.
	Interval time.Duration

	// Interval at which to fetch new feature flag definitions, 5min by default
	DefaultFeatureFlagsPollingInterval time.Duration

	// Timeout for fetching feature flags, 3 seconds by default
	FeatureFlagRequestTimeout time.Duration

	// Calculate when feature flag definitions should be polled next. Setting this property
	// will override DefaultFeatureFlagsPollingInterval.
	NextFeatureFlagsPollingTick func() time.Duration

	// The HTTP transport used by the client, this allows an application to
	// redefine how requests are being sent at the HTTP level (for example,
	// to change the connection pooling policy).
	// If none is specified the client uses `http.DefaultTransport`.
	Transport http.RoundTripper

	// The logger used by the client to output info or error messages when that
	// are generated by background operations.
	// If none is specified the client uses a standard logger that outputs to
	// `os.Stderr`.
	Logger Logger

	// Properties that will be included in every event sent by the client.
	// This is useful for adding common metadata like service name or app version across all events.
	// If a property conflict occurs, the value from DefaultEventProperties will overwrite any existing value.
	DefaultEventProperties Properties

	// The callback object that will be used by the client to notify the
	// application when messages sends to the backend API succeeded or failed.
	Callback Callback

	// The maximum number of messages that will be sent in one API call.
	// Messages will be sent when they've been queued up to the maximum batch
	// size or when the flushing interval timer triggers.
	// Note that the API will still enforce a 500KB limit on each HTTP request
	// which is independent from the number of embedded messages.
	BatchSize int

	// When set to true the client will send more frequent and detailed messages
	// to its logger.
	Verbose bool

	// The retry policy used by the client to resend requests that have failed.
	// The function is called with how many times the operation has been retried
	// and is expected to return how long the client should wait before trying
	// again.
	// If not set the client will fallback to use a default retry policy.
	RetryAfter func(int) time.Duration

	// A function called by the client to get the current time, `time.Now` is
	// used by default.
	// This field is not exported and only exposed internally to let unit tests
	// mock the current time.
	now func() time.Time

	// The maximum number of goroutines that will be spawned by a client to send
	// requests to the backend API.
	// This field is not exported and only exposed internally to let unit tests
	// mock the current time.
	maxConcurrentRequests int
}

// This constant sets the default endpoint to which client instances send
// messages if none was explictly set.
const DefaultEndpoint = "https://app.posthog.com"

// This constant sets the default flush interval used by client instances if
// none was explicitly set.
const DefaultInterval = 5 * time.Second

// Specifies the default interval at which to fetch new feature flags
const DefaultFeatureFlagsPollingInterval = 5 * time.Minute

// Specifies the default timeout for fetching feature flags
const DefaultFeatureFlagRequestTimeout = 3 * time.Second

// This constant sets the default batch size used by client instances if none
// was explicitly set.
const DefaultBatchSize = 250

// Verifies that fields that don't have zero-values are set to valid values,
// returns an error describing the problem if a field was invalid.
func (c *Config) validate() error {
	if c.Interval < 0 {
		return ConfigError{
			Reason: "negative time intervals are not supported",
			Field:  "Interval",
			Value:  c.Interval,
		}
	}

	if c.BatchSize < 0 {
		return ConfigError{
			Reason: "negative batch sizes are not supported",
			Field:  "BatchSize",
			Value:  c.BatchSize,
		}
	}

	return nil
}

// Given a config object as argument the function will set all zero-values to
// their defaults and return the modified object.
func makeConfig(c Config) Config {
	if len(c.Endpoint) == 0 {
		c.Endpoint = DefaultEndpoint
	}

	if c.Interval == 0 {
		c.Interval = DefaultInterval
	}

	if c.DefaultFeatureFlagsPollingInterval == 0 {
		c.DefaultFeatureFlagsPollingInterval = DefaultFeatureFlagsPollingInterval
	}

	if c.FeatureFlagRequestTimeout == 0 {
		c.FeatureFlagRequestTimeout = DefaultFeatureFlagRequestTimeout
	}

	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}

	if c.Logger == nil {
		c.Logger = newDefaultLogger()
	}

	if c.BatchSize == 0 {
		c.BatchSize = DefaultBatchSize
	}

	if c.RetryAfter == nil {
		c.RetryAfter = DefaultBacko().Duration
	}

	if c.now == nil {
		c.now = time.Now
	}

	if c.maxConcurrentRequests == 0 {
		c.maxConcurrentRequests = 1000
	}

	return c
}
