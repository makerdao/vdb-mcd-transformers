// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package vat_frob_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Vat frob repository", func() {
	var (
		db                *postgres.DB
		vatFrobRepository vat_frob.VatFrobRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatFrobRepository = vat_frob.VatFrobRepository{}
		vatFrobRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatFrobModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatFrobChecked,
			LogEventTableName:        "maker.vat_frob",
			TestModel:                test_data.VatFrobModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &vatFrobRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat frob", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = vatFrobRepository.Create(headerID, []interface{}{test_data.VatFrobModel})
			Expect(err).NotTo(HaveOccurred())
			var dbVatFrob vat_frob.VatFrobModel
			err = db.Get(&dbVatFrob, `SELECT urn, v, w, dink, dart, log_idx, tx_idx, raw_log FROM maker.vat_frob WHERE header_id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.VatFrobModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			urnID, err := shared.GetOrCreateUrn(test_data.VatFrobModel.Urn, ilkID, db)
			Expect(dbVatFrob.Urn).To(Equal(strconv.Itoa(urnID)))
			Expect(dbVatFrob.V).To(Equal(test_data.VatFrobModel.V))
			Expect(dbVatFrob.W).To(Equal(test_data.VatFrobModel.W))
			Expect(dbVatFrob.Dink).To(Equal(test_data.VatFrobModel.Dink))
			Expect(dbVatFrob.Dart).To(Equal(test_data.VatFrobModel.Dart))
			Expect(dbVatFrob.LogIndex).To(Equal(test_data.VatFrobModel.LogIndex))
			Expect(dbVatFrob.TransactionIndex).To(Equal(test_data.VatFrobModel.TransactionIndex))
			Expect(dbVatFrob.Raw).To(MatchJSON(test_data.VatFrobModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatFrobChecked,
			Repository:              &vatFrobRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
