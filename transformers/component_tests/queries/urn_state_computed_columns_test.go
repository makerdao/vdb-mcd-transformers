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
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Urn state computed columns", func() {
	var (
		fakeGuy                = fakes.RandomString(42)
		blockOne, timestampOne int
		headerOne              core.Header
		logId                  int64
		vatRepository          vat.VatStorageRepository
		catRepository          cat.CatStorageRepository
		jugRepository          jug.JugStorageRepository
		headerRepository       repositories.HeaderRepository
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		fakeHeaderSyncLog := test_data.CreateTestLog(headerOne.Id, db)
		logId = fakeHeaderSyncLog.ID

		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	Describe("urn_state_ilk", func() {
		It("returns the ilk for an urn", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, diffID, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			fakeGuy := "fakeAddress"
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepository)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
					FROM api.urn_state_ilk(
					(SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
					FROM api.get_urn($1, $2, $3)))`, test_helpers.FakeIlk.Identifier, fakeGuy, headerOne.BlockNumber)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("urn_state_frobs", func() {
		It("returns frobs for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepository)

			frobRepo := vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent := test_data.VatFrobModelWithPositiveDart()
			frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
			frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			frobEvent.ColumnValues[constants.LogFK] = logId
			insertFrobErr := frobRepo.Create([]shared.InsertionModel{frobEvent})
			Expect(insertFrobErr).NotTo(HaveOccurred())

			var actualFrobs test_helpers.FrobEvent
			getFrobsErr := db.Get(&actualFrobs,
				`SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_state_frobs(
                        (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
                         FROM api.all_urns($1))
                    )`, blockOne)
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
				urnSetupData := test_helpers.GetUrnSetupData()
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepository)

				frobRepo := vat_frob.VatFrobRepository{}
				frobRepo.SetDB(db)

				frobEventOne = test_data.VatFrobModelWithPositiveDart()
				frobEventOne.ForeignKeyValues[constants.UrnFK] = fakeGuy
				frobEventOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				frobEventOne.ColumnValues[constants.HeaderFK] = headerOne.Id
				frobEventOne.ColumnValues[constants.LogFK] = logId
				insertFrobErrOne := frobRepo.Create([]shared.InsertionModel{frobEventOne})
				Expect(insertFrobErrOne).NotTo(HaveOccurred())

				// insert more recent frob for same urn
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepository)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				frobEventTwo = test_data.VatFrobModelWithNegativeDink()
				frobEventTwo.ForeignKeyValues[constants.UrnFK] = fakeGuy
				frobEventTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
				frobEventTwo.ColumnValues[constants.HeaderFK] = headerTwo.Id
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
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepository)

			biteEvent := generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerOne.Id, logId, db)
			insertBiteErr := event.PersistModels([]event.InsertionModel{biteEvent}, db)
			Expect(insertBiteErr).NotTo(HaveOccurred())

			var actualBites test_helpers.BiteEvent
			getBitesErr := db.Get(&actualBites, `
				SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_state_bites(
				    (SELECT (urn_identifier, ilk_identifier, block_height, ink, art, created, updated)::api.urn_state
				    FROM api.all_urns($1)))`,
				blockOne)
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
				urnSetupData := test_helpers.GetUrnSetupData()
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(urnSetupData, diffID, headerOne.Id, urnMetadata, vatRepository)

				biteEventOne = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerOne.Id, logId, db)
				insertBiteOneErr := event.PersistModels([]event.InsertionModel{biteEventOne}, db)
				Expect(insertBiteOneErr).NotTo(HaveOccurred())

				// insert more recent bite for same urn
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepository)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				biteEventTwo = generateBite(test_helpers.FakeIlk.Hex, fakeGuy, headerTwo.Id, logTwoId, db)
				insertBiteTwoErr := event.PersistModels([]event.InsertionModel{biteEventTwo}, db)
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
