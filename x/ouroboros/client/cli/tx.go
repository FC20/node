package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ouroboros-crypto/node/x/ouroboros/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	ouroborosTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Ouroboros transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ouroborosTxCmd.AddCommand(flags.PostCommands(
		GetCmdAmnesty(cdc),
	)...)

	return ouroborosTxCmd
}

func GetCmdAmnesty(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "amnesty",
		Short: "Amnesty all the validators",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgAmnesty(cliCtx.GetFromAddress())
			err := msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
