package clip

import (
	"fmt"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	Dog        = "dog"
	Vow        = "vow"
	Spotter    = "spotter"
	Calc       = "calc"
	Buf        = "buf"
	Tail       = "tail"
	Cusp       = "cusp"
	Chip       = "chip"
	Tip        = "tip"
	Chost      = "chost"
	Kicks      = "kicks"
	Active     = "active"
	ActiveSale = "active_sale"

	Packed = "packed_storage_values"

	SalePos = "sale_pos"
	SaleTab = "sale_tab"
	SaleLot = "sale_lot"
	SaleUsr = "sale_usr"
	SaleTic = "sale_tic"
	SaleTop = "sale_top"

	insertClipDogQuery        = `INSERT INTO maker.clip_dog (diff_id, header_id, address_id, dog) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipVowQuery        = `INSERT INTO maker.clip_vow (diff_id, header_id, address_id, vow) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipSpotterQuery    = `INSERT INTO maker.clip_spotter (diff_id, header_id, address_id, spotter) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipCalcQuery       = `INSERT INTO maker.clip_calc (diff_id, header_id, address_id, calc) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipBufQuery        = `INSERT INTO maker.clip_buf (diff_id, header_id, address_id, buf) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipTailQuery       = `INSERT INTO maker.clip_tail (diff_id, header_id, address_id, tail) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipCuspQuery       = `INSERT INTO maker.clip_cusp (diff_id, header_id, address_id, cusp) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipChipQuery       = `INSERT INTO maker.clip_chip (diff_id, header_id, address_id, chip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipTipQuery        = `INSERT INTO maker.clip_tip (diff_id, header_id, address_id, tip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipChostQuery      = `INSERT INTO maker.clip_chost (diff_id, header_id, address_id, chost) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipKicksQuery      = `INSERT INTO maker.clip_kicks (diff_id, header_id, address_id, kicks) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipActiveQuery     = `INSERT INTO maker.clip_active (diff_id, header_id, address_id, active) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertClipActiveSaleQuery = `INSERT INTO maker.clip_active_sales (diff_id, header_id, address_id, sale_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	insertSalePosQuery = `INSERT INTO maker.clip_sale_pos (diff_id, header_id, address_id, sale_id, pos) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertSaleTabQuery = `INSERT INTO maker.clip_sale_tab (diff_id, header_id, address_id, sale_id, tab) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertSaleLotQuery = `INSERT INTO maker.clip_sale_lot (diff_id, header_id, address_id, sale_id, lot) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertSaleUsrQuery = `INSERT INTO maker.clip_sale_usr (diff_id, header_id, address_id, sale_id, usr) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertSaleTicQuery = `INSERT INTO maker.clip_sale_tic (diff_id, header_id, address_id, sale_id, tic) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	insertSaleTopQuery = `INSERT INTO maker.clip_sale_top (diff_id, header_id, address_id, sale_id, top) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	ContractAddress   string
	contractAddressID int64
	db                *postgres.DB
}

func (repo *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dog:
		return repo.insertDog(diffID, headerID, value.(string))
	case Vow:
		return repo.insertVow(diffID, headerID, value.(string))
	case Spotter:
		return repo.insertSpotter(diffID, headerID, value.(string))
	case Calc:
		return repo.insertCalc(diffID, headerID, value.(string))
	case Buf:
		return repo.insertBuf(diffID, headerID, value.(string))
	case Tail:
		return repo.insertTail(diffID, headerID, value.(string))
	case Cusp:
		return repo.insertCusp(diffID, headerID, value.(string))
	case Packed:
		return repo.insertPackedValueRecord(diffID, headerID, metadata, value.(map[int]string))
	case Chost:
		return repo.insertChost(diffID, headerID, value.(string))
	case Kicks:
		return repo.insertKicks(diffID, headerID, value.(string))
	case Active:
		return repo.insertActive(diffID, headerID, value.(string))
	case ActiveSale:
		return repo.insertActiveSale(diffID, headerID, value.(string))
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repo.ContractAddress, value.(string), repo.db)
	case SalePos:
		return repo.insertSalePos(diffID, headerID, metadata, value.(string))
	case SaleTab:
		return repo.insertSaleTab(diffID, headerID, metadata, value.(string))
	case SaleLot:
		return repo.insertSaleLot(diffID, headerID, metadata, value.(string))
	case SaleTop:
		return repo.insertSaleTop(diffID, headerID, metadata, value.(string))
	default:
		return fmt.Errorf("unrecognized clip contract storage name: %s", metadata.Name)
	}
}

func (repo *StorageRepository) SetDB(db *postgres.DB) {
	repo.db = db
}

func (repo *StorageRepository) insertDog(diffID, headerID int64, dog string) error {
	dogAddressID, addressErr := repository.GetOrCreateAddress(repo.db, dog)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip dog: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipDogQuery,
		strconv.FormatInt(dogAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s dog %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, dog, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertVow(diffID, headerID int64, vow string) error {
	vowAddressID, addressErr := repository.GetOrCreateAddress(repo.db, vow)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip vow: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipVowQuery,
		strconv.FormatInt(vowAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s vow %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, vow, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSpotter(diffID, headerID int64, spotter string) error {
	spotterAddressID, addressErr := repository.GetOrCreateAddress(repo.db, spotter)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip spotter: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipSpotterQuery,
		strconv.FormatInt(spotterAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s spotter %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, spotter, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertCalc(diffID, headerID int64, calc string) error {
	calcAddressID, addressErr := repository.GetOrCreateAddress(repo.db, calc)
	if addressErr != nil {
		return fmt.Errorf("error inserting clip calc: %w", addressErr)
	}
	insertErr := shared.InsertRecordWithAddress(
		diffID,
		headerID,
		insertClipCalcQuery,
		strconv.FormatInt(calcAddressID, 10),
		repo.ContractAddress,
		repo.db)
	if insertErr != nil {
		msgToFormat := "error inserting clip %s calc %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, repo.ContractAddress, calc, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertBuf(diffID, headerID int64, buf string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipBufQuery, diffID, headerID, addressID, buf)
	if err != nil {
		return fmt.Errorf("error inserting clip buf %s from diff ID %d: %w", buf, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertTail(diffID, headerID int64, tail string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipTailQuery, diffID, headerID, addressID, tail)
	if err != nil {
		return fmt.Errorf("error inserting clip tail %s from diff ID %d: %w", tail, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertCusp(diffID, headerID int64, cusp string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipCuspQuery, diffID, headerID, addressID, cusp)
	if err != nil {
		return fmt.Errorf("error inserting clip cusp %s from diff ID %d: %w", cusp, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertChip(diffID, headerID int64, chip string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipChipQuery, diffID, headerID, addressID, chip)
	if err != nil {
		return fmt.Errorf("error inserting clip chip %s from diff ID %d: %w", chip, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertTip(diffID, headerID int64, tip string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipTipQuery, diffID, headerID, addressID, tip)
	if err != nil {
		return fmt.Errorf("error inserting clip tip %s from diff ID %d: %w", tip, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertChost(diffID, headerID int64, chost string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipChostQuery, diffID, headerID, addressID, chost)
	if err != nil {
		return fmt.Errorf("error inserting clip chost %s from diff ID %d: %w", chost, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertKicks(diffID, headerID int64, kicks string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipKicksQuery, diffID, headerID, addressID, kicks)
	if err != nil {
		return fmt.Errorf("error inserting clip kicks %s from diff ID %d: %w", kicks, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertActive(diffID, headerID int64, active string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipActiveQuery, diffID, headerID, addressID, active)
	if err != nil {
		return fmt.Errorf("error inserting clip active %s from diff ID %d: %w", active, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertActiveSale(diffID, headerID int64, saleID string) error {
	addressID, addressErr := repo.ContractAddressID()
	if addressErr != nil {
		return fmt.Errorf("could not retrieve address id for %s, error: %w", repo.ContractAddress, addressErr)
	}

	_, err := repo.db.Exec(insertClipActiveSaleQuery, diffID, headerID, addressID, saleID)
	if err != nil {
		return fmt.Errorf("error inserting clip active sale %s from diff ID %d: %w", saleID, diffID, err)
	}
	return nil
}

func (repo *StorageRepository) insertPackedValueRecord(diffID, headerID int64, metadata types.ValueMetadata, packedValues map[int]string) error {
	for order, value := range packedValues {
		var insertErr error
		switch metadata.PackedNames[order] {
		case Chip:
			insertErr = repo.insertChip(diffID, headerID, value)
		case Tip:
			insertErr = repo.insertTip(diffID, headerID, value)
		case SaleUsr:
			insertErr = repo.insertSaleUsr(diffID, headerID, metadata, value)
		case SaleTic:
			insertErr = repo.insertSaleTic(diffID, headerID, metadata, value)
		default:
			return fmt.Errorf("unrecognized clip contract storage name in packed values: %s", metadata.Name)
		}
		if insertErr != nil {
			return fmt.Errorf("error inserting clip packed value from diff ID %d: %w", diffID, insertErr)
		}
	}
	return nil
}

func (repo *StorageRepository) insertSalePos(diffID, headerID int64, metadata types.ValueMetadata, pos string) error {
	saleID, err := getSaleID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting saleID for clip sale pos: %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertSalePosQuery, saleID, pos, repo.ContractAddress, repo.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting saleID %s pos %s from diff ID %d", saleID, pos, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSaleTab(diffID, headerID int64, metadata types.ValueMetadata, tab string) error {
	saleID, err := getSaleID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting saleID for clip sale tab : %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertSaleTabQuery, saleID, tab, repo.ContractAddress, repo.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting saleID %s tab %s from diff ID %d", saleID, tab, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSaleLot(diffID, headerID int64, metadata types.ValueMetadata, lot string) error {
	saleID, err := getSaleID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting saleID for clip sale lot : %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertSaleLotQuery, saleID, lot, repo.ContractAddress, repo.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting saleID %s lot %s from diff ID %d", saleID, lot, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSaleUsr(diffID, headerID int64, metadata types.ValueMetadata, usr string) error {
	saleID, err := getSaleID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting saleID for clip sale usr : %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertSaleUsrQuery, saleID, usr, repo.ContractAddress, repo.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting saleID %s usr %s from diff ID %d", saleID, usr, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSaleTic(diffID, headerID int64, metadata types.ValueMetadata, tic string) error {
	saleID, err := getSaleID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting saleID for clip sale tic : %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertSaleTicQuery, saleID, tic, repo.ContractAddress, repo.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting saleID %s tic %s from diff ID %d", saleID, tic, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) insertSaleTop(diffID, headerID int64, metadata types.ValueMetadata, top string) error {
	saleID, err := getSaleID(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting saleID for clip sale top : %w", err)
	}
	insertErr := shared.InsertRecordWithAddressAndBidID(diffID, headerID, insertSaleTopQuery, saleID, top, repo.ContractAddress, repo.db)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting saleID %s top %s from diff ID %d", saleID, top, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repo *StorageRepository) ContractAddressID() (int64, error) {
	if repo.contractAddressID == 0 {
		addressID, addressErr := repository.GetOrCreateAddress(repo.db, repo.ContractAddress)
		repo.contractAddressID = addressID
		return repo.contractAddressID, addressErr
	}
	return repo.contractAddressID, nil
}

func getSaleID(keys map[types.Key]string) (string, error) {
	sale, ok := keys[constants.SaleId]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.SaleId}
	}
	return sale, nil
}
