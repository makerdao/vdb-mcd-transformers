package flap

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery    = `INSERT INTO maker.flap_vat (diff_id, header_id, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertGemQuery    = `INSERT INTO maker.flap_gem (diff_id, header_id, address_id, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBegQuery    = `INSERT INTO maker.flap_beg (diff_id, header_id, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTTLQuery    = `INSERT INTO maker.flap_ttl (diff_id, header_id, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTauQuery    = `INSERT INTO maker.flap_tau (diff_id, header_id, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertKicksQuery  = `INSERT INTO maker.flap_kicks (diff_id, header_id, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertLiveQuery   = `INSERT INTO maker.flap_live (diff_id, header_id, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidBidQuery = `INSERT INTO maker.flap_bid_bid (diff_id, header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidLotQuery = `INSERT INTO maker.flap_bid_lot (diff_id, header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidGuyQuery = `INSERT INTO maker.flap_bid_guy (diff_id, header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidTicQuery = `INSERT INTO maker.flap_bid_tic (diff_id, header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidEndQuery = `INSERT INTO maker.flap_bid_end (diff_id, header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case storage.Gem:
		return repository.insertGem(diffID, headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(diffID, headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(diffID, headerID, value.(string))
	case storage.Live:
		return repository.insertLive(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case storage.BidBid:
		return repository.insertBidBid(diffID, headerID, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(diffID, headerID, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flap contract storage name: %s", metadata.Name))
	}
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertVatQuery, vat, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertGem(diffID, headerID int64, gem string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertGemQuery, gem, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap gem %s from diff ID %d: %w", gem, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertBegQuery, beg, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap beg %s from diff ID %d: %w", beg, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertTTL(diffID, headerID int64, ttl string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertTTLQuery, ttl, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap ttl %s from diff ID %d: %w", ttl, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertTau(diffID, headerID int64, tau string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertTauQuery, tau, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap tau %s from diff ID %d: %w", tau, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, InsertKicksQuery, kicks, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap kicks %s from diff ID %d: %w", kicks, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertLive(diffID, headerID int64, live string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertLiveQuery, live, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flap live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertBidBid(diffID, headerID int64, metadata types.ValueMetadata, bid string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flap bid bid: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertBidBidQuery, bidID, bid, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flap bid %s bid %s from diff ID %d", bidID, bid, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flap bid lot: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertBidLotQuery, bidID, lot, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flap bid %s lot %s from diff ID %d", bidID, lot, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidGuy(diffID, headerID int64, metadata types.ValueMetadata, guy string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flap bid guy: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertBidGuyQuery, bidID, guy, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flap bid %s guy %s from diff ID %d", bidID, guy, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flap bid tic: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertBidTicQuery, bidID, tic, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flap bid %s tic %s from diff ID %d", bidID, tic, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidEnd(diffID, headerID int64, metadata types.ValueMetadata, end string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flap bid end: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertBidEndQuery, bidID, end, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flap bid %s end %s from diff ID %d", bidID, end, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case storage.Ttl:
			insertErr = repository.insertTTL(diffID, headerID, value)
		case storage.Tau:
			insertErr = repository.insertTau(diffID, headerID, value)
		case storage.BidGuy:
			insertErr = repository.insertBidGuy(diffID, headerID, metadata, value)
		case storage.BidTic:
			insertErr = repository.insertBidTic(diffID, headerID, metadata, value)
		case storage.BidEnd:
			insertErr = repository.insertBidEnd(diffID, headerID, metadata, value)
		default:
			panic(fmt.Sprintf("unrecognized flap contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return fmt.Errorf("error inserting flap packed value from diff ID %d: %w", diffID, insertErr)
		}
	}
	return nil
}

func getBidID(keys map[types.Key]string) (string, error) {
	bid, ok := keys[constants.BidId]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bid, nil
}
