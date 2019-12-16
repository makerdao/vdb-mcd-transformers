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
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"

	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dent"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/tend"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/yank"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	vdbStorageFactory "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
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

	EmptyMetadatas []vdbStorage.ValueMetadata

	FakeIlkRateMetadata = vdbStorage.GetValueMetadata(vat.IlkRate, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	FakeIlkArtMetadata  = vdbStorage.GetValueMetadata(vat.IlkArt, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	FakeIlkSpotMetadata = vdbStorage.GetValueMetadata(vat.IlkSpot, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	FakeIlkLineMetadata = vdbStorage.GetValueMetadata(vat.IlkLine, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	FakeIlkDustMetadata = vdbStorage.GetValueMetadata(vat.IlkDust, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	fakeIlkChopMetadata = vdbStorage.GetValueMetadata(cat.IlkChop, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	fakeIlkLumpMetadata = vdbStorage.GetValueMetadata(cat.IlkLump, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	fakeIlkFlipMetadata = vdbStorage.GetValueMetadata(cat.IlkFlip, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	fakeIlkRhoMetadata  = vdbStorage.GetValueMetadata(jug.IlkRho, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	fakeIlkTaxMetadata  = vdbStorage.GetValueMetadata(jug.IlkDuty, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)
	fakeIlkPipMetadata  = vdbStorage.GetValueMetadata(spot.IlkPip, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Address)
	fakeIlkMatMetadata  = vdbStorage.GetValueMetadata(spot.IlkMat, map[vdbStorage.Key]string{constants.Ilk: FakeIlk.Hex}, vdbStorage.Uint256)

	FakeIlkVatMetadatas = []vdbStorage.ValueMetadata{
		FakeIlkRateMetadata,
		FakeIlkArtMetadata,
		FakeIlkSpotMetadata,
		FakeIlkLineMetadata,
		FakeIlkDustMetadata,
	}
	FakeIlkCatMetadatas = []vdbStorage.ValueMetadata{
		fakeIlkChopMetadata,
		fakeIlkLumpMetadata,
		fakeIlkFlipMetadata,
	}
	FakeIlkJugMetadatas = []vdbStorage.ValueMetadata{
		fakeIlkRhoMetadata,
		fakeIlkTaxMetadata,
	}
	FakeIlkSpotMetadatas = []vdbStorage.ValueMetadata{
		fakeIlkPipMetadata,
		fakeIlkMatMetadata,
	}

	anotherFakeIlkRateMetadata = vdbStorage.GetValueMetadata(vat.IlkRate, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkArtMetadata  = vdbStorage.GetValueMetadata(vat.IlkArt, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkSpotMetadata = vdbStorage.GetValueMetadata(vat.IlkSpot, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkLineMetadata = vdbStorage.GetValueMetadata(vat.IlkLine, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkDustMetadata = vdbStorage.GetValueMetadata(vat.IlkDust, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkChopMetadata = vdbStorage.GetValueMetadata(cat.IlkChop, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkLumpMetadata = vdbStorage.GetValueMetadata(cat.IlkLump, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkFlipMetadata = vdbStorage.GetValueMetadata(cat.IlkFlip, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Address)
	anotherFakeIlkRhoMetadata  = vdbStorage.GetValueMetadata(jug.IlkRho, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkTaxMetadata  = vdbStorage.GetValueMetadata(jug.IlkDuty, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)
	anotherFakeIlkPipMetadata  = vdbStorage.GetValueMetadata(spot.IlkPip, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Address)
	anotherFakeIlkMatMetadata  = vdbStorage.GetValueMetadata(spot.IlkMat, map[vdbStorage.Key]string{constants.Ilk: AnotherFakeIlk.Hex}, vdbStorage.Uint256)

	AnotherFakeIlkVatMetadatas = []vdbStorage.ValueMetadata{
		anotherFakeIlkRateMetadata,
		anotherFakeIlkArtMetadata,
		anotherFakeIlkSpotMetadata,
		anotherFakeIlkLineMetadata,
		anotherFakeIlkDustMetadata,
	}
	AnotherFakeIlkCatMetadatas = []vdbStorage.ValueMetadata{
		anotherFakeIlkChopMetadata,
		anotherFakeIlkLumpMetadata,
		anotherFakeIlkFlipMetadata,
	}
	AnotherFakeIlkJugMetadatas = []vdbStorage.ValueMetadata{
		anotherFakeIlkRhoMetadata,
		anotherFakeIlkTaxMetadata,
	}
	AnotherFakeIlkSpotMetadatas = []vdbStorage.ValueMetadata{
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
	BlockNumber   string `db:"block_number"`
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

func GetIlkValues(seed int) map[string]interface{} {
	valuesMap := make(map[string]interface{})
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

func IlkStateFromValues(ilk, updated, created string, ilkValues map[string]interface{}) IlkState {
	parsedCreated, _ := strconv.ParseInt(created, 10, 64)
	parsedUpdated, _ := strconv.ParseInt(updated, 10, 64)
	createdTimestamp := time.Unix(parsedCreated, 0).UTC().Format(time.RFC3339)
	updatedTimestamp := time.Unix(parsedUpdated, 0).UTC().Format(time.RFC3339)

	ilkIdentifier := shared.DecodeHexToText(ilk)
	return IlkState{
		IlkIdentifier: ilkIdentifier,
		Rate:          ilkValues[vat.IlkRate].(string),
		Art:           ilkValues[vat.IlkArt].(string),
		Spot:          ilkValues[vat.IlkSpot].(string),
		Line:          ilkValues[vat.IlkLine].(string),
		Dust:          ilkValues[vat.IlkDust].(string),
		Chop:          ilkValues[cat.IlkChop].(string),
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

func CreateVatRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []vdbStorage.ValueMetadata, repository vat.VatStorageRepository) {
	insertValues(db, &repository, header, valuesMap, metadatas)
}

func CreateCatRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []vdbStorage.ValueMetadata, repository cat.CatStorageRepository) {
	insertValues(db, &repository, header, valuesMap, metadatas)
}

func CreateJugRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []vdbStorage.ValueMetadata, repository jug.JugStorageRepository) {
	insertValues(db, &repository, header, valuesMap, metadatas)
}

func CreateSpotRecords(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []vdbStorage.ValueMetadata, repository spot.SpotStorageRepository) {
	insertValues(db, &repository, header, valuesMap, metadatas)
}

// Creates urn by creating necessary state diffs and the corresponding header
func CreateUrn(db *postgres.DB, setupData map[string]interface{}, header core.Header, metadata UrnMetadata, vatRepo vat.VatStorageRepository) {
	// This also creates the ilk if it doesn't exist
	urnMetadata := []vdbStorage.ValueMetadata{metadata.UrnInk, metadata.UrnArt}
	insertValues(db, &vatRepo, header, setupData, urnMetadata)
}

func CreateIlk(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, vatMetadatas, catMetadatas, jugMetadatas, spotMetadatas []vdbStorage.ValueMetadata) {
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
	CreateVatRecords(db, header, valuesMap, vatMetadatas, vatRepo)
	CreateCatRecords(db, header, valuesMap, catMetadatas, catRepo)
	CreateJugRecords(db, header, valuesMap, jugMetadatas, jugRepo)
	CreateSpotRecords(db, header, valuesMap, spotMetadatas, spotRepo)
}

func GetUrnSetupData() map[string]interface{} {
	urnData := make(map[string]interface{})
	urnData[vat.UrnInk] = rand.Int()
	urnData[vat.UrnArt] = rand.Int()
	return urnData
}

func GetUrnMetadata(ilk, urn string) UrnMetadata {
	return UrnMetadata{
		UrnInk: vdbStorage.GetValueMetadata(vat.UrnInk,
			map[vdbStorage.Key]string{constants.Ilk: ilk, constants.Guy: urn}, vdbStorage.Uint256),
		UrnArt: vdbStorage.GetValueMetadata(vat.UrnArt,
			map[vdbStorage.Key]string{constants.Ilk: ilk, constants.Guy: urn}, vdbStorage.Uint256),
	}
}

type UrnMetadata struct {
	UrnInk vdbStorage.ValueMetadata
	UrnArt vdbStorage.ValueMetadata
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
	Expect(actual.UrnIdentifier).To(Equal(expected.UrnIdentifier))
	Expect(actual.IlkIdentifier).To(Equal(expected.IlkIdentifier))
	Expect(actual.BlockHeight).To(Equal(expected.BlockHeight))
	Expect(actual.Ink).To(Equal(expected.Ink))
	Expect(actual.Art).To(Equal(expected.Art))
	Expect(actual.Created).To(Equal(expected.Created))
	Expect(actual.Updated).To(Equal(expected.Updated))
}

func getCommonBidMetadatas(bidId string) []vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	packedNames := map[int]string{0: storage.BidGuy, 1: storage.BidTic, 2: storage.BidEnd}
	packedTypes := map[int]vdbStorage.ValueType{0: vdbStorage.Address, 1: vdbStorage.Uint48, 2: vdbStorage.Uint48}
	return []vdbStorage.ValueMetadata{
		vdbStorage.GetValueMetadata(storage.Kicks, nil, vdbStorage.Uint256),
		vdbStorage.GetValueMetadata(storage.BidBid, keys, vdbStorage.Uint256),
		vdbStorage.GetValueMetadata(storage.BidLot, keys, vdbStorage.Uint256),
		vdbStorage.GetValueMetadataForPackedSlot(storage.Packed, keys, vdbStorage.PackedSlot, packedNames, packedTypes),
	}
}

func GetFlopMetadatas(bidId string) []vdbStorage.ValueMetadata {
	return getCommonBidMetadatas(bidId)
}

func GetFlapMetadatas(bidId string) []vdbStorage.ValueMetadata {
	return getCommonBidMetadatas(bidId)
}

func GetCdpManagerMetadatas(cdpi string) []vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.Cdpi: cdpi}
	return []vdbStorage.ValueMetadata{
		vdbStorage.GetValueMetadata(cdp_manager.Cdpi, nil, vdbStorage.Uint256),
		vdbStorage.GetValueMetadata(cdp_manager.Urns, keys, vdbStorage.Address),
		vdbStorage.GetValueMetadata(cdp_manager.Owns, keys, vdbStorage.Address),
		vdbStorage.GetValueMetadata(cdp_manager.Ilks, keys, vdbStorage.Bytes32),
	}
}

func GetFlipMetadatas(bidId string) []vdbStorage.ValueMetadata {
	keys := map[vdbStorage.Key]string{constants.BidId: bidId}
	return append(getCommonBidMetadatas(bidId),
		vdbStorage.GetValueMetadata(storage.Ilk, nil, vdbStorage.Bytes32),
		vdbStorage.GetValueMetadata(storage.BidUsr, keys, vdbStorage.Address),
		vdbStorage.GetValueMetadata(storage.BidGal, keys, vdbStorage.Address),
		vdbStorage.GetValueMetadata(storage.BidTab, keys, vdbStorage.Uint256))
}

func GetCdpManagerStorageValues(seed int, ilkHex string, urnGuy string, cdpi int) map[string]interface{} {
	valuesMap := make(map[string]interface{})
	valuesMap[cdp_manager.Cdpi] = strconv.Itoa(cdpi)
	valuesMap[cdp_manager.Urns] = urnGuy
	valuesMap[cdp_manager.Owns] = "address1" + strconv.Itoa(seed)
	valuesMap[cdp_manager.Ilks] = ilkHex
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

func insertValues(db *postgres.DB, repo vdbStorageFactory.Repository, header core.Header, valuesMap map[string]interface{}, metadatas []vdbStorage.ValueMetadata) {
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

		persistedDiff := test_helpers.CreateDiffRecord(db, header, common.Hash{}, key, valueForDiffRecord)

		err := repo.Create(persistedDiff.ID, header.Id, metadata, valueForStorageRecord)
		Expect(err).NotTo(HaveOccurred())
	}
}

func CreateFlop(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flopMetadatas []vdbStorage.ValueMetadata, contractAddress string) {
	flopRepo := flop.FlopStorageRepository{ContractAddress: contractAddress}
	flopRepo.SetDB(db)
	insertValues(db, &flopRepo, header, valuesMap, flopMetadatas)
}

func CreateFlap(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flapMetadatas []vdbStorage.ValueMetadata, contractAddress string) {
	flapRepo := flap.FlapStorageRepository{ContractAddress: contractAddress}
	flapRepo.SetDB(db)
	insertValues(db, &flapRepo, header, valuesMap, flapMetadatas)
}

func CreateFlip(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, flipMetadatas []vdbStorage.ValueMetadata, contractAddress string) {
	flipRepo := flip.FlipStorageRepository{ContractAddress: contractAddress}
	flipRepo.SetDB(db)
	insertValues(db, &flipRepo, header, valuesMap, flipMetadatas)
}

func CreateManagedCdp(db *postgres.DB, header core.Header, valuesMap map[string]interface{}, metadatas []vdbStorage.ValueMetadata) error {
	cdpManagerRepo := cdp_manager.CdpManagerStorageRepository{}
	cdpManagerRepo.SetDB(db)
	_, err := shared.GetOrCreateUrn(valuesMap[cdp_manager.Urns].(string), valuesMap[cdp_manager.Ilks].(string), db)
	if err != nil {
		return err
	}
	insertValues(db, &cdpManagerRepo, header, valuesMap, metadatas)
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
	ilkId, ilkErr := shared.GetOrCreateIlk(setupData.IlkHex, setupData.DB)
	if ilkErr != nil {
		return 0, 0, ilkErr
	}

	urnId, urnErr := shared.GetOrCreateUrn(setupData.UrnGuy, setupData.IlkHex, setupData.DB)
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
	addressID, addressErr := shared.GetOrCreateAddress(input.ContractAddress, input.DB)
	Expect(addressErr).NotTo(HaveOccurred())
	dealLog := test_data.CreateTestLog(input.DealHeaderId, input.DB)
	dealModel := test_data.CopyModel(test_data.DealModel)
	dealModel.ColumnValues[deal.Id] = strconv.Itoa(input.BidId)
	dealModel.ColumnValues[event.HeaderFK] = input.DealHeaderId
	dealModel.ColumnValues[event.LogFK] = dealLog.ID
	dealModel.ColumnValues[constants.AddressColumn] = addressID
	deals := []event.InsertionModel{dealModel}
	return event.PersistModels(deals, input.DB)
}

func CreateFlipKick(contractAddress string, bidId int, headerId, logId int64, usr string, db *postgres.DB) error {
	addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
	Expect(addressErr).NotTo(HaveOccurred())
	flipKickModel := test_data.CopyModel(test_data.FlipKickModel())
	flipKickModel.ColumnValues[event.HeaderFK] = headerId
	flipKickModel.ColumnValues[event.LogFK] = logId
	flipKickModel.ColumnValues[event.AddressFK] = addressId
	flipKickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
	flipKickModel.ColumnValues[constants.UsrColumn] = usr
	return event.PersistModels([]event.InsertionModel{flipKickModel}, db)
}

func CreateFlapKick(contractAddress string, bidId int, headerId, logId int64, db *postgres.DB) error {
	flapKickModel := test_data.FlapKickModel()
	addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
	Expect(addressErr).NotTo(HaveOccurred())
	flapKickModel.ColumnValues[event.HeaderFK] = headerId
	flapKickModel.ColumnValues[event.LogFK] = logId
	flapKickModel.ColumnValues[event.AddressFK] = addressId
	flapKickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
	return event.PersistModels([]event.InsertionModel{flapKickModel}, db)
}

func CreateFlopKick(contractAddress string, bidId int, headerId, logId int64, db *postgres.DB) error {
	addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
	Expect(addressErr).NotTo(HaveOccurred())
	flopKickModel := test_data.FlopKickModel()
	flopKickModel.ColumnValues[event.HeaderFK] = headerId
	flopKickModel.ColumnValues[event.LogFK] = logId
	flopKickModel.ColumnValues[event.AddressFK] = addressId
	flopKickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
	return event.PersistModels([]event.InsertionModel{flopKickModel}, db)
}

func CreateTend(input TendCreationInput) (err error) {
	addressID, addressErr := shared.GetOrCreateAddress(input.ContractAddress, input.DB)
	Expect(addressErr).NotTo(HaveOccurred())
	tendModel := test_data.TendModel()
	tendLog := test_data.CreateTestLog(input.TendHeaderId, input.DB)
	tendModel.ColumnValues[tend.Id] = strconv.Itoa(input.BidId)
	tendModel.ColumnValues[tend.Lot] = strconv.Itoa(input.Lot)
	tendModel.ColumnValues[tend.Bid] = strconv.Itoa(input.BidAmount)
	tendModel.ColumnValues[constants.AddressColumn] = addressID
	tendModel.ColumnValues[event.HeaderFK] = input.TendHeaderId
	tendModel.ColumnValues[event.LogFK] = tendLog.ID
	return event.PersistModels([]event.InsertionModel{tendModel}, input.DB)
}

func CreateDent(input DentCreationInput) (err error) {
	addressID, addressErr := shared.GetOrCreateAddress(input.ContractAddress, input.DB)
	Expect(addressErr).NotTo(HaveOccurred())
	dentModel := test_data.DentModel()
	dentModel.ColumnValues[dent.Id] = strconv.Itoa(input.BidId)
	dentModel.ColumnValues[dent.Lot] = strconv.Itoa(input.Lot)
	dentModel.ColumnValues[dent.Bid] = strconv.Itoa(input.BidAmount)
	dentModel.ColumnValues[constants.AddressColumn] = addressID
	dentModel.ColumnValues[event.HeaderFK] = input.DentHeaderId
	dentModel.ColumnValues[event.LogFK] = input.DentLogId
	return event.PersistModels([]event.InsertionModel{dentModel}, input.DB)
}

func CreateYank(input YankCreationInput) (err error) {
	addressID, addressErr := shared.GetOrCreateAddress(input.ContractAddress, input.DB)
	Expect(addressErr).NotTo(HaveOccurred())
	yankModel := test_data.CopyModel(test_data.YankModel)
	yankModel.ColumnValues[yank.BidId] = strconv.Itoa(input.BidId)
	yankModel.ColumnValues[constants.AddressColumn] = addressID
	yankModel.ColumnValues[event.HeaderFK] = input.YankHeaderId
	yankModel.ColumnValues[event.LogFK] = input.YankLogId
	return event.PersistModels([]event.InsertionModel{yankModel}, input.DB)
}

func CreateTick(input TickCreationInput) (err error) {
	addressID, addressErr := shared.GetOrCreateAddress(input.ContractAddress, input.DB)
	Expect(addressErr).NotTo(HaveOccurred())
	tickLog := test_data.CreateTestLog(input.TickHeaderId, input.DB)
	tickModel := test_data.CopyModel(test_data.TickModel)
	tickModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(input.BidId)
	tickModel.ColumnValues[constants.AddressColumn] = addressID
	tickModel.ColumnValues[event.HeaderFK] = input.TickHeaderId
	tickModel.ColumnValues[event.LogFK] = tickLog.ID
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
