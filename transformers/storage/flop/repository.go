package flop

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
	insertFlopVatQuery   = `INSERT INTO maker.flop_vat (diff_id, header_id, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopGemQuery   = `INSERT INTO maker.flop_gem (diff_id, header_id, address_id, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopBegQuery   = `INSERT INTO maker.flop_beg (diff_id, header_id, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopPadQuery   = `INSERT INTO maker.flop_pad (diff_id, header_id, address_id, pad) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopTTLQuery   = `INSERT INTO maker.flop_ttl (diff_id, header_id, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopTauQuery   = `INSERT INTO maker.flop_tau (diff_id, header_id, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlopKicksQuery = `INSERT INTO maker.flop_kicks (diff_id, header_id, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopLiveQuery  = `INSERT INTO maker.flop_live (diff_id, header_id, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlopVowQuery   = `INSERT INTO maker.flop_vow (diff_id, header_id, address_id, vow) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	InsertFlopBidBidQuery = `INSERT INTO maker.flop_bid_bid (diff_id, header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidLotQuery = `INSERT INTO Maker.flop_bid_lot (diff_id, header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidGuyQuery = `INSERT INTO Maker.flop_bid_guy (diff_id, header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidTicQuery = `INSERT INTO Maker.flop_bid_tic (diff_id, header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlopBidEndQuery = `INSERT INTO Maker.flop_bid_end (diff_id, header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case storage.Gem:
		return repository.insertGem(diffID, headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(diffID, headerID, value.(string))
	case storage.Pad:
		return repository.insertPad(diffID, headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(diffID, headerID, value.(string))
	case storage.Live:
		return repository.insertLive(diffID, headerID, value.(string))
	case storage.Vow:
		return repository.insertVow(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case storage.Packed:
		return repository.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	case storage.BidBid:
		return repository.insertBidBid(diffID, headerID, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(diffID, headerID, metadata, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized flop contract storage name: %s", metadata.Name))
	}
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopVatQuery, vat, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertGem(diffID, headerID int64, gem string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopGemQuery, gem, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop gem %s from diff ID %d: %w", gem, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopBegQuery, beg, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop beg %s from diff ID %d: %w", beg, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertPad(diffID, headerID int64, pad string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopPadQuery, pad, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop pad %s from diff ID %d: %w", pad, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertTTL(diffID, headerID int64, ttl string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopTTLQuery, ttl, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop ttl %s from diff ID %d: %w", ttl, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertTau(diffID, headerID int64, tau string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopTauQuery, tau, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop tau %s from diff ID %d: %w", tau, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, InsertFlopKicksQuery, kicks, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop kicks %s from diff ID %d: %w", kicks, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertLive(diffID, headerID int64, live string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopLiveQuery, live, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	err := shared.InsertRecordWithAddress(diffID, headerID, insertFlopVowQuery, vow, repository.ContractAddress, repository.db)
	if err != nil {
		return fmt.Errorf("error inserting flop vow %s from diff ID %d: %w", vow, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertBidBid(diffID, headerID int64, metadata types.ValueMetadata, bid string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for bid bid: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, InsertFlopBidBidQuery, bidID, bid, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flop bid %s bid %s from diff ID %d", bidID, bid, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for bid lot: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, InsertFlopBidLotQuery, bidID, lot, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flop bid %s lot %s from diff ID %d", bidID, lot, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidGuy(diffID, headerID int64, metadata types.ValueMetadata, guy string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for bid guy: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, InsertFlopBidGuyQuery, bidID, guy, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flop bid %s guy %s from diff ID %d", bidID, guy, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flop bid tic: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, InsertFlopBidTicQuery, bidID, tic, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flop bid %s tic %s from diff ID %d", bidID, tic, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidEnd(diffID, headerID int64, metadata types.ValueMetadata, end string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flop bid end: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, InsertFlopBidEndQuery, bidID, end, repository.ContractAddress, repository.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting flop bid %s end %s from diff ID %d", bidID, end, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	var insertErr error
	for order, value := range packedValues {
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
			panic(fmt.Sprintf("unrecognized flop contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return fmt.Errorf("error inserting flop packed value from diff ID %d: %w", diffID, insertErr)
		}
	}
	return nil
}

func getBidID(keys map[types.Key]string) (string, error) {
	bidID, ok := keys[constants.BidId]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bidID, nil
}
