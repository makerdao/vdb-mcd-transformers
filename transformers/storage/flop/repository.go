package flop

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	insertFlopVatQuery   = `INSERT INTO maker.flop_vat (block_number, block_hash, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopGemQuery   = `INSERT INTO maker.flop_gem (block_number, block_hash, address_id, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopBegQuery   = `INSERT INTO maker.flop_beg (block_number, block_hash, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopPadQuery   = `INSERT INTO maker.flop_pad (block_number, block_hash, address_id, pad) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopTtlQuery   = `INSERT INTO maker.flop_ttl (block_number, block_hash, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopTauQuery   = `INSERT INTO maker.flop_tau (block_number, block_hash, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopKicksQuery = `INSERT INTO maker.flop_kicks (block_number, block_hash, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopLiveQuery  = `INSERT INTO maker.flop_live (block_number, block_hash, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	InsertFlopBidBidQuery = `INSERT INTO maker.flop_bid_bid (block_number, block_hash, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidLotQuery = `INSERT INTO Maker.flop_bid_lot (block_number, block_hash, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidGuyQuery = `INSERT INTO Maker.flop_bid_guy (block_number, block_hash, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidTicQuery = `INSERT INTO Maker.flop_bid_tic (block_number, block_hash, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidEndQuery = `INSERT INTO Maker.flop_bid_end (block_number, block_hash, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type FlopStorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *FlopStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(blockNumber, blockHash, value.(string))
	case storage.Gem:
		return repository.insertGem(blockNumber, blockHash, value.(string))
	case storage.Beg:
		return repository.insertBeg(blockNumber, blockHash, value.(string))
	case storage.Pad:
		return repository.insertPad(blockNumber, blockHash, value.(string))
	case storage.Kicks:
		return repository.insertKicks(blockNumber, blockHash, value.(string))
	case storage.Live:
		return repository.insertLive(blockNumber, blockHash, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(blockNumber, blockHash, metadata, value.(map[int]string))
	case storage.BidBid:
		return repository.insertBidBid(blockNumber, blockHash, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(blockNumber, blockHash, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized flop contract storage name: %s", metadata.Name))
	}
}

func (repository *FlopStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlopStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopVatQuery, vat)
}

func (repository *FlopStorageRepository) insertGem(blockNumber int, blockHash string, gem string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopGemQuery, gem)
}

func (repository *FlopStorageRepository) insertBeg(blockNumber int, blockHash string, beg string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopBegQuery, beg)
}

func (repository *FlopStorageRepository) insertPad(blockNumber int, blockHash string, pad string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopPadQuery, pad)
}

func (repository *FlopStorageRepository) insertTtl(blockNumber int, blockHash string, ttl string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopTtlQuery, ttl)
}

func (repository *FlopStorageRepository) insertTau(blockNumber int, blockHash string, tau string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopTauQuery, tau)
}

func (repository *FlopStorageRepository) insertKicks(blockNumber int, blockHash string, kicks string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, InsertFlopKicksQuery, kicks)
}

func (repository *FlopStorageRepository) insertLive(blockNumber int, blockHash string, live string) error {
	return repository.insertRecordWithAddress(blockNumber, blockHash, insertFlopLiveQuery, live)
}

func (repository *FlopStorageRepository) insertBidBid(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(blockNumber, blockHash, InsertFlopBidBidQuery, bidId, bid)
}

func (repository *FlopStorageRepository) insertBidLot(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(blockNumber, blockHash, InsertFlopBidLotQuery, bidId, lot)
}

func (repository *FlopStorageRepository) insertBidGuy(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(blockNumber, blockHash, InsertFlopBidGuyQuery, bidId, guy)
}

func (repository *FlopStorageRepository) insertBidTic(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(blockNumber, blockHash, InsertFlopBidTicQuery, bidId, tic)
}

func (repository *FlopStorageRepository) insertBidEnd(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(blockNumber, blockHash, InsertFlopBidEndQuery, bidId, end)
}

func (repository *FlopStorageRepository) insertPackedValueRecord(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
	var insertErr error
	for order, value := range packedValues {
		switch metadata.PackedNames[order] {
		case storage.Ttl:
			insertErr = repository.insertTtl(blockNumber, blockHash, value)
		case storage.Tau:
			insertErr = repository.insertTau(blockNumber, blockHash, value)
		case storage.BidGuy:
			insertErr = repository.insertBidGuy(blockNumber, blockHash, metadata, value)
		case storage.BidTic:
			insertErr = repository.insertBidTic(blockNumber, blockHash, metadata, value)
		case storage.BidEnd:
			insertErr = repository.insertBidEnd(blockNumber, blockHash, metadata, value)
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

func (repository *FlopStorageRepository) insertRecordWithAddress(blockNumber int, blockHash, query, value string) error {
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
	_, insertErr := tx.Exec(query, blockNumber, blockHash, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flop field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlopStorageRepository) insertRecordWithAddressAndBidId(blockNumber int, blockHash, query, bidId, value string) error {
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
	_, insertErr := tx.Exec(query, blockNumber, blockHash, addressId, bidId, value)
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
