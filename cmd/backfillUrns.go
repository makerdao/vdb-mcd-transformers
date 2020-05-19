package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/frob"
	"github.com/makerdao/vdb-mcd-transformers/backfill/grab"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type backFillInitializer func(core.BlockChain, repository.EventsRepository, repository.StorageRepository) backfill.BackFiller

var (
	eventsToBackFill []string
	frobEvent        = "frob"
	grabEvent        = "grab"
	initializers     = map[string]backFillInitializer{
		frobEvent: frob.NewFrobBackFiller,
		grabEvent: grab.NewGrabBackFiller,
	}
	maxEvents     = 2
	minEvents     = 1
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
	validationErr := validateEventsToBackfill()
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
		backFiller := initializer(blockChain, eventRepository, storageRepository)
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

func validateEventsToBackfill() error {
	lenEvents := len(eventsToBackFill)
	if lenEvents < minEvents {
		return fmt.Errorf("must specify at least %d event(s)", minEvents)
	}
	if lenEvents > maxEvents {
		return fmt.Errorf("only %d events are allowed", maxEvents)
	}
	if lenEvents == 1 {
		if eventsToBackFill[0] != frobEvent && eventsToBackFill[0] != grabEvent {
			return errors.New("only frob and/or grab are allowed")
		}
	}
	if lenEvents == 2 {
		if (eventsToBackFill[0] != frobEvent && eventsToBackFill[0] != grabEvent) || (eventsToBackFill[0] != frobEvent && eventsToBackFill[1] != grabEvent) {
			return errors.New("only frob and/or grab are allowed")
		}
		if eventsToBackFill[0] == eventsToBackFill[1] {
			return errors.New("same event included twice")
		}
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
