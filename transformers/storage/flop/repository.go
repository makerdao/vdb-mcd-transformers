package flop

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertFlopVatQuery   = `INSERT INTO maker.flop_vat (header_id, address_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopGemQuery   = `INSERT INTO maker.flop_gem (header_id, address_id, gem) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopBegQuery   = `INSERT INTO maker.flop_beg (header_id, address_id, beg) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopPadQuery   = `INSERT INTO maker.flop_pad (header_id, address_id, pad) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopTtlQuery   = `INSERT INTO maker.flop_ttl (header_id, address_id, ttl) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopTauQuery   = `INSERT INTO maker.flop_tau (header_id, address_id, tau) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertFlopKicksQuery = `INSERT INTO maker.flop_kicks (header_id, address_id, kicks) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopLiveQuery  = `INSERT INTO maker.flop_live (header_id, address_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	InsertFlopBidBidQuery = `INSERT INTO maker.flop_bid_bid (header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopBidLotQuery = `INSERT INTO Maker.flop_bid_lot (header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopBidGuyQuery = `INSERT INTO Maker.flop_bid_guy (header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopBidTicQuery = `INSERT INTO Maker.flop_bid_tic (header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopBidEndQuery = `INSERT INTO Maker.flop_bid_end (header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type FlopStorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *FlopStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(headerID, value.(string))
	case storage.Gem:
		return repository.insertGem(headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(headerID, value.(string))
	case storage.Pad:
		return repository.insertPad(headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(headerID, value.(string))
	case storage.Live:
		return repository.insertLive(headerID, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(headerID, metadata, value.(map[int]string))
	case storage.BidBid:
		return repository.insertBidBid(headerID, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(headerID, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized flop contract storage name: %s", metadata.Name))
	}
}

func (repository *FlopStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlopStorageRepository) insertVat(headerID int64, vat string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopVatQuery, vat)
}

func (repository *FlopStorageRepository) insertGem(headerID int64, gem string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopGemQuery, gem)
}

func (repository *FlopStorageRepository) insertBeg(headerID int64, beg string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopBegQuery, beg)
}

func (repository *FlopStorageRepository) insertPad(headerID int64, pad string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopPadQuery, pad)
}

func (repository *FlopStorageRepository) insertTtl(headerID int64, ttl string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopTtlQuery, ttl)
}

func (repository *FlopStorageRepository) insertTau(headerID int64, tau string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopTauQuery, tau)
}

func (repository *FlopStorageRepository) insertKicks(headerID int64, kicks string) error {
	return repository.insertRecordWithAddress(headerID, InsertFlopKicksQuery, kicks)
}

func (repository *FlopStorageRepository) insertLive(headerID int64, live string) error {
	return repository.insertRecordWithAddress(headerID, insertFlopLiveQuery, live)
}

func (repository *FlopStorageRepository) insertBidBid(headerID int64, metadata utils.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlopBidBidQuery, bidId, bid)
}

func (repository *FlopStorageRepository) insertBidLot(headerID int64, metadata utils.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlopBidLotQuery, bidId, lot)
}

func (repository *FlopStorageRepository) insertBidGuy(headerID int64, metadata utils.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlopBidGuyQuery, bidId, guy)
}

func (repository *FlopStorageRepository) insertBidTic(headerID int64, metadata utils.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlopBidTicQuery, bidId, tic)
}

func (repository *FlopStorageRepository) insertBidEnd(headerID int64, metadata utils.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlopBidEndQuery, bidId, end)
}

func (repository *FlopStorageRepository) insertPackedValueRecord(headerID int64, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
	var insertErr error
	for order, value := range packedValues {
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
			panic(fmt.Sprintf("unrecognized flop contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func getBidId(keys map[utils.Key]string) (string, error) {
	bidId, ok := keys[constants.BidId]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bidId, nil
}

func (repository *FlopStorageRepository) insertRecordWithAddress(headerID int64, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressId, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flop address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flop field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlopStorageRepository) insertRecordWithAddressAndBidId(headerID int64, query, bidId, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressId, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, repository.db)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flop address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, headerID, addressId, bidId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("flop field with address for bid id %s", bidId)
			return shared.FormatRollbackError(errorString, insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}
