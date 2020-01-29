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
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
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
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		fakeEventLog := test_data.CreateTestLog(headerOne.Id, db)
		logId = fakeEventLog.ID

		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
	})

	Describe("urn_state_ilk", func() {
		It("returns the ilk for an urn", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			fakeGuy := "fakeAddress"
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

			expectedIlk := test_helpers.IlkSnapshotFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

			var result test_helpers.IlkSnapshot
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
			test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

			frobEvent := test_data.VatFrobModelWithPositiveDart()
			urnID, urnErr := shared.GetOrCreateUrn(fakeGuy, test_helpers.FakeIlk.Hex, db)
			Expect(urnErr).NotTo(HaveOccurred())
			frobEvent.ColumnValues[event.HeaderFK] = headerOne.Id
			frobEvent.ColumnValues[constants.UrnColumn] = urnID
			frobEvent.ColumnValues[event.LogFK] = logId
			insertFrobErr := event.PersistModels([]event.InsertionModel{frobEvent}, db)
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
				Dink:          frobEvent.ColumnValues[constants.DinkColumn].(string),
				Dart:          frobEvent.ColumnValues[constants.DartColumn].(string),
			}

			Expect(actualFrobs).To(Equal(expectedFrobs))
		})

		Describe("result pagination", func() {
			var frobEventOne, frobEventTwo event.InsertionModel

			BeforeEach(func() {
				urnSetupData := test_helpers.GetUrnSetupData()
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

				frobEventOne = test_data.VatFrobModelWithPositiveDart()
				urnID, urnErr := shared.GetOrCreateUrn(fakeGuy, test_helpers.FakeIlk.Hex, db)
				Expect(urnErr).NotTo(HaveOccurred())
				frobEventOne.ColumnValues[constants.UrnColumn] = urnID
				frobEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
				frobEventOne.ColumnValues[event.LogFK] = logId
				insertFrobErrOne := event.PersistModels([]event.InsertionModel{frobEventOne}, db)
				Expect(insertFrobErrOne).NotTo(HaveOccurred())

				// insert more recent frob for same urn
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepository)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				frobEventTwo = test_data.VatFrobModelWithNegativeDink()
				frobEventTwo.ColumnValues[constants.UrnColumn] = urnID
				frobEventTwo.ColumnValues[event.HeaderFK] = headerTwo.Id
				frobEventTwo.ColumnValues[event.LogFK] = logTwoId
				insertFrobErrTwo := event.PersistModels([]event.InsertionModel{frobEventTwo}, db)
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
					Dink:          frobEventTwo.ColumnValues[constants.DinkColumn].(string),
					Dart:          frobEventTwo.ColumnValues[constants.DartColumn].(string),
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
					Dink:          frobEventOne.ColumnValues[constants.DinkColumn].(string),
					Dart:          frobEventOne.ColumnValues[constants.DartColumn].(string),
				}
				Expect(actualFrobs).To(ConsistOf(expectedFrobs))
			})
		})
	})

	Describe("urn_state_bites", func() {
		It("returns bites for an urn_state", func() {
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

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
				Ink:           biteEvent.ColumnValues[constants.InkColumn].(string),
				Art:           biteEvent.ColumnValues[constants.ArtColumn].(string),
				Tab:           biteEvent.ColumnValues[constants.TabColumn].(string),
			}

			Expect(actualBites).To(Equal(expectedBites))
		})

		Describe("result pagination", func() {
			var biteEventOne, biteEventTwo event.InsertionModel

			BeforeEach(func() {
				urnSetupData := test_helpers.GetUrnSetupData()
				urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
				test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

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
					Ink:           biteEventTwo.ColumnValues[constants.InkColumn].(string),
					Art:           biteEventTwo.ColumnValues[constants.ArtColumn].(string),
					Tab:           biteEventTwo.ColumnValues[constants.TabColumn].(string),
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
					Ink:           biteEventOne.ColumnValues[constants.InkColumn].(string),
					Art:           biteEventOne.ColumnValues[constants.ArtColumn].(string),
					Tab:           biteEventOne.ColumnValues[constants.TabColumn].(string),
				}
				Expect(actualBites).To(ConsistOf(expectedBite))
			})
		})
	})
})
