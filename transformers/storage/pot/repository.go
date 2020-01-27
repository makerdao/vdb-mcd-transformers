package pot

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertPotUserPieQuery = `INSERT INTO maker.pot_user_pie (diff_id, header_id, "user", pie) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertPotPieQuery     = `INSERT INTO maker.pot_pie (diff_id, header_id, pie) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertPotDsrQuery     = `INSERT INTO maker.pot_dsr (diff_id, header_id, dsr) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertPotChiQuery     = `INSERT INTO maker.pot_chi (diff_id, header_id, chi) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertPotVatQuery     = `INSERT INTO maker.pot_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertPotVowQuery     = `INSERT INTO maker.pot_vow (diff_id, header_id, vow) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertPotRhoQuery     = `INSERT INTO maker.pot_rho (diff_id, header_id, rho) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertPotLiveQuery    = `INSERT INTO maker.pot_live (diff_id, header_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type PotStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

type diffInserter func(int64, *sqlx.Tx) error

func (repository PotStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case UserPie:
		return repository.insertUserPie(diffID, headerID, metadata, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case Pie:
		return repository.insertPie(diffID, headerID, value.(string))
	case Dsr:
		return repository.insertDsr(diffID, headerID, value.(string))
	case Chi:
		return repository.insertChi(diffID, headerID, value.(string))
	case Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case Vow:
		return repository.insertVow(diffID, headerID, value.(string))
	case Rho:
		return repository.insertRho(diffID, headerID, value.(string))
	case Live:
		return repository.insertLive(diffID, headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized pot contract storage name: %s", metadata.Name))
	}
}

func (repository *PotStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository PotStorageRepository) insertUserPie(diffID, headerID int64, metadata types.ValueMetadata, pie string) error {
	user, err := getUser(metadata.Keys)
	if err != nil {
		return err
	}

	var insertUserPieDiff diffInserter = func(addressID int64, tx *sqlx.Tx) error {
		_, writeErr := tx.Exec(insertPotUserPieQuery, diffID, headerID, addressID, pie)
		return writeErr
	}
	return repository.insertRecordWithAddress(user, UserPie, insertUserPieDiff)
}

func (repository PotStorageRepository) insertPie(diffID, headerID int64, pie string) error {
	_, err := repository.db.Exec(insertPotPieQuery, diffID, headerID, pie)
	return err
}

func (repository PotStorageRepository) insertDsr(diffID, headerID int64, dsr string) error {
	_, err := repository.db.Exec(insertPotDsrQuery, diffID, headerID, dsr)
	return err
}

func (repository PotStorageRepository) insertChi(diffID, headerID int64, chi string) error {
	_, err := repository.db.Exec(insertPotChiQuery, diffID, headerID, chi)
	return err
}

func (repository PotStorageRepository) insertVat(diffID, headerID int64, vat string) error {
	var insertVatDiff diffInserter = func(addressID int64, tx *sqlx.Tx) error {
		_, writeErr := tx.Exec(insertPotVatQuery, diffID, headerID, addressID)
		return writeErr
	}
	return repository.insertRecordWithAddress(vat, Vat, insertVatDiff)
}

func (repository PotStorageRepository) insertVow(diffID, headerID int64, vow string) error {
	var insertVowDiff diffInserter = func(addressID int64, tx *sqlx.Tx) error {
		_, writeErr := tx.Exec(insertPotVowQuery, diffID, headerID, addressID)
		return writeErr
	}
	return repository.insertRecordWithAddress(vow, Vow, insertVowDiff)
}

func (repository PotStorageRepository) insertRho(diffID, headerID int64, rho string) error {
	_, err := repository.db.Exec(insertPotRhoQuery, diffID, headerID, rho)
	return err
}

func (repository PotStorageRepository) insertLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertPotLiveQuery, diffID, headerID, live)
	return err
}

func (repository *PotStorageRepository) insertRecordWithAddress(address, variableName string, insertDiffRecord diffInserter) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressID, addressErr := shared.GetOrCreateAddressInTransaction(address, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("address", addressErr.Error())
		}
		return addressErr
	}

	writeErr := insertDiffRecord(addressID, tx)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func getUser(keys map[types.Key]string) (string, error) {
	ilk, ok := keys[constants.MsgSender]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.MsgSender}
	}
	return ilk, nil
}
