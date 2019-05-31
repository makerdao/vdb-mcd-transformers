package vat_fork_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_fork"
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
		modelWithDifferentLogIdx := test_data.VatForkModelWithNegativeDart
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatForkChecked,
			LogEventTableName:        "maker.vat_fork",
			TestModel:                test_data.VatForkModelWithNegativeDart,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &vatForkRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat fork", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = vatForkRepository.Create(headerID, []interface{}{test_data.VatForkModelWithNegativeDart})
			Expect(err).NotTo(HaveOccurred())

			var dbVatFork vat_fork.VatForkModel
			err = db.Get(&dbVatFork, `SELECT ilk_id, src, dst, dink, dart, log_idx, tx_idx, raw_log FROM maker.vat_fork WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbVatFork.Src).To(Equal(test_data.VatForkModelWithNegativeDart.Src))
			Expect(dbVatFork.Dst).To(Equal(test_data.VatForkModelWithNegativeDart.Dst))
			Expect(dbVatFork.Dink).To(Equal(test_data.VatForkModelWithNegativeDart.Dink))
			Expect(dbVatFork.Dart).To(Equal(test_data.VatForkModelWithNegativeDart.Dart))
			Expect(dbVatFork.LogIndex).To(Equal(test_data.VatForkModelWithNegativeDart.LogIndex))
			Expect(dbVatFork.TransactionIndex).To(Equal(test_data.VatForkModelWithNegativeDart.TransactionIndex))
			Expect(dbVatFork.Raw).To(MatchJSON(test_data.VatForkModelWithNegativeDart.Raw))
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
