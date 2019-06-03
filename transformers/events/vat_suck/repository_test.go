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

package vat_suck_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_suck"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("VatSuck Repository", func() {
	var (
		db         *postgres.DB
		repository vat_suck.VatSuckRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_suck.VatSuckRepository{}
		repository.SetDB(db)
	})

	type VatSuckDBResult struct {
		vat_suck.VatSuckModel
		Id       int
		HeaderId int64 `db:"header_id"`
	}

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatSuckModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatSuckChecked,
			LogEventTableName:        "maker.vat_suck",
			TestModel:                test_data.VatSuckModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("persists vat suck records", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerId, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			anotherVatSuck := test_data.VatSuckModel
			anotherVatSuck.LogIndex = test_data.VatSuckModel.LogIndex + 1
			err = repository.Create(headerId, []interface{}{test_data.VatSuckModel, anotherVatSuck})

			var dbResult []VatSuckDBResult
			err = db.Select(&dbResult, `SELECT * from maker.vat_suck where header_id = $1`, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbResult)).To(Equal(2))
			Expect(dbResult[0].U).To(Equal(test_data.VatSuckModel.U))
			Expect(dbResult[0].V).To(Equal(test_data.VatSuckModel.V))
			Expect(dbResult[0].Rad).To(Equal(test_data.VatSuckModel.Rad))
			Expect(dbResult[0].LogIndex).To(Equal(test_data.VatSuckModel.LogIndex))
			Expect(dbResult[1].LogIndex).To(Equal(test_data.VatSuckModel.LogIndex + 1))
			Expect(dbResult[0].TransactionIndex).To(Equal(test_data.VatSuckModel.TransactionIndex))
			Expect(dbResult[0].Raw).To(MatchJSON(test_data.VatSuckModel.Raw))
			Expect(dbResult[0].HeaderId).To(Equal(headerId))
		})
	})

	Describe("MarkCheckedHeader", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatSuckChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
