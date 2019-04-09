package test_helpers

import (
	"database/sql"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
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

	FakeIlkRateMetadata = GetMetadata(vat.IlkRate, FakeIlk, utils.Uint256)
	FakeIlkArtMetadata  = GetMetadata(vat.IlkArt, FakeIlk, utils.Uint256)
	FakeIlkSpotMetadata = GetMetadata(vat.IlkSpot, FakeIlk, utils.Uint256)
	FakeIlkLineMetadata = GetMetadata(vat.IlkLine, FakeIlk, utils.Uint256)
	FakeIlkDustMetadata = GetMetadata(vat.IlkDust, FakeIlk, utils.Uint256)
	fakeIlkChopMetadata = GetMetadata(cat.IlkChop, FakeIlk, utils.Uint256)
	fakeIlkLumpMetadata = GetMetadata(cat.IlkLump, FakeIlk, utils.Uint256)
	fakeIlkFlipMetadata = GetMetadata(cat.IlkFlip, FakeIlk, utils.Address)
	fakeIlkRhoMetadata  = GetMetadata(jug.IlkRho, FakeIlk, utils.Uint256)
	fakeIlkTaxMetadata  = GetMetadata(jug.IlkTax, FakeIlk, utils.Uint256)
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

	anotherFakeIlkRateMetadata = GetMetadata(vat.IlkRate, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkArtMetadata  = GetMetadata(vat.IlkArt, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkSpotMetadata = GetMetadata(vat.IlkSpot, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkLineMetadata = GetMetadata(vat.IlkLine, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkDustMetadata = GetMetadata(vat.IlkDust, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkChopMetadata = GetMetadata(cat.IlkChop, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkLumpMetadata = GetMetadata(cat.IlkLump, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkFlipMetadata = GetMetadata(cat.IlkFlip, AnotherFakeIlk, utils.Address)
	anotherFakeIlkRhoMetadata  = GetMetadata(jug.IlkRho, AnotherFakeIlk, utils.Uint256)
	anotherFakeIlkTaxMetadata  = GetMetadata(jug.IlkTax, AnotherFakeIlk, utils.Uint256)

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

type IlkState struct {
	Ilk     string
	Rate    string
	Art     string
	Spot    string
	Line    string
	Dust    string
	Chop    string
	Lump    string
	Flip    string
	Rho     string
	Tax     string
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
	valuesMap[jug.IlkTax] = strconv.Itoa(9 + seed)

	return valuesMap
}

func IlkStateFromValues(ilk, updated, created string, ilkValues map[string]string) IlkState {
	return IlkState{
		Ilk:     ilk,
		Rate:    ilkValues[vat.IlkRate],
		Art:     ilkValues[vat.IlkArt],
		Spot:    ilkValues[vat.IlkSpot],
		Line:    ilkValues[vat.IlkLine],
		Dust:    ilkValues[vat.IlkDust],
		Chop:    ilkValues[cat.IlkChop],
		Lump:    ilkValues[cat.IlkLump],
		Flip:    ilkValues[cat.IlkFlip],
		Rho:     ilkValues[jug.IlkRho],
		Tax:     ilkValues[jug.IlkTax],
		Created: sql.NullString{String: created, Valid: true},
		Updated: sql.NullString{String: updated, Valid: true},
	}
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
	UrnId       string
	IlkId       string
	BlockHeight int
	Ink         string
	Art         string
	Ratio       sql.NullString
	Safe        bool
	Created     sql.NullString
	Updated     sql.NullString
}

func AssertUrn(actual, expected UrnState) {
	Expect(actual.UrnId).To(Equal(expected.UrnId))
	Expect(actual.IlkId).To(Equal(expected.IlkId))
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

type FrobEvent struct {
	IlkId string
	UrnId string
	Dink  string
	Dart  string
}
