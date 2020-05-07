package cmd

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startingBlock int

// backfillUrnsCmd represents the backfillUrns command
var backfillUrnsCmd = &cobra.Command{
	Use:   "backfillUrns",
	Short: "Backfill diffs for urns, looking up diffs based on associated events",
	Long: `Fetch diffs when events indicate the state of an Urn changed at a given block.
You can optionally pass a starting block number to backfill since a given block.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := backfillUrns()
		if err != nil {
			logrus.Fatalf("error backfilling urns: %s", err.Error())
		}
		logrus.Println("Backfilling urns completed successfully")
		return
	},
}

func init() {
	rootCmd.AddCommand(backfillUrnsCmd)
	backfillUrnsCmd.Flags().IntVarP(&startingBlock, "starting-block", "s", 0, "starting block for backfilling diffs derived from urn events")

}

func backfillUrns() error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	eventRepository := backfill.NewEventsRepository(&db)
	urnsRepository := backfill.NewUrnsRepository(&db)
	backfiller := backfill.NewUrnBackFiller(blockChain, eventRepository, urnsRepository)

	return backfiller.BackfillUrns(startingBlock)
}
