package main

import (
	"github.com/ouroboros-crypto/node"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	distrCli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	paramsCli "github.com/cosmos/cosmos-sdk/x/params/client/cli"
	ouroborosCli "github.com/ouroboros-crypto/node/x/ouroboros/client/cli"
	ouroborosRest "github.com/ouroboros-crypto/node/x/ouroboros/client/rest"
	addrs "github.com/ouroboros-crypto/node/x/ouroboros/types"
	paraminingCli "github.com/ouroboros-crypto/node/x/paramining/client/cli"
	paraminingRest "github.com/ouroboros-crypto/node/x/paramining/client/rest"
	structureCli "github.com/ouroboros-crypto/node/x/structure/client/cli"
	structureRest "github.com/ouroboros-crypto/node/x/structure/client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	storeAcc = "acc"
	storeNS  = "ouroboros"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := Node_Github.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(addrs.AccAddr, addrs.AccPub)
	config.SetBech32PrefixForValidator(addrs.ValAddr, addrs.ValPub)
	config.SetBech32PrefixForConsensusNode(addrs.ConsAddr, addrs.ConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "ouroboroscli",
		Short: "Ouroboros Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(Node_Github.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
	)

	executor := cli.PrepareMainCmd(rootCmd, "OURO", Node_Github.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)

	structureRest.RegisterRoutes(rs.CliCtx, rs.Mux)
	paraminingRest.RegisterRoutes(rs.CliCtx, rs.Mux)
	ouroborosRest.RegisterRoutes(rs.CliCtx, rs.Mux)

	Node_Github.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)

	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		client.LineBreak,
		ouroborosCli.GetQueryCmd(cdc),
		client.LineBreak,
		structureCli.GetQueryCmd(cdc),
		client.LineBreak,
		paraminingCli.GetQueryCmd(cdc),
		client.LineBreak,
	)

	// add modules' query commands
	Node_Github.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	paramsProposalCmd := client.PostCommands(paramsCli.GetCmdSubmitProposal(cdc))[0]
	distrProposalCmd := client.PostCommands(distrCli.GetCmdSubmitProposal(cdc))[0]

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		paramsProposalCmd,
		client.LineBreak,
		distrProposalCmd,
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		client.LineBreak,
		paraminingCli.GetTxCmd(cdc),
		client.LineBreak,
		ouroborosCli.GetTxCmd(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	Node_Github.ModuleBasics.AddTxCommands(txCmd, cdc)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)

	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")

	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}

	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}

	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}