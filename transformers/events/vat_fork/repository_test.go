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

package vat_fork_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fork"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Vat fork repository", func() {
	var (
		db                *postgres.DB
		vatForkRepository vat_fork.VatForkRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatForkRepository = vat_fork.VatForkRepository{}
		vatForkRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatForkModelWithNegativeDinkDart
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatForkChecked,
			LogEventTableName:        "maker.vat_fork",
			TestModel:                test_data.VatForkModelWithNegativeDinkDart,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &vatForkRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat fork", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = vatForkRepository.Create(headerID, []interface{}{test_data.VatForkModelWithNegativeDinkDart})
			Expect(err).NotTo(HaveOccurred())

			var dbVatFork vat_fork.VatForkModel
			err = db.Get(&dbVatFork, `SELECT ilk_id, src, dst, dink, dart, log_idx, tx_idx, raw_log FROM maker.vat_fork WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			ilkID, err := shared.GetOrCreateIlk(test_data.VatForkModelWithNegativeDinkDart.Ilk, db)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbVatFork.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbVatFork.Src).To(Equal(test_data.VatForkModelWithNegativeDinkDart.Src))
			Expect(dbVatFork.Dst).To(Equal(test_data.VatForkModelWithNegativeDinkDart.Dst))
			Expect(dbVatFork.Dink).To(Equal(test_data.VatForkModelWithNegativeDinkDart.Dink))
			Expect(dbVatFork.Dart).To(Equal(test_data.VatForkModelWithNegativeDinkDart.Dart))
			Expect(dbVatFork.LogIndex).To(Equal(test_data.VatForkModelWithNegativeDinkDart.LogIndex))
			Expect(dbVatFork.TransactionIndex).To(Equal(test_data.VatForkModelWithNegativeDinkDart.TransactionIndex))
			Expect(dbVatFork.Raw).To(MatchJSON(test_data.VatForkModelWithNegativeDinkDart.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatForkChecked,
			Repository:              &vatForkRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
