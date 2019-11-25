package dsr

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Repository struct {
	db *postgres.DB
}

func (repository Repository) Create(models []event.InsertionModel) error {
	return event.Create(models, repository.db)
}

func (repository *Repository) SetDB(db *postgres.DB) {
	repository.db = db
}
