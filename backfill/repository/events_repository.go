package repository

import (
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Frob struct {
	HeaderID int `db:"header_id"`
	UrnID    int `db:"urn_id"`
	Dink     string
	Dart     string
}

type Grab struct {
	HeaderID int `db:"header_id"`
	UrnID    int `db:"urn_id"`
	Dink     string
	Dart     string
}

type EventsRepository interface {
	GetFrobs(startingBlock int) ([]Frob, error)
	GetGrabs(startingBlock int) ([]Grab, error)
	GetHeaderByID(id int) (core.Header, error)
}

type eventsRepository struct {
	db *postgres.DB
}

func NewEventsRepository(db *postgres.DB) EventsRepository {
	return eventsRepository{db: db}
}

func (e eventsRepository) GetFrobs(startingBlock int) ([]Frob, error) {
	var frobs []Frob
	err := e.db.Select(&frobs, `SELECT header_id, urn_id, dink, dart
		FROM maker.vat_frob
		JOIN public.headers ON vat_frob.header_id = headers.id
		WHERE headers.block_number >= $1
		ORDER BY headers.block_number ASC`, startingBlock)
	return frobs, err
}

func (e eventsRepository) GetGrabs(startingBlock int) ([]Grab, error) {
	var grabs []Grab
	err := e.db.Select(&grabs, `SELECT header_id, urn_id, dink, dart
		FROM maker.vat_grab
		JOIN public.headers ON vat_grab.header_id = headers.id
		WHERE headers.block_number >= $1`, startingBlock)
	return grabs, err
}

func (e eventsRepository) GetHeaderByID(id int) (core.Header, error) {
	var header core.Header
	headerErr := e.db.Get(&header, `SELECT id, block_number, hash, raw, block_timestamp FROM headers WHERE id = $1`, id)
	return header, headerErr
}
