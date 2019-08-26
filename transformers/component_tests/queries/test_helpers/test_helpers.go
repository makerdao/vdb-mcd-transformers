package test_helpers

import (
	"database/sql"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_tick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
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
		Expect(!expected.Ratio.Valid)
	}

	Expect(actual.Safe).To(Equal(expected.Safe))
	Expect(actual.Created).To(Equal(expected.Created))
	Expect(actual.Updated).To(Equal(expected.Updated))
}

func GetFlopMetadatas(bidId string) []utils.StorageValueMetadata {
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

func GetFlapMetadatas(bidId string) []utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return append(GetFlopMetadatas(bidId), utils.GetStorageValueMetadata(storage.BidGal, keys, utils.Address))
}

func GetFlipMetadatas(bidId string) []utils.StorageValueMetadata {
	keys := map[utils.Key]string{constants.BidId: bidId}
	return append(GetFlapMetadatas(bidId),
		utils.GetStorageValueMetadata(storage.Ilk, nil, utils.Bytes32),
		utils.GetStorageValueMetadata(storage.BidUsr, keys, utils.Address),
		utils.GetStorageValueMetadata(storage.BidTab, keys, utils.Uint256))
}

func GetFlopStorageValues(seed int, bidId int) map[string]interface{} {
	packedValues := map[int]string{0: "address1" + strconv.Itoa(seed), 1: strconv.Itoa(1 + seed), 2: strconv.Itoa(2 + seed)}
	valuesMap := make(map[string]interface{})
	valuesMap[storage.Kicks] = strconv.Itoa(bidId)
	valuesMap[storage.BidBid] = strconv.Itoa(3 + seed)
	valuesMap[storage.BidLot] = strconv.Itoa(4 + seed)
	valuesMap[storage.Packed] = packedValues

	return valuesMap
}

func GetFlapStorageValues(seed int, bidId int) map[string]interface{} {
	valuesMap := GetFlopStorageValues(seed, bidId)
	valuesMap[storage.BidGal] = "address2" + strconv.Itoa(seed)
	return valuesMap
}

func GetFlipStorageValues(seed int, ilk string, bidId int) map[string]interface{} {
	valuesMap := GetFlapStorageValues(seed, bidId)
	valuesMap[storage.Ilk] = ilk
	valuesMap[storage.BidUsr] = "address3" + strconv.Itoa(seed)
	valuesMap[storage.BidTab] = strconv.Itoa(5 + seed)
	return valuesMap
}

func createBid(repo repository.StorageRepository, header core.Header, valuesMap map[string]interface{}, metadatas []utils.StorageValueMetadata) {
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
	createBid(&flopRepo, header, valuesMap, flopMetadatas)
}

func CreateFlap(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flapMetadatas []utils.StorageValueMetadata, contractAddress string) {
	flapRepo := flap.FlapStorageRepository{ContractAddress: contractAddress}
	flapRepo.SetDB(db)
	createBid(&flapRepo, header, valuesMap, flapMetadatas)
}

func CreateFlip(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flipMetadatas []utils.StorageValueMetadata, contractAddress string) {
	flipRepo := flip.FlipStorageRepository{ContractAddress: contractAddress}
	flipRepo.SetDB(db)
	createBid(&flipRepo, header, valuesMap, flipMetadatas)
}

func FlopBidFromValues(bidId, dealt, updated, created string, bidValues map[string]interface{}) FlopBid {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	parsedUpdated, _ := strconv.ParseInt(updated, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)
	updatedTimestamp := time.Unix(parsedUpdated, 0).UTC().Format(time.RFC3339)
	packedValues := bidValues[storage.Packed].(map[int]string)

	return FlopBid{
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

func FlapBidFromValues(bidId, dealt, updated, created string, bidValues map[string]interface{}) FlapBid {
	return FlapBid{
		FlopBid: FlopBidFromValues(bidId, dealt, updated, created, bidValues),
		Gal:     bidValues[storage.BidGal].(string),
	}
}

func FlipBidFromValues(bidId, ilkId, urnId, dealt, updated, created string, bidValues map[string]interface{}) FlipBid {
	return FlipBid{
		FlapBid: FlapBidFromValues(bidId, dealt, updated, created, bidValues),
		IlkId:   ilkId,
		UrnId:   urnId,
		Tab:     bidValues[storage.BidTab].(string),
	}
}

type FlopBid struct {
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

type FlapBid struct {
	FlopBid
	Gal string
}

type FlipBid struct {
	FlapBid
	IlkId string `db:"ilk_id"`
	UrnId string `db:"urn_id"`
	Tab   string
}

func SetUpFlipBidContext(setupData FlipBidContextInput) (ilkId, urnId int, err error) {
	ilkId, ilkErr := shared.GetOrCreateIlk(setupData.IlkHex, setupData.Db)
	if ilkErr != nil {
		return 0, 0, ilkErr
	}

	urnId, urnErr := shared.GetOrCreateUrn(setupData.UrnGuy, setupData.IlkHex, setupData.Db)
	if urnErr != nil {
		return 0, 0, urnErr
	}

	flipKickErr := CreateFlipKick(setupData.ContractAddress, setupData.BidId, setupData.FlipKickHeaderId, setupData.UrnGuy, setupData.FlipKickRepo)
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

func CreateDeal(input DealCreationInput) (err error) {
	dealModel := test_data.DealModel
	dealModel.ColumnValues["contract_address"] = input.ContractAddress
	dealModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	deals := []shared.InsertionModel{dealModel}
	return input.DealRepo.Create(input.DealHeaderId, deals)
}

func CreateFlipKick(contractAddress string, bidId int, headerId int64, usr string, repo flip_kick.FlipKickRepository) error {
	flipKickModel := test_data.FlipKickModel
	flipKickModel.ContractAddress = contractAddress
	flipKickModel.BidId = strconv.Itoa(bidId)
	flipKickModel.Usr = usr
	return repo.Create(headerId, []interface{}{flipKickModel})
}

func CreateTend(input TendCreationInput) (err error) {
	tendModel := test_data.TendModel
	tendModel.ColumnValues["contract_address"] = input.ContractAddress
	tendModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	tendModel.ColumnValues["lot"] = strconv.Itoa(input.Lot)
	tendModel.ColumnValues["bid"] = strconv.Itoa(input.BidAmount)
	if input.LogIndex != 0 {
		tendModel.ColumnValues["log_idx"] = input.LogIndex
	}
	if input.TxIndex != 0 {
		tendModel.ColumnValues["tx_idx"] = input.TxIndex
	}
	return input.TendRepo.Create(input.TendHeaderId, []shared.InsertionModel{tendModel})
}

func CreateDent(input DentCreationInput) (err error) {
	dentModel := test_data.DentModel
	dentModel.ColumnValues["contract_address"] = input.ContractAddress
	dentModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	dentModel.ColumnValues["lot"] = strconv.Itoa(input.Lot)
	dentModel.ColumnValues["bid"] = strconv.Itoa(input.BidAmount)
	if input.LogIndex != 0 {
		dentModel.ColumnValues["log_idx"] = input.LogIndex
	}
	if input.TxIndex != 0 {
		dentModel.ColumnValues["tx_idx"] = input.TxIndex
	}
	return input.DentRepo.Create(input.DentHeaderId, []shared.InsertionModel{dentModel})
}

func CreateYank(input YankCreationInput) (err error) {
	yankModel := test_data.YankModel
	yankModel.ColumnValues["contract_address"] = input.ContractAddress
	yankModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	return input.YankRepo.Create(input.YankHeaderId, []shared.InsertionModel{yankModel})
}

func CreateFlipTick(input FlipTickCreationInput) (err error) {
	flipTickModel := test_data.FlipTickModel
	flipTickModel.ColumnValues["contract_address"] = input.ContractAddress
	flipTickModel.ColumnValues["bid_id"] = strconv.Itoa(input.BidId)
	return input.FlipTickRepo.Create(input.FlipTickHeaderId, []shared.InsertionModel{flipTickModel})
}

type YankCreationInput struct {
	ContractAddress string
	BidId           int
	YankRepo        yank.YankRepository
	YankHeaderId    int64
}

type TendCreationInput struct {
	ContractAddress string
	BidId           int
	Lot             int
	BidAmount       int
	TxIndex         int
	LogIndex        int
	TendRepo        tend.TendRepository
	TendHeaderId    int64
}

type DentCreationInput struct {
	ContractAddress string
	BidId           int
	Lot             int
	BidAmount       int
	TxIndex         int
	LogIndex        int
	DentRepo        dent.DentRepository
	DentHeaderId    int64
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

type FlipTickCreationInput struct {
	BidId            int
	ContractAddress  string
	FlipTickRepo     flip_tick.FlipTickRepository
	FlipTickHeaderId int64
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
}

type BiteEvent struct {
	IlkIdentifier string `db:"ilk_identifier"`
	UrnIdentifier string `db:"urn_identifier"`
	Ink           string
	Art           string
	Tab           string
}

type SinQueueEvent struct {
	Era string
	Act string
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
