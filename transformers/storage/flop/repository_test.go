package flop_test

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
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Flop storage repository", func() {
	var (
		db              *postgres.DB
		repo            flop.FlopStorageRepository
		fakeBlockHash   string
		fakeBlockNumber int
	)

	BeforeEach(func() {
		fakeBlockNumber = rand.Int()
		fakeBlockHash = fakes.FakeHash.Hex()
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = flop.FlopStorageRepository{ContractAddress: "0x668001c75a9c02d6b10c7a17dbd8aa4afff95037"}
		repo.SetDB(db)
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := utils.StorageValueMetadata{Name: "unrecognized"}
		flopCreate := func() {
			repo.Create(fakeBlockNumber, fakeBlockHash, unrecognizedMetadata, "")
		}

		Expect(flopCreate).Should(Panic())
	})

	Describe("Vat", func() {
		vatMetadata := utils.StorageValueMetadata{Name: storage.Vat}

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Vat,
			Value:            FakeAddress,
			StorageTableName: "maker.flop_vat",
			Repository:       &repo,
			Metadata:         vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("Gem", func() {
		gemMetadata := utils.StorageValueMetadata{Name: storage.Gem}
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Gem,
			Value:            FakeAddress,
			StorageTableName: "maker.flop_gem",
			Repository:       &repo,
			Metadata:         gemMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("Beg", func() {
		begMetadata := utils.StorageValueMetadata{Name: storage.Beg}
		fakeBeg := strconv.Itoa(rand.Int())

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Beg,
			Value:            fakeBeg,
			StorageTableName: "maker.flop_beg",
			Repository:       &repo,
			Metadata:         begMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			createErr := repo.Create(fakeBlockNumber, fakeBlockHash, begMetadata, "")
			Expect(createErr).To(HaveOccurred())
			Expect(createErr.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("Pad", func() {
		padMetadata := utils.StorageValueMetadata{Name: storage.Pad}
		fakePad := strconv.Itoa(rand.Int())

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Pad,
			Value:            fakePad,
			StorageTableName: "maker.flop_pad",
			Repository:       &repo,
			Metadata:         padMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			createErr := repo.Create(fakeBlockNumber, fakeBlockHash, padMetadata, "")
			Expect(createErr).To(HaveOccurred())
			Expect(createErr.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
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
			createErr := repo.Create(fakeBlockNumber, fakeBlockHash, ttlAndTauMetadata, values)
			Expect(createErr).NotTo(HaveOccurred())

			var ttlResult VariableRes
			getResErr := db.Get(&ttlResult, `SELECT block_number, block_hash, ttl AS value FROM maker.flop_ttl`)
			Expect(getResErr).NotTo(HaveOccurred())
			AssertVariable(ttlResult, fakeBlockNumber, fakeBlockHash, fakeTtl)
		})

		It("persists a tau record", func() {
			createErr := repo.Create(fakeBlockNumber, fakeBlockHash, ttlAndTauMetadata, values)
			Expect(createErr).NotTo(HaveOccurred())

			var tauResult VariableRes
			getResErr := db.Get(&tauResult, `SELECT block_number, block_hash, tau AS value FROM maker.flop_tau`)
			Expect(getResErr).NotTo(HaveOccurred())
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
			createErr := repo.Create(fakeBlockNumber, fakeBlockHash, ttlAndTauMetadata, badValues)
			Expect(createErr).To(HaveOccurred())
			Expect(createErr.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("Kicks", func() {
		var kicksMetadata = utils.StorageValueMetadata{Name: storage.Kicks}
		var fakeKicks = strconv.Itoa(rand.Int())
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Kicks,
			StorageTableName: "maker.flop_kicks",
			Repository:       &repo,
			Metadata:         kicksMetadata,
			Value:            fakeKicks,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("Live", func() {
		var liveMetadata = utils.StorageValueMetadata{Name: storage.Live}
		var fakeKicks = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Live,
			StorageTableName: "maker.flop_live",
			Repository:       &repo,
			Metadata:         liveMetadata,
			Value:            fakeKicks,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("Bid", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := utils.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[utils.Key]string{},
				Type: utils.Uint256,
			}
			createErr := repo.Create(fakeBlockNumber, fakeBlockHash, badMetadata, "")
			Expect(createErr).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("bid_bid", func() {
			var fakeBidValue = strconv.Itoa(rand.Int())
			var bidBidMetadata = utils.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "bid",
				Value:            fakeBidValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flop_bid_bid",
				Repository:       &repo,
				Metadata:         bidBidMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

			It("triggers an update to the flop table", func() {
				err := repo.Create(fakeBlockNumber, fakeBlockHash, bidBidMetadata, fakeBidValue)
				Expect(err).NotTo(HaveOccurred())

				var flop FlopRes
				queryErr := db.Get(&flop, `SELECT block_number, block_hash, bid_id, bid FROM maker.flop`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flop.BlockNumber).To(Equal(fakeBlockNumber))
				Expect(flop.BlockHash).To(Equal(fakeBlockHash))
				Expect(flop.BidId).To(Equal(fakeBidId))
				Expect(flop.Bid).To(Equal(fakeBidValue))
			})
		})

		Describe("bid_lot", func() {
			var fakeLotValue = strconv.Itoa(rand.Int())
			var bidLotMetadata = utils.StorageValueMetadata{
				Name: storage.BidLot,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Uint256,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "lot",
				Value:            fakeLotValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flop_bid_lot",
				Repository:       &repo,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

			It("triggers an update to the flop table", func() {
				err := repo.Create(fakeBlockNumber, fakeBlockHash, bidLotMetadata, fakeLotValue)
				Expect(err).NotTo(HaveOccurred())

				var flop FlopRes
				queryErr := db.Get(&flop, `SELECT block_number, block_hash, bid_id, lot FROM maker.flop`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flop.BlockNumber).To(Equal(fakeBlockNumber))
				Expect(flop.BlockHash).To(Equal(fakeBlockHash))
				Expect(flop.BidId).To(Equal(fakeBidId))
				Expect(flop.Lot).To(Equal(fakeLotValue))
			})
		})

		Describe("bid_guy, bid_tic and bid_end packed storage", func() {
			packedNames := make(map[int]string)
			packedNames[0] = storage.BidGuy
			packedNames[1] = storage.BidTic
			packedNames[2] = storage.BidEnd
			var bidGuyTicEndMetadata = utils.StorageValueMetadata{
				Name:        storage.Packed,
				Keys:        map[utils.Key]string{constants.BidId: fakeBidId},
				PackedNames: packedNames,
			}

			Describe("Create", func() {
				fakeGuy := FakeAddress
				fakeTic := strconv.Itoa(rand.Intn(100))
				fakeEnd := strconv.Itoa(rand.Intn(100))
				values := make(map[int]string)
				values[0] = fakeGuy
				values[1] = fakeTic
				values[2] = fakeEnd

				BeforeEach(func() {
					err := repo.Create(fakeBlockNumber, fakeBlockHash, bidGuyTicEndMetadata, values)
					Expect(err).NotTo(HaveOccurred())
				})

				It("persists bid guy record", func() {
					var guyResult MappingRes
					selectErr := db.Get(&guyResult, `SELECT block_number, block_hash, bid_id AS key, guy AS value FROM maker.flop_bid_guy`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, fakeBlockNumber, fakeBlockHash, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					selectErr := db.Get(&ticResult, `SELECT block_number, block_hash, bid_id AS key, tic AS value FROM maker.flop_bid_tic`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, fakeBlockNumber, fakeBlockHash, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					selectErr := db.Get(&endResult, `SELECT block_number, block_hash, bid_id AS key, "end" AS value FROM maker.flop_bid_end`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, fakeBlockNumber, fakeBlockHash, fakeBidId, fakeEnd)
				})

				It("triggers an update to the flop table with the latest guy, tic, and end values", func() {
					var flop FlopRes
					queryErr := db.Get(&flop, `SELECT block_number, block_hash, bid_id, guy, tic, "end" FROM maker.flop`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(flop.BlockNumber).To(Equal(fakeBlockNumber))
					Expect(flop.BlockHash).To(Equal(fakeBlockHash))
					Expect(flop.BidId).To(Equal(fakeBidId))
					Expect(flop.Guy).To(Equal(fakeGuy))
					Expect(flop.Tic).To(Equal(fakeTic))
					Expect(flop.End).To(Equal(fakeEnd))
				})
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[1] = ""
				err := repo.Create(fakeBlockNumber, fakeBlockHash, bidGuyTicEndMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for integer"))
			})
		})
	})
})
