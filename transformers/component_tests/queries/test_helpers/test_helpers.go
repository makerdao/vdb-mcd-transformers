// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test_helpers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	mcdShared "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	vdbStorageFactory "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/gomega"
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

	EmptyMetadatas []types.ValueMetadata

	FakeIlkRateMetadata = types.GetValueMetadata(vat.IlkRate, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	FakeIlkArtMetadata  = types.GetValueMetadata(vat.IlkArt, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	FakeIlkSpotMetadata = types.GetValueMetadata(vat.IlkSpot, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	FakeIlkLineMetadata = types.GetValueMetadata(vat.IlkLine, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	FakeIlkDustMetadata = types.GetValueMetadata(vat.IlkDust, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkChopMetadata = types.GetValueMetadata(cat.IlkChop, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkDunkMetadata = types.GetValueMetadata(cat.IlkDunk, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkLumpMetadata = types.GetValueMetadata(cat.IlkLump, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkFlipMetadata = types.GetValueMetadata(cat.IlkFlip, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkRhoMetadata  = types.GetValueMetadata(jug.IlkRho, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkTaxMetadata  = types.GetValueMetadata(jug.IlkDuty, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)
	fakeIlkPipMetadata  = types.GetValueMetadata(spot.IlkPip, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Address)
	fakeIlkMatMetadata  = types.GetValueMetadata(spot.IlkMat, map[types.Key]string{constants.Ilk: FakeIlk.Hex}, types.Uint256)

	FakeIlkVatMetadatas = []types.ValueMetadata{
		FakeIlkRateMetadata,
		FakeIlkArtMetadata,
		FakeIlkSpotMetadata,
		FakeIlkLineMetadata,
		FakeIlkDustMetadata,
	}
	FakeIlkCatMetadatas = []types.ValueMetadata{
		fakeIlkChopMetadata,
		fakeIlkDunkMetadata,
		fakeIlkLumpMetadata,
		fakeIlkFlipMetadata,
	}
	FakeIlkJugMetadatas = []types.ValueMetadata{
		fakeIlkRhoMetadata,
		fakeIlkTaxMetadata,
	}
	FakeIlkSpotMetadatas = []types.ValueMetadata{
		fakeIlkPipMetadata,
		fakeIlkMatMetadata,
	}

	anotherFakeIlkRateMetadata = types.GetValueMetadata(vat.IlkRate, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkArtMetadata  = types.GetValueMetadata(vat.IlkArt, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkSpotMetadata = types.GetValueMetadata(vat.IlkSpot, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkLineMetadata = types.GetValueMetadata(vat.IlkLine, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkDustMetadata = types.GetValueMetadata(vat.IlkDust, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkChopMetadata = types.GetValueMetadata(cat.IlkChop, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkDunkMetadata = types.GetValueMetadata(cat.IlkDunk, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkLumpMetadata = types.GetValueMetadata(cat.IlkLump, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkFlipMetadata = types.GetValueMetadata(cat.IlkFlip, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Address)
	anotherFakeIlkRhoMetadata  = types.GetValueMetadata(jug.IlkRho, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkTaxMetadata  = types.GetValueMetadata(jug.IlkDuty, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)
	anotherFakeIlkPipMetadata  = types.GetValueMetadata(spot.IlkPip, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Address)
	anotherFakeIlkMatMetadata  = types.GetValueMetadata(spot.IlkMat, map[types.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, types.Uint256)

	AnotherFakeIlkVatMetadatas = []types.ValueMetadata{
		anotherFakeIlkRateMetadata,
		anotherFakeIlkArtMetadata,
		anotherFakeIlkSpotMetadata,
		anotherFakeIlkLineMetadata,
		anotherFakeIlkDustMetadata,
	}
	AnotherFakeIlkCatMetadatas = []types.ValueMetadata{
		anotherFakeIlkChopMetadata,
		anotherFakeIlkDunkMetadata,
		anotherFakeIlkLumpMetadata,
		anotherFakeIlkFlipMetadata,
	}
	AnotherFakeIlkJugMetadatas = []types.ValueMetadata{
		anotherFakeIlkRhoMetadata,
		anotherFakeIlkTaxMetadata,
	}
	AnotherFakeIlkSpotMetadatas = []types.ValueMetadata{
		anotherFakeIlkPipMetadata,
		anotherFakeIlkMatMetadata,
	}
)

type TestIlk struct {
	Hex        string
	Identifier string
}

type IlkSnapshot struct {
	IlkIdentifier string `db:"ilk_identifier"`
	BlockNumber   string `db:"block_number"`
	Rate          string
	Art           string
	Spot          string
	Line          string
	Dust          string
	Chop          string
	Dunk          string
	Lump          string
	Flip          string
	Rho           string
	Duty          string
	Pip           string
	Mat           string
	Created       sql.NullString
	Updated       sql.NullString
}

func GetIlkValues(seed int) map[string]interface{} {
	valuesMap := make(map[string]interface{})
	valuesMap[vat.IlkRate] = strconv.Itoa(1 + seed)
	valuesMap[vat.IlkArt] = strconv.Itoa(2 + seed)
	valuesMap[vat.IlkSpot] = strconv.Itoa(3 + seed)
	valuesMap[vat.IlkLine] = strconv.Itoa(4 + seed)
	valuesMap[vat.IlkDust] = strconv.Itoa(5 + seed)
	valuesMap[cat.IlkChop] = strconv.Itoa(6 + seed)
	valuesMap[cat.IlkDunk] = strconv.Itoa(7 + seed)
	valuesMap[cat.IlkLump] = strconv.Itoa(8 + seed)
	valuesMap[cat.IlkFlip] = "an address" + strconv.Itoa(seed)
	valuesMap[jug.IlkRho] = strconv.Itoa(9 + seed)
	valuesMap[jug.IlkDuty] = strconv.Itoa(10 + seed)
	valuesMap[spot.IlkPip] = "an address2" + strconv.Itoa(seed)
	valuesMap[spot.IlkMat] = strconv.Itoa(11 + seed)

	return valuesMap
}

func IlkSnapshotFromValues(ilk, updated, created string, ilkValues map[string]interface{}) IlkSnapshot {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	parsedUpdated, _ := strconv.ParseInt(updated, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)
	updatedTimestamp := time.Unix(parsedUpdated, 0).UTC().Format(time.RFC3339)

	ilkIdentifier := shared.DecodeHexToText(ilk)
	return IlkSnapshot{
		IlkIdentifier: ilkIdentifier,
		Rate:          ilkValues[vat.IlkRate].(string),
		Art:           ilkValues[vat.IlkArt].(string),
		Spot:          ilkValues[vat.IlkSpot].(string),
		Line:          ilkValues[vat.IlkLine].(string),
		Dust:          ilkValues[vat.IlkDust].(string),
		Chop:          ilkValues[cat.IlkChop].(string),
		Dunk:          ilkValues[cat.IlkDunk].(string),
		Lump:          ilkValues[cat.IlkLump].(string),
		Flip:          ilkValues[cat.IlkFlip].(string),
		Rho:           ilkValues[jug.IlkRho].(string),
		Duty:          ilkValues[jug.IlkDuty].(string),
		Pip:           ilkValues[spot.IlkPip].(string),
		Mat:           ilkValues[spot.IlkMat].(string),
		Created:       sql.NullString{String: createdTimestamp, Valid: true},
		Updated:       sql.NullString{String: updatedTimestamp, Valid: true},
	}
}

func CreateVatRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []types.ValueMetadata, repository vat.StorageRepository) {
	InsertValues(db, &repository, header, valuesMap, metadatas)
}

func CreateCatRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []types.ValueMetadata, repository cat.StorageRepository) {
	InsertValues(db, &repository, header, valuesMap, metadatas)
}

func CreateJugRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []types.ValueMetadata, repository jug.StorageRepository) {
	InsertValues(db, &repository, header, valuesMap, metadatas)
}

func CreateSpotRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []types.ValueMetadata, repository spot.StorageRepository) {
	InsertValues(db, &repository, header, valuesMap, metadatas)
}

// Creates urn by creating necessary state diffs and the corresponding header
func CreateUrn(db *postgres.DB, setupData map[string]interface{}, header core.Header, metadata UrnMetadata, vatRepo vat.StorageRepository) {
	// This also creates the ilk if it doesn't exist
	urnMetadata := []types.ValueMetadata{metadata.UrnInk, metadata.UrnArt}
	InsertValues(db, &vatRepo, header, setupData, urnMetadata)
}

func CreateIlk(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, vatMetadatas, catMetadatas, jugMetadatas, spotMetadatas []types.ValueMetadata) {
	var (
		vatRepo  vat.StorageRepository
		catRepo  cat.StorageRepository
		jugRepo  jug.StorageRepository
		spotRepo spot.StorageRepository
	)
	vatRepo.SetDB(db)
	catRepo.SetDB(db)
	jugRepo.SetDB(db)
	spotRepo.SetDB(db)
	CreateVatRecords(db, header, valuesMap, vatMetadatas, vatRepo)
	CreateCatRecords(db, header, valuesMap, catMetadatas, catRepo)
	CreateJugRecords(db, header, valuesMap, jugMetadatas, jugRepo)
	CreateSpotRecords(db, header, valuesMap, spotMetadatas, spotRepo)
	CreateVatInit(db, header.Id, vatMetadatas[0].Keys[constants.Ilk])
}

func CreateVatInit(db *postgres.DB, headerID int64, ilkHex string) event.InsertionModel {
	ilkID, ilkErr := mcdShared.GetOrCreateIlk(ilkHex, db)
	Expect(ilkErr).NotTo(HaveOccurred())

	logID := test_data.CreateTestLog(headerID, db).ID

	vatInit := test_data.VatInitModel()
	vatInit.ColumnValues[constants.IlkColumn] = ilkID
	vatInit.ColumnValues[event.HeaderFK] = headerID
	vatInit.ColumnValues[event.LogFK] = logID

	insertErr := event.PersistModels([]event.InsertionModel{vatInit}, db)
	Expect(insertErr).NotTo(HaveOccurred())

	return vatInit
}

func GetUrnSetupData() map[string]interface{} {
	urnData := make(map[string]interface{})
	urnData[vat.UrnInk] = rand.Int()
	urnData[vat.UrnArt] = rand.Int()
	return urnData
}

func GetUrnMetadata(ilk, urn string) UrnMetadata {
	return UrnMetadata{
		UrnInk: types.GetValueMetadata(vat.UrnInk,
			map[types.Key]string{constants.Ilk: ilk, constants.Guy: urn}, types.Uint256),
		UrnArt: types.GetValueMetadata(vat.UrnArt,
			map[types.Key]string{constants.Ilk: ilk, constants.Guy: urn}, types.Uint256),
	}
}

type UrnMetadata struct {
	UrnInk types.ValueMetadata
	UrnArt types.ValueMetadata
}

type UrnState struct {
	UrnIdentifier string `db:"urn_identifier"`
	IlkIdentifier string `db:"ilk_identifier"`
	BlockHeight   int    `db:"block_height"`
	Ink           string
	Art           string
	Created       sql.NullString
	Updated       sql.NullString
}

func AssertUrn(actual, expected UrnState) {
	Expect(actual.UrnIdentifier).To(Equal(expected.UrnIdentifier), "Urn Identifier")
	Expect(actual.IlkIdentifier).To(Equal(expected.IlkIdentifier), "Ilk Identifier")
	Expect(actual.BlockHeight).To(Equal(expected.BlockHeight), "Block Height")
	Expect(actual.Ink).To(Equal(expected.Ink), "Ink")
	Expect(actual.Art).To(Equal(expected.Art), "Art")
	Expect(actual.Created).To(Equal(expected.Created), "Created")
	Expect(actual.Updated).To(Equal(expected.Updated), "Updated")
}

func GetCommonBidMetadatas(bidId string) []types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidId}
	packedNames := map[int]string{0: storage.BidGuy, 1: storage.BidTic, 2: storage.BidEnd}
	packedTypes := map[int]types.ValueType{0: types.Address, 1: types.Uint48, 2: types.Uint48}
	return []types.ValueMetadata{
		types.GetValueMetadata(storage.Kicks, nil, types.Uint256),
		types.GetValueMetadata(storage.BidBid, keys, types.Uint256),
		types.GetValueMetadata(storage.BidLot, keys, types.Uint256),
		types.GetValueMetadataForPackedSlot(storage.Packed, keys, types.PackedSlot, packedNames, packedTypes),
	}
}

func GetFlopMetadatas(bidId string) []types.ValueMetadata {
	return GetCommonBidMetadatas(bidId)
}

func GetFlapMetadatas(bidId string) []types.ValueMetadata {
	return GetCommonBidMetadatas(bidId)
}

func GetCdpManagerMetadatas(cdpi string) []types.ValueMetadata {
	keys := map[types.Key]string{constants.Cdpi: cdpi}
	return []types.ValueMetadata{
		types.GetValueMetadata(cdp_manager.Cdpi, nil, types.Uint256),
		types.GetValueMetadata(cdp_manager.Urns, keys, types.Address),
		types.GetValueMetadata(cdp_manager.Owns, keys, types.Address),
		types.GetValueMetadata(cdp_manager.Ilks, keys, types.Bytes32),
	}
}

func GetFlipMetadatas(bidId string) []types.ValueMetadata {
	keys := map[types.Key]string{constants.BidId: bidId}
	return append(GetCommonBidMetadatas(bidId),
		types.GetValueMetadata(storage.Ilk, nil, types.Bytes32),
		types.GetValueMetadata(storage.BidUsr, keys, types.Address),
		types.GetValueMetadata(storage.BidGal, keys, types.Address),
		types.GetValueMetadata(storage.BidTab, keys, types.Uint256))
}

func GetCdpManagerStorageValues(seed int, ilkHex string, urnGuy string, cdpi int) map[string]interface{} {
	valuesMap := make(map[string]interface{})
	valuesMap[cdp_manager.Cdpi] = strconv.Itoa(cdpi)
	valuesMap[cdp_manager.Urns] = urnGuy
	valuesMap[cdp_manager.Owns] = "address1" + strconv.Itoa(seed)
	valuesMap[cdp_manager.Ilks] = ilkHex
	return valuesMap
}

func GetCommonBidStorageValues(seed, bidId int) map[string]interface{} {
	packedValues := map[int]string{0: "address1" + strconv.Itoa(seed), 1: strconv.Itoa(1 + seed), 2: strconv.Itoa(2 + seed)}
	valuesMap := make(map[string]interface{})
	valuesMap[storage.Kicks] = strconv.Itoa(bidId)
	valuesMap[storage.BidBid] = strconv.Itoa(3 + seed)
	valuesMap[storage.BidLot] = strconv.Itoa(4 + seed)
	valuesMap[storage.Packed] = packedValues

	return valuesMap
}

func GetFlopStorageValues(seed, bidId int) map[string]interface{} {
	return GetCommonBidStorageValues(seed, bidId)
}

func GetFlapStorageValues(seed, bidId int) map[string]interface{} {
	return GetCommonBidStorageValues(seed, bidId)
}

func GetFlipStorageValues(seed int, ilk string, bidId int) map[string]interface{} {
	valuesMap := GetCommonBidStorageValues(seed, bidId)
	valuesMap[storage.Ilk] = ilk
	valuesMap[storage.BidGal] = "address2" + strconv.Itoa(seed)
	valuesMap[storage.BidUsr] = "address3" + strconv.Itoa(seed)
	valuesMap[storage.BidTab] = strconv.Itoa(5 + seed)
	return valuesMap
}

func InsertValues(db *postgres.DB, repo vdbStorageFactory.Repository, header core.Header, valuesMap map[string]interface{}, metadatas []types.ValueMetadata) {
	for _, metadata := range metadatas {
		value := valuesMap[metadata.Name]
		key := common.HexToHash(test_data.RandomString(32))
		var valueForDiffRecord common.Hash
		var valueForStorageRecord interface{}

		switch v := value.(type) {
		case string:
			valueForDiffRecord = common.HexToHash(v)
			valueForStorageRecord = v
		case int:
			valueForDiffRecord = common.HexToHash(strconv.Itoa(v))
			valueForStorageRecord = strconv.Itoa(v)
		case map[int]string:
			values := make([]string, 0, len(v))
			for _, value := range v {
				values = append(values, value)
			}
			valueForDiffRecord = common.HexToHash(strings.Join(values, ""))
			valueForStorageRecord = v
		default:
			panic(fmt.Sprintf("valuesMap value type not recognized %v", v))
		}

		persistedDiff := test_helpers.CreateDiffRecord(db, header, common.Address{}, key, valueForDiffRecord)

		err := repo.Create(persistedDiff.ID, header.Id, metadata, valueForStorageRecord)
		Expect(err).NotTo(HaveOccurred())
	}
}

func CreateFlop(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flopMetadatas []types.ValueMetadata, contractAddress string) {
	flopRepo := flop.StorageRepository{ContractAddress: contractAddress}
	flopRepo.SetDB(db)
	InsertValues(db, &flopRepo, header, valuesMap, flopMetadatas)
}

func CreateFlap(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flapMetadatas []types.ValueMetadata, contractAddress string) {
	flapRepo := flap.StorageRepository{ContractAddress: contractAddress}
	flapRepo.SetDB(db)
	InsertValues(db, &flapRepo, header, valuesMap, flapMetadatas)
}

func CreateFlip(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flipMetadatas []types.ValueMetadata, contractAddress string) {
	flipRepo := flip.StorageRepository{ContractAddress: contractAddress}
	flipRepo.SetDB(db)
	InsertValues(db, &flipRepo, header, valuesMap, flipMetadatas)
}

func CreateManagedCdp(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []types.ValueMetadata) error {
	cdpManagerRepo := cdp_manager.StorageRepository{}
	cdpManagerRepo.SetDB(db)
	_, err := mcdShared.GetOrCreateUrn(valuesMap[cdp_manager.Urns].(string), valuesMap[cdp_manager.Ilks].(string), db)
	if err != nil {
		return err
	}
	InsertValues(db, &cdpManagerRepo, header, valuesMap, metadatas)
	return nil
}

func ManagedCdpFromValues(ilkIdentifier, created string, cdpValues map[string]interface{}) ManagedCdp {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)

	return ManagedCdp{
		Usr:           cdpValues[cdp_manager.Owns].(string),
		Id:            cdpValues[cdp_manager.Cdpi].(string),
		UrnIdentifier: cdpValues[cdp_manager.Urns].(string),
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

func FlipBidFromValues(bidId, ilkId, urnId, dealt, updated, created, flipAddress string, bidValues map[string]interface{}) FlipBid {
	return FlipBid{
		commonBid:   commonBidFromValues(bidId, dealt, updated, created, bidValues),
		IlkId:       ilkId,
		UrnId:       urnId,
		Gal:         bidValues[storage.BidGal].(string),
		Tab:         bidValues[storage.BidTab].(string),
		FlipAddress: flipAddress,
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
	IlkId       string `db:"ilk_id"`
	UrnId       string `db:"urn_id"`
	Gal         string
	Tab         string
	FlipAddress string `db:"flip_address"`
	BlockHeight string `db:"block_height"`
}

func SetUpFlipBidContext(setupData FlipBidContextInput) (ilkId, urnId int64, err error) {
	ilkId, ilkErr := mcdShared.GetOrCreateIlk(setupData.IlkHex, setupData.DB)
	if ilkErr != nil {
		return 0, 0, ilkErr
	}

	urnId, urnErr := mcdShared.GetOrCreateUrn(setupData.UrnGuy, setupData.IlkHex, setupData.DB)
	if urnErr != nil {
		return 0, 0, urnErr
	}

	flipKickLog := test_data.CreateTestLog(setupData.FlipKickHeaderId, setupData.DB)
	flipKickErr := CreateFlipKick(setupData.ContractAddress, setupData.BidId, setupData.FlipKickHeaderId, flipKickLog.ID, setupData.UrnGuy, setupData.DB)
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
	flapKickLog := test_data.CreateTestLog(setupData.FlapKickHeaderId, setupData.DB)
	flapKickErr := CreateFlapKick(setupData.ContractAddress, setupData.BidId, setupData.FlapKickHeaderId, flapKickLog.ID, setupData.DB)
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
	flopKickLog := test_data.CreateTestLog(setupData.FlopKickHeaderId, setupData.DB)
	flopKickErr := CreateFlopKick(setupData.ContractAddress, setupData.BidId, setupData.FlopKickHeaderId, flopKickLog.ID, setupData.DB)
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
	addressID, addressErr := repository.GetOrCreateAddress(input.DB, input.ContractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	dealLog := test_data.CreateTestLog(input.DealHeaderId, input.DB)
	dealModel := test_data.DealModel()
	dealModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(input.BidId)
	dealModel.ColumnValues[event.HeaderFK] = input.DealHeaderId
	dealModel.ColumnValues[event.LogFK] = dealLog.ID
	dealModel.ColumnValues[event.AddressFK] = addressID
	test_data.AssignMessageSenderID(test_data.DealEventLog, dealModel, input.DB)
	deals := []event.InsertionModel{dealModel}
	return event.PersistModels(deals, input.DB)
}

func CreateFlipKick(contractAddress string, bidId int, headerId, logId int64, usr string, db *postgres.DB) error {
	addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	flipKickModel := test_data.FlipKickModel()
	flipKickModel.ColumnValues[event.HeaderFK] = headerId
	flipKickModel.ColumnValues[event.LogFK] = logId
	flipKickModel.ColumnValues[event.AddressFK] = addressId
	flipKickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
	flipKickModel.ColumnValues[constants.UsrColumn] = usr
	return event.PersistModels([]event.InsertionModel{flipKickModel}, db)
}

func CreateFlapKick(contractAddress string, bidId int, headerId, logId int64, db *postgres.DB) error {
	flapKickModel := test_data.FlapKickModel()
	addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	flapKickModel.ColumnValues[event.HeaderFK] = headerId
	flapKickModel.ColumnValues[event.LogFK] = logId
	flapKickModel.ColumnValues[event.AddressFK] = addressId
	flapKickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
	return event.PersistModels([]event.InsertionModel{flapKickModel}, db)
}

func CreateFlopKick(contractAddress string, bidId int, headerId, logId int64, db *postgres.DB) error {
	addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	flopKickModel := test_data.FlopKickModel()
	flopKickModel.ColumnValues[event.HeaderFK] = headerId
	flopKickModel.ColumnValues[event.LogFK] = logId
	flopKickModel.ColumnValues[event.AddressFK] = addressId
	flopKickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
	return event.PersistModels([]event.InsertionModel{flopKickModel}, db)
}

func CreateTend(input TendCreationInput) (err error) {
	addressID, addressErr := repository.GetOrCreateAddress(input.DB, input.ContractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	tendModel := test_data.TendModel()
	tendLog := test_data.CreateTestLog(input.TendHeaderId, input.DB)
	tendModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(input.BidId)
	tendModel.ColumnValues[constants.LotColumn] = strconv.Itoa(input.Lot)
	tendModel.ColumnValues[constants.BidColumn] = strconv.Itoa(input.BidAmount)
	tendModel.ColumnValues[event.AddressFK] = addressID
	tendModel.ColumnValues[event.HeaderFK] = input.TendHeaderId
	tendModel.ColumnValues[event.LogFK] = tendLog.ID
	test_data.AssignMessageSenderID(test_data.TendEventLog, tendModel, input.DB)
	return event.PersistModels([]event.InsertionModel{tendModel}, input.DB)
}

func CreateDent(input DentCreationInput) (err error) {
	addressID, addressErr := repository.GetOrCreateAddress(input.DB, input.ContractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	dentModel := test_data.DentModel()
	dentModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(input.BidId)
	dentModel.ColumnValues[constants.LotColumn] = strconv.Itoa(input.Lot)
	dentModel.ColumnValues[constants.BidColumn] = strconv.Itoa(input.BidAmount)
	dentModel.ColumnValues[event.AddressFK] = addressID
	dentModel.ColumnValues[event.HeaderFK] = input.DentHeaderId
	dentModel.ColumnValues[event.LogFK] = input.DentLogId
	test_data.AssignMessageSenderID(test_data.DentEventLog, dentModel, input.DB)
	return event.PersistModels([]event.InsertionModel{dentModel}, input.DB)
}

func CreateYank(input YankCreationInput) (err error) {
	addressID, addressErr := repository.GetOrCreateAddress(input.DB, input.ContractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	yankModel := test_data.YankModel()
	yankModel.ColumnValues[event.AddressFK] = addressID
	yankModel.ColumnValues[event.HeaderFK] = input.YankHeaderId
	yankModel.ColumnValues[event.LogFK] = input.YankLogId
	test_data.AssignMessageSenderID(test_data.YankEventLog, yankModel, input.DB)
	yankModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(input.BidId)
	return event.PersistModels([]event.InsertionModel{yankModel}, input.DB)
}

func CreateTick(input TickCreationInput) (err error) {
	addressID, addressErr := repository.GetOrCreateAddress(input.DB, input.ContractAddress)
	Expect(addressErr).NotTo(HaveOccurred())
	tickLog := test_data.CreateTestLog(input.TickHeaderId, input.DB)
	tickModel := test_data.TickModel()
	tickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(input.BidId)
	tickModel.ColumnValues[event.AddressFK] = addressID
	tickModel.ColumnValues[event.HeaderFK] = input.TickHeaderId
	tickModel.ColumnValues[event.LogFK] = tickLog.ID
	test_data.AssignMessageSenderID(test_data.FlipTickEventLog, tickModel, input.DB)
	return event.PersistModels([]event.InsertionModel{tickModel}, input.DB)
}

type YankCreationInput struct {
	DB              *postgres.DB
	ContractAddress string
	BidId           int
	YankHeaderId    int64
	YankLogId       int64
}

type TendCreationInput struct {
	DB              *postgres.DB
	ContractAddress string
	BidId           int
	Lot             int
	BidAmount       int
	TendHeaderId    int64
	TendLogId       int64
}

type DentCreationInput struct {
	DB              *postgres.DB
	ContractAddress string
	BidId           int
	Lot             int
	BidAmount       int
	DentHeaderId    int64
	DentLogId       int64
}

type DealCreationInput struct {
	DB              *postgres.DB
	BidId           int
	ContractAddress string
	DealHeaderId    int64
}

type FlipBidContextInput struct {
	DealCreationInput
	Dealt            bool
	IlkHex           string
	UrnGuy           string
	FlipKickHeaderId int64
}

type FlapBidCreationInput struct {
	DealCreationInput
	Dealt            bool
	FlapKickHeaderId int64
}

type TickCreationInput struct {
	DB              *postgres.DB
	BidId           int
	ContractAddress string
	TickHeaderId    int64
	TickLogId       int64
}

type FlopBidCreationInput struct {
	DealCreationInput
	Dealt            bool
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
