package vat_grab_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_grab"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Vat grab repository", func() {
	var (
		db                *postgres.DB
		vatGrabRepository vat_grab.VatGrabRepository
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatGrabRepository = vat_grab.VatGrabRepository{}
		vatGrabRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatGrabModelWithPositiveDink
		modelWithDifferentLogIdx.LogIndex++

		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatGrabChecked,
			LogEventTableName:        "maker.vat_grab",
			TestModel:                test_data.VatGrabModelWithPositiveDink,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &vatGrabRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat grab event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = vatGrabRepository.Create(headerID, []interface{}{test_data.VatGrabModelWithPositiveDink})
			Expect(err).NotTo(HaveOccurred())
			var dbVatGrab vat_grab.VatGrabModel
			err = db.Get(&dbVatGrab, `SELECT urn_id, v, w, dink, dart, log_idx, tx_idx, raw_log FROM maker.vat_grab WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.VatGrabModelWithPositiveDink.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			urnID, err := shared.GetOrCreateUrn(test_data.VatGrabModelWithPositiveDink.Urn, ilkID, db)
			Expect(dbVatGrab.Urn).To(Equal(strconv.Itoa(urnID)))
			Expect(dbVatGrab.V).To(Equal(test_data.VatGrabModelWithPositiveDink.V))
			Expect(dbVatGrab.W).To(Equal(test_data.VatGrabModelWithPositiveDink.W))
			Expect(dbVatGrab.Dink).To(Equal(test_data.VatGrabModelWithPositiveDink.Dink))
			Expect(dbVatGrab.Dart).To(Equal(test_data.VatGrabModelWithPositiveDink.Dart))
			Expect(dbVatGrab.LogIndex).To(Equal(test_data.VatGrabModelWithPositiveDink.LogIndex))
			Expect(dbVatGrab.TransactionIndex).To(Equal(test_data.VatGrabModelWithPositiveDink.TransactionIndex))
			Expect(dbVatGrab.Raw).To(MatchJSON(test_data.VatGrabModelWithPositiveDink.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatGrabChecked,
			Repository:              &vatGrabRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
