package flap

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
)

const (
	insertVatQuery    = `INSERT INTO maker.flap_vat (block_number, block_hash, contract_address, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertGemQuery    = `INSERT INTO maker.flap_gem (block_number, block_hash, contract_address, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBegQuery    = `INSERT INTO maker.flap_beg (block_number, block_hash, contract_address, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTtlQuery    = `INSERT INTO maker.flap_ttl (block_number, block_hash, contract_address, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTauQuery    = `INSERT INTO maker.flap_tau (block_number, block_hash, contract_address, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertKicksQuery  = `INSERT INTO maker.flap_kicks (block_number, block_hash, contract_address, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertLiveQuery   = `INSERT INTO maker.flap_live (block_number, block_hash, contract_address, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidBidQuery = `INSERT INTO maker.flap_bid_bid (block_number, block_hash, contract_address, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidLotQuery = `INSERT INTO maker.flap_bid_lot (block_number, block_hash, contract_address, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidGuyQuery = `INSERT INTO maker.flap_bid_guy (block_number, block_hash, contract_address, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidTicQuery = `INSERT INTO maker.flap_bid_tic (block_number, block_hash, contract_address, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidEndQuery = `INSERT INTO maker.flap_bid_end (block_number, block_hash, contract_address, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidGalQuery = `INSERT INTO maker.flap_bid_gal (block_number, block_hash, contract_address, bid_id, gal) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type FlapStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *FlapStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
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
	case storage.BidBid:
		return repository.insertBidBid(blockNumber, blockHash, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(blockNumber, blockHash, metadata, value.(string))
	case storage.BidGal:
		return repository.insertBidGal(blockNumber, blockHash, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(blockNumber, blockHash, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flap contract storage name: %s", metadata.Name))
	}
}

func (repository *FlapStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlapStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, blockNumber, blockHash, repository.ContractAddress, vat)
	return err
}

func (repository *FlapStorageRepository) insertGem(blockNumber int, blockHash string, gem string) error {
	_, err := repository.db.Exec(insertGemQuery, blockNumber, blockHash, repository.ContractAddress, gem)
	return err
}

func (repository *FlapStorageRepository) insertBeg(blockNumber int, blockHash string, beg string) error {
	_, err := repository.db.Exec(insertBegQuery, blockNumber, blockHash, repository.ContractAddress, beg)
	return err
}

func (repository *FlapStorageRepository) insertTtl(blockNumber int, blockHash string, ttl string) error {
	_, writeErr := repository.db.Exec(insertTtlQuery, blockNumber, blockHash, repository.ContractAddress, ttl)
	return writeErr
}

func (repository *FlapStorageRepository) insertTau(blockNumber int, blockHash string, tau string) error {
	_, writeErr := repository.db.Exec(insertTauQuery, blockNumber, blockHash, repository.ContractAddress, tau)
	return writeErr
}

func (repository *FlapStorageRepository) insertKicks(blockNumber int, blockHash string, kicks string) error {
	_, err := repository.db.Exec(InsertKicksQuery, blockNumber, blockHash, repository.ContractAddress, kicks)
	return err
}

func (repository *FlapStorageRepository) insertLive(blockNumber int, blockHash string, live string) error {
	_, err := repository.db.Exec(insertLiveQuery, blockNumber, blockHash, repository.ContractAddress, live)
	return err
}

func (repository *FlapStorageRepository) insertBidBid(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, err = repository.db.Exec(insertBidBidQuery, blockNumber, blockHash, repository.ContractAddress, bidId, bid)
	return err
}

func (repository *FlapStorageRepository) insertBidLot(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, err = repository.db.Exec(insertBidLotQuery, blockNumber, blockHash, repository.ContractAddress, bidId, lot)
	return err
}

func (repository *FlapStorageRepository) insertBidGuy(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, err = repository.db.Exec(insertBidGuyQuery, blockNumber, blockHash, repository.ContractAddress, bidId, guy)
	return err
}

func (repository *FlapStorageRepository) insertBidTic(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertBidTicQuery, blockNumber, blockHash, repository.ContractAddress, bidId, tic)
	return writeErr
}

func (repository *FlapStorageRepository) insertBidEnd(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertBidEndQuery, blockNumber, blockHash, repository.ContractAddress, bidId, end)
	return writeErr
}

func (repository *FlapStorageRepository) insertBidGal(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, gal string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertBidGalQuery, blockNumber, blockHash, repository.ContractAddress, bidId, gal)
	return writeErr
}

func (repository *FlapStorageRepository) insertPackedValueRecord(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
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
