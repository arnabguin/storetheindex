package config

import (
	"time"
)

// Discovery holds addresses of peers to from which to receive index
// advertisements, which peers to allow and block, and time related settings
// for provider discovery
type Discovery struct {
	// FilterIPs, when true, removes any private, loopback, or unspecified IP
	// addresses from provider and publisher addresses.
	FilterIPs bool
	// LotusGateway is the host or host:port for a lotus gateway used to
	// verify providers on the blockchain.
	LotusGateway string
	// Policy configures which providers are allowed and blocked, rate-limited,
	// and allow to publish on behalf of others.
	Policy Policy
	// PollInterval is the amount of time to wait without getting any updates
	// for a provider, before sending a request for the latest advertisement.
	// If there is no response after at least one poll attempt, then the
	// provider is considered inactive and is not returned in find responses.
	// Values are a number ending in "s", "m", "h" for seconds. minutes, hours.
	PollInterval Duration
	// PollRetryAfter is the amount of time from one poll attempt, without a
	// response, to the next poll attempt, and is also the time between checks
	// for providers to poll. This value must be smaller than PollStopAfter for
	// there to be more than one poll attempt for a provider.
	PollRetryAfter Duration
	// PollStopAfter is the amount of time, from the start of polling, to
	// continuing polling for the latest advertisment without getting a
	// response. After this time elapses with no updates from the provider, the
	// provider's data is removed from the indexer and must be re-fetched if
	// the provider returns.
	PollStopAfter Duration
	// PollOverrides configures polling for specific providers.
	PollOverrides []Polling
	// RediscoverWait is the amount of time that must pass before a provider
	// can be discovered following a previous discovery attempt. A value of 0
	// means there is no wait time.
	RediscoverWait Duration
	// Timeout is the maximum amount of time that the indexer will spend trying
	// to discover and verify a new provider.
	Timeout Duration
}

// Polling is a set of polling parameters that is applied to a specific
// provider. The values override the matching Poll values in the Discovery
// config.
type Polling struct {
	// ProviderID identifies the provider that this override applies to.
	ProviderID string
	// Interval overrides Discovery.PollInterval.
	Interval Duration
	// RetryAfter overrides Discovery.PollRetryAfter.
	RetryAfter Duration
	// StopAfter overrides Discovery.PollStopAfter.
	StopAfter Duration
}

// NewDiscovery returns Discovery with values set to their defaults.
func NewDiscovery() Discovery {
	return Discovery{
		LotusGateway:   "https://api.chain.love",
		Policy:         NewPolicy(),
		PollInterval:   Duration(24 * time.Hour),
		PollRetryAfter: Duration(5 * time.Hour),
		PollStopAfter:  Duration(7 * 24 * time.Hour),
		RediscoverWait: Duration(5 * time.Minute),
		Timeout:        Duration(2 * time.Minute),
	}
}

// populateUnset replaces zero-values in the config with default values.
func (c *Discovery) populateUnset() {
	def := NewDiscovery()

	if c.PollInterval == 0 {
		c.PollInterval = def.PollInterval
	}
	if c.PollRetryAfter == 0 {
		c.PollRetryAfter = def.PollRetryAfter
	}
	if c.PollStopAfter == 0 {
		c.PollStopAfter = def.PollStopAfter
	}
	if c.Timeout == 0 {
		c.Timeout = def.Timeout
	}
}
