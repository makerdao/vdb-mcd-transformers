package flap_test

import (
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Flap storage repository", func() {
	var (
		db              *postgres.DB
		repository      flap.FlapStorageRepository
		fakeBlockNumber int
		fakeHash        string
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = flap.FlapStorageRepository{}
		repository.SetDB(db)
		fakeBlockNumber = rand.Int()
		fakeHash = fakes.FakeHash.Hex()
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := utils.StorageValueMetadata{Name: "unrecognized"}
		flapCreate := func() {
			repository.Create(fakeBlockNumber, fakeHash, unrecognizedMetadata, "")
		}

		Expect(flapCreate).Should(Panic())
	})

	Describe("vat", func() {
		var vatMetadata = utils.StorageValueMetadata{Name: flap.Vat}
		var fakeAddress = FakeAddress

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			FieldName:        flap.Vat,
			Value:            fakeAddress,
			StorageTableName: "maker.flap_vat",
			Repository:       &repository,
			Metadata:         vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("gem", func() {
		var gemMetadata = utils.StorageValueMetadata{Name: flap.Gem}
		var fakeAddress = FakeAddress
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			FieldName:        flap.Gem,
			Value:            fakeAddress,
			StorageTableName: "maker.flap_gem",
			Repository:       &repository,
			Metadata:         gemMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("beg", func() {
		var begMetadata = utils.StorageValueMetadata{Name: flap.Beg}
		var fakeBeg = strconv.Itoa(rand.Int())
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			FieldName:        flap.Beg,
			StorageTableName: "maker.flap_beg",
			Repository:       &repository,
			Metadata:         begMetadata,
			Value:            fakeBeg,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			err := repository.Create(fakeBlockNumber, fakeHash, begMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("ttl and tau", func() {
		packedNames := make(map[int]string)
		packedNames[0] = flap.Ttl
		packedNames[1] = flap.Tau
		var ttlAndTauMetadata = utils.StorageValueMetadata{
			Name:        flap.Packed,
			PackedNames: packedNames,
		}

		var fakeTtl = strconv.Itoa(rand.Intn(100))
		var fakeTau = strconv.Itoa(rand.Intn(100))
		values := make(map[int]string)
		values[0] = fakeTtl
		values[1] = fakeTau

		It("persists a ttl record", func() {
			err := repository.Create(fakeBlockNumber, fakeHash, ttlAndTauMetadata, values)
			Expect(err).NotTo(HaveOccurred())

			var ttlResult VariableRes
			err = db.Get(&ttlResult, `SELECT block_number, block_hash, ttl AS value FROM maker.flap_ttl`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(ttlResult, fakeBlockNumber, fakeHash, fakeTtl)

			var tauResult VariableRes
			err = db.Get(&tauResult, `SELECT block_number, block_hash, tau AS value FROM maker.flap_tau`)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(tauResult, fakeBlockNumber, fakeHash, fakeTau)
		})

		It("panics if the packed name is not recognized", func() {
			packedNames := make(map[int]string)
			packedNames[0] = "notRecognized"

			var badMetadata = utils.StorageValueMetadata{
				Name:        flap.Packed,
				PackedNames: packedNames,
			}

			createFunc := func() {
				repository.Create(fakeBlockNumber, fakeHash, badMetadata, values)
			}
			Expect(createFunc).To(Panic())
		})

		It("returns an error if inserting fails", func() {
			badValues := make(map[int]string)
			badValues[0] = ""
			err := repository.Create(fakeBlockNumber, fakeHash, ttlAndTauMetadata, badValues)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for integer"))
		})
	})

	Describe("kicks", func() {
		var kicksMetadata = utils.StorageValueMetadata{Name: flap.Kicks}
		var fakeKicks = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			FieldName:        flap.Kicks,
			StorageTableName: "maker.flap_kicks",
			Repository:       &repository,
			Metadata:         kicksMetadata,
			Value:            fakeKicks,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			err := repository.Create(fakeBlockNumber, fakeHash, kicksMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("live", func() {
		var liveMetadata = utils.StorageValueMetadata{Name: flap.Live}
		var fakeLive = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			FieldName:        flap.Live,
			StorageTableName: "maker.flap_live",
			Repository:       &repository,
			Metadata:         liveMetadata,
			Value:            fakeLive,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			err := repository.Create(fakeBlockNumber, fakeHash, liveMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("bids", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("for mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := utils.StorageValueMetadata{
				Name: flap.BidBid,
				Keys: map[utils.Key]string{},
				Type: utils.Uint256,
			}
			err := repository.Create(fakeBlockNumber, fakeHash, badMetadata, "")
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("bid_bid", func() {
			var fakeBidValue = strconv.Itoa(rand.Int())
			var bidBidMetadata = utils.StorageValueMetadata{
				Name: flap.BidBid,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "bid",
				Value:            fakeBidValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flap_bid_bid",
				Repository:       &repository,
				Metadata:         bidBidMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("bid_lot", func() {
			var fakeLotValue = strconv.Itoa(rand.Int())
			var bidLotMetadata = utils.StorageValueMetadata{
				Name: flap.BidLot,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "lot",
				Value:            fakeLotValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flap_bid_lot",
				Repository:       &repository,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("bid_guy", func() {
			var fakeGuyValue = FakeAddress
			var bidGuyMetadata = utils.StorageValueMetadata{
				Name: flap.BidGuy,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "guy",
				Value:            fakeGuyValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flap_bid_guy",
				Repository:       &repository,
				Metadata:         bidGuyMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("bid_tic", func() {
			var fakeTicValue = strconv.Itoa(rand.Intn(100))
			var bidTicMetadata = utils.StorageValueMetadata{
				Name: flap.BidTic,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint48,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "tic",
				Value:            fakeTicValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flap_bid_tic",
				Repository:       &repository,
				Metadata:         bidTicMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("bid_end", func() {
			var fakeEndValue = strconv.Itoa(rand.Intn(100))
			var bidEndMetadata = utils.StorageValueMetadata{
				Name: flap.BidEnd,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint48,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "\"end\"",
				Value:            fakeEndValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flap_bid_end",
				Repository:       &repository,
				Metadata:         bidEndMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("bid_gal", func() {
			var fakeGalValue = FakeAddress
			var bidGalMetadata = utils.StorageValueMetadata{
				Name: flap.BidGal,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				FieldName:        "gal",
				Value:            fakeGalValue,
				BidId:            fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flap_bid_gal",
				Repository:       &repository,
				Metadata:         bidGalMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})
	})
})
