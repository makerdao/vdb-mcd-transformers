package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type EventsRepository struct {
	GetFrobsError               error
	GetFrobsFrobsToReturn       []backfill.Frob
	GetFrobsPassedStartingBlock int
	GetFrobsPassedUrnIDs        []int
	GetGrabsError               error
	GetGrabsGrabsToReturn       []backfill.Grab
	GetGrabsPassedStartingBlock int
	GetHeaderByIDError          error
	GetHeaderByIDHeaderToReturn core.Header
	GetHeaderByIDPassedIDs      []int
}

func (e *EventsRepository) GetFrobs(urnID, startingBlock int) ([]backfill.Frob, error) {
	e.GetFrobsPassedUrnIDs = append(e.GetFrobsPassedUrnIDs, urnID)
	e.GetFrobsPassedStartingBlock = startingBlock
	return e.GetFrobsFrobsToReturn, e.GetFrobsError
}

func (e *EventsRepository) GetGrabs(startingBlock int) ([]backfill.Grab, error) {
	e.GetGrabsPassedStartingBlock = startingBlock
	return e.GetGrabsGrabsToReturn, e.GetGrabsError
}

func (e *EventsRepository) GetHeaderByID(id int) (core.Header, error) {
	e.GetHeaderByIDPassedIDs = append(e.GetHeaderByIDPassedIDs, id)
	return e.GetHeaderByIDHeaderToReturn, e.GetHeaderByIDError
}
