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
	insertFlipTTLQuery   = `INSERT INTO maker.flip_ttl (diff_id, header_id, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipTauQuery   = `INSERT INTO maker.flip_tau (diff_id, header_id, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipKicksQuery = `INSERT INTO maker.flip_kicks (diff_id, header_id, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFlipCatQuery   = `INSERT INTO maker.flip_cat (diff_id, header_id, address_id, cat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	InsertFlipBidBidQuery = `INSERT INTO maker.flip_bid_bid (diff_id, header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidLotQuery = `INSERT INTO maker.flip_bid_lot (diff_id, header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidGuyQuery = `INSERT INTO maker.flip_bid_guy (diff_id, header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidTicQuery = `INSERT INTO maker.flip_bid_tic (diff_id, header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidEndQuery = `INSERT INTO maker.flip_bid_end (diff_id, header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidUsrQuery = `INSERT INTO maker.flip_bid_usr (diff_id, header_id, address_id, bid_id, usr) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidGalQuery = `INSERT INTO maker.flip_bid_gal (diff_id, header_id, address_id, bid_id, gal) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertFlipBidTabQuery = `INSERT INTO maker.flip_bid_tab (diff_id, header_id, address_id, bid_id, tab) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case storage.Ilk:
		return repository.insertIlk(diffID, headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(diffID, headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(diffID, headerID, value.(string))
	case storage.Cat:
		return repository.insertCat(diffID, headerID, value.(string))
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

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipVatQuery,
		vat,
		repository.ContractAddress,
		repository.db)
	if err != nil {
		msgToFormat := "error inserting flip %s vat %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, vat, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repository *StorageRepository) insertIlk(diffID, headerID int64, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repository.db)
	if ilkErr != nil {
		return fmt.Errorf("error getting or creating ilk for flip ilk: %w", ilkErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipIlkQuery,
		strconv.FormatInt(ilkID, 10),
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s ilk %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, ilk, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipBegQuery,
		beg,
		repository.ContractAddress,
		repository.db)
	if err != nil {
		msgToFormat := "error inserting flip %s beg %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, beg, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repository *StorageRepository) insertTTL(diffID, headerID int64, ttl string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipTTLQuery,
		ttl,
		repository.ContractAddress,
		repository.db)
	if err != nil {
		msgToFormat := "error inserting flip %s ttl %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, ttl, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repository *StorageRepository) insertTau(diffID, headerID int64, tau string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipTauQuery,
		tau,
		repository.ContractAddress,
		repository.db)
	if err != nil {
		msgToFormat := "error inserting flip %s tau %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, tau, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repository *StorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		InsertFlipKicksQuery,
		kicks,
		repository.ContractAddress,
		repository.db)
	if err != nil {
		msgToFormat := "error inserting flip %s kicks %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, kicks, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repository *StorageRepository) insertCat(diffID, headerID int64, cat string) error {
	catAddressID, addressErr := shared.GetOrCreateAddress(cat, repository.db)
	if addressErr != nil {
		return fmt.Errorf("error inserting flip cat: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipCatQuery,
		strconv.FormatInt(catAddressID, 10),
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s cat %s from diff ID %d: %w"
		return fmt.Errorf(msgToFormat, repository.ContractAddress, cat, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidBid(diffID, headerID int64, metadata types.ValueMetadata, bid string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid bid: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidBidQuery,
		bidID,
		bid,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s bid %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, bid, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid lot: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidLotQuery,
		bidID,
		lot,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s lot %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, lot, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidGuy(diffID, headerID int64, metadata types.ValueMetadata, guy string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid guy: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidGuyQuery,
		bidID,
		guy,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s guy %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, guy, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid tic: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidTicQuery,
		bidID,
		tic,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s tic %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, tic, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidEnd(diffID, headerID int64, metadata types.ValueMetadata, end string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid end: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidEndQuery,
		bidID,
		end,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s end %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, end, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidUsr(diffID, headerID int64, metadata types.ValueMetadata, usr string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid usr: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidUsrQuery,
		bidID,
		usr,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s usr %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, usr, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidGal(diffID, headerID int64, metadata types.ValueMetadata, gal string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid gal: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidGalQuery,
		bidID,
		gal,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s gal %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, gal, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertBidTab(diffID, headerID int64, metadata types.ValueMetadata, tab string) error {
	bidID, err := getBidID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting bid ID for flip bid tab: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(
		diffID,
		headerID,
		InsertFlipBidTabQuery,
		bidID,
		tab,
		repository.ContractAddress,
		repository.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s tab %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repository.ContractAddress, bidID, tab, diffID)
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
			panic(fmt.Sprintf("unrecognized flip contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return fmt.Errorf("error inserting flip packed value from diff ID %d: %w", diffID, insertErr)
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
