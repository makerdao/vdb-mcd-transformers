package test_helpers

import (
	"database/sql"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/events/tick"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	vdbStorage "github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
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
		Hex:        "0x464b450000000000000000000000000000000000000000000000000000000000",
		Identifier: "FKE",
	}

	AnotherFakeIlk = TestIlk{
		Hex:        "0x464b453200000000000000000000000000000000000000000000000000000000",
		Identifier: "FKE2",
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
	fakeIlkPipMetadata  = utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Address)
	fakeIlkMatMetadata  = utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: FakeIlk.Hex}, utils.Uint256)

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
	FakeIlkSpotMetadatas = []utils.StorageValueMetadata{
		fakeIlkPipMetadata,
		fakeIlkMatMetadata,
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
	anotherFakeIlkPipMetadata  = utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Address)
	anotherFakeIlkMatMetadata  = utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, utils.Uint256)

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
	AnotherFakeIlkSpotMetadatas = []utils.StorageValueMetadata{
		anotherFakeIlkPipMetadata,
		anotherFakeIlkMatMetadata,
	}
)

type TestIlk struct {
	Hex        string
	Identifier string
}

type IlkState struct {
	IlkIdentifier string `db:"ilk_identifier"`
	Rate          string
	Art           string
	Spot          string
	Line          string
	Dust          string
	Chop          string
	Lump          string
	Flip          string
	Rho           string
	Duty          string
	Pip           string
	Mat           string
	Created       sql.NullString
	Updated       sql.NullString
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
	valuesMap[spot.IlkPip] = "an address2" + strconv.Itoa(seed)
	valuesMap[spot.IlkMat] = strconv.Itoa(10 + seed)

	return valuesMap
}

func IlkStateFromValues(ilk, updated, created string, ilkValues map[string]string) IlkState {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	parsedUpdated, _ := strconv.ParseInt(updated, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)
	updatedTimestamp := time.Unix(parsedUpdated, 0).UTC().Format(time.RFC3339)

	ilkIdentifier := shared.DecodeHexToText(ilk)
	return IlkState{
		IlkIdentifier: ilkIdentifier,
		Rate:          ilkValues[vat.IlkRate],
		Art:           ilkValues[vat.IlkArt],
		Spot:          ilkValues[vat.IlkSpot],
		Line:          ilkValues[vat.IlkLine],
		Dust:          ilkValues[vat.IlkDust],
		Chop:          ilkValues[cat.IlkChop],
		Lump:          ilkValues[cat.IlkLump],
		Flip:          ilkValues[cat.IlkFlip],
		Rho:           ilkValues[jug.IlkRho],
		Duty:          ilkValues[jug.IlkDuty],
		Pip:           ilkValues[spot.IlkPip],
		Mat:           ilkValues[spot.IlkMat],
		Created:       sql.NullString{String: createdTimestamp, Valid: true},
		Updated:       sql.NullString{String: updatedTimestamp, Valid: true},
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

func CreateSpotRecords(header core.Header, valuesMap map[string]string, metadatas []utils.StorageValueMetadata, repository spot.SpotStorageRepository) {
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

func CreateIlk(db *postgres.DB, header core.Header, valuesMap map[string]string, vatMetadatas, catMetadatas, jugMetadatas, spotMetadatas []utils.StorageValueMetadata) {
	var (
		vatRepo  vat.VatStorageRepository
		catRepo  cat.CatStorageRepository
		jugRepo  jug.JugStorageRepository
		spotRepo spot.SpotStorageRepository
	)
	vatRepo.SetDB(db)
	catRepo.SetDB(db)
	jugRepo.SetDB(db)
	spotRepo.SetDB(db)
	CreateVatRecords(header, valuesMap, vatMetadatas, vatRepo)
	CreateCatRecords(header, valuesMap, catMetadatas, catRepo)
	CreateJugRecords(header, valuesMap, jugMetadatas, jugRepo)
	CreateSpotRecords(header, valuesMap, spotMetadatas, spotRepo)
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
	UrnIdentifier string `db:"urn_identifier"`
	IlkIdentifier string `db:"ilk_identifier"`
	BlockHeight   int    `db:"block_height"`
	Ink           string
	Art           string
	Ratio         sql.NullString
	Safe          bool
	Created       sql.NullString
	Updated       sql.NullString
}

func AssertUrn(actual, expected UrnState) {
	Expect(actual.UrnIdentifier).To(Equal(expected.UrnIdentifier))
	Expect(actual.IlkIdentifier).To(Equal(expected.IlkIdentifier))
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
		Expect(!expected.Ratio.Valid).To(BeTrue())
	}

	Expect(actual.Safe).To(Equal(expected.Safe))
	Expect(actual.Created).To(Equal(expected.Created))
	Expect(actual.Updated).To(Equal(expected.Updated))
}

func getCommonBidMetadatas(bidId string) []utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	packedNames := map[int]string{0: storage.BidGuy, 1: storage.BidTic, 2: storage.BidEnd}
	packedTypes := map[int]utils.ValueType{0: utils.Address, 1: utils.Uint48, 2: utils.Uint48}
	return []utils.StorageValueMetadata{
		utils.GetStorageValueMetadata(storage.Kicks, nil, utils.Uint256),
		utils.GetStorageValueMetadata(storage.BidBid, keys, utils.Uint256),
		utils.GetStorageValueMetadata(storage.BidLot, keys, utils.Uint256),
		utils.GetStorageValueMetadataForPackedSlot(storage.Packed, keys, utils.PackedSlot, packedNames, packedTypes),
	}
}

func GetFlopMetadatas(bidId string) []utils.StorageValueMetadata {
	return getCommonBidMetadatas(bidId)
}

func GetFlapMetadatas(bidId string) []utils.StorageValueMetadata {
	return getCommonBidMetadatas(bidId)
}

func GetCdpManagerMetadatas(cdpi string) []utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.Cdpi: cdpi}
	return []utils.StorageValueMetadata{
		utils.GetStorageValueMetadata(cdp_manager.CdpManagerCdpi, nil, utils.Uint256),
		utils.GetStorageValueMetadata(cdp_manager.CdpManagerUrns, keys, utils.Address),
		utils.GetStorageValueMetadata(cdp_manager.CdpManagerOwns, keys, utils.Address),
		utils.GetStorageValueMetadata(cdp_manager.CdpManagerIlks, keys, utils.Bytes32),
	}
}

func GetFlipMetadatas(bidId string) []utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return append(getCommonBidMetadatas(bidId),
		utils.GetStorageValueMetadata(storage.Ilk, nil, utils.Bytes32),
		utils.GetStorageValueMetadata(storage.BidUsr, keys, utils.Address),
		utils.GetStorageValueMetadata(storage.BidGal, keys, utils.Address),
		utils.GetStorageValueMetadata(storage.BidTab, keys, utils.Uint256))
}

func GetCdpManagerStorageValues(seed int, ilkHex string, urnGuy string, cdpi int) map[string]interface{} {
	valuesMap := make(map[string]interface{})
	valuesMap[cdp_manager.CdpManagerCdpi] = strconv.Itoa(cdpi)
	valuesMap[cdp_manager.CdpManagerUrns] = urnGuy
	valuesMap[cdp_manager.CdpManagerOwns] = "address1" + strconv.Itoa(seed)
	valuesMap[cdp_manager.CdpManagerIlks] = ilkHex
	return valuesMap
}

func getCommonBidStorageValues(seed, bidId int) map[string]interface{} {
	packedValues := map[int]string{0: "address1" + strconv.Itoa(seed), 1: strconv.Itoa(1 + seed), 2: strconv.Itoa(2 + seed)}
	valuesMap := make(map[string]interface{})
	valuesMap[storage.Kicks] = strconv.Itoa(bidId)
	valuesMap[storage.BidBid] = strconv.Itoa(3 + seed)
	valuesMap[storage.BidLot] = strconv.Itoa(4 + seed)
	valuesMap[storage.Packed] = packedValues

	return valuesMap
}

func GetFlopStorageValues(seed, bidId int) map[string]interface{} {
	return getCommonBidStorageValues(seed, bidId)
}

func GetFlapStorageValues(seed, bidId int) map[string]interface{} {
	return getCommonBidStorageValues(seed, bidId)
}

func GetFlipStorageValues(seed int, ilk string, bidId int) map[string]interface{} {
	valuesMap := getCommonBidStorageValues(seed, bidId)
	valuesMap[storage.Ilk] = ilk
	valuesMap[storage.BidGal] = "address2" + strconv.Itoa(seed)
	valuesMap[storage.BidUsr] = "address3" + strconv.Itoa(seed)
	valuesMap[storage.BidTab] = strconv.Itoa(5 + seed)
	return valuesMap
}

func insertValues(repo vdbStorage.Repository, header core.Header, valuesMap map[string]interface{}, metadatas []utils.StorageValueMetadata) {
	blockHash := header.Hash
	blockNumber := int(header.BlockNumber)

	for _, metadata := range metadatas {
		value := valuesMap[metadata.Name]
		err := repo.Create(blockNumber, blockHash, metadata, value)

		Expect(err).NotTo(HaveOccurred())
	}
}

func CreateFlop(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flopMetadatas []utils.StorageValueMetadata, contractAddress string) {
	flopRepo := flop.FlopStorageRepository{ContractAddress: contractAddress}
	flopRepo.SetDB(db)
	insertValues(&flopRepo, header, valuesMap, flopMetadatas)
}

func CreateFlap(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flapMetadatas []utils.StorageValueMetadata, contractAddress string) {
	flapRepo := flap.FlapStorageRepository{ContractAddress: contractAddress}
	flapRepo.SetDB(db)
	insertValues(&flapRepo, header, valuesMap, flapMetadatas)
}

func CreateFlip(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flipMetadatas []utils.StorageValueMetadata, contractAddress string) {
	flipRepo := flip.FlipStorageRepository{ContractAddress: contractAddress}
	flipRepo.SetDB(db)
	insertValues(&flipRepo, header, valuesMap, flipMetadatas)
}

func CreateManagedCdp(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []utils.StorageValueMetadata) error {
	cdpManagerRepo := cdp_manager.CdpManagerStorageRepository{}
	cdpManagerRepo.SetDB(db)
	_, err := shared.GetOrCreateUrn(valuesMap[cdp_manager.CdpManagerUrns].(string), valuesMap[cdp_manager.CdpManagerIlks].(string), db)
	if err != nil {
		return err
	}
	insertValues(&cdpManagerRepo, header, valuesMap, metadatas)
	return nil
}

func ManagedCdpFromValues(ilkIdentifier, created string, cdpValues map[string]interface{}) ManagedCdp {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)

	return ManagedCdp{
		Usr:           cdpValues[cdp_manager.CdpManagerOwns].(string),
		Id:            cdpValues[cdp_manager.CdpManagerCdpi].(string),
		UrnIdentifier: cdpValues[cdp_manager.CdpManagerUrns].(string),
		IlkIdentifier: ilkIdentifier,
		Created:       sql.NullString{String: createdTimestamp, Valid: true},
	}
}

func commonBidFromValues(bidId, dealt, updated, created string, bidValues map[string]interface{}) commonBid {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	parsedUpdated, _ := strconv.ParseInt(updated, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)
	updatedTimestamp := time.Unix(parsedUpdated, 0).UTC().Format(time.RFC3339)
	packedValues := bidValues[storage.Packed].(map[int]string)

	return commonBid{
		BidId:   bidId,
		Guy:     packedValues[0],
		Tic:     packedValues[1],
		End:     packedValues[2],
		Lot:     bidValues[storage.BidLot].(string),
		Bid:     bidValues[storage.BidBid].(string),
		Dealt:   dealt,
		Created: sql.NullString{String: createdTimestamp, Valid: true},
		Updated: sql.NullString{String: updatedTimestamp, Valid: true},
	}
}

func FlopBidFromValues(bidId, dealt, updated, created string, bidValues map[string]interface{}) FlopBid {
	return FlopBid{
		commonBid: commonBidFromValues(bidId, dealt, updated, created, bidValues),
	}
}

func FlapBidFromValues(bidId, dealt, updated, created string, bidValues map[string]interface{}) FlapBid {
	return FlapBid{
		commonBid: commonBidFromValues(bidId, dealt, updated, created, bidValues),
	}
}

func FlipBidFromValues(bidId, ilkId, urnId, dealt, updated, created string, bidValues map[string]interface{}) FlipBid {
	return FlipBid{
		commonBid: commonBidFromValues(bidId, dealt, updated, created, bidValues),
		IlkId:     ilkId,
		UrnId:     urnId,
		Gal:       bidValues[storage.BidGal].(string),
		Tab:       bidValues[storage.BidTab].(string),
	}
}

type ManagedCdp struct {
	Id            string `db:"cdpi"`
	Usr           string
	UrnIdentifier string `db:"urn_identifier"`
	IlkIdentifier string `db:"ilk_identifier"`
	Created       sql.NullString
}

type commonBid struct {
	BidId   string `db:"bid_id"`
	Guy     string
	Tic     string
	End     string
	Lot     string
	Bid     string
	Dealt   string
	Created sql.NullString
	Updated sql.NullString
}

type FlopBid struct {
	commonBid
}

type FlapBid struct {
	commonBid
}

type FlipBid struct {
	commonBid
	IlkId string `db:"ilk_id"`
	UrnId string `db:"urn_id"`
	Gal   string
	Tab   string
}

func SetUpFlipBidContext(setupData FlipBidContextInput) (ilkId, urnId int64, err error) {
	ilkId, ilkErr := shared.GetOrCreateIlk(setupData.IlkHex, setupData.Db)
	if ilkErr != nil {
		return 0, 0, ilkErr
	}

	urnId, urnErr := shared.GetOrCreateUrn(setupData.UrnGuy, setupData.IlkHex, setupData.Db)
	if urnErr != nil {
		return 0, 0, urnErr
	}

	flipKickLog := test_data.CreateTestLog(setupData.FlipKickHeaderId, setupData.Db)
	flipKickErr := CreateFlipKick(setupData.ContractAddress, setupData.BidId, setupData.FlipKickHeaderId, flipKickLog.ID, setupData.UrnGuy, setupData.FlipKickRepo)
	if flipKickErr != nil {
		return 0, 0, flipKickErr
	}

	if setupData.Dealt {
		dealErr := CreateDeal(setupData.DealCreationInput)
		if dealErr != nil {
			return 0, 0, dealErr
		}
	}

	return ilkId, urnId, nil
}

func SetUpFlapBidContext(setupData FlapBidCreationInput) (err error) {
	flapKickLog := test_data.CreateTestLog(setupData.FlapKickHeaderId, setupData.Db)
	flapKickErr := CreateFlapKick(setupData.ContractAddress, setupData.BidId, setupData.FlapKickHeaderId, flapKickLog.ID, setupData.FlapKickRepo)
	if flapKickErr != nil {
		return flapKickErr
	}

	if setupData.Dealt {
		dealErr := CreateDeal(setupData.DealCreationInput)
		if dealErr != nil {
			return dealErr
		}
	}

	return nil
}

func SetUpFlopBidContext(setupData FlopBidCreationInput) (err error) {
	flopKickLog := test_data.CreateTestLog(setupData.FlopKickHeaderId, setupData.Db)
	flopKickErr := CreateFlopKick(setupData.ContractAddress, setupData.BidId, setupData.FlopKickHeaderId, flopKickLog.ID, setupData.FlopKickRepo)
	if flopKickErr != nil {
		return flopKickErr
	}

	if setupData.Dealt {
		dealErr := CreateDeal(setupData.DealCreationInput)
		if dealErr != nil {
			return dealErr
		}
	}
	return nil
}

func CreateDeal(input DealCreationInput) (err error) {
	dealLog := test_data.CreateTestLog(input.DealHeaderId, input.Db)
	dealModel := test_data.CopyModel(test_data.DealModel)
	dealModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	dealModel.ColumnValues["tx_idx"] = rand.Int31()
	dealModel.ForeignKeyValues[constants.AddressFK] = input.ContractAddress
	dealModel.ColumnValues[constants.HeaderFK] = input.DealHeaderId
	dealModel.ColumnValues[constants.LogFK] = dealLog.ID
	deals := []shared.InsertionModel{dealModel}
	return input.DealRepo.Create(deals)
}

func CreateFlipKick(contractAddress string, bidId int, headerId, logId int64, usr string, repo flip_kick.FlipKickRepository) error {
	flipKickModel := test_data.CopyModel(test_data.FlipKickModel())
	flipKickModel.ForeignKeyValues[constants.AddressFK] = contractAddress
	flipKickModel.ColumnValues["bid_id"] = strconv.Itoa(bidId)
	flipKickModel.ColumnValues["usr"] = usr
	flipKickModel.ColumnValues[constants.HeaderFK] = headerId
	flipKickModel.ColumnValues[constants.LogFK] = logId
	return repo.Create([]shared.InsertionModel{flipKickModel})
}

func CreateFlapKick(contractAddress string, bidId int, headerId, logId int64, repo flap_kick.FlapKickRepository) error {
	flapKickModel := test_data.FlapKickModel()
	flapKickModel.ForeignKeyValues[constants.AddressFK] = contractAddress
	flapKickModel.ColumnValues[constants.HeaderFK] = headerId
	flapKickModel.ColumnValues[constants.LogFK] = logId
	flapKickModel.ColumnValues["bid_id"] = strconv.Itoa(bidId)
	return repo.Create([]shared.InsertionModel{flapKickModel})
}

func CreateFlopKick(contractAddress string, bidId int, headerId, logId int64, repo flop_kick.FlopKickRepository) error {
	flopKickModel := test_data.FlopKickModel()
	flopKickModel.ForeignKeyValues[constants.AddressFK] = contractAddress
	flopKickModel.ColumnValues["bid_id"] = strconv.Itoa(bidId)
	flopKickModel.ColumnValues[constants.HeaderFK] = headerId
	flopKickModel.ColumnValues[constants.LogFK] = logId
	return repo.Create([]shared.InsertionModel{flopKickModel})
}

func CreateTend(input TendCreationInput) (err error) {
	tendModel := test_data.CopyModel(test_data.TendModel)
	tendModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	tendModel.ColumnValues["lot"] = strconv.Itoa(input.Lot)
	tendModel.ColumnValues["bid"] = strconv.Itoa(input.BidAmount)
	tendModel.ForeignKeyValues[constants.AddressFK] = input.ContractAddress
	tendModel.ColumnValues[constants.HeaderFK] = input.TendHeaderId
	tendModel.ColumnValues[constants.LogFK] = input.TendLogId
	return input.TendRepo.Create([]shared.InsertionModel{tendModel})
}

func CreateDent(input DentCreationInput) (err error) {
	dentModel := test_data.CopyModel(test_data.DentModel)
	dentModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	dentModel.ColumnValues["lot"] = strconv.Itoa(input.Lot)
	dentModel.ColumnValues["bid"] = strconv.Itoa(input.BidAmount)
	dentModel.ForeignKeyValues[constants.AddressFK] = input.ContractAddress
	dentModel.ColumnValues[constants.HeaderFK] = input.DentHeaderId
	dentModel.ColumnValues[constants.LogFK] = input.DentLogId
	return input.DentRepo.Create([]shared.InsertionModel{dentModel})
}

func CreateYank(input YankCreationInput) (err error) {
	yankModel := test_data.CopyModel(test_data.YankModel)
	yankModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	yankModel.ColumnValues["tx_idx"] = rand.Int31()
	yankModel.ForeignKeyValues[constants.AddressFK] = input.ContractAddress
	yankModel.ColumnValues[constants.HeaderFK] = input.YankHeaderId
	yankModel.ColumnValues[constants.LogFK] = input.YankLogId
	return input.YankRepo.Create([]shared.InsertionModel{yankModel})
}

func CreateTick(input TickCreationInput) (err error) {
	tickModel := test_data.CopyModel(test_data.TickModel)
	tickModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	tickModel.ColumnValues["tx_idx"] = rand.Int31()
	tickModel.ForeignKeyValues[constants.AddressFK] = input.ContractAddress
	tickModel.ColumnValues[constants.HeaderFK] = input.TickHeaderId
	tickModel.ColumnValues[constants.LogFK] = input.TickLogId
	return input.TickRepo.Create([]shared.InsertionModel{tickModel})
}

type YankCreationInput struct {
	ContractAddress string
	BidId           int
	YankRepo        yank.YankRepository
	YankHeaderId    int64
	YankLogId       int64
}

type TendCreationInput struct {
	ContractAddress string
	BidId           int
	Lot             int
	BidAmount       int
	TendRepo        tend.TendRepository
	TendHeaderId    int64
	TendLogId       int64
}

type DentCreationInput struct {
	ContractAddress string
	BidId           int
	Lot             int
	BidAmount       int
	DentRepo        dent.DentRepository
	DentHeaderId    int64
	DentLogId       int64
}

type DealCreationInput struct {
	Db              *postgres.DB
	BidId           int
	ContractAddress string
	DealRepo        deal.DealRepository
	DealHeaderId    int64
}

type FlipBidContextInput struct {
	DealCreationInput
	Dealt            bool
	IlkHex           string
	UrnGuy           string
	FlipKickRepo     flip_kick.FlipKickRepository
	FlipKickHeaderId int64
}

type FlapBidCreationInput struct {
	DealCreationInput
	Dealt            bool
	FlapKickRepo     flap_kick.FlapKickRepository
	FlapKickHeaderId int64
}

type TickCreationInput struct {
	BidId           int
	ContractAddress string
	TickRepo        tick.TickRepository
	TickHeaderId    int64
	TickLogId       int64
}

type FlopBidCreationInput struct {
	DealCreationInput
	Dealt            bool
	FlopKickRepo     flop_kick.FlopKickRepository
	FlopKickHeaderId int64
}

type BidEvent struct {
	BidId           string `db:"bid_id"`
	Lot             string
	BidAmount       string `db:"bid_amount"`
	Act             string
	ContractAddress string `db:"contract_address"`
}

type IlkFileEvent struct {
	IlkIdentifier sql.NullString `db:"ilk_identifier"`
	What          string
	Data          string
}

type FrobEvent struct {
	IlkIdentifier string `db:"ilk_identifier"`
	UrnIdentifier string `db:"urn_identifier"`
	Dink          string
	Dart          string
	Rate          string `db:"ilk_rate"`
}

type BiteEvent struct {
	IlkIdentifier string `db:"ilk_identifier"`
	UrnIdentifier string `db:"urn_identifier"`
	Ink           string
	Art           string
	Tab           string
}

type SinQueueEvent struct {
	Era         string
	Act         string
	BlockHeight string `db:"block_height"`
}

type PokeEvent struct {
	IlkId string `db:"ilk_id"`
	Val   string
	Spot  string
}

func GetExpectedTimestamp(epoch int) string {
	return time.Unix(int64(epoch), 0).UTC().Format(time.RFC3339)
}

func GetValidNullString(val string) sql.NullString {
	return sql.NullString{
		String: val,
		Valid:  true,
	}
}

func GetEmptyNullString() sql.NullString {
	return sql.NullString{}
}

func GetRandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
