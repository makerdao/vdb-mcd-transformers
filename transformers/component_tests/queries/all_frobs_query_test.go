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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Frobs query", func() {
	var (
		db         *postgres.DB
		frobRepo   vat_frob.VatFrobRepository
		headerRepo repositories.HeaderRepository
		fakeUrn    = test_data.RandomString(5)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		frobRepo = vat_frob.VatFrobRepository{}
		frobRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("urn_frobs", func() {
		It("returns frobs for relevant ilk/urn", func() {
			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)

			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			frobBlockOne := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobBlockOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobBlockOne.ForeignKeyValues[constants.UrnFK] = fakeUrn
			frobBlockOne.ColumnValues["dink"] = strconv.Itoa(rand.Int())
			frobBlockOne.ColumnValues["dart"] = strconv.Itoa(rand.Int())

			irrelevantFrob := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			irrelevantFrob.ForeignKeyValues[constants.IlkFK] = test_helpers.AnotherFakeIlk.Hex
			irrelevantFrob.ForeignKeyValues[constants.UrnFK] = fakeUrn
			irrelevantFrob.ColumnValues["dink"] = strconv.Itoa(rand.Int())
			irrelevantFrob.ColumnValues["dart"] = strconv.Itoa(rand.Int())
			irrelevantFrob.ColumnValues["tx_idx"] = frobBlockOne.ColumnValues["tx_idx"].(uint) + 1

			err = frobRepo.Create(headerOneId, []shared.InsertionModel{frobBlockOne, irrelevantFrob})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())

			frobBlockTwo := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobBlockTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobBlockTwo.ForeignKeyValues[constants.UrnFK] = fakeUrn
			frobBlockTwo.ColumnValues["dink"] = strconv.Itoa(rand.Int())
			frobBlockTwo.ColumnValues["dart"] = strconv.Itoa(rand.Int())

			err = frobRepo.Create(headerTwoId, []shared.InsertionModel{frobBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			err = db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_frobs($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn,
					Dink: frobBlockOne.ColumnValues["dink"].(string), Dart: frobBlockOne.ColumnValues["dart"].(string)},
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn,
					Dink: frobBlockTwo.ColumnValues["dink"].(string), Dart: frobBlockTwo.ColumnValues["dart"].(string)},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.urn_frobs()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.urn_frobs() does not exist"))
		})

		It("fails if only one argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.urn_frobs($1::text)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.urn_frobs(text) does not exist"))
		})
	})

	Describe("all_frobs", func() {
		It("returns all frobs for a whole ilk", func() {
			headerOne := fakes.GetFakeHeader(1)

			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			frobOne := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobOne.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobOne.ForeignKeyValues[constants.UrnFK] = fakeUrn
			frobOne.ColumnValues["dink"] = strconv.Itoa(rand.Int())
			frobOne.ColumnValues["dart"] = strconv.Itoa(rand.Int())

			anotherUrn := "anotherUrn"
			frobTwo := test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobTwo.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
			frobTwo.ForeignKeyValues[constants.UrnFK] = anotherUrn
			frobTwo.ColumnValues["dink"] = strconv.Itoa(rand.Int())
			frobTwo.ColumnValues["dart"] = strconv.Itoa(rand.Int())
			frobTwo.ColumnValues["tx_idx"] = frobOne.ColumnValues["tx_idx"].(uint) + 1

			err = frobRepo.Create(headerOneId, []shared.InsertionModel{frobOne, frobTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			err = db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart FROM api.all_frobs($1)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn,
					Dink: frobOne.ColumnValues["dink"].(string), Dart: frobOne.ColumnValues["dart"].(string)},
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: anotherUrn,
					Dink: frobTwo.ColumnValues["dink"].(string), Dart: frobTwo.ColumnValues["dart"].(string)},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.all_frobs()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.all_frobs() does not exist"))
		})
	})
})
