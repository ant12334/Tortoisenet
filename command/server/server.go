package server

import (
	"fmt"

	"github.com/0xPolygon/polygon-edge/command"
	"github.com/0xPolygon/polygon-edge/command/helper"
	"github.com/0xPolygon/polygon-edge/command/server/config"
	"github.com/0xPolygon/polygon-edge/command/server/export"
	"github.com/0xPolygon/polygon-edge/server"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:     "server",
		Short:   "The default command that starts the Polygon Edge client, by bootstrapping all modules together",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	helper.RegisterGRPCAddressFlag(serverCmd)
	helper.RegisterLegacyGRPCAddressFlag(serverCmd)
	helper.RegisterJSONRPCFlag(serverCmd)

	registerSubcommands(serverCmd)
	setFlags(serverCmd)

	return serverCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// server export
		export.GetCommand(),
	)
}

func setFlags(cmd *cobra.Command) {
	defaultConfig := config.DefaultConfig()

	cmd.Flags().StringVar(
		&params.rawConfig.LogLevel,
		command.LogLevelFlag,
		defaultConfig.LogLevel,
		"the log level for console output",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.GenesisPath,
		genesisPathFlag,
		defaultConfig.GenesisPath,
		"the genesis file used for starting the chain",
	)

	cmd.Flags().StringVar(
		&params.configPath,
		configFlag,
		"",
		"the path to the CLI config. Supports .json, .hcl, .yaml, .yml",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.DataDir,
		dataDirFlag,
		defaultConfig.DataDir,
		"the data directory used for storing Polygon Edge client data",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Network.Libp2pAddr,
		libp2pAddressFlag,
		defaultConfig.Network.Libp2pAddr,
		"the address and port for the libp2p service",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Telemetry.PrometheusAddr,
		prometheusAddressFlag,
		"",
		"the address and port for the prometheus instrumentation service (address:port). "+
			"If only port is defined (:port) it will bind to 0.0.0.0:port",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Network.NatAddr,
		natFlag,
		"",
		"the external IP address without port, as can be seen by peers",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.Network.DNSAddr,
		dnsFlag,
		"",
		"the host DNS address which can be used by a remote peer for connection",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.BlockGasTarget,
		blockGasTargetFlag,
		defaultConfig.BlockGasTarget,
		"the target block gas limit for the chain. If omitted, the value of the parent block is used",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.SecretsConfigPath,
		secretsConfigFlag,
		"",
		"the path to the SecretsManager config file. Used for Hashicorp Vault. "+
			"If omitted, the local FS secrets manager is used",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.RestoreFile,
		restoreFlag,
		"",
		"the path to the archive blockchain data to restore on initialization",
	)

	cmd.Flags().BoolVar(
		&params.rawConfig.ShouldSeal,
		sealFlag,
		defaultConfig.ShouldSeal,
		"the flag indicating that the client should seal blocks",
	)

	cmd.Flags().BoolVar(
		&params.rawConfig.Network.NoDiscover,
		command.NoDiscoverFlag,
		defaultConfig.Network.NoDiscover,
		"prevent the client from discovering other peers",
	)

	cmd.Flags().Int64Var(
		&params.rawConfig.Network.MaxPeers,
		maxPeersFlag,
		-1,
		"the client's max number of peers allowed",
	)
	// override default usage value
	cmd.Flag(maxPeersFlag).DefValue = fmt.Sprintf("%d", defaultConfig.Network.MaxPeers)

	cmd.Flags().Int64Var(
		&params.rawConfig.Network.MaxInboundPeers,
		maxInboundPeersFlag,
		-1,
		"the client's max number of inbound peers allowed",
	)
	// override default usage value
	cmd.Flag(maxInboundPeersFlag).DefValue = fmt.Sprintf("%d", defaultConfig.Network.MaxInboundPeers)
	cmd.MarkFlagsMutuallyExclusive(maxPeersFlag, maxInboundPeersFlag)

	cmd.Flags().Int64Var(
		&params.rawConfig.Network.MaxOutboundPeers,
		maxOutboundPeersFlag,
		-1,
		"the client's max number of outbound peers allowed",
	)
	// override default usage value
	cmd.Flag(maxOutboundPeersFlag).DefValue = fmt.Sprintf("%d", defaultConfig.Network.MaxOutboundPeers)
	cmd.MarkFlagsMutuallyExclusive(maxPeersFlag, maxOutboundPeersFlag)

	cmd.Flags().IntVar(
		&params.rawConfig.Network.GossipMessageSize,
		gossipMessageSizeFlag,
		pubsub.DefaultMaxMessageSize,
		"the maximum size of a gossip message",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.PriceLimit,
		priceLimitFlag,
		defaultConfig.TxPool.PriceLimit,
		fmt.Sprintf(
			"the minimum gas price limit to enforce for acceptance into the pool (default %d)",
			defaultConfig.TxPool.PriceLimit,
		),
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.MaxSlots,
		maxSlotsFlag,
		defaultConfig.TxPool.MaxSlots,
		"maximum slots in the pool",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.MaxAccountEnqueued,
		maxEnqueuedFlag,
		defaultConfig.TxPool.MaxAccountEnqueued,
		"maximum number of enqueued transactions per account",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.TxPool.TxGossipBatchSize,
		txGossipBatchSizeFlag,
		defaultConfig.TxPool.TxGossipBatchSize,
		"maximum number of transactions in gossip message",
	)

	cmd.Flags().StringArrayVar(
		&params.rawConfig.CorsAllowedOrigins,
		corsOriginFlag,
		defaultConfig.Headers.AccessControlAllowOrigins,
		"the CORS header indicating whether any JSON-RPC response can be shared with the specified origin",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.JSONRPCBatchRequestLimit,
		jsonRPCBatchRequestLimitFlag,
		defaultConfig.JSONRPCBatchRequestLimit,
		"max length to be considered when handling json-rpc batch requests, value of 0 disables it",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.JSONRPCBlockRangeLimit,
		jsonRPCBlockRangeLimitFlag,
		defaultConfig.JSONRPCBlockRangeLimit,
		"max block range to be considered when executing json-rpc requests "+
			"that consider fromBlock/toBlock values (e.g. eth_getLogs), value of 0 disables it",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.LogFilePath,
		logFileLocationFlag,
		defaultConfig.LogFilePath,
		"write all logs to the file at specified location instead of writing them to console",
	)

	cmd.Flags().BoolVar(
		&params.rawConfig.UseTLS,
		useTLSFlag,
		defaultConfig.UseTLS,
		"start json rpc endpoint with tls enabled",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.TLSCertFile,
		tlsCertFileLocationFlag,
		defaultConfig.TLSCertFile,
		"path to TLS cert file, if no file is provided then cert file is loaded from secrets manager",
	)

	cmd.Flags().StringVar(
		&params.rawConfig.TLSKeyFile,
		tlsKeyFileLocationFlag,
		defaultConfig.TLSKeyFile,
		"path to TLS key file, if no file is provided then key file is loaded from secrets manager",
	)

	cmd.Flags().BoolVar(
		&params.rawConfig.Relayer,
		relayerFlag,
		defaultConfig.Relayer,
		"start the state sync relayer service (PolyBFT only)",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.ConcurrentRequestsDebug,
		concurrentRequestsDebugFlag,
		defaultConfig.ConcurrentRequestsDebug,
		"maximal number of concurrent requests for debug endpoints",
	)

	cmd.Flags().Uint64Var(
		&params.rawConfig.WebSocketReadLimit,
		webSocketReadLimitFlag,
		defaultConfig.WebSocketReadLimit,
		"maximum size in bytes for a message read from the peer by websocket",
	)

	cmd.Flags().DurationVar(
		&params.rawConfig.MetricsInterval,
		metricsIntervalFlag,
		defaultConfig.MetricsInterval,
		"the interval (in seconds) at which special metrics are generated. a value of zero means the metrics are disabled",
	)

	{ // event tracker
		cmd.Flags().Uint64Var(
			&params.rawConfig.EventTracker.SyncBatchSize,
			trackerSyncBatchSizeFlag,
			defaultConfig.EventTracker.SyncBatchSize,
			`defines a batch size of blocks that will be gotten from tracked chain,
			when tracker is out of sync and needs to sync a number of blocks.
			(e.g., SyncBatchSize = 10, trackers last processed block is 10, latest block on tracked chain is 100,
			it will get blocks 11-20, get logs from confirmed blocks of given batch, remove processed confirm logs
			from memory, and continue to the next batch)`,
		)

		cmd.Flags().Uint64Var(
			&params.rawConfig.EventTracker.NumBlockConfirmations,
			trackerNumBlockConfirmationsFlag,
			defaultConfig.EventTracker.NumBlockConfirmations,
			"minimal number of child blocks required for the parent block to be considered final on tracked chain",
		)

		cmd.Flags().Uint64Var(
			&params.rawConfig.EventTracker.NumOfBlocksToReconcile,
			trackerNumOfBlocksToReconcileFlag,
			defaultConfig.EventTracker.NumBlockConfirmations,
			`defines how many blocks we will sync up from the latest block on tracked chain. 
			If a node that has a tracker, was offline for days, months, a year, it is going to miss a lot of blocks potentially. 
			In the meantime, we expect the rest of nodes to have collected the desired events and did their 
			logic with them, continuing consensus and relayer stuff. 
			In order to not waste too much unnecessary time in syncing all those blocks, with NumOfBlocksToReconcile, 
			we tell the tracker to sync only latestBlock.Number - NumOfBlocksToReconcile number of blocks.`,
		)
	}

	setDevFlags(cmd)
}

func setDevFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(
		&params.isDevMode,
		devFlag,
		false,
		"should the client start in dev mode (default false)",
	)

	_ = cmd.Flags().MarkHidden(devFlag)

	cmd.Flags().Uint64Var(
		&params.devInterval,
		devIntervalFlag,
		0,
		"the client's dev notification interval in seconds (default 1)",
	)

	_ = cmd.Flags().MarkHidden(devIntervalFlag)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	// Set the grpc and json ip:port bindings
	// The config file will have precedence over --flag
	params.setRawGRPCAddress(helper.GetGRPCAddress(cmd))
	params.setRawJSONRPCAddress(helper.GetJSONRPCAddress(cmd))
	params.setJSONLogFormat(helper.GetJSONLogFormat(cmd))

	// Check if the config file has been specified
	// Config file settings will override JSON-RPC and GRPC address values
	if isConfigFileSpecified(cmd) {
		if err := params.initConfigFromFile(); err != nil {
			return err
		}
	}

	if err := params.initRawParams(); err != nil {
		return err
	}

	return nil
}

func isConfigFileSpecified(cmd *cobra.Command) bool {
	return cmd.Flags().Changed(configFlag)
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)

	if err := runServerLoop(params.generateConfig(), outputter); err != nil {
		outputter.SetError(err)
		outputter.WriteOutput()

		return
	}
}

func runServerLoop(
	config *server.Config,
	outputter command.OutputFormatter,
) error {
	serverInstance, err := server.NewServer(config)
	if err != nil {
		return err
	}

	return helper.HandleSignals(serverInstance.Close, outputter)
}
