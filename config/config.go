// Copyright (C) 2024, AllianceBlock. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/profiler"
	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/config"
	"github.com/ava-labs/hypersdk/gossiper"
	"github.com/ava-labs/hypersdk/trace"
	"github.com/ava-labs/hypersdk/vm"

	nconsts "github.com/nuklai/nuklaivm/consts"
	"github.com/nuklai/nuklaivm/version"
)

var _ vm.Config = (*Config)(nil)

const (
	defaultContinuousProfilerFrequency = 1 * time.Minute
	defaultContinuousProfilerMaxFiles  = 10
	defaultStoreTransactions           = true
)

type Config struct {
	*config.Config

	// Concurrency
	AuthVerificationCores     int `json:"authVerificationCores"`
	RootGenerationCores       int `json:"rootGenerationCores"`
	TransactionExecutionCores int `json:"transactionExecutionCores"`

	// Gossip
	GossipMaxSize       int   `json:"gossipMaxSize"`
	GossipProposerDiff  int   `json:"gossipProposerDiff"`
	GossipProposerDepth int   `json:"gossipProposerDepth"`
	NoGossipBuilderDiff int   `json:"noGossipBuilderDiff"`
	VerifyTimeout       int64 `json:"verifyTimeout"`

	// Tracing
	TraceEnabled    bool    `json:"traceEnabled"`
	TraceSampleRate float64 `json:"traceSampleRate"`

	// Profiling
	ContinuousProfilerDir string `json:"continuousProfilerDir"` // "*" is replaced with rand int

	// Streaming settings
	StreamingBacklogSize int `json:"streamingBacklogSize"`

	// Mempool
	MempoolSize           int      `json:"mempoolSize"`
	MempoolSponsorSize    int      `json:"mempoolSponsorSize"`
	MempoolExemptSponsors []string `json:"mempoolExemptSponsors"`

	// Misc
	VerifyAuth        bool          `json:"verifyAuth"`
	StoreTransactions bool          `json:"storeTransactions"`
	TestMode          bool          `json:"testMode"` // makes gossip/building manual
	LogLevel          logging.Level `json:"logLevel"`

	// State Sync
	StateSyncServerDelay time.Duration `json:"stateSyncServerDelay"` // for testing

	loaded               bool
	nodeID               ids.NodeID
	parsedExemptSponsors []codec.Address
}

func New(nodeID ids.NodeID, b []byte) (*Config, error) {
	c := &Config{nodeID: nodeID}
	c.setDefault()
	if len(b) > 0 {
		if err := json.Unmarshal(b, c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config %s: %w", string(b), err)
		}
		c.loaded = true
	}

	// Parse any exempt sponsors (usually used when a single account is
	// broadcasting many txs at once)
	c.parsedExemptSponsors = make([]codec.Address, len(c.MempoolExemptSponsors))
	for i, sponsor := range c.MempoolExemptSponsors {
		p, err := codec.ParseAddressBech32(nconsts.HRP, sponsor)
		if err != nil {
			return nil, err
		}
		c.parsedExemptSponsors[i] = p
	}
	return c, nil
}

func (c *Config) setDefault() {
	c.LogLevel = c.Config.GetLogLevel()
	gcfg := gossiper.DefaultProposerConfig()
	c.GossipMaxSize = gcfg.GossipMaxSize
	c.GossipProposerDiff = gcfg.GossipProposerDiff
	c.GossipProposerDepth = gcfg.GossipProposerDepth
	c.NoGossipBuilderDiff = gcfg.NoGossipBuilderDiff
	c.VerifyTimeout = gcfg.VerifyTimeout
	c.AuthVerificationCores = c.Config.GetAuthVerificationCores()
	c.RootGenerationCores = c.Config.GetRootGenerationCores()
	c.TransactionExecutionCores = c.Config.GetTransactionExecutionCores()
	c.MempoolSize = c.Config.GetMempoolSize()
	c.MempoolSponsorSize = c.Config.GetMempoolSponsorSize()
	c.StateSyncServerDelay = c.Config.GetStateSyncServerDelay()
	c.StreamingBacklogSize = c.Config.GetStreamingBacklogSize()
	c.VerifyAuth = c.Config.GetVerifyAuth()
	c.StoreTransactions = defaultStoreTransactions
}

func (c *Config) GetLogLevel() logging.Level                { return c.LogLevel }
func (c *Config) GetTestMode() bool                         { return c.TestMode }
func (c *Config) GetAuthVerificationCores() int             { return c.AuthVerificationCores }
func (c *Config) GetRootGenerationCores() int               { return c.RootGenerationCores }
func (c *Config) GetTransactionExecutionCores() int         { return c.TransactionExecutionCores }
func (c *Config) GetMempoolSize() int                       { return c.MempoolSize }
func (c *Config) GetMempoolSponsorSize() int                { return c.MempoolSponsorSize }
func (c *Config) GetMempoolExemptSponsors() []codec.Address { return c.parsedExemptSponsors }
func (c *Config) GetTraceConfig() *trace.Config {
	return &trace.Config{
		Enabled:         c.TraceEnabled,
		TraceSampleRate: c.TraceSampleRate,
		AppName:         nconsts.Name,
		Agent:           c.nodeID.String(),
		Version:         version.Version.String(),
	}
}
func (c *Config) GetStateSyncServerDelay() time.Duration { return c.StateSyncServerDelay }
func (c *Config) GetStreamingBacklogSize() int           { return c.StreamingBacklogSize }
func (c *Config) GetContinuousProfilerConfig() *profiler.Config {
	if len(c.ContinuousProfilerDir) == 0 {
		return &profiler.Config{Enabled: false}
	}
	// Replace all instances of "*" with nodeID. This is useful when
	// running multiple instances of nuklaivm on the same machine.
	c.ContinuousProfilerDir = strings.ReplaceAll(c.ContinuousProfilerDir, "*", c.nodeID.String())
	return &profiler.Config{
		Enabled:     true,
		Dir:         c.ContinuousProfilerDir,
		Freq:        defaultContinuousProfilerFrequency,
		MaxNumFiles: defaultContinuousProfilerMaxFiles,
	}
}
func (c *Config) GetVerifyAuth() bool        { return c.VerifyAuth }
func (c *Config) GetStoreTransactions() bool { return c.StoreTransactions }
func (c *Config) Loaded() bool               { return c.loaded }
