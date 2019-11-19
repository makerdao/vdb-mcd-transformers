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

package queries

import (
	"math/big"
	"math/rand"
	"strconv"

	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"

	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/bite"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
)

var _ = Describe("Urn state computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		headerId, logId  int64
		vatRepository    vat.VatStorageRepository
		catRepository    cat.CatStorageRepository
		jugRepository    jug.JugStorageRepository
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		fakeHeaderSyncLog := test_data.CreateTestLog(headerId, db)
		logId = fakeHeaderSyncLog.ID

		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("urn_state_ilk", func() {
		It("returns the ilk for an urn", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			fakeGuy := "fakeAddress"
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			ilkRate, convertRateErr := strconv.Atoi(ilkValues[vat.IlkRate])
			Expect(convertRateErr).NotTo(HaveOccurred())
			urnSetupData.Rate = ilkRate
			ilkSpot, convertSpotErr := strconv.Atoi(ilkValues[vat.IlkSpot])
			Expect(convertSpotErr).NotTo(HaveOccurred())
			urnSetupData.Spot = ilkSpot
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
					FROM api.urn_state_ilk(
					(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
					FROM api.get_urn($1, $2, $3)))`, test_helpers.FakeIlk.Identifier, fakeGuy, fakeHeader.BlockNumber)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("urn_state_frobs", func() {
		It("returns frobs for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModelWithPositiveDart()
			frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
			frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobEvent.ColumnValues[constants.HeaderFK] = headerId
			frobEvent.ColumnValues[constants.LogFK] = logId
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs test_helpers.FrobEvent
			getFrobsErr := db.Get(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
                        (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
                         FROM api.all_urns($1))
                    )`, fakeBlock)
			Expect(getFrobsErr).NotTo(HaveOccurred())

			expectedFrobs := test_helpers.FrobEvent{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Dink:          frobEvent.ColumnValues["dink"].(string),
				Dart:          frobEvent.ColumnValues["dart"].(string),
			}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})

		Describe("result pagination", func() {
			var frobEventOne, frobEventTwo shared.InsertionModel

			BeforeEach(func() {
				urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
				urnSetupData.Header.Hash = fakeHeader.Hash
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

				frobRepo := vat_frob.VatFrobRepository{}
				frobRepo.SetDB(db)

				frobEventOne = test_data.VatFrobModelWithPositiveDart()
				frobEventOne.ForeignKeyValues[constants.UrnFK] = fakeGuy
				frobEventOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				frobEventOne.ColumnValues[constants.HeaderFK] = headerId
				frobEventOne.ColumnValues[constants.LogFK] = logId
				insertFrobErrOne := frobRepo.Create([]shared.InsertionModel{frobEventOne})
				Expect(insertFrobErrOne).NotTo(HaveOccurred())

				// insert more recent frob for same urn
				laterBlock := fakeBlock + 1
				fakeHeaderTwo := fakes.GetFakeHeader(int64(laterBlock))
				headerTwoId, insertHeaderTwoErr := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
				Expect(insertHeaderTwoErr).NotTo(HaveOccurred())
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID

				frobEventTwo = test_data.VatFrobModelWithNegativeDink()
				frobEventTwo.ForeignKeyValues[constants.UrnFK] = fakeGuy
				frobEventTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				frobEventTwo.ColumnValues[constants.HeaderFK] = headerTwoId
				frobEventTwo.ColumnValues[constants.LogFK] = logTwoId
				insertFrobErrTwo := frobRepo.Create([]shared.InsertionModel{frobEventTwo})
				Expect(insertFrobErrTwo).NotTo(HaveOccurred())
			})

			It("limits results to latest block number if max_results argument is provided", func() {
				maxResults := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3)`, test_helpers.FakeIlk.Identifier, fakeGuy, maxResults)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				expectedFrob := test_helpers.FrobEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Dink:          frobEventTwo.ColumnValues["dink"].(string),
					Dart:          frobEventTwo.ColumnValues["dart"].(string),
				}
				Expect(actualFrobs).To(ConsistOf(expectedFrob))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualFrobs []test_helpers.FrobEvent
				getFrobsErr := db.Select(&actualFrobs,
					`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, fakeGuy, maxResults, resultOffset)
				Expect(getFrobsErr).NotTo(HaveOccurred())

				expectedFrobs := test_helpers.FrobEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Dink:          frobEventOne.ColumnValues["dink"].(string),
					Dart:          frobEventOne.ColumnValues["dart"].(string),
				}
				Expect(actualFrobs).To(ConsistOf(expectedFrobs))
			})
		})
	})

	Describe("urn_state_bites", func() {
		It("returns bites for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			biteRepo := bite.Repository{}
			biteRepo.SetDB(db)
			biteEvent := generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerId, logId, db)
			insertBiteErr := biteRepo.Create([]event.InsertionModel{biteEvent})
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites test_helpers.BiteEvent
			getBitesErr := db.Get(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
				    (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
				    FROM api.all_urns($1)))`,
				fakeBlock)
			Expect(getBitesErr).NotTo(HaveOccurred())

			expectedBites := test_helpers.BiteEvent{
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
				UrnIdentifier: fakeGuy,
				Ink:           biteEvent.ColumnValues["ink"].(string),
				Art:           biteEvent.ColumnValues["art"].(string),
				Tab:           biteEvent.ColumnValues["tab"].(string),
			}

			Expect(actualBites).To(Equal(expectedBites))
		})

		Describe("result pagination", func() {
			var biteEventOne, biteEventTwo event.InsertionModel

			BeforeEach(func() {
				urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
				urnSetupData.Header.Hash = fakeHeader.Hash
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

				biteRepo := bite.Repository{}
				biteRepo.SetDB(db)

				biteEventOne = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerId, logId, db)
				insertBiteOneErr := biteRepo.Create([]event.InsertionModel{biteEventOne})
				Expect(insertBiteOneErr).NotTo(HaveOccurred())

				// insert more recent bite for same urn
				laterBlock := fakeBlock + 1
				fakeHeaderTwo := fakes.GetFakeHeader(int64(laterBlock))
				headerTwoId, insertHeaderTwoErr := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
				Expect(insertHeaderTwoErr).NotTo(HaveOccurred())
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID

				biteEventTwo = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerTwoId, logTwoId, db)
				insertBiteTwoErr := biteRepo.Create([]event.InsertionModel{biteEventTwo})
				Expect(insertBiteTwoErr).NotTo(HaveOccurred())
			})

			It("limits results to latest block number if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
					SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3)`,
					test_helpers.FakeIlk.Identifier, fakeGuy, maxResults)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Ink:           biteEventTwo.ColumnValues["ink"].(string),
					Art:           biteEventTwo.ColumnValues["art"].(string),
					Tab:           biteEventTwo.ColumnValues["tab"].(string),
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				getBitesErr := db.Select(&actualBites, `
					SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
						(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
						 FROM api.get_urn($1, $2)), $3, $4)`,
					test_helpers.FakeIlk.Identifier, fakeGuy, maxResults, resultOffset)
				Expect(getBitesErr).NotTo(HaveOccurred())

				expectedBite := test_helpers.BiteEvent{
					IlkIdentifier: test_helpers.FakeIlk.Identifier,
					UrnIdentifier: fakeGuy,
					Ink:           biteEventOne.ColumnValues["ink"].(string),
					Art:           biteEventOne.ColumnValues["art"].(string),
					Tab:           biteEventOne.ColumnValues["tab"].(string),
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})
		})
	})
})

func randomizeBite(bite event.InsertionModel) event.InsertionModel {
	bite.ColumnValues["ink"] = big.NewInt(rand.Int63()).String()
	bite.ColumnValues["art"] = big.NewInt(rand.Int63()).String()
	bite.ColumnValues["tab"] = big.NewInt(rand.Int63()).String()
	return bite
}
