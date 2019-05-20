package test_helpers

import (
	"database/sql"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
	"time"
)

var (
	FakeIlk = TestIlk{
		Hex:  "464b450000000000000000000000000000000000000000000000000000000000",
		Name: "FKE",
	}

	AnotherFakeIlk = TestIlk{
		Hex:  "464b453200000000000000000000000000000000000000000000000000000000",
		Name: "FKE2",
	}

	EmptyMetadatas []utils.StorageValueMetadata

	FakeIlkRateMetadata = utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	FakeIlkArtMetadata  = utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	FakeIlkSpotMetadata = utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	FakeIlkLineMetadata = utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	FakeIlkDustMetadata = utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	fakeIlkChopMetadata = utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	fakeIlkLumpMetadata = utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	fakeIlkFlipMetadata = utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	fakeIlkRhoMetadata  = utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)
	fakeIlkTaxMetadata  = utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)

	FakeIlkVatMetadatas = []utils.StorageValueMetadata{
		FakeIlkRateMetadata,
		FakeIlkArtMetadata,
		FakeIlkSpotMetadata,
		FakeIlkLineMetadata,
		FakeIlkDustMetadata,
	}
	FakeIlkCatMetadatas = []utils.StorageValueMetadata{
		fakeIlkChopMetadata,
		fakeIlkLumpMetadata,
		fakeIlkFlipMetadata,
	}
	FakeIlkJugMetadatas = []utils.StorageValueMetadata{
		fakeIlkRhoMetadata,
		fakeIlkTaxMetadata,
	}

	anotherFakeIlkRateMetadata = utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkArtMetadata  = utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkSpotMetadata = utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkLineMetadata = utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkDustMetadata = utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkChopMetadata = utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkLumpMetadata = utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkFlipMetadata = utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Address)
	anotherFakeIlkRhoMetadata  = utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)
	anotherFakeIlkTaxMetadata  = utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)

	AnotherFakeIlkVatMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkRateMetadata,
		anotherFakeIlkArtMetadata,
		anotherFakeIlkSpotMetadata,
		anotherFakeIlkLineMetadata,
		anotherFakeIlkDustMetadata,
	}
	AnotherFakeIlkCatMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkChopMetadata,
		anotherFakeIlkLumpMetadata,
		anotherFakeIlkFlipMetadata,
	}
	AnotherFakeIlkJugMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkRhoMetadata,
		anotherFakeIlkTaxMetadata,
	}
)

type TestIlk struct {
	Hex  string
	Name string
}

type IlkState struct {
	IlkName string `db:"ilk_name"`
	Rate    string
	Art     string
	Spot    string
	Line    string
	Dust    string
	Chop    string
	Lump    string
	Flip    string
	Rho     string
	Duty    string
	Created sql.NullString
	Updated sql.NullString
}

func GetIlkValues(seed int) map[string]string {
	valuesMap := make(map[string]string)
	valuesMap[vat.IlkRate] = strconv.Itoa(1 + seed)
	valuesMap[vat.IlkArt] = strconv.Itoa(2 + seed)
	valuesMap[vat.IlkSpot] = strconv.Itoa(3 + seed)
	valuesMap[vat.IlkLine] = strconv.Itoa(4 + seed)
	valuesMap[vat.IlkDust] = strconv.Itoa(5 + seed)
	valuesMap[cat.IlkChop] = strconv.Itoa(6 + seed)
	valuesMap[cat.IlkLump] = strconv.Itoa(7 + seed)
	valuesMap[cat.IlkFlip] = "an address" + strconv.Itoa(seed)
	valuesMap[jug.IlkRho] = strconv.Itoa(8 + seed)
	valuesMap[jug.IlkDuty] = strconv.Itoa(9 + seed)

	return valuesMap
}

func IlkStateFromValues(ilk, updated, created string, ilkValues map[string]string) IlkState {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	parsedUpdated, _ := strconv.ParseInt(updated, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)
	updatedTimestamp := time.Unix(parsedUpdated, 0).UTC().Format(time.RFC3339)

	ilkName, err := shared.DecodeIlkName(ilk)
	Expect(err).NotTo(HaveOccurred())
	return IlkState{
		IlkName: ilkName,
		Rate:    ilkValues[vat.IlkRate],
		Art:     ilkValues[vat.IlkArt],
		Spot:    ilkValues[vat.IlkSpot],
		Line:    ilkValues[vat.IlkLine],
		Dust:    ilkValues[vat.IlkDust],
		Chop:    ilkValues[cat.IlkChop],
		Lump:    ilkValues[cat.IlkLump],
		Flip:    ilkValues[cat.IlkFlip],
		Rho:     ilkValues[jug.IlkRho],
		Duty:    ilkValues[jug.IlkDuty],
		Created: sql.NullString{String: createdTimestamp, Valid: true},
		Updated: sql.NullString{String: updatedTimestamp, Valid: true},
	}
}

func CreateVatRecords(header core.Header, valuesMap map[string]string, metadatas []utils.StorageValueMetadata, repository vat.VatStorageRepository) {
	blockHash := header.Hash
	blockNumber := int(header.BlockNumber)

	for _, metadata := range metadatas {
		value := valuesMap[metadata.Name]
		err := repository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}
}

func CreateCatRecords(header core.Header, valuesMap map[string]string, metadatas []utils.StorageValueMetadata, repository cat.CatStorageRepository) {
	blockHash := header.Hash
	blockNumber := int(header.BlockNumber)

	for _, metadata := range metadatas {
		value := valuesMap[metadata.Name]
		err := repository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}
}

func CreateJugRecords(header core.Header, valuesMap map[string]string, metadatas []utils.StorageValueMetadata, repository jug.JugStorageRepository) {
	blockHash := header.Hash
	blockNumber := int(header.BlockNumber)

	for _, metadata := range metadatas {
		value := valuesMap[metadata.Name]
		err := repository.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}
}

func GetExpectedRatio(ink, spot, art, rate int) float64 {
	inkXspot := float64(ink) * float64(spot)
	artXrate := float64(art) * float64(rate)
	return inkXspot / artXrate
}

// Creates urn by creating necessary state diffs and the corresponding header
func CreateUrn(setupData UrnSetupData, metadata UrnMetadata, vatRepo vat.VatStorageRepository, headerRepo repositories.HeaderRepository) {
	blockNo := int(setupData.Header.BlockNumber)
	hash := setupData.Header.Hash

	// This also creates the ilk if it doesn't exist
	err := vatRepo.Create(blockNo, hash, metadata.UrnInk, strconv.Itoa(setupData.Ink))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.UrnArt, strconv.Itoa(setupData.Art))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.IlkSpot, strconv.Itoa(setupData.Spot))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.IlkRate, strconv.Itoa(setupData.Rate))
	Expect(err).NotTo(HaveOccurred())

	_, err = headerRepo.CreateOrUpdateHeader(setupData.Header)
	if err == repositories.ErrValidHeaderExists {
		// In some tests, the header has been created in other operations
		err = nil
	}
	Expect(err).NotTo(HaveOccurred())
}

func CreateIlk(db *postgres.DB, header core.Header, valuesMap map[string]string, vatMetadatas, catMetadatas, jugMetadatas []utils.StorageValueMetadata) {
	var (
		vatRepo vat.VatStorageRepository
		catRepo cat.CatStorageRepository
		jugRepo jug.JugStorageRepository
	)
	vatRepo.SetDB(db)
	catRepo.SetDB(db)
	jugRepo.SetDB(db)
	CreateVatRecords(header, valuesMap, vatMetadatas, vatRepo)
	CreateCatRecords(header, valuesMap, catMetadatas, catRepo)
	CreateJugRecords(header, valuesMap, jugMetadatas, jugRepo)
}

// Does not return values computed by the query (ratio, safe, updated, created)
func GetUrnSetupData(block, timestamp int) UrnSetupData {
	fakeHeader := fakes.GetFakeHeader(int64(block))
	fakeHeader.Timestamp = strconv.Itoa(timestamp)
	fakeHeader.Hash = test_data.RandomString(5)

	return UrnSetupData{
		Header: fakeHeader,
		Ink:    rand.Int(),
		Art:    rand.Int(),
		Spot:   rand.Int(),
		Rate:   rand.Int(),
	}
}

type UrnSetupData struct {
	Header core.Header
	Ink    int
	Art    int
	Spot   int
	Rate   int
}

func GetUrnMetadata(ilk, urn string) UrnMetadata {
	return UrnMetadata{
		UrnInk: utils.GetStorageValueMetadata(vat.UrnInk,
			map[utils.Key]string{constants.Ilk: ilk, constants.Guy: urn}, utils.Uint256),
		UrnArt: utils.GetStorageValueMetadata(vat.UrnArt,
			map[utils.Key]string{constants.Ilk: ilk, constants.Guy: urn}, utils.Uint256),
		IlkSpot: utils.GetStorageValueMetadata(vat.IlkSpot,
			map[utils.Key]string{constants.Ilk: ilk}, utils.Uint256),
		IlkRate: utils.GetStorageValueMetadata(vat.IlkRate,
			map[utils.Key]string{constants.Ilk: ilk}, utils.Uint256),
	}
}

type UrnMetadata struct {
	UrnInk  utils.StorageValueMetadata
	UrnArt  utils.StorageValueMetadata
	IlkSpot utils.StorageValueMetadata
	IlkRate utils.StorageValueMetadata
}

type UrnState struct {
	UrnGuy      string `db:"urn_guy"`
	IlkName     string `db:"ilk_name"`
	BlockHeight int    `db:"block_height"`
	Ink         string
	Art         string
	Ratio       sql.NullString
	Safe        bool
	Created     sql.NullString
	Updated     sql.NullString
}

func AssertUrn(actual, expected UrnState) {
	Expect(actual.UrnGuy).To(Equal(expected.UrnGuy))
	Expect(actual.IlkName).To(Equal(expected.IlkName))
	Expect(actual.BlockHeight).To(Equal(expected.BlockHeight))
	Expect(actual.Ink).To(Equal(expected.Ink))
	Expect(actual.Art).To(Equal(expected.Art))

	if actual.Ratio.Valid {
		actualRatio, err := strconv.ParseFloat(actual.Ratio.String, 64)
		Expect(err).NotTo(HaveOccurred())
		expectedRatio, err := strconv.ParseFloat(expected.Ratio.String, 64)
		Expect(err).NotTo(HaveOccurred())
		Expect(actualRatio).To(BeNumerically("~", expectedRatio))
	} else {
		Expect(!expected.Ratio.Valid)
	}

	Expect(actual.Safe).To(Equal(expected.Safe))
	Expect(actual.Created).To(Equal(expected.Created))
	Expect(actual.Updated).To(Equal(expected.Updated))
}

type FileEvent struct {
	Id      string
	IlkName sql.NullString `db:"ilk_name"`
	What    string
	Data    string
}

type FrobEvent struct {
	IlkName string `db:"ilk_name"`
	UrnGuy  string `db:"urn_guy"`
	Dink    string
	Dart    string
}

type BiteEvent struct {
	IlkName string `db:"ilk_name"`
	UrnGuy  string `db:"urn_guy"`
	Ink     string
	Art     string
	Tab     string
}

type LogValue struct {
	Val             string
	BlockNumber     uint64 `db:"block_number"`
	TxIdx           uint   `db:"tx_idx"`
	ContractAddress string `db:"contract_address"`
}

func GetExpectedTimestamp(epoch int) string {
	return time.Unix(int64(epoch), 0).UTC().Format(time.RFC3339)
}
