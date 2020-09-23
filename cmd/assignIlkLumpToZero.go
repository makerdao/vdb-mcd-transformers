package cmd

import (
	"github.com/makerdao/vdb-mcd-transformers/zero_value_diff"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	blockNumber int
	ilks        []string
)

var assignIlkLumpToZeroCmd = &cobra.Command{
	Use:   "assignIlkLumpToZero",
	Short: "Create a diff with a zero ilk lump value.",
	Long: `Inserts artificial cat_ilk_lump storage_diffs with storage_values of zero for the given ilks at the given block
height. The storage_diffs will then be picked up in the execute process to be transformed into maker.cat_ilk_lump records.
The block-number (b) and ilks (i) arguments are required.

Note: the ilk argument(s) passed in need to be the hex representation of the ilk.

Note: this command should only be run for the block height when the Cat contract upgrade change takes place, i.e. this
should only be run for the block where the first cat_ilk_dunk was set for the given ilks.
`,
	Example: `./vdb-mcd-transformers assignIlkLumpToZero --blockHeight=10769102
		--ilks 0x4241542d41000000000000000000000000000000000000000000000000000000
		--ilks 0x4554482d41000000000000000000000000000000000000000000000000000000
	`,
	PreRun: setViperConfigs,
	Run: func(cmd *cobra.Command, args []string) {
		err := zeroOutIlkLump()
		if err != nil {
			logrus.Fatalf("error inserting a zero ilk lump diff: %s", err.Error())
		}
		logrus.Infof("Successfully created a zero-value ilk lump diffs for block %d", blockNumber)
		return
	},
}

func init() {
	rootCmd.AddCommand(assignIlkLumpToZeroCmd)
	assignIlkLumpToZeroCmd.Flags().IntVarP(&blockNumber, "block-number", "b", -1, "the block associated with artificial ilk lump diffs")
	assignIlkLumpToZeroCmd.MarkFlagRequired("block-number")
	assignIlkLumpToZeroCmd.Flags().StringSliceVarP(&ilks, "ilks", "i", []string{}, "the ilk to set lump to zero")
	assignIlkLumpToZeroCmd.MarkFlagRequired("ilk")
}

func zeroOutIlkLump() error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	generator := zero_value_diff.NewZeroValueDiffGenerator(db)
	return generator.CreateZeroValueIlkLumpDiff(blockNumber, ilks)
}
