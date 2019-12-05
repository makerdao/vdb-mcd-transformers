package flip

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertFlipVatQuery   = `INSERT INTO maker.flip_vat (header_id, address_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlipIlkQuery   = `INSERT INTO maker.flip_ilk (header_id, address_id, ilk_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlipBegQuery   = `INSERT INTO maker.flip_beg (header_id, address_id, beg) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlipTtlQuery   = `INSERT INTO maker.flip_ttl (header_id, address_id, ttl) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlipTauQuery   = `INSERT INTO maker.flip_tau (header_id, address_id, tau) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertFlipKicksQuery = `INSERT INTO maker.flip_kicks (header_id, address_id, kicks) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	InsertFlipBidBidQuery = `INSERT INTO maker.flip_bid_bid (header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidLotQuery = `INSERT INTO maker.flip_bid_lot (header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidGuyQuery = `INSERT INTO maker.flip_bid_guy (header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidTicQuery = `INSERT INTO maker.flip_bid_tic (header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidEndQuery = `INSERT INTO maker.flip_bid_end (header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidUsrQuery = `INSERT INTO maker.flip_bid_usr (header_id, address_id, bid_id, usr) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidGalQuery = `INSERT INTO maker.flip_bid_gal (header_id, address_id, bid_id, gal) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertFlipBidTabQuery = `INSERT INTO maker.flip_bid_tab (header_id, address_id, bid_id, tab) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type FlipStorageRepository struct {
	ContractAddress string
	db              *postgres.DB
}

func (repository *FlipStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case storage.Vat:
		return repository.insertVat(headerID, value.(string))
	case storage.Ilk:
		return repository.insertIlk(headerID, value.(string))
	case storage.Beg:
		return repository.insertBeg(headerID, value.(string))
	case storage.Kicks:
		return repository.insertKicks(headerID, value.(string))
	case storage.BidBid:
		return repository.insertBidBid(headerID, metadata, value.(string))
	case storage.BidLot:
		return repository.insertBidLot(headerID, metadata, value.(string))
	case storage.BidUsr:
		return repository.insertBidUsr(headerID, metadata, value.(string))
	case storage.BidGal:
		return repository.insertBidGal(headerID, metadata, value.(string))
	case storage.BidTab:
		return repository.insertBidTab(headerID, metadata, value.(string))
	case storage.Packed:
		return repository.insertPackedValueRecord(headerID, metadata, value.(map[int]string))
	default:
		panic(fmt.Sprintf("unrecognized flip contract storage name: %s", metadata.Name))
	}
}

func (repository *FlipStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlipStorageRepository) insertVat(headerID int64, vat string) error {
	return repository.insertRecordWithAddress(headerID, insertFlipVatQuery, vat)
}

func (repository *FlipStorageRepository) insertIlk(headerID int64, ilk string) error {
	ilkID, ilkErr := shared.GetOrCreateIlk(ilk, repository.db)
	if ilkErr != nil {
		return ilkErr
	}

	return repository.insertRecordWithAddress(headerID, insertFlipIlkQuery, strconv.FormatInt(ilkID, 10))
}

func (repository *FlipStorageRepository) insertBeg(headerID int64, beg string) error {
	return repository.insertRecordWithAddress(headerID, insertFlipBegQuery, beg)
}

func (repository *FlipStorageRepository) insertTtl(headerID int64, ttl string) error {
	return repository.insertRecordWithAddress(headerID, insertFlipTtlQuery, ttl)
}

func (repository *FlipStorageRepository) insertTau(headerID int64, tau string) error {
	return repository.insertRecordWithAddress(headerID, insertFlipTauQuery, tau)
}

func (repository *FlipStorageRepository) insertKicks(headerID int64, kicks string) error {
	return repository.insertRecordWithAddress(headerID, InsertFlipKicksQuery, kicks)
}

func (repository *FlipStorageRepository) insertBidBid(headerID int64, metadata utils.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidBidQuery, bidId, bid)
}

func (repository *FlipStorageRepository) insertBidLot(headerID int64, metadata utils.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidLotQuery, bidId, lot)
}

func (repository *FlipStorageRepository) insertBidGuy(headerID int64, metadata utils.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidGuyQuery, bidId, guy)
}

func (repository *FlipStorageRepository) insertBidTic(headerID int64, metadata utils.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidTicQuery, bidId, tic)
}

func (repository *FlipStorageRepository) insertBidEnd(headerID int64, metadata utils.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidEndQuery, bidId, end)
}

func (repository *FlipStorageRepository) insertBidUsr(headerID int64, metadata utils.StorageValueMetadata, usr string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidUsrQuery, bidId, usr)
}

func (repository *FlipStorageRepository) insertBidGal(headerID int64, metadata utils.StorageValueMetadata, gal string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidGalQuery, bidId, gal)
}

func (repository *FlipStorageRepository) insertBidTab(headerID int64, metadata utils.StorageValueMetadata, tab string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertRecordWithAddressAndBidId(headerID, InsertFlipBidTabQuery, bidId, tab)
}

func (repository *FlipStorageRepository) insertPackedValueRecord(headerID int64, metadata utils.StorageValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case storage.Ttl:
			insertErr = repository.insertTtl(headerID, value)
		case storage.Tau:
			insertErr = repository.insertTau(headerID, value)
		case storage.BidGuy:
			insertErr = repository.insertBidGuy(headerID, metadata, value)
		case storage.BidTic:
			insertErr = repository.insertBidTic(headerID, metadata, value)
		case storage.BidEnd:
			insertErr = repository.insertBidEnd(headerID, metadata, value)
		default:
			panic(fmt.Sprintf("unrecognized flip contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func (repository *FlipStorageRepository) insertRecordWithAddress(headerID int64, query, value string) error {
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
	_, insertErr := tx.Exec(query, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flip field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlipStorageRepository) insertRecordWithAddressAndBidId(headerID int64, query, bidId, value string) error {
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
	_, insertErr := tx.Exec(query, headerID, addressId, bidId, value)
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

func getBidId(keys map[utils.Key]string) (string, error) {
	bidId, ok := keys[constants.BidId]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bidId, nil
}
