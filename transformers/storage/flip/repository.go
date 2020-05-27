package flip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertFlipVatQuery   = `INSERT INTO maker.flip_vat (diff_id, header_id, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipIlkQuery   = `INSERT INTO maker.flip_ilk (diff_id, header_id, address_id, ilk_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipBegQuery   = `INSERT INTO maker.flip_beg (diff_id, header_id, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipTtlQuery   = `INSERT INTO maker.flip_ttl (diff_id, header_id, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipTauQuery   = `INSERT INTO maker.flip_tau (diff_id, header_id, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipKicksQuery = `INSERT INTO maker.flip_kicks (diff_id, header_id, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	InsertFlipBidBidQuery = `INSERT INTO maker.flip_bid_bid (diff_id, header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidLotQuery = `INSERT INTO maker.flip_bid_lot (diff_id, header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidGuyQuery = `INSERT INTO maker.flip_bid_guy (diff_id, header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidTicQuery = `INSERT INTO maker.flip_bid_tic (diff_id, header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidEndQuery = `INSERT INTO maker.flip_bid_end (diff_id, header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidUsrQuery = `INSERT INTO maker.flip_bid_usr (diff_id, header_id, address_id, bid_id, usr) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidGalQuery = `INSERT INTO maker.flip_bid_gal (diff_id, header_id, address_id, bid_id, gal) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidTabQuery = `INSERT INTO maker.flip_bid_tab (diff_id, header_id, address_id, bid_id, tab) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type FlipStorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *FlipStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case storage.Ilk:
		return repository.insertIlk(diffID, headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(diffID, headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case storage.BidBid:
		return repository.insertBidBid(diffID, headerID, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(diffID, headerID, metadata, value.(string))
	case storage.BidUsr:
		return repository.insertBidUsr(diffID, headerID, metadata, value.(string))
	case storage.BidGal:
		return repository.insertBidGal(diffID, headerID, metadata, value.(string))
	case storage.BidTab:
		return repository.insertBidTab(diffID, headerID, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flip contract storage name: %s", metadata.Name))
	}
}

func (repository *FlipStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlipStorageRepository) insertVat(diffID, headerID int64, vat string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlipVatQuery, vat, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertIlk(diffID, headerID int64, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repository.db)
	if ilkErr != nil {
		return ilkErr
	}

	return shared.InsertRecordWithAddress(diffID, headerID, insertFlipIlkQuery, strconv.FormatInt(ilkID, 10), repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlipBegQuery, beg, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertTtl(diffID, headerID int64, ttl string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlipTtlQuery, ttl, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertTau(diffID, headerID int64, tau string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, insertFlipTauQuery, tau, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	return shared.InsertRecordWithAddress(diffID, headerID, InsertFlipKicksQuery, kicks, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidBid(diffID, headerID int64, metadata types.ValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidBidQuery, bidId, bid, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidLotQuery, bidId, lot, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidGuy(diffID, headerID int64, metadata types.ValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidGuyQuery, bidId, guy, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidTicQuery, bidId, tic, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidEnd(diffID, headerID int64, metadata types.ValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidEndQuery, bidId, end, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidUsr(diffID, headerID int64, metadata types.ValueMetadata, usr string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidUsrQuery, bidId, usr, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidGal(diffID, headerID int64, metadata types.ValueMetadata, gal string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidGalQuery, bidId, gal, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertBidTab(diffID, headerID int64, metadata types.ValueMetadata, tab string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return shared.InsertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidTabQuery, bidId, tab, repository.ContractAddress, repository.db)
}

func (repository *FlipStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
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
			panic(fmt.Sprintf("unrecognized flip contract storage name in packed values: %s", metadata.Name))
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
