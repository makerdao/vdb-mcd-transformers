package pot

import (
	"fmt"

	shared2 "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
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

type StorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repo StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case UserPie:
		return repo.insertUserPie(diffID, headerID, metadata, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	case Pie:
		return repo.insertPie(diffID, headerID, value.(string))
	case Dsr:
		return repo.insertDsr(diffID, headerID, value.(string))
	case Chi:
		return repo.insertChi(diffID, headerID, value.(string))
	case Vat:
		return repo.insertVat(diffID, headerID, value.(string))
	case Vow:
		return repo.insertVow(diffID, headerID, value.(string))
	case Rho:
		return repo.insertRho(diffID, headerID, value.(string))
	case Live:
		return repo.insertLive(diffID, headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized pot contract storage name: %s", metadata.Name))
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo StorageRepository) insertUserPie(diffID, headerID int64, metadata types.ValueMetadata, pie string) error {
	user, err := getUser(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting user for pot user pie: %w", err)
	}
	insertErr := shared2.InsertRecordWithAddress(diffID, headerID, insertPotUserPieQuery, pie, user, repo.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting pot user %s pie %s from diff ID %d: %w", user, pie, diffID, insertErr)
	}
	return nil
}

func (repo StorageRepository) insertPie(diffID, headerID int64, pie string) error {
	_, err := repo.db.Exec(insertPotPieQuery, diffID, headerID, pie)
	if err != nil {
		return fmt.Errorf("error inserting pot pie %s from diff ID %d: %w", pie, diffID, err)
	}
	return nil
}

func (repo StorageRepository) insertDsr(diffID, headerID int64, dsr string) error {
	_, err := repo.db.Exec(insertPotDsrQuery, diffID, headerID, dsr)
	if err != nil {
		return fmt.Errorf("error inserting pot dsr %s from diff ID %d: %w", dsr, diffID, err)
	}
	return nil
}

func (repo StorageRepository) insertChi(diffID, headerID int64, chi string) error {
	_, err := repo.db.Exec(insertPotChiQuery, diffID, headerID, chi)
	if err != nil {
		return fmt.Errorf("error inserting pot chi %s from diff ID %d: %w", chi, diffID, err)
	}
	return nil
}

func (repo StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	err := repo.insertAddressID(diffID, headerID, insertPotVatQuery, vat)
	if err != nil {
		return fmt.Errorf("error inserting pot vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repo StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	err := repo.insertAddressID(diffID, headerID, insertPotVowQuery, vow)
	if err != nil {
		return fmt.Errorf("error inserting pot vow %s from diff ID %d: %w", vow, diffID, err)
	}
	return nil
}

func (repo StorageRepository) insertRho(diffID, headerID int64, rho string) error {
	_, err := repo.db.Exec(insertPotRhoQuery, diffID, headerID, rho)
	if err != nil {
		return fmt.Errorf("error inserting pot rho %s from diff ID %d: %w", rho, diffID, err)
	}
	return nil
}

func (repo StorageRepository) insertLive(diffID, headerID int64, live string) error {
	_, err := repo.db.Exec(insertPotLiveQuery, diffID, headerID, live)
	if err != nil {
		return fmt.Errorf("error inserting pot live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertAddressID(diffID, headerID int64, query, address string) error {
	tx, txErr := repo.db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction: %w", txErr)
	}
	addressID, addressErr := repository.GetOrCreateAddressInTransaction(tx, address)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("address", addressErr)
		}
		return fmt.Errorf("error getting or creating address: %w", addressErr)
	}

	_, insertErr := tx.Exec(query, diffID, headerID, addressID)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("value that is address_id", insertErr)
		}
		msg := fmt.Sprintf("error inserting pot address_id value for %s from diff ID %d", address, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
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
