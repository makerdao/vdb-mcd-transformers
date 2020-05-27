package pot_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/pot"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pot storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 pot.StorageRepository
		fakeAddress          = "0x" + fakes.RandomString(40)
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = pot.StorageRepository{ContractAddress: test_data.PotAddress()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Variable", func() {
		It("panics if the metadata name is not recognized", func() {
			unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}
			repoCreate := func() {
				repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
			}

			Expect(repoCreate).Should(Panic())
		})
	})

	Describe("User pie", func() {
		It("writes a row", func() {
			userPieMetadata := types.GetValueMetadata(pot.UserPie, map[types.Key]string{constants.MsgSender: fakeAddress}, types.Uint256)

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
			userPieMetadata := types.GetValueMetadata(pot.UserPie, map[types.Key]string{constants.MsgSender: fakeAddress}, types.Uint256)
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
			malformedUserPieMetadata := types.GetValueMetadata(pot.UserPie, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedUserPieMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.MsgSender}))
		})
	})

	Describe("Pie", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: pot.Pie,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.PotPieTable,
			Repository:     &repo,
			Metadata:       pot.PieMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("dsr", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: pot.Dsr,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.PotDsrTable,
			Repository:     &repo,
			Metadata:       pot.DsrMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("chi", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: pot.Chi,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.PotChiTable,
			Repository:     &repo,
			Metadata:       pot.ChiMetadata,
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
			ValueFieldName: pot.Rho,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.PotRhoTable,
			Repository:     &repo,
			Metadata:       pot.RhoMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("live", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: pot.Live,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.PotLiveTable,
			Repository:     &repo,
			Metadata:       pot.LiveMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})
	Describe("Wards mapping", func() {
		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result WardsMappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			err := db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := shared.GetOrCreateAddress(fakeUserAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())
			Expect(result.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))
			AssertMapping(result.MappingRes, diffID, fakeHeaderID, strconv.FormatInt(userAddressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing user", func() {
			malformedWardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedWardsMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.User}))
		})
	})
})
