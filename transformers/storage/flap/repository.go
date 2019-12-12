package flap

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery    = `INSERT INTO maker.flap_vat (diff_id, header_id, address_id, vat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertGemQuery    = `INSERT INTO maker.flap_gem (diff_id, header_id, address_id, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBegQuery    = `INSERT INTO maker.flap_beg (diff_id, header_id, address_id, beg) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTtlQuery    = `INSERT INTO maker.flap_ttl (diff_id, header_id, address_id, ttl) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertTauQuery    = `INSERT INTO maker.flap_tau (diff_id, header_id, address_id, tau) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertKicksQuery  = `INSERT INTO maker.flap_kicks (diff_id, header_id, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertLiveQuery   = `INSERT INTO maker.flap_live (diff_id, header_id, address_id, live) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertBidBidQuery = `INSERT INTO maker.flap_bid_bid (diff_id, header_id, address_id, bid_id, bid) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidLotQuery = `INSERT INTO maker.flap_bid_lot (diff_id, header_id, address_id, bid_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidGuyQuery = `INSERT INTO maker.flap_bid_guy (diff_id, header_id, address_id, bid_id, guy) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidTicQuery = `INSERT INTO maker.flap_bid_tic (diff_id, header_id, address_id, bid_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertBidEndQuery = `INSERT INTO maker.flap_bid_end (diff_id, header_id, address_id, bid_id, "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type FlapStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *FlapStorageRepository) Create(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, value interface{}) error {
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

func (repository *FlapStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *FlapStorageRepository) insertVat(diffID, headerID int64, vat string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertVatQuery, vat)
}

func (repository *FlapStorageRepository) insertGem(diffID, headerID int64, gem string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertGemQuery, gem)
}

func (repository *FlapStorageRepository) insertBeg(diffID, headerID int64, beg string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertBegQuery, beg)
}

func (repository *FlapStorageRepository) insertTtl(diffID, headerID int64, ttl string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertTtlQuery, ttl)
}

func (repository *FlapStorageRepository) insertTau(diffID, headerID int64, tau string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertTauQuery, tau)
}

func (repository *FlapStorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	return repository.insertRecordWithAddress(diffID, headerID, InsertKicksQuery, kicks)
}

func (repository *FlapStorageRepository) insertLive(diffID, headerID int64, live string) error {
	return repository.insertRecordWithAddress(diffID, headerID, insertLiveQuery, live)
}

func (repository *FlapStorageRepository) insertBidBid(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, bid string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(diffID, headerID, insertBidBidQuery, bidId, bid)
}

func (repository *FlapStorageRepository) insertBidLot(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, lot string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(diffID, headerID, insertBidLotQuery, bidId, lot)
}

func (repository *FlapStorageRepository) insertBidGuy(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, guy string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(diffID, headerID, insertBidGuyQuery, bidId, guy)
}

func (repository *FlapStorageRepository) insertBidTic(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, tic string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(diffID, headerID, insertBidTicQuery, bidId, tic)
}

func (repository *FlapStorageRepository) insertBidEnd(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, end string) error {
	bidId, err := getBidId(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertRecordWithAddressAndBidId(diffID, headerID, insertBidEndQuery, bidId, end)
}

func (repository *FlapStorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata vdbStorage.StorageValueMetadata, packedValues map[int]string) error {
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
			panic(fmt.Sprintf("unrecognized flap contract storage name in packed values: %s", metadata.Name))
		}
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func getBidId(keys map[vdbStorage.Key]string) (string, error) {
	bid, ok := keys[constants.BidId]
	if !ok {
		return "", vdbStorage.ErrMetadataMalformed{MissingData: constants.BidId}
	}
	return bid, nil
}

func (repository *FlapStorageRepository) insertRecordWithAddress(diffID, headerID int64, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressId, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap field with address", insertErr.Error())
		}
		return insertErr
	}

	return tx.Commit()
}

func (repository *FlapStorageRepository) insertRecordWithAddressAndBidId(diffID, headerID int64, query, bidId, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressId, addressErr := shared.GetOrCreateAddressInTransaction(repository.ContractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("flap address", addressErr.Error())
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, bidId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("flap field with address for bid id %s", bidId)
			return shared.FormatRollbackError(errorString, insertErr.Error())
		}
		return insertErr
	}
	return tx.Commit()
}
