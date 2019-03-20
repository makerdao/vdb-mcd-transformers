package test_helpers

import (
	"database/sql"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/drip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var (
	FakeIlk        = "fakeIlk"
	AnotherFakeIlk = "anotherFakeIlk"

	EmptyMetadatas []utils.StorageValueMetadata

	FakeIlkTakeMetadata = GetMetadata(vat.IlkTake, FakeIlk, utils.Uint256)
	FakeIlkRateMetadata = GetMetadata(vat.IlkRate, FakeIlk, utils.Uint256)
	fakeIlkInkMetadata  = GetMetadata(vat.IlkInk, FakeIlk, utils.Uint256)
	fakeIlkArtMetadata  = GetMetadata(vat.IlkArt, FakeIlk, utils.Uint256)
	fakeIlkSpotMetadata = GetMetadata(pit.IlkSpot, FakeIlk, utils.Uint256)
	fakeIlkLineMetadata = GetMetadata(pit.IlkLine, FakeIlk, utils.Uint256)
	fakeIlkChopMetadata = GetMetadata(cat.IlkChop, FakeIlk, utils.Uint256)
	fakeIlkLumpMetadata = GetMetadata(cat.IlkLump, FakeIlk, utils.Uint256)
	fakeIlkFlipMetadata = GetMetadata(cat.IlkFlip, FakeIlk, utils.Address)
	fakeIlkRhoMetadata  = GetMetadata(drip.IlkRho, FakeIlk, utils.Uint256)
	fakeIlkTaxMetadata  = GetMetadata(drip.IlkTax, FakeIlk, utils.Uint256)
	FakeIlkVatMetadatas = []utils.StorageValueMetadata{
		FakeIlkTakeMetadata,
		FakeIlkRateMetadata,
		fakeIlkInkMetadata,
		fakeIlkArtMetadata,
	}
	FakeIlkPitMetadatas = []utils.StorageValueMetadata{
		fakeIlkSpotMetadata,
		fakeIlkLineMetadata,
	}
	FakeIlkCatMetadatas = []utils.StorageValueMetadata{
		fakeIlkChopMetadata,
		fakeIlkLumpMetadata,
		fakeIlkFlipMetadata,
	}
	FakeIlkDripMetadatas = []utils.StorageValueMetadata{
		fakeIlkRhoMetadata,
		fakeIlkTaxMetadata,
	}

	anotherFakeIlkTakeMetadata = GetMetadata(vat.IlkTake, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkRateMetadata = GetMetadata(vat.IlkRate, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkInkMetadata  = GetMetadata(vat.IlkInk, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkArtMetadata  = GetMetadata(vat.IlkArt, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkSpotMetadata = GetMetadata(pit.IlkSpot, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkLineMetadata = GetMetadata(pit.IlkLine, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkChopMetadata = GetMetadata(cat.IlkChop, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkLumpMetadata = GetMetadata(cat.IlkLump, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkFlipMetadata = GetMetadata(cat.IlkFlip, AnotherFakeIlk, utils.Address)
	anotherFakeIlkRhoMetadata  = GetMetadata(drip.IlkRho, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkTaxMetadata  = GetMetadata(drip.IlkTax, AnotherFakeIlk, utils.Uint256)

	AnotherFakeIlkVatMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkTakeMetadata,
		anotherFakeIlkRateMetadata,
		anotherFakeIlkInkMetadata,
		anotherFakeIlkArtMetadata,
	}
	AnotherFakeIlkPitMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkSpotMetadata,
		anotherFakeIlkLineMetadata,
	}
	AnotherFakeIlkCatMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkChopMetadata,
		anotherFakeIlkLumpMetadata,
		anotherFakeIlkFlipMetadata,
	}
	AnotherFakeIlkDripMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkRhoMetadata,
		anotherFakeIlkTaxMetadata,
	}

	FakeIlkMetadatas = fakeIlkMetadatas()
)

func fakeIlkMetadatas() map[string][]utils.StorageValueMetadata {
	m := make(map[string][]utils.StorageValueMetadata)
	m["vat"] = FakeIlkVatMetadatas
	m["pit"] = FakeIlkPitMetadatas
	m["cat"] = FakeIlkCatMetadatas
	m["drip"] = FakeIlkDripMetadatas
	return m
}

type IlkState struct {
	Ilk     string
	Take    string
	Rate    string
	Ink     string
	Art     string
	Spot    string
	Line    string
	Chop    string
	Lump    string
	Flip    string
	Rho     string
	Tax     string
	Created sql.NullString
	Updated sql.NullString
}

func GetIlkState(seed int) map[string]string {
	valuesMap := make(map[string]string)
	valuesMap[vat.IlkTake] = strconv.Itoa(1 + seed)
	valuesMap[vat.IlkRate] = strconv.Itoa(2 + seed)
	valuesMap[vat.IlkInk] = strconv.Itoa(3 + seed)
	valuesMap[vat.IlkArt] = strconv.Itoa(4 + seed)
	valuesMap[pit.IlkSpot] = strconv.Itoa(5 + seed)
	valuesMap[pit.IlkLine] = strconv.Itoa(6 + seed)
	valuesMap[cat.IlkChop] = strconv.Itoa(7 + seed)
	valuesMap[cat.IlkLump] = strconv.Itoa(8 + seed)
	valuesMap[cat.IlkFlip] = "an address" + strconv.Itoa(seed)
	valuesMap[drip.IlkRho] = strconv.Itoa(9 + seed)
	valuesMap[drip.IlkTax] = strconv.Itoa(10 + seed)

	return valuesMap
}

func GetMetadata(fieldType, ilk string, valueType utils.ValueType) utils.StorageValueMetadata {
	return utils.GetStorageValueMetadata(fieldType, map[utils.Key]string{constants.Ilk: ilk}, valueType)
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

func CreatePitRecords(header core.Header, valuesMap map[string]string, metadatas []utils.StorageValueMetadata, repository pit.PitStorageRepository) {
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

func CreateDripRecords(header core.Header, valuesMap map[string]string, metadatas []utils.StorageValueMetadata, repository drip.DripStorageRepository) {
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
func CreateUrn(setupData UrnSetupData, metadata UrnMetadata, vatRepo vat.VatStorageRepository,
	pitRepo pit.PitStorageRepository, headerRepo repositories.HeaderRepository) {

	blockNo := int(setupData.Header.BlockNumber)
	hash := setupData.Header.Hash

	// This also creates the ilk if it doesn't exist
	err := vatRepo.Create(blockNo, hash, metadata.UrnInk, strconv.Itoa(setupData.Ink))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.UrnArt, strconv.Itoa(setupData.Art))
	Expect(err).NotTo(HaveOccurred())

	err = pitRepo.Create(blockNo, hash, metadata.IlkSpot, strconv.Itoa(setupData.Spot))
	Expect(err).NotTo(HaveOccurred())

	err = vatRepo.Create(blockNo, hash, metadata.IlkRate, strconv.Itoa(setupData.Rate))
	Expect(err).NotTo(HaveOccurred())

	_, err = headerRepo.CreateOrUpdateHeader(setupData.Header)
	Expect(err).NotTo(HaveOccurred())
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
		IlkSpot: utils.GetStorageValueMetadata(pit.IlkSpot,
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
	UrnId   string
	IlkId   string
	Ink     string
	Art     string
	Ratio   sql.NullString
	Safe    bool
	Created sql.NullString
	Updated sql.NullString
	// Frobs and bites collections, and ilk object, are missing
}
