package fork

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/sirupsen/logrus"
)

type urnDinkDart struct {
	Ilk      string
	Guy      string
	Dart     string
	Dink     string
	HeaderID int
}

type forkBackFiller struct {
	blockChain        core.BlockChain
	dartDinkRetriever shared.DartDinkRetriever
	eventsRepository  repository.EventsRepository
	storageRepository repository.StorageRepository
}

func NewForkBackFiller(blockChain core.BlockChain, eventsRepository repository.EventsRepository, storageRepository repository.StorageRepository, dartDinkRetriever shared.DartDinkRetriever) backfill.BackFiller {
	return forkBackFiller{
		blockChain:        blockChain,
		dartDinkRetriever: dartDinkRetriever,
		eventsRepository:  eventsRepository,
		storageRepository: storageRepository,
	}
}

func (backFiller forkBackFiller) BackFill(startingBlock int) error {
	forks, forksErr := backFiller.eventsRepository.GetForks(startingBlock)
	if forksErr != nil {
		return fmt.Errorf("error getting forks: %w", forksErr)
	}

	logrus.Infof("getting storage for %d forks", len(forks))
	for _, fork := range forks {
		srcErr := backFiller.retrieveDiffsForUrn(getDartDinkForForkGuy(fork.Src, fork))
		if srcErr != nil {
			return fmt.Errorf("error retrieving src for fork: %w", srcErr)
		}

		dstErr := backFiller.retrieveDiffsForUrn(getDartDinkForForkGuy(fork.Dst, fork))
		if dstErr != nil {
			return fmt.Errorf("error retrieving dst for fork: %w", dstErr)
		}
	}
	logrus.Infof("finished getting storage for %d forks", len(forks))
	return nil
}

func (backFiller forkBackFiller) retrieveDiffsForUrn(data urnDinkDart) error {
	urnID, urnErr := backFiller.storageRepository.GetOrCreateUrn(data.Guy, data.Ilk)
	if urnErr != nil {
		return fmt.Errorf("error getting or creating urn ID: %w", urnErr)
	}
	dartDink := shared.DartDink{
		Dart:     data.Dart,
		Dink:     data.Dink,
		HeaderID: data.HeaderID,
		UrnID:    int(urnID),
	}
	err := backFiller.dartDinkRetriever.RetrieveDartDinkDiffs(dartDink)
	if err != nil {
		return fmt.Errorf("error getching and persisting diffs: %w", err)
	}
	return nil
}

func getDartDinkForForkGuy(guy string, fork repository.Fork) urnDinkDart {
	return urnDinkDart{
		Ilk:      fork.Ilk,
		Guy:      guy,
		Dart:     fork.Dart,
		Dink:     fork.Dink,
		HeaderID: fork.HeaderID,
	}
}
