package flip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
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
		repo = flip.FlipStorageRepository{ContractAddress: test_data.EthFlipAddress()}
		repo.SetDB(db)
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := utils.StorageValueMetadata{Name: "unrecognized"}
		flipCreate := func() {
			_ = repo.Create(fakeBlockNumber, fakeBlockHash, unrecognizedMetadata, "")
		}

		Expect(flipCreate).Should(Panic())
	})

	It("rolls back the record and address insertions if there's a failure", func() {
		var begMetadata = utils.StorageValueMetadata{Name: storage.Beg}
		err := repo.Create(fakeBlockNumber, fakeBlockHash, begMetadata, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

		var addressCount int
		countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
		Expect(countErr).NotTo(HaveOccurred())
		Expect(addressCount).To(Equal(0))
	})

	Describe("Variable", func() {
		Describe("Vat", func() {
			vatMetadata := utils.StorageValueMetadata{Name: storage.Vat}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				ValueFieldName:   storage.Vat,
				Value:            FakeAddress,
				StorageTableName: "maker.flip_vat",
				Repository:       &repo,
				Metadata:         vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("Ilk", func() {
			It("writes row", func() {
				ilkMetadata := utils.StorageValueMetadata{Name: storage.Ilk}

				insertErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMetadata, FakeIlk)

				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				getErr := db.Get(&result, `SELECT block_number, block_hash, ilk_id AS value FROM maker.flip_ilk`)
				Expect(getErr).NotTo(HaveOccurred())
				ilkID, ilkErr := shared.GetOrCreateIlk(FakeIlk, db)
				Expect(ilkErr).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID))
			})

			It("does not duplicate row", func() {
				ilkMetadata := utils.StorageValueMetadata{Name: storage.Ilk}
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMetadata, FakeIlk)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMetadata, FakeIlk)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.flip_ilk`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Describe("Beg", func() {
			begMetadata := utils.StorageValueMetadata{Name: storage.Beg}
			fakeBeg := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				ValueFieldName:   storage.Beg,
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
				ValueFieldName:   storage.Kicks,
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
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "bid",
				Value:            fakeBidValue,
				Key:              fakeBidId,
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
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "lot",
				Value:            fakeLotValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_lot",
				Repository:       &repo,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("BidGuy, BidTic and BidEnd packed storage", func() {
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
					selectErr := db.Get(&guyResult, `SELECT block_number, block_hash, bid_id AS key, guy AS value FROM maker.flip_bid_guy`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, fakeBlockNumber, fakeBlockHash, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					selectErr := db.Get(&ticResult, `SELECT block_number, block_hash, bid_id AS key, tic AS value FROM maker.flip_bid_tic`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, fakeBlockNumber, fakeBlockHash, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					selectErr := db.Get(&endResult, `SELECT block_number, block_hash, bid_id AS key, "end" AS value FROM maker.flip_bid_end`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, fakeBlockNumber, fakeBlockHash, fakeBidId, fakeEnd)
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

		Describe("BidUsr", func() {
			bidUsrMetadata := utils.StorageValueMetadata{
				Name: storage.BidUsr,
				Keys: map[utils.Key]string{constants.BidId: fakeBidId},
				Type: utils.Address,
			}
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "usr",
				Value:            FakeAddress,
				Key:              fakeBidId,
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
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "gal",
				Value:            FakeAddress,
				Key:              fakeBidId,
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
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "tab",
				Value:            fakeTabValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_tab",
				Repository:       &repo,
				Metadata:         bidTabMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})
	})
})
