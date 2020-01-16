package flip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
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

func (repository *FlipStorageRepository) Create(diffID, headerID int64, metadata vdbStorage.ValueMetadata, value interface{}) error {
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
	return repository.insertRecordWithAddress(diffID, headerID, insertFlipVatQuery, vat)
}

func (repository *FlipStorageRepository) insertIlk(diffID, headerID int64, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repository.db)
	if ilkErr != nil {
		return ilkErr
	}

	return repository.insertRecordWithAddress(diffID, headerID, insertFlipIlkQuery, strconv.FormatInt(ilkID, 10))
}

func (repository *FlipStorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertFlipBegQuery, beg)
}

func (repository *FlipStorageRepository) insertTtl(diffID, headerID int64, ttl string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertFlipTtlQuery, ttl)
}

func (repository *FlipStorageRepository) insertTau(diffID, headerID int64, tau string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertFlipTauQuery, tau)
}

func (repository *FlipStorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	return repository.insertRecordWithAddress(diffID, headerID, InsertFlipKicksQuery, kicks)
}

func (repository *FlipStorageRepository) insertBidBid(diffID, headerID int64, metadata vdbStorage.ValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidBidQuery, bidId, bid)
}

func (repository *FlipStorageRepository) insertBidLot(diffID, headerID int64, metadata vdbStorage.ValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidLotQuery, bidId, lot)
}

func (repository *FlipStorageRepository) insertBidGuy(diffID, headerID int64, metadata vdbStorage.ValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidGuyQuery, bidId, guy)
}

func (repository *FlipStorageRepository) insertBidTic(diffID, headerID int64, metadata vdbStorage.ValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidTicQuery, bidId, tic)
}

func (repository *FlipStorageRepository) insertBidEnd(diffID, headerID int64, metadata vdbStorage.ValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidEndQuery, bidId, end)
}

func (repository *FlipStorageRepository) insertBidUsr(diffID, headerID int64, metadata vdbStorage.ValueMetadata, usr string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidUsrQuery, bidId, usr)
}

func (repository *FlipStorageRepository) insertBidGal(diffID, headerID int64, metadata vdbStorage.ValueMetadata, gal string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidGalQuery, bidId, gal)
}

func (repository *FlipStorageRepository) insertBidTab(diffID, headerID int64, metadata vdbStorage.ValueMetadata, tab string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(diffID, headerID, InsertFlipBidTabQuery, bidId, tab)
}

func (repository *FlipStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata vdbStorage.ValueMetadata, packedValues map[int]string) error {
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

func (repository *FlipStorageRepository) insertRecordWithAddress(diffID, headerID int64, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressId, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlipStorageRepository) insertRecordWithAddressAndBidId(diffID, headerID int64, query, bidId, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressId, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, bidId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("flip field with address for bid id %s", bidId)
			return shared.FormatRollbackError(errorString, insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}

func getBidId(keys map[vdbStorage.Key]string) (string, error) {
	bidId, ok := keys[constants.BidId]
	if !ok {
		return "", vdbStorage.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bidId, nil
}
