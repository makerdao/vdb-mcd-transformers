package flap_test

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared"
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
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
)

var _ = Describe("Flap storage repository", func() {
	var (
		db                  *postgres.DB
		repository          flap.FlapStorageRepository
		fakeBlockNumber     int
		fakeHash            string
		flapContractAddress = "flapContractAddress"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = flap.FlapStorageRepository{
			ContractAddress: flapContractAddress,
		}
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

	It("rolls back the record and address insertions if there's a failure", func() {
		var begMetadata = utils.StorageValueMetadata{Name: storage.Beg}
		err := repository.Create(fakeBlockNumber, fakeHash, begMetadata, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

		var addressCount int
		countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
		Expect(countErr).NotTo(HaveOccurred())
		Expect(addressCount).To(Equal(0))
	})

	It("rolls back the records and address insertions for records with bid_ids if there's a failure", func() {
		var bidBidMetadata = utils.StorageValueMetadata{
			Name: storage.BidBid,
			Keys: map[utils.Key]string{constants.BidId: strconv.Itoa(rand.Int())},
			Type: utils.Uint256,
		}

		err := repository.Create(fakeBlockNumber, fakeHash, bidBidMetadata, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

		var addressCount int
		countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
		Expect(countErr).NotTo(HaveOccurred())
		Expect(addressCount).To(Equal(0))
	})

	It("gets or creates the address record", func() {
		var begMetadata = utils.StorageValueMetadata{Name: storage.Beg}
		createErr := repository.Create(fakeBlockNumber, fakeHash, begMetadata, strconv.Itoa(rand.Int()))
		Expect(createErr).NotTo(HaveOccurred())

		addressId, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		var ttlContractAddressId int64
		query := "SELECT address_id FROM maker.flap_beg LIMIT 1"
		getErr := db.Get(&ttlContractAddressId, query)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(ttlContractAddressId).To(Equal(addressId))
	})

	Describe("vat", func() {
		var vatMetadata = utils.StorageValueMetadata{Name: storage.Vat}
		var fakeAddress = FakeAddress

		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Vat,
			Value:            fakeAddress,
			StorageTableName: "maker.flap_vat",
			Repository:       &repository,
			Metadata:         vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("gem", func() {
		var gemMetadata = utils.StorageValueMetadata{Name: storage.Gem}
		var fakeAddress = FakeAddress
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Gem,
			Value:            fakeAddress,
			StorageTableName: "maker.flap_gem",
			Repository:       &repository,
			Metadata:         gemMetadata,
		}

		shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
	})

	Describe("beg", func() {
		var begMetadata = utils.StorageValueMetadata{Name: storage.Beg}
		var fakeBeg = strconv.Itoa(rand.Int())
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Beg,
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
				Name:        storage.Packed,
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
		var kicksMetadata = utils.StorageValueMetadata{Name: storage.Kicks}
		var fakeKicks = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Kicks,
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
		var liveMetadata = utils.StorageValueMetadata{Name: storage.Live}
		var fakeLive = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageVariableBehaviorInputs{
			ValueFieldName:   storage.Live,
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
				Name: storage.BidBid,
				Keys: map[utils.Key]string{},
				Type: utils.Uint256,
			}
			err := repository.Create(fakeBlockNumber, fakeHash, badMetadata, "")
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.BidId}))
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
				StorageTableName: "maker.flap_bid_bid",
				Repository:       &repository,
				Metadata:         bidBidMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

			It("triggers an update to the flap table", func() {
				err := repository.Create(fakeBlockNumber, fakeHash, bidBidMetadata, fakeBidValue)
				Expect(err).NotTo(HaveOccurred())

				var flap FlapRes
				queryErr := db.Get(&flap, `SELECT block_number, block_hash, bid_id, bid FROM maker.flap`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flap.BlockNumber).To(Equal(fakeBlockNumber))
				Expect(flap.BlockHash).To(Equal(fakeHash))
				Expect(flap.BidId).To(Equal(fakeBidId))
				Expect(flap.Bid).To(Equal(fakeBidValue))
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
				StorageTableName: "maker.flap_bid_lot",
				Repository:       &repository,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)

			It("triggers an update to the flap table", func() {
				err := repository.Create(fakeBlockNumber, fakeHash, bidLotMetadata, fakeLotValue)
				Expect(err).NotTo(HaveOccurred())

				var flap FlapRes
				queryErr := db.Get(&flap, `SELECT block_number, block_hash, bid_id, lot FROM maker.flap`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flap.BlockNumber).To(Equal(fakeBlockNumber))
				Expect(flap.BlockHash).To(Equal(fakeHash))
				Expect(flap.BidId).To(Equal(fakeBidId))
				Expect(flap.Lot).To(Equal(fakeLotValue))
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
					err := repository.Create(fakeBlockNumber, fakeHash, bidGuyTicEndMetadata, values)
					Expect(err).NotTo(HaveOccurred())
				})

				It("persists bid guy record", func() {
					var guyResult MappingRes
					selectErr := db.Get(&guyResult, `SELECT block_number, block_hash, bid_id AS key, guy AS value FROM maker.flap_bid_guy`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, fakeBlockNumber, fakeHash, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					selectErr := db.Get(&ticResult, `SELECT block_number, block_hash, bid_id AS key, tic AS value FROM maker.flap_bid_tic`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, fakeBlockNumber, fakeHash, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					selectErr := db.Get(&endResult, `SELECT block_number, block_hash, bid_id AS key, "end" AS value FROM maker.flap_bid_end`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, fakeBlockNumber, fakeHash, fakeBidId, fakeEnd)
				})

				It("triggers an update to the flap table with the latest guy, tic, and end values", func() {
					var flap FlapRes
					queryErr := db.Get(&flap, `SELECT block_number, block_hash, bid_id, guy, tic, "end" FROM maker.flap`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(flap.BlockNumber).To(Equal(fakeBlockNumber))
					Expect(flap.BlockHash).To(Equal(fakeHash))
					Expect(flap.BidId).To(Equal(fakeBidId))
					Expect(flap.Guy).To(Equal(fakeGuy))
					Expect(flap.Tic).To(Equal(fakeTic))
					Expect(flap.End).To(Equal(fakeEnd))
				})
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[1] = ""
				err := repository.Create(fakeBlockNumber, fakeHash, bidGuyTicEndMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for integer"))
			})
		})
	})
})
