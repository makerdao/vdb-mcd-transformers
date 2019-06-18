package vat_grab

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatGrabRepository struct {
	db *postgres.DB
}

func (repository VatGrabRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		vatGrab, ok := model.(VatGrabModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VatGrabModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(vatGrab.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		urnID, urnErr := shared.GetOrCreateUrnInTransaction(vatGrab.Urn, ilkID, tx)
		if urnErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback", rollbackErr)
			}
			return urnErr
		}

		_, execErr := tx.Exec(
			`INSERT into maker.vat_grab (header_id, urn_id, v, w, dink, dart, log_idx, tx_idx, raw_log)
	   VALUES($1, $2, $3, $4, $5::NUMERIC, $6::NUMERIC, $7, $8, $9)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET urn_id = $2, v = $3, w = $4, dink = $5, dart = $6, raw_log = $9;`,
			headerID, urnID, vatGrab.V, vatGrab.W, vatGrab.Dink, vatGrab.Dart, vatGrab.LogIndex, vatGrab.TransactionIndex, vatGrab.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}
	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatGrabChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository VatGrabRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, constants.VatGrabChecked)
}

func (repository *VatGrabRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
