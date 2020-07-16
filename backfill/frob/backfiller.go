package frob

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/sirupsen/logrus"
)

type frobBackFiller struct {
	blockChain        core.BlockChain
	dartDinkRetriever shared.DartDinkRetriever
	eventsRepository  repository.EventsRepository
	storageRepository repository.StorageRepository
}

func NewFrobBackFiller(blockChain core.BlockChain, eventsRepository repository.EventsRepository, storageRepository repository.StorageRepository, dartDinkRetriever shared.DartDinkRetriever) backfill.BackFiller {
	return frobBackFiller{
		blockChain:        blockChain,
		dartDinkRetriever: dartDinkRetriever,
		eventsRepository:  eventsRepository,
		storageRepository: storageRepository,
	}
}

func (backFiller frobBackFiller) BackFill(startingBlock int) error {
	frobs, frobsErr := backFiller.eventsRepository.GetFrobs(startingBlock)
	if frobsErr != nil {
		return fmt.Errorf("error getting frobs: %w", frobsErr)
	}

	logrus.Infof("getting storage for %d frobs", len(frobs))
	for i, frob := range frobs {
		dartDink := shared.DartDink{
			Dart:     frob.Dart,
			Dink:     frob.Dink,
			HeaderID: frob.HeaderID,
			UrnID:    frob.UrnID,
		}
		err := backFiller.dartDinkRetriever.RetrieveDartDinkDiffs(dartDink)
		if err != nil {
			return fmt.Errorf("error fetching and persisting diffs for frob %d: %w", i, err)
		}
	}
	logrus.Infof("finished getting storage for %d frobs", len(frobs))
	return nil
}
