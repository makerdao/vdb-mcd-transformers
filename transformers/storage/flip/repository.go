package flip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
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

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repo.insertVat(diffID, headerID, value.(string))
	case storage.Ilk:
		return repo.insertIlk(diffID, headerID, value.(string))
	case storage.Beg:
		return repo.insertBeg(diffID, headerID, value.(string))
	case storage.Kicks:
		return repo.insertKicks(diffID, headerID, value.(string))
	case storage.Cat:
		return repo.insertCat(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	case storage.BidBid:
		return repo.insertBidBid(diffID, headerID, metadata, value.(string))
	case storage.BidLot:
		return repo.insertBidLot(diffID, headerID, metadata, value.(string))
	case storage.BidUsr:
		return repo.insertBidUsr(diffID, headerID, metadata, value.(string))
	case storage.BidGal:
		return repo.insertBidGal(diffID, headerID, metadata, value.(string))
	case storage.BidTab:
		return repo.insertBidTab(diffID, headerID, metadata, value.(string))
	case storage.Packed:
		return repo.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	default:
		return fmt.Errorf("unrecognized flip contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipVatQuery,
		vat,
		repo.ContractAddress,
		repo.db)
	if err != nil {
		msgToFormat := "error inserting flip %s vat %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, vat, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repo *StorageRepository) insertIlk(diffID, headerID int64, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repo.db)
	if ilkErr != nil {
		return fmt.Errorf("error getting or creating ilk for flip ilk: %w", ilkErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipIlkQuery,
		strconv.FormatInt(ilkID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s ilk %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, ilk, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipBegQuery,
		beg,
		repo.ContractAddress,
		repo.db)
	if err != nil {
		msgToFormat := "error inserting flip %s beg %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, beg, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repo *StorageRepository) insertTTL(diffID, headerID int64, ttl string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipTTLQuery,
		ttl,
		repo.ContractAddress,
		repo.db)
	if err != nil {
		msgToFormat := "error inserting flip %s ttl %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, ttl, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repo *StorageRepository) insertTau(diffID, headerID int64, tau string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipTauQuery,
		tau,
		repo.ContractAddress,
		repo.db)
	if err != nil {
		msgToFormat := "error inserting flip %s tau %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, tau, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repo *StorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	err := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		InsertFlipKicksQuery,
		kicks,
		repo.ContractAddress,
		repo.db)
	if err != nil {
		msgToFormat := "error inserting flip %s kicks %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, kicks, diffID)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (repo *StorageRepository) insertCat(diffID, headerID int64, cat string) error {
	catAddressID, addressErr := repository.GetOrCreateAddress(repo.db, cat)
	if addressErr != nil {
		return fmt.Errorf("error inserting flip cat: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertFlipCatQuery,
		strconv.FormatInt(catAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s cat %s from diff ID %d: %w"
		return fmt.Errorf(msgToFormat, repo.ContractAddress, cat, diffID, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidBid(diffID, headerID int64, metadata types.ValueMetadata, bid string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s bid %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, bid, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s lot %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, lot, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidGuy(diffID, headerID int64, metadata types.ValueMetadata, guy string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s guy %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, guy, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s tic %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, tic, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidEnd(diffID, headerID int64, metadata types.ValueMetadata, end string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s end %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, end, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidUsr(diffID, headerID int64, metadata types.ValueMetadata, usr string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s usr %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, usr, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidGal(diffID, headerID int64, metadata types.ValueMetadata, gal string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s gal %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, gal, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBidTab(diffID, headerID int64, metadata types.ValueMetadata, tab string) error {
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
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting flip %s bid %s tab %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, bidID, tab, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case storage.Ttl:
			insertErr = repo.insertTTL(diffID, headerID, value)
		case storage.Tau:
			insertErr = repo.insertTau(diffID, headerID, value)
		case storage.BidGuy:
			insertErr = repo.insertBidGuy(diffID, headerID, metadata, value)
		case storage.BidTic:
			insertErr = repo.insertBidTic(diffID, headerID, metadata, value)
		case storage.BidEnd:
			insertErr = repo.insertBidEnd(diffID, headerID, metadata, value)
		default:
			return fmt.Errorf("unrecognized flip contract storage name in packed values: %s", metadata.Name)
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
