package cdp_manager_test

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("CDP Manager storage repository", func() {
	var (
		db              *postgres.DB
		repository      cdp_manager.CdpManagerStorageRepository
		fakeBlockNumber int
		fakeHash        string
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = cdp_manager.CdpManagerStorageRepository{}
		repository.SetDB(db)
		fakeBlockNumber = rand.Int()
		fakeHash = fakes.FakeHash.Hex()
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := utils.StorageValueMetadata{Name: "unrecognized"}
		repoCreate := func() {
			repository.Create(fakeBlockNumber, fakeHash, unrecognizedMetadata, "")
		}

		Expect(repoCreate).Should(Panic())
	})

	Describe("vat", func() {
		var vatMetadata = utils.StorageValueMetadata{Name: cdp_manager.CdpManagerVat}
		var fakeAddress = FakeAddress

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   cdp_manager.CdpManagerVat,
			Value:            fakeAddress,
			StorageTableName: "maker.cdp_manager_vat",
			Repository:       &repository,
			Metadata:         vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("cdpi", func() {
		var (
			cdpiMetadata  = utils.StorageValueMetadata{Name: cdp_manager.CdpManagerCdpi}
			fakeCdpi      = strconv.Itoa(rand.Int())
			fakeTimestamp int
			header        core.Header
		)

		BeforeEach(func() {
			fakeTimestamp = int(rand.Int31())
			header = fakes.GetFakeHeaderWithTimestamp(int64(fakeTimestamp), int64(fakeBlockNumber))
			headerRepo := repositories.NewHeaderRepository(db)
			_, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())
		})

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   cdp_manager.CdpManagerCdpi,
			Value:            fakeCdpi,
			StorageTableName: "maker.cdp_manager_cdpi",
			Repository:       &repository,
			Metadata:         cdpiMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

		It("triggers an update to the managed_cdp table", func() {
			createdTimestamp := time.Unix(int64(fakeTimestamp), 0).UTC().Format(time.RFC3339)
			expectedTimeCreated := sql.NullString{String: createdTimestamp, Valid: true}
			err := repository.Create(fakeBlockNumber, fakeHash, cdpiMetadata, fakeCdpi)
			Expect(err).NotTo(HaveOccurred())

			var cdp test_helpers.ManagedCdp
			queryErr := db.Get(&cdp, `SELECT cdpi, created FROM api.managed_cdp`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(cdp.Id).To(Equal(fakeCdpi))
			Expect(cdp.Created).To(Equal(expectedTimeCreated))
		})
	})

	Describe("cdpi mapping tables", func() {
		fakeCdpi := strconv.Itoa(rand.Int())

		It("returns an error if mapping metadata is missing the key", func() {
			badMetadata := utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerUrns,
				Keys: map[utils.Key]string{},
				Type: utils.Address,
			}
			err := repository.Create(fakeBlockNumber, fakeHash, badMetadata, "")
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Cdpi}))
		})

		Describe("urns", func() {
			var fakeUrnsValue = FakeAddress
			var urnsMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerUrns,
				Keys: map[utils.Key]string{constants.Cdpi: fakeCdpi},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Cdpi),
				ValueFieldName:   "urn",
				Key:              fakeCdpi,
				Value:            fakeUrnsValue,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_urns",
				Repository:       &repository,
				Metadata:         urnsMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

			It("triggers an update to the managed_cdp table", func() {
				err := repository.Create(fakeBlockNumber, fakeHash, urnsMetadata, fakeUrnsValue)
				Expect(err).NotTo(HaveOccurred())

				var cdp test_helpers.ManagedCdp
				queryErr := db.Get(&cdp, `SELECT cdpi, urn_identifier FROM api.managed_cdp`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(cdp.Id).To(Equal(fakeCdpi))
				Expect(cdp.UrnIdentifier).To(Equal(fakeUrnsValue))
			})
		})

		Describe("list_prev", func() {
			var fakePrevValue = strconv.Itoa(rand.Int())
			var prevMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerListPrev,
				Keys: map[utils.Key]string{constants.Cdpi: fakeCdpi},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Cdpi),
				ValueFieldName:   "prev",
				Key:              fakeCdpi,
				Value:            fakePrevValue,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_list_prev",
				Repository:       &repository,
				Metadata:         prevMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("list_next", func() {
			var fakeNextValue = strconv.Itoa(rand.Int())
			var nextMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerListNext,
				Keys: map[utils.Key]string{constants.Cdpi: fakeCdpi},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Cdpi),
				ValueFieldName:   "next",
				Key:              fakeCdpi,
				Value:            fakeNextValue,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_list_next",
				Repository:       &repository,
				Metadata:         nextMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("owns", func() {
			var fakeOwner = FakeAddress
			var ownsMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerOwns,
				Keys: map[utils.Key]string{constants.Cdpi: fakeCdpi},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Cdpi),
				ValueFieldName:   "owner",
				Key:              fakeCdpi,
				Value:            fakeOwner,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_owns",
				Repository:       &repository,
				Metadata:         ownsMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

			It("triggers an update to the managed_cdp table", func() {
				err := repository.Create(fakeBlockNumber, fakeHash, ownsMetadata, fakeOwner)
				Expect(err).NotTo(HaveOccurred())

				var cdp test_helpers.ManagedCdp
				queryErr := db.Get(&cdp, `SELECT cdpi, usr FROM api.managed_cdp`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(cdp.Id).To(Equal(fakeCdpi))
				Expect(cdp.Usr).To(Equal(fakeOwner))
			})
		})

		Describe("ilks", func() {
			var (
				ilksMetadata = utils.StorageValueMetadata{
					Name: cdp_manager.CdpManagerIlks,
					Keys: map[utils.Key]string{constants.Cdpi: fakeCdpi},
					Type: utils.Bytes32,
				}
				fakeIlksValue   = test_helpers.FakeIlk.Hex
				fakeBlockNumber = rand.Int()
				fakeHash        = fakes.FakeHash.Hex()
			)

			It("persists a record", func() {
				createErr := repository.Create(fakeBlockNumber, fakeHash, ilksMetadata, fakeIlksValue)
				Expect(createErr).NotTo(HaveOccurred())

				var result MappingRes
				readErr := db.Get(&result, "SELECT block_number, block_hash, cdpi AS key, ilk_id AS value FROM maker.cdp_manager_ilks")
				Expect(readErr).NotTo(HaveOccurred())

				ilkId, ilkErr := shared.GetOrCreateIlk(fakeIlksValue, db)
				Expect(ilkErr).NotTo(HaveOccurred())

				AssertMapping(result, fakeBlockNumber, fakeHash, fakeCdpi, strconv.FormatInt(ilkId, 10))
			})

			It("doesn't duplicate a record", func() {
				err := repository.Create(fakeBlockNumber, fakeHash, ilksMetadata, fakeIlksValue)
				Expect(err).NotTo(HaveOccurred())

				err = repository.Create(fakeBlockNumber, fakeHash, ilksMetadata, fakeIlksValue)
				Expect(err).NotTo(HaveOccurred())

				var count int
				err = db.Get(&count, "SELECT COUNT(*) FROM maker.cdp_manager_ilks")
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("triggers an update to the managed_cdp table", func() {
				err := repository.Create(fakeBlockNumber, fakeHash, ilksMetadata, fakeIlksValue)
				Expect(err).NotTo(HaveOccurred())

				var cdp test_helpers.ManagedCdp
				queryErr := db.Get(&cdp, `SELECT cdpi, ilk_identifier FROM api.managed_cdp`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(cdp.Id).To(Equal(fakeCdpi))
				Expect(cdp.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
			})
		})
	})

	Describe("owner mapping tables", func() {
		fakeOwner := FakeAddress

		Describe("first", func() {
			var fakeFirstValue = strconv.Itoa(rand.Int())
			var firstMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerFirst,
				Keys: map[utils.Key]string{constants.Owner: fakeOwner},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Owner),
				ValueFieldName:   "first",
				Key:              fakeOwner,
				Value:            fakeFirstValue,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_first",
				Repository:       &repository,
				Metadata:         firstMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("last", func() {
			var fakeLastValue = strconv.Itoa(rand.Int())
			var lastMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerLast,
				Keys: map[utils.Key]string{constants.Owner: fakeOwner},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Owner),
				ValueFieldName:   "last",
				Key:              fakeOwner,
				Value:            fakeLastValue,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_last",
				Repository:       &repository,
				Metadata:         lastMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("count", func() {
			var fakeCountValue = strconv.Itoa(rand.Int())
			var countMetadata = utils.StorageValueMetadata{
				Name: cdp_manager.CdpManagerCount,
				Keys: map[utils.Key]string{constants.Owner: fakeOwner},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.Owner),
				ValueFieldName:   "count",
				Key:              fakeOwner,
				Value:            fakeCountValue,
				IsAMapping:       true,
				StorageTableName: "maker.cdp_manager_count",
				Repository:       &repository,
				Metadata:         countMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})
	})
})
