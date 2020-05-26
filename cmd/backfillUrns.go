package cmd

import (
	"fmt"
	"sync"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/fork"
	"github.com/makerdao/vdb-mcd-transformers/backfill/frob"
	"github.com/makerdao/vdb-mcd-transformers/backfill/grab"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type backFillInitializer func(core.BlockChain, repository.EventsRepository, repository.StorageRepository, shared.DartDinkRetriever) backfill.BackFiller

var (
	eventsToBackFill []string
	initializers     = map[string]backFillInitializer{
		backfill.ForkEvent: fork.NewForkBackFiller,
		backfill.FrobEvent: frob.NewFrobBackFiller,
		backfill.GrabEvent: grab.NewGrabBackFiller,
	}
	startingBlock int
)

// backfillUrnsCmd represents the backfillUrns command
var backfillUrnsCmd = &cobra.Command{
	Use:   "backfillUrns",
	Short: "Backfill diffs for urns, looking up diffs based on associated events",
	Long: `Fetch diffs when events indicate the state of an Urn changed at a given block.
Optionally pass a starting block number to backfill since a given block.
Optionally pass events to watch (frob, grab) to backfill based off of certain events.`,
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
	backfillUrnsCmd.Flags().StringSliceVarP(&eventsToBackFill, "events-to-backfill", "e", []string{"frob", "grab"}, "events to back-fill")
}

func backfillUrns() error {
	validationErr := backfill.ValidateArgs(eventsToBackFill)
	if validationErr != nil {
		return fmt.Errorf("invalid events-to-backfill: %w", validationErr)
	}

	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	eventRepository := repository.NewEventsRepository(&db)
	storageRepository := repository.NewStorageRepository(&db)

	var wg sync.WaitGroup
	done := make(chan bool)
	errs := make(chan error)

	for _, e := range eventsToBackFill {
		initializer := initializers[e]
		headerRepository := repositories.NewHeaderRepository(&db)
		dartDinkRetriever := shared.NewDartDinkRetriever(blockChain, eventRepository, headerRepository, storageRepository)
		backFiller := initializer(blockChain, eventRepository, storageRepository, dartDinkRetriever)
		wg.Add(1)
		go backFillEvents(backFiller, startingBlock, errs, &wg)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case err := <-errs:
		logrus.Errorf("error executing back-fill: %w", err)
		return err
	}

	return nil
}

func backFillEvents(backFiller backfill.BackFiller, startingBlock int, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	err := backFiller.BackFill(startingBlock)
	if err != nil {
		errs <- err
	}
	return
}
