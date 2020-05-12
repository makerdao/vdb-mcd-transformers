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
	eventsRepository  repository.EventsRepository
	storageRepository repository.StorageRepository
}

func NewFrobBackFiller(blockChain core.BlockChain, eventsRepository repository.EventsRepository, storageRepository repository.StorageRepository) backfill.BackFiller {
	return frobBackFiller{
		blockChain:        blockChain,
		eventsRepository:  eventsRepository,
		storageRepository: storageRepository,
	}
}

func (backFiller frobBackFiller) BackFill(startingBlock int) error {
	frobs, frobsErr := backFiller.eventsRepository.GetFrobs(startingBlock)
	if frobsErr != nil {
		return fmt.Errorf("error getting frobs: %w", frobsErr)
	}

	logrus.Infof("getting storage for %d grabs", len(frobs))
	for i, frob := range frobs {
		err := shared.FetchAndPersistDartDinkDiffs(frob.Dart, frob.Dink, frob.UrnID, frob.HeaderID, backFiller.eventsRepository, backFiller.storageRepository, backFiller.blockChain)
		if err != nil {
			return fmt.Errorf("error fetching and persisting diffs for frob %d: %w", i, err)
		}
	}
	return nil
}
