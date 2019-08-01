package flip_test

import (
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flip storage repository", func() {
	var (
		db              *postgres.DB
		repo            flip.FlipStorageRepository
		fakeBlockHash   = fakes.FakeHash.Hex()
		fakeBlockNumber int
	)

	BeforeEach(func() {
		fakeBlockNumber = rand.Int()
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = flip.FlipStorageRepository{ContractAddress: constants.EthFlipContractAddressA()}
		repo.SetDB(db)
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := utils.StorageValueMetadata{Name: "unrecognized"}
		flipCreate := func() {
			_ = repo.Create(fakeBlockNumber, fakeBlockHash, unrecognizedMetadata, "")
		}

		Expect(flipCreate).Should(Panic())
	})

	Describe("Variable", func() {
		Describe("Vat", func() {
			vatMetadata := utils.StorageValueMetadata{Name: storage.Vat}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        storage.Vat,
				Value:            FakeAddress,
				StorageTableName: "maker.flip_vat",
				Repository:       &repo,
				Metadata:         vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("Ilk", func() {
			ilkMetadata := utils.StorageValueMetadata{Name: storage.Ilk}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        storage.Ilk,
				Value:            FakeAddress,
				StorageTableName: "maker.flip_ilk",
				Repository:       &repo,
				Metadata:         ilkMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("Beg", func() {
			begMetadata := utils.StorageValueMetadata{Name: storage.Beg}
			fakeBeg := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        storage.Beg,
				Value:            fakeBeg,
				StorageTableName: "maker.flip_beg",
				Repository:       &repo,
				Metadata:         begMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("Ttl and Tau", func() {
			packedNames := make(map[int]string)
			packedNames[0] = storage.Ttl
			packedNames[1] = storage.Tau
			var ttlAndTauMetadata = utils.StorageValueMetadata{
				Name:        storage.Packed,
				PackedNames: packedNames,
			}

			var fakeTtl = strconv.Itoa(rand.Intn(100))
			var fakeTau = strconv.Itoa(rand.Intn(100))
			values := make(map[int]string)
			values[0] = fakeTtl
			values[1] = fakeTau

			It("persists a ttl record", func() {
				err := repo.Create(fakeBlockNumber, fakeBlockHash, ttlAndTauMetadata, values)
				Expect(err).NotTo(HaveOccurred())

				var ttlResult VariableRes
				err = db.Get(&ttlResult, `SELECT block_number, block_hash, ttl AS value FROM maker.flip_ttl`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(ttlResult, fakeBlockNumber, fakeBlockHash, fakeTtl)

				var tauResult VariableRes
				err = db.Get(&tauResult, `SELECT block_number, block_hash, tau AS value FROM maker.flip_tau`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(tauResult, fakeBlockNumber, fakeBlockHash, fakeTau)
			})

			It("panics if the packed name is not recognized", func() {
				packedNames := make(map[int]string)
				packedNames[0] = "notRecognized"

				var badMetadata = utils.StorageValueMetadata{
					Name:        storage.Packed,
					PackedNames: packedNames,
				}

				createFunc := func() {
					_ = repo.Create(fakeBlockNumber, fakeBlockHash, badMetadata, values)
				}
				Expect(createFunc).To(Panic())
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[0] = ""
				err := repo.Create(fakeBlockNumber, fakeBlockHash, ttlAndTauMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
			})

		})

		Describe("Kicks", func() {
			kicksMetadata := utils.StorageValueMetadata{Name: storage.Kicks}
			fakeKicks := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        storage.Kicks,
				Value:            fakeKicks,
				StorageTableName: "maker.flip_kicks",
				Repository:       &repo,
				Metadata:         kicksMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})
	})

	Describe("Bid", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := utils.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[utils.Key]string{},
				Type: utils.Uint256,
			}
			err := repo.Create(fakeBlockNumber, fakeBlockHash, badMetadata, "")
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("BidBid", func() {
			fakeBidValue := strconv.Itoa(rand.Int())
			bidBidMetadata := utils.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "bid",
				Value:            fakeBidValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_bid",
				Repository:       &repo,
				Metadata:         bidBidMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidLot", func() {
			fakeLotValue := strconv.Itoa(rand.Int())
			bidLotMetadata := utils.StorageValueMetadata{
				Name: storage.BidLot,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "lot",
				Value:            fakeLotValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_lot",
				Repository:       &repo,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidGuy", func() {
			bidGuyMetadata := utils.StorageValueMetadata{
				Name: storage.BidGuy,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "guy",
				Value:            FakeAddress,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_guy",
				Repository:       &repo,
				Metadata:         bidGuyMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidTic", func() {
			fakeTicValue := strconv.Itoa(rand.Intn(100))
			bidTicMetadata := utils.StorageValueMetadata{
				Name: storage.BidTic,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint48,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "tic",
				Value:            fakeTicValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_tic",
				Repository:       &repo,
				Metadata:         bidTicMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidEnd", func() {
			fakeEndValue := strconv.Itoa(rand.Intn(100))
			bidEndMetadata := utils.StorageValueMetadata{
				Name: storage.BidEnd,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint48,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "\"end\"",
				Value:            fakeEndValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_end",
				Repository:       &repo,
				Metadata:         bidEndMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidUsr", func() {
			bidUsrMetadata := utils.StorageValueMetadata{
				Name: storage.BidUsr,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "usr",
				Value:            FakeAddress,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_usr",
				Repository:       &repo,
				Metadata:         bidUsrMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidGal", func() {
			bidGalMetadata := utils.StorageValueMetadata{
				Name: storage.BidGal,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "gal",
				Value:            FakeAddress,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_gal",
				Repository:       &repo,
				Metadata:         bidGalMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidTab", func() {
			fakeTabValue := strconv.Itoa(rand.Int())
			bidTabMetadata := utils.StorageValueMetadata{
				Name: storage.BidTab,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "tab",
				Value:            fakeTabValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_tab",
				Repository:       &repo,
				Metadata:         bidTabMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})
	})
})
