package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type EventsRepository struct {
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

func (e *EventsRepository) GetFrobs(startingBlock int) ([]repository.Frob, error) {
	e.GetFrobsPassedStartingBlock = startingBlock
	return e.GetFrobsFrobsToReturn, e.GetFrobsError
}

func (e *EventsRepository) GetGrabs(startingBlock int) ([]repository.Grab, error) {
	e.GetGrabsPassedStartingBlock = startingBlock
	return e.GetGrabsGrabsToReturn, e.GetGrabsError
}

func (e *EventsRepository) GetHeaderByID(id int) (core.Header, error) {
	e.GetHeaderByIDPassedIDs = append(e.GetHeaderByIDPassedIDs, id)
	return e.GetHeaderByIDHeaderToReturn, e.GetHeaderByIDError
}
