package cmd

import (
	"errors"
	"fmt"
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sync"
)

var (
	startingBlock    int
	eventsToBackFill []string
	minEvents        = 1
	maxEvents        = 2
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
	eventRepository := backfill.NewEventsRepository(&db)
	storageRepository := backfill.NewStorageRepository(&db)

	lenEvents := len(eventsToBackFill)
	if lenEvents == 1 {
		if eventsToBackFill[0] == "frob" {
			backFiller := backfill.NewFrobBackFiller(blockChain, eventRepository, storageRepository)
			return backFiller.BackFillFrobStorage(startingBlock)
		} else {
			backFiller := backfill.NewGrabBackFiller(blockChain, eventRepository, storageRepository)
			return backFiller.BackFillGrabStorage(startingBlock)
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(2)

		done := make(chan bool)
		errs := make(chan error)

		go backFillFrobEvents(blockChain, eventRepository, storageRepository, errs, &wg)
		go backFillGrabEvents(blockChain, eventRepository, storageRepository, errs, &wg)

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
		if eventsToBackFill[0] != "frob" && eventsToBackFill[0] != "grab" {
			return errors.New("only frob and/or grab are allowed")
		}
	}
	if lenEvents == 2 {
		if (eventsToBackFill[0] != "frob" && eventsToBackFill[0] != "grab") || (eventsToBackFill[0] != "frob" && eventsToBackFill[1] != "grab") {
			return errors.New("only frob and/or grab are allowed")
		}
		if eventsToBackFill[0] == eventsToBackFill[1] {
			return errors.New("same event included twice")
		}
	}
	return nil
}

func backFillFrobEvents(blockChain core.BlockChain, eventsRepository backfill.EventsRepository, storageRepository backfill.StorageRepository, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	backFiller := backfill.NewFrobBackFiller(blockChain, eventsRepository, storageRepository)
	err := backFiller.BackFillFrobStorage(startingBlock)
	if err != nil {
		errs <- err
		return
	}
}

func backFillGrabEvents(blockChain core.BlockChain, eventsRepository backfill.EventsRepository, storageRepository backfill.StorageRepository, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	backFiller := backfill.NewGrabBackFiller(blockChain, eventsRepository, storageRepository)
	err := backFiller.BackFillGrabStorage(startingBlock)
	if err != nil {
		errs <- err
		return
	}
}
