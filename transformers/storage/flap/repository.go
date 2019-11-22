package flap

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery    = `INSERT INTO maker.flap_vat (header_id, address_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertGemQuery    = `INSERT INTO maker.flap_gem (header_id, address_id, gem) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertBegQuery    = `INSERT INTO maker.flap_beg (header_id, address_id, beg) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertTtlQuery    = `INSERT INTO maker.flap_ttl (header_id, address_id, ttl) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertTauQuery    = `INSERT INTO maker.flap_tau (header_id, address_id, tau) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertKicksQuery  = `INSERT INTO maker.flap_kicks (header_id, address_id, kicks) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertLiveQuery   = `INSERT INTO maker.flap_live (header_id, address_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertBidBidQuery = `INSERT INTO maker.flap_bid_bid (header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidLotQuery = `INSERT INTO maker.flap_bid_lot (header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidGuyQuery = `INSERT INTO maker.flap_bid_guy (header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidTicQuery = `INSERT INTO maker.flap_bid_tic (header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidEndQuery = `INSERT INTO maker.flap_bid_end (header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type FlapStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *FlapStorageRepository) Create(headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(headerID, value.(string))
	case storage.Gem:
		return repository.insertGem(headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(headerID, value.(string))
	case storage.Live:
		return repository.insertLive(headerID, value.(string))
	case storage.BidBid:
		return repository.insertBidBid(headerID, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(headerID, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(headerID, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flap contract storage name: %s", metadata.Name))
	}
}

func (repository *FlapStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlapStorageRepository) insertVat(headerID int64, vat string) error {
	return repository.insertRecordWithAddress(headerID, insertVatQuery, vat)
}

func (repository *FlapStorageRepository) insertGem(headerID int64, gem string) error {
	return repository.insertRecordWithAddress(headerID, insertGemQuery, gem)
}

func (repository *FlapStorageRepository) insertBeg(headerID int64, beg string) error {
	return repository.insertRecordWithAddress(headerID, insertBegQuery, beg)
}

func (repository *FlapStorageRepository) insertTtl(headerID int64, ttl string) error {
	return repository.insertRecordWithAddress(headerID, insertTtlQuery, ttl)
}

func (repository *FlapStorageRepository) insertTau(headerID int64, tau string) error {
	return repository.insertRecordWithAddress(headerID, insertTauQuery, tau)
}

func (repository *FlapStorageRepository) insertKicks(headerID int64, kicks string) error {
	return repository.insertRecordWithAddress(headerID, InsertKicksQuery, kicks)
}

func (repository *FlapStorageRepository) insertLive(headerID int64, live string) error {
	return repository.insertRecordWithAddress(headerID, insertLiveQuery, live)
}

func (repository *FlapStorageRepository) insertBidBid(headerID int64, metadata utils.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, insertBidBidQuery, bidId, bid)
}

func (repository *FlapStorageRepository) insertBidLot(headerID int64, metadata utils.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, insertBidLotQuery, bidId, lot)
}

func (repository *FlapStorageRepository) insertBidGuy(headerID int64, metadata utils.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, insertBidGuyQuery, bidId, guy)
}

func (repository *FlapStorageRepository) insertBidTic(headerID int64, metadata utils.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, insertBidTicQuery, bidId, tic)
}

func (repository *FlapStorageRepository) insertBidEnd(headerID int64, metadata utils.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, insertBidEndQuery, bidId, end)
}

func (repository *FlapStorageRepository) insertPackedValueRecord(headerID int64, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case storage.Ttl:
			insertErr = repository.insertTtl(headerID, value)
		case storage.Tau:
			insertErr = repository.insertTau(headerID, value)
		case storage.BidGuy:
			insertErr = repository.insertBidGuy(headerID, metadata, value)
		case storage.BidTic:
			insertErr = repository.insertBidTic(headerID, metadata, value)
		case storage.BidEnd:
			insertErr = repository.insertBidEnd(headerID, metadata, value)
		default:
			panic(fmt.Sprintf("unrecognized flap contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func getBidId(keys map[utils.Key]string) (string, error) {
	bid, ok := keys[constants.BidId]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bid, nil
}

func (repository *FlapStorageRepository) insertRecordWithAddress(headerID int64, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressId, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlapStorageRepository) insertRecordWithAddressAndBidId(headerID int64, query, bidId, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressId, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, headerID, addressId, bidId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("flap field with address for bid id %s", bidId)
			return shared.FormatRollbackError(errorString, insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}
