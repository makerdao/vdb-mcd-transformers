package grab

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/sirupsen/logrus"
)

type grabBackFiller struct {
	blockChain        core.BlockChain
	eventsRepository  repository.EventsRepository
	storageRepository repository.StorageRepository
}

func NewGrabBackFiller(blockChain core.BlockChain, eventsRepository repository.EventsRepository, storageRepository repository.StorageRepository) backfill.BackFiller {
	return grabBackFiller{
		blockChain:        blockChain,
		eventsRepository:  eventsRepository,
		storageRepository: storageRepository,
	}
}

func (backFiller grabBackFiller) BackFill(startingBlock int) error {
	grabs, grabsErr := backFiller.eventsRepository.GetGrabs(startingBlock)
	if grabsErr != nil {
		return fmt.Errorf("error getting grab events: %w", grabsErr)
	}

	logrus.Infof("getting storage for %d grabs", len(grabs))
	for i, grab := range grabs {
		dartDink := shared.DartDink{
			Dart:     grab.Dart,
			Dink:     grab.Dink,
			HeaderID: grab.HeaderID,
			UrnID:    grab.UrnID,
		}
		err := shared.FetchAndPersistDartDinkDiffs(dartDink, backFiller.eventsRepository, backFiller.storageRepository, backFiller.blockChain)
		if err != nil {
			return fmt.Errorf("error fetching and persisting diffs for grab %d: %w", i, err)
		}
	}
	return nil
}
