package wards

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var insertWardsQuery = `INSERT INTO maker.wards (diff_id, header_id, address_id, usr, wards) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

func InsertWards(diffID, headerID int64, metadata storage.ValueMetadata, contractAddress, wards string, db *postgres.DB) error {
	user, userErr := getUser(metadata.Keys)
	if userErr != nil {
		return userErr
	}

	tx, txErr := db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressID, addressErr := shared.GetOrCreateAddress(contractAddress, db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("contract address", addressErr.Error())
		}
		return addressErr
	}

	userAddressID, userAddressErr := shared.GetOrCreateAddress(user, db)
	if userAddressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("wards user address", userAddressErr.Error())
		}
		return addressErr
	}

	_, insertErr := tx.Exec(insertWardsQuery, diffID, headerID, addressID, userAddressID, wards)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("wards record with address", insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}

func getUser(keys map[storage.Key]string) (string, error) {
	user, ok := keys[constants.User]
	if !ok {
		return "", storage.ErrMetadataMalformed{MissingData: constants.User}
	}
	return user, nil
}
