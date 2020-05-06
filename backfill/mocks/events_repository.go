package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type EventsRepository struct {
	GetFrobsPassedUrnIDs        []int
	GetFrobsPassedStartingBlock int
	GetFrobsFrobsToReturn       []backfill.Frob
	GetFrobsError               error
	GetHeaderByIDPassedIDs      []int
	GetHeaderByIDHeaderToReturn core.Header
	GetHeaderByIDError          error
}

func (e *EventsRepository) GetFrobs(urnID, startingBlock int) ([]backfill.Frob, error) {
	e.GetFrobsPassedUrnIDs = append(e.GetFrobsPassedUrnIDs, urnID)
	e.GetFrobsPassedStartingBlock = startingBlock
	return e.GetFrobsFrobsToReturn, e.GetFrobsError
}

func (e *EventsRepository) GetHeaderByID(id int) (core.Header, error) {
	e.GetHeaderByIDPassedIDs = append(e.GetHeaderByIDPassedIDs, id)
	return e.GetHeaderByIDHeaderToReturn, e.GetHeaderByIDError
}
