package flop

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	insertFlopVatQuery   = `INSERT INTO maker.flop_vat (block_number, block_hash, contract_address, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopGemQuery   = `INSERT INTO maker.flop_gem (block_number, block_hash, contract_address, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopBegQuery   = `INSERT INTO maker.flop_beg (block_number, block_hash, contract_address, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopTtlQuery   = `INSERT INTO maker.flop_ttl (block_number, block_hash, contract_address, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopTauQuery   = `INSERT INTO maker.flop_tau (block_number, block_hash, contract_address, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopKicksQuery = `INSERT INTO maker.flop_kicks (block_number, block_hash, contract_address, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopLiveQuery  = `INSERT INTO maker.flop_live (block_number, block_hash, contract_address, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	InsertFlopBidBidQuery = `INSERT INTO maker.flop_bid_bid (block_number, block_hash, contract_address, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidLotQuery = `INSERT INTO Maker.flop_bid_lot (block_number, block_hash, contract_address, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidGuyQuery = `INSERT INTO Maker.flop_bid_guy (block_number, block_hash, contract_address, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidTicQuery = `INSERT INTO Maker.flop_bid_tic (block_number, block_hash, contract_address, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidEndQuery = `INSERT INTO Maker.flop_bid_end (block_number, block_hash, contract_address, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
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
	_, writeErr := repository.db.Exec(insertFlopVatQuery, blockNumber, blockHash, repository.ContractAddress, vat)
	return writeErr
}

func (repository *FlopStorageRepository) insertGem(blockNumber int, blockHash string, gem string) error {
	_, writeErr := repository.db.Exec(insertFlopGemQuery, blockNumber, blockHash, repository.ContractAddress, gem)
	return writeErr
}

func (repository *FlopStorageRepository) insertBeg(blockNumber int, blockHash string, beg string) error {
	_, writeErr := repository.db.Exec(insertFlopBegQuery, blockNumber, blockHash, repository.ContractAddress, beg)
	return writeErr
}

func (repository *FlopStorageRepository) insertTtl(blockNumber int, blockHash string, ttl string) error {
	_, writeErr := repository.db.Exec(insertFlopTtlQuery, blockNumber, blockHash, repository.ContractAddress, ttl)
	return writeErr
}

func (repository *FlopStorageRepository) insertTau(blockNumber int, blockHash string, tau string) error {
	_, writeErr := repository.db.Exec(insertFlopTauQuery, blockNumber, blockHash, repository.ContractAddress, tau)
	return writeErr
}

func (repository *FlopStorageRepository) insertKicks(blockNumber int, blockHash string, kicks string) error {
	_, writeErr := repository.db.Exec(InsertFlopKicksQuery, blockNumber, blockHash, repository.ContractAddress, kicks)
	return writeErr
}

func (repository *FlopStorageRepository) insertLive(blockNumber int, blockHash string, live string) error {
	_, writeErr := repository.db.Exec(insertFlopLiveQuery, blockNumber, blockHash, repository.ContractAddress, live)
	return writeErr
}

func (repository *FlopStorageRepository) insertBidBid(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(InsertFlopBidBidQuery, blockNumber, blockHash, repository.ContractAddress, bidId, bid)
	return writeErr
}

func (repository *FlopStorageRepository) insertBidLot(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(InsertFlopBidLotQuery, blockNumber, blockHash, repository.ContractAddress, bidId, lot)
	return writeErr
}

func (repository *FlopStorageRepository) insertBidGuy(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(InsertFlopBidGuyQuery, blockNumber, blockHash, repository.ContractAddress, bidId, guy)
	return writeErr
}

func (repository *FlopStorageRepository) insertBidTic(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(InsertFlopBidTicQuery, blockNumber, blockHash, repository.ContractAddress, bidId, tic)
	return writeErr
}

func (repository *FlopStorageRepository) insertBidEnd(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(InsertFlopBidEndQuery, blockNumber, blockHash, repository.ContractAddress, bidId, end)
	return writeErr
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
