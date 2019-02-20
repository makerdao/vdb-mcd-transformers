package vat_tune

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatTuneRepository struct {
	db *postgres.DB
}

func (repository VatTuneRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		vatTune, ok := model.(VatTuneModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VatTuneModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(vatTune.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		urnID, urnErr := shared.GetOrCreateUrnInTransaction(vatTune.Urn, ilkID, tx)
		if urnErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback", rollbackErr)
			}
			return urnErr
		}

		_, execErr := tx.Exec(
			`INSERT into maker.vat_tune (header_id, urn, v, w, dink, dart, tx_idx, log_idx, raw_log)
	   VALUES($1, $2, $3, $4, $5::NUMERIC, $6::NUMERIC, $7, $8, $9)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET urn = $2, v = $3, w = $4, dink = $5, dart = $6, raw_log = $9;`,
			headerID, urnID, vatTune.V, vatTune.W, vatTune.Dink, vatTune.Dart, vatTune.TransactionIndex, vatTune.LogIndex, vatTune.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatTuneChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository VatTuneRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, constants.VatTuneChecked)
}

func (repository VatTuneRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.VatTuneChecked)
}

func (repository VatTuneRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.VatTuneChecked)
}

func (repository *VatTuneRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
