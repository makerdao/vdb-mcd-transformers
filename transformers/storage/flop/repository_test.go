package flop_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flop storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 flop.FlopStorageRepository
		blockNumber          int64
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = flop.FlopStorageRepository{ContractAddress: "0x668001c75a9c02d6b10c7a17dbd8aa4afff95037"}
		repo.SetDB(db)
		blockNumber = rand.Int63()
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := vdbStorage.StorageValueMetadata{Name: "unrecognized"}
		flopCreate := func() {
			repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
		}

		Expect(flopCreate).Should(Panic())
	})

	Describe("Vat", func() {
		vatMetadata := vdbStorage.StorageValueMetadata{Name: storage.Vat}

		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   storage.Vat,
			Value:            FakeAddress,
			StorageTableName: "maker.flop_vat",
			Repository:       &repo,
			Metadata:         vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Gem", func() {
		gemMetadata := vdbStorage.StorageValueMetadata{Name: storage.Gem}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   storage.Gem,
			Value:            FakeAddress,
			StorageTableName: "maker.flop_gem",
			Repository:       &repo,
			Metadata:         gemMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Beg", func() {
		begMetadata := vdbStorage.StorageValueMetadata{Name: storage.Beg}
		fakeBeg := strconv.Itoa(rand.Int())

		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   storage.Beg,
			Value:            fakeBeg,
			StorageTableName: "maker.flop_beg",
			Repository:       &repo,
			Metadata:         begMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			createErr := repo.Create(diffID, fakeHeaderID, begMetadata, "")
			Expect(createErr).To(HaveOccurred())
			Expect(createErr.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("Pad", func() {
		padMetadata := vdbStorage.StorageValueMetadata{Name: storage.Pad}
		fakePad := strconv.Itoa(rand.Int())

		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   storage.Pad,
			Value:            fakePad,
			StorageTableName: "maker.flop_pad",
			Repository:       &repo,
			Metadata:         padMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			createErr := repo.Create(diffID, fakeHeaderID, padMetadata, "")
			Expect(createErr).To(HaveOccurred())
			Expect(createErr.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("Ttl and Tau", func() {
		packedNames := make(map[int]string)
		packedNames[0] = storage.Ttl
		packedNames[1] = storage.Tau
		var ttlAndTauMetadata = vdbStorage.StorageValueMetadata{
			Name:        storage.Packed,
			PackedNames: packedNames,
		}

		var fakeTtl = strconv.Itoa(rand.Intn(100))
		var fakeTau = strconv.Itoa(rand.Intn(100))
		values := make(map[int]string)
		values[0] = fakeTtl
		values[1] = fakeTau

		It("persists a ttl record", func() {
			diffID = CreateFakeDiffRecord(db)

			createErr := repo.Create(diffID, fakeHeaderID, ttlAndTauMetadata, values)
			Expect(createErr).NotTo(HaveOccurred())

			var ttlResult VariableRes
			getResErr := db.Get(&ttlResult, `SELECT diff_id, header_id, ttl AS value FROM maker.flop_ttl`)
			Expect(getResErr).NotTo(HaveOccurred())
			AssertVariable(ttlResult, diffID, fakeHeaderID, fakeTtl)
		})

		It("persists a tau record", func() {
			diffID = CreateFakeDiffRecord(db)

			createErr := repo.Create(diffID, fakeHeaderID, ttlAndTauMetadata, values)
			Expect(createErr).NotTo(HaveOccurred())

			var tauResult VariableRes
			getResErr := db.Get(&tauResult, `SELECT diff_id, header_id, tau AS value FROM maker.flop_tau`)
			Expect(getResErr).NotTo(HaveOccurred())
			AssertVariable(tauResult, diffID, fakeHeaderID, fakeTau)
		})

		It("panics if the packed name is not recognized", func() {
			packedNames := make(map[int]string)
			packedNames[0] = "notRecognized"

			var badMetadata = vdbStorage.StorageValueMetadata{
				Name:        storage.Packed,
				PackedNames: packedNames,
			}

			createFunc := func() {
				_ = repo.Create(diffID, fakeHeaderID, badMetadata, values)
			}
			Expect(createFunc).To(Panic())
		})

		It("returns an error if inserting fails", func() {
			badValues := make(map[int]string)
			badValues[0] = ""
			createErr := repo.Create(diffID, fakeHeaderID, ttlAndTauMetadata, badValues)
			Expect(createErr).To(HaveOccurred())
			Expect(createErr.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("Kicks", func() {
		var kicksMetadata = vdbStorage.StorageValueMetadata{Name: storage.Kicks}
		var fakeKicks = strconv.Itoa(rand.Int())
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   storage.Kicks,
			StorageTableName: "maker.flop_kicks",
			Repository:       &repo,
			Metadata:         kicksMetadata,
			Value:            fakeKicks,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Live", func() {
		var liveMetadata = vdbStorage.StorageValueMetadata{Name: storage.Live}
		var fakeKicks = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName:   storage.Live,
			StorageTableName: "maker.flop_live",
			Repository:       &repo,
			Metadata:         liveMetadata,
			Value:            fakeKicks,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Bid", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := vdbStorage.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[vdbStorage.Key]string{},
				Type: vdbStorage.Uint256,
			}
			createErr := repo.Create(diffID, fakeHeaderID, badMetadata, "")
			Expect(createErr).To(MatchError(vdbStorage.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("bid_bid", func() {
			var fakeBidValue = strconv.Itoa(rand.Int())
			var bidBidMetadata = vdbStorage.StorageValueMetadata{
				Name: storage.BidBid,
				Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
				Type: vdbStorage.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "bid",
				Value:            fakeBidValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flop_bid_bid",
				Repository:       &repo,
				Metadata:         bidBidMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			It("triggers an update to the flop table", func() {
				diffID = CreateFakeDiffRecord(db)

				err := repo.Create(diffID, fakeHeaderID, bidBidMetadata, fakeBidValue)
				Expect(err).NotTo(HaveOccurred())

				var flop FlopRes
				queryErr := db.Get(&flop, `SELECT block_number, bid_id, bid FROM maker.flop`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flop.BlockNumber).To(Equal(blockNumber))
				Expect(flop.BidId).To(Equal(fakeBidId))
				Expect(flop.Bid).To(Equal(fakeBidValue))
			})
		})

		Describe("bid_lot", func() {
			var fakeLotValue = strconv.Itoa(rand.Int())
			var bidLotMetadata = vdbStorage.StorageValueMetadata{
				Name: storage.BidLot,
				Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
				Type: vdbStorage.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "lot",
				Value:            fakeLotValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flop_bid_lot",
				Repository:       &repo,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			It("triggers an update to the flop table", func() {
				diffID = CreateFakeDiffRecord(db)

				err := repo.Create(diffID, fakeHeaderID, bidLotMetadata, fakeLotValue)
				Expect(err).NotTo(HaveOccurred())

				var flop FlopRes
				queryErr := db.Get(&flop, `SELECT block_number, bid_id, lot FROM maker.flop`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flop.BlockNumber).To(Equal(blockNumber))
				Expect(flop.BidId).To(Equal(fakeBidId))
				Expect(flop.Lot).To(Equal(fakeLotValue))
			})
		})

		Describe("bid_guy, bid_tic and bid_end packed storage", func() {
			packedNames := make(map[int]string)
			packedNames[0] = storage.BidGuy
			packedNames[1] = storage.BidTic
			packedNames[2] = storage.BidEnd
			var bidGuyTicEndMetadata = vdbStorage.StorageValueMetadata{
				Name:        storage.Packed,
				Keys:        map[vdbStorage.Key]string{constants.BidId: fakeBidId},
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
					diffID = CreateFakeDiffRecord(db)

					err := repo.Create(diffID, fakeHeaderID, bidGuyTicEndMetadata, values)
					Expect(err).NotTo(HaveOccurred())
				})

				It("persists bid guy record", func() {
					var guyResult MappingRes
					selectErr := db.Get(&guyResult, `SELECT diff_id, header_id, bid_id AS key, guy AS value FROM maker.flop_bid_guy`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, diffID, fakeHeaderID, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					selectErr := db.Get(&ticResult, `SELECT diff_id, header_id, bid_id AS key, tic AS value FROM maker.flop_bid_tic`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, diffID, fakeHeaderID, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					selectErr := db.Get(&endResult, `SELECT diff_id, header_id, bid_id AS key, "end" AS value FROM maker.flop_bid_end`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, diffID, fakeHeaderID, fakeBidId, fakeEnd)
				})

				It("triggers an update to the flop table with the latest guy, tic, and end values", func() {
					var flop FlopRes
					queryErr := db.Get(&flop, `SELECT block_number, bid_id, guy, tic, "end" FROM maker.flop`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(flop.BlockNumber).To(Equal(blockNumber))
					Expect(flop.BidId).To(Equal(fakeBidId))
					Expect(flop.Guy).To(Equal(fakeGuy))
					Expect(flop.Tic).To(Equal(fakeTic))
					Expect(flop.End).To(Equal(fakeEnd))
				})
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[1] = ""
				err := repo.Create(diffID, fakeHeaderID, bidGuyTicEndMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for integer"))
			})
		})
	})
})
