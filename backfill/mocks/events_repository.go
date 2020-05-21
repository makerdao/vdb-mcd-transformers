package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type EventsRepository struct {
	GetForksError               error
	GetForksForksToReturn       []repository.Fork
	GetForksPassedStartingBlock int
	GetFrobsError               error
	GetFrobsFrobsToReturn       []repository.Frob
	GetFrobsPassedStartingBlock int
	GetGrabsError               error
	GetGrabsGrabsToReturn       []repository.Grab
	GetGrabsPassedStartingBlock int
	GetHeaderByIDError          error
	GetHeaderByIDHeaderToReturn core.Header
	GetHeaderByIDPassedIDs      []int
}

func (mock *EventsRepository) GetForks(startingBlock int) ([]repository.Fork, error) {
	mock.GetForksPassedStartingBlock = startingBlock
	return mock.GetForksForksToReturn, mock.GetForksError
}

func (mock *EventsRepository) GetFrobs(startingBlock int) ([]repository.Frob, error) {
	mock.GetFrobsPassedStartingBlock = startingBlock
	return mock.GetFrobsFrobsToReturn, mock.GetFrobsError
}

func (mock *EventsRepository) GetGrabs(startingBlock int) ([]repository.Grab, error) {
	mock.GetGrabsPassedStartingBlock = startingBlock
	return mock.GetGrabsGrabsToReturn, mock.GetGrabsError
}

func (mock *EventsRepository) GetHeaderByID(id int) (core.Header, error) {
	mock.GetHeaderByIDPassedIDs = append(mock.GetHeaderByIDPassedIDs, id)
	return mock.GetHeaderByIDHeaderToReturn, mock.GetHeaderByIDError
}
