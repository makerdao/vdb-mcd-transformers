package repository

import (
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Fork struct {
	HeaderID int64 `db:"header_id"`
	Ilk      string
	Src      string
	Dst      string
	Dink     string
	Dart     string
}

type Frob struct {
	HeaderID int64 `db:"header_id"`
	UrnID    int64 `db:"urn_id"`
	Dink     string
	Dart     string
}

type Grab struct {
	HeaderID int64 `db:"header_id"`
	UrnID    int64 `db:"urn_id"`
	Dink     string
	Dart     string
}

type EventsRepository interface {
	GetForks(startingBlock int) ([]Fork, error)
	GetFrobs(startingBlock int) ([]Frob, error)
	GetGrabs(startingBlock int) ([]Grab, error)
}

type eventsRepository struct {
	db *postgres.DB
}

func NewEventsRepository(db *postgres.DB) EventsRepository {
	return eventsRepository{db: db}
}

func (e eventsRepository) GetForks(startingBlock int) ([]Fork, error) {
	var forks []Fork
	err := e.db.Select(&forks, `SELECT header_id, ilks.ilk, src, dst, dink, dart
		FROM maker.vat_fork
		    JOIN maker.ilks on vat_fork.ilk_id = ilks.id
			JOIN public.headers ON vat_fork.header_id = headers.id
		WHERE headers.block_number >= $1
		ORDER BY headers.block_number ASC`, startingBlock)
	return forks, err
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
		WHERE headers.block_number >= $1
		ORDER BY headers.block_number ASC`, startingBlock)
	return grabs, err
}
