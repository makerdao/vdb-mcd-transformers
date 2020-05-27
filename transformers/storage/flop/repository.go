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
	insertFlopTtlQuery   = `INSERT INTO maker.flop_ttl (diff_id, header_id, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
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

type FlopStorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *FlopStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
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

func (repository *FlopStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlopStorageRepository) insertVat(diffID, headerID int64, vat string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopVatQuery, vat, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertGem(diffID, headerID int64, gem string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopGemQuery, gem, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopBegQuery, beg, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertPad(diffID, headerID int64, pad string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopPadQuery, pad, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertTtl(diffID, headerID int64, ttl string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopTtlQuery, ttl, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertTau(diffID, headerID int64, tau string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopTauQuery, tau, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, InsertFlopKicksQuery, kicks, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertLive(diffID, headerID int64, live string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopLiveQuery, live, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertVow(diffID, headerID int64, vow string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlopVowQuery, vow, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertBidBid(diffID, headerID int64, metadata types.ValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlopBidBidQuery, bidId, bid, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertBidLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlopBidLotQuery, bidId, lot, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertBidGuy(diffID, headerID int64, metadata types.ValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlopBidGuyQuery, bidId, guy, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertBidTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlopBidTicQuery, bidId, tic, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertBidEnd(diffID, headerID int64, metadata types.ValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlopBidEndQuery, bidId, end, repository.ContractAddress, repository.db)
}

func (repository *FlopStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	var insertErr error
	for order, value := range packedValues {
		switch metadata.PackedNames[order] {
		case storage.Ttl:
			insertErr = repository.insertTtl(diffID, headerID, value)
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
			return insertErr
		}
	}
	return nil
}

func getBidId(keys map[types.Key]string) (string, error) {
	bidId, ok := keys[constants.BidId]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bidId, nil
}
