package pot_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/pot"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pot storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 pot.PotStorageRepository
		fakeAddress          = "0x" + fakes.RandomString(20)
		fakeUint256          = strconv.Itoa(rand.Int())
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = pot.PotStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		diffID = CreateFakeDiffRecord(db)
	})

	Describe("User pie", func() {
		It("writes a row", func() {
			userPieMetadata := storage.GetValueMetadata(pot.UserPie, map[storage.Key]string{constants.MsgSender: fakeAddress}, storage.Uint256)

			insertErr := repo.Create(diffID, fakeHeaderID, userPieMetadata, fakeUint256)
			Expect(insertErr).NotTo(HaveOccurred())

			var result MappingRes
			dbErr := db.Get(&result, `SELECT diff_id, header_id, "user" AS key, pie AS value FROM maker.pot_user_pie`)
			Expect(dbErr).NotTo(HaveOccurred())
			addressID, addressErr := shared.GetOrCreateAddress(fakeAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			userPieMetadata := storage.GetValueMetadata(pot.UserPie, map[storage.Key]string{constants.MsgSender: fakeAddress}, storage.Uint256)
			setupErr := repo.Create(diffID, fakeHeaderID, userPieMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			insertErr := repo.Create(diffID, fakeHeaderID, userPieMetadata, fakeUint256)
			Expect(insertErr).NotTo(HaveOccurred())

			Expect(insertErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT COUNT(*) FROM maker.pot_user_pie`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing ilk", func() {
			malformedUserPieMetadata := storage.GetValueMetadata(pot.UserPie, nil, storage.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedUserPieMetadata, fakeUint256)
			Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.MsgSender}))
		})
	})

	Describe("Pie", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   pot.Pie,
			Value:            fakeUint256,
			StorageTableName: "maker.pot_pie",
			Repository:       &repo,
			Metadata:         pot.PieMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("dsr", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   pot.Dsr,
			Value:            fakeUint256,
			StorageTableName: "maker.pot_dsr",
			Repository:       &repo,
			Metadata:         pot.DsrMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("chi", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   pot.Chi,
			Value:            fakeUint256,
			StorageTableName: "maker.pot_chi",
			Repository:       &repo,
			Metadata:         pot.ChiMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("vat", func() {
		It("persists a record", func() {
			err := repo.Create(diffID, fakeHeaderID, pot.VatMetadata, fakeAddress)
			Expect(err).NotTo(HaveOccurred())

			addressID, addressErr := shared.GetOrCreateAddress(fakeAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			var result VariableRes
			err = db.Get(&result, `SELECT diff_id, header_id, vat AS value FROM maker.pot_vat`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
		})

		It("doesn't duplicate a record", func() {
			setupErr := repo.Create(diffID, fakeHeaderID, pot.VatMetadata, fakeAddress)
			Expect(setupErr).NotTo(HaveOccurred())

			createErr := repo.Create(diffID, fakeHeaderID, pot.VatMetadata, fakeAddress)
			Expect(createErr).NotTo(HaveOccurred())

			var count int
			query := `SELECT COUNT(*) FROM maker.pot_vat`
			queryErr := db.Get(&count, query)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("vow", func() {
		It("persists a record", func() {
			err := repo.Create(diffID, fakeHeaderID, pot.VowMetadata, fakeAddress)
			Expect(err).NotTo(HaveOccurred())

			addressID, addressErr := shared.GetOrCreateAddress(fakeAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			var result VariableRes
			err = db.Get(&result, `SELECT diff_id, header_id, vow AS value FROM maker.pot_vow`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
		})

		It("doesn't duplicate a record", func() {
			setupErr := repo.Create(diffID, fakeHeaderID, pot.VowMetadata, fakeAddress)
			Expect(setupErr).NotTo(HaveOccurred())

			createErr := repo.Create(diffID, fakeHeaderID, pot.VowMetadata, fakeAddress)
			Expect(createErr).NotTo(HaveOccurred())

			var count int
			query := `SELECT COUNT(*) FROM maker.pot_vow`
			queryErr := db.Get(&count, query)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("rho", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   pot.Rho,
			Value:            fakeUint256,
			StorageTableName: "maker.pot_rho",
			Repository:       &repo,
			Metadata:         pot.RhoMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("live", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   pot.Live,
			Value:            fakeUint256,
			StorageTableName: "maker.pot_live",
			Repository:       &repo,
			Metadata:         pot.LiveMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})
})
