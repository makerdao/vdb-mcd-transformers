package flip_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flip storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 flip.FlipStorageRepository
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = flip.FlipStorageRepository{ContractAddress: test_data.EthFlipAddress()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := vdbStorage.ValueMetadata{Name: "unrecognized"}
		flipCreate := func() {
			_ = repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
		}

		Expect(flipCreate).Should(Panic())
	})

	It("rolls back the record and address insertions if there's a failure", func() {
		var begMetadata = vdbStorage.ValueMetadata{Name: storage.Beg}
		err := repo.Create(diffID, fakeHeaderID, begMetadata, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

		var addressCount int
		countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
		Expect(countErr).NotTo(HaveOccurred())
		Expect(addressCount).To(Equal(0))
	})

	Describe("Variable", func() {
		Describe("Vat", func() {
			vatMetadata := vdbStorage.ValueMetadata{Name: storage.Vat}
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:   storage.Vat,
				Value:            FakeAddress,
				StorageTableName: "maker.flip_vat",
				Repository:       &repo,
				Metadata:         vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Ilk", func() {
			It("writes row", func() {
				diffID = CreateFakeDiffRecord(db)

				ilkMetadata := vdbStorage.ValueMetadata{Name: storage.Ilk}
				insertErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				getErr := db.Get(&result, `SELECT diff_id, header_id, ilk_id AS value FROM maker.flip_ilk`)
				Expect(getErr).NotTo(HaveOccurred())
				ilkID, ilkErr := shared.GetOrCreateIlk(FakeIlk, db)
				Expect(ilkErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10))
			})

			It("does not duplicate row", func() {
				diffID = CreateFakeDiffRecord(db)

				ilkMetadata := vdbStorage.ValueMetadata{Name: storage.Ilk}
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.flip_ilk`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Describe("Beg", func() {
			begMetadata := vdbStorage.ValueMetadata{Name: storage.Beg}
			fakeBeg := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:   storage.Beg,
				Value:            fakeBeg,
				StorageTableName: "maker.flip_beg",
				Repository:       &repo,
				Metadata:         begMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Ttl and Tau", func() {
			packedNames := make(map[int]string)
			packedNames[0] = storage.Ttl
			packedNames[1] = storage.Tau
			var ttlAndTauMetadata = vdbStorage.ValueMetadata{
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

				err := repo.Create(diffID, fakeHeaderID, ttlAndTauMetadata, values)
				Expect(err).NotTo(HaveOccurred())

				var ttlResult VariableRes
				err = db.Get(&ttlResult, `SELECT diff_id, header_id, ttl AS value FROM maker.flip_ttl`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(ttlResult, diffID, fakeHeaderID, fakeTtl)

				var tauResult VariableRes
				err = db.Get(&tauResult, `SELECT diff_id, header_id, tau AS value FROM maker.flip_tau`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(tauResult, diffID, fakeHeaderID, fakeTau)
			})

			It("panics if the packed name is not recognized", func() {
				packedNames := make(map[int]string)
				packedNames[0] = "notRecognized"

				var badMetadata = vdbStorage.ValueMetadata{
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
				err := repo.Create(diffID, fakeHeaderID, ttlAndTauMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
			})

		})

		Describe("Kicks", func() {
			kicksMetadata := vdbStorage.ValueMetadata{Name: storage.Kicks}
			fakeKicks := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:   storage.Kicks,
				Value:            fakeKicks,
				StorageTableName: "maker.flip_kicks",
				Repository:       &repo,
				Metadata:         kicksMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})

	Describe("Bid", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := vdbStorage.ValueMetadata{
				Name: storage.BidBid,
				Keys: map[vdbStorage.Key]string{},
				Type: vdbStorage.Uint256,
			}
			err := repo.Create(diffID, fakeHeaderID, badMetadata, "")
			Expect(err).To(MatchError(vdbStorage.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("BidBid", func() {
			fakeBidValue := strconv.Itoa(rand.Int())
			bidBidMetadata := vdbStorage.ValueMetadata{
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
				StorageTableName: "maker.flip_bid_bid",
				Repository:       &repo,
				Metadata:         bidBidMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("BidLot", func() {
			fakeLotValue := strconv.Itoa(rand.Int())
			bidLotMetadata := vdbStorage.ValueMetadata{
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
				StorageTableName: "maker.flip_bid_lot",
				Repository:       &repo,
				Metadata:         bidLotMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("BidGuy, BidTic and BidEnd packed storage", func() {
			packedNames := make(map[int]string)
			packedNames[0] = storage.BidGuy
			packedNames[1] = storage.BidTic
			packedNames[2] = storage.BidEnd
			var bidGuyTicEndMetadata = vdbStorage.ValueMetadata{
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
					selectErr := db.Get(&guyResult, `SELECT diff_id, header_id, bid_id AS key, guy AS value FROM maker.flip_bid_guy`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, diffID, fakeHeaderID, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					selectErr := db.Get(&ticResult, `SELECT diff_id, header_id, bid_id AS key, tic AS value FROM maker.flip_bid_tic`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, diffID, fakeHeaderID, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					selectErr := db.Get(&endResult, `SELECT diff_id, header_id, bid_id AS key, "end" AS value FROM maker.flip_bid_end`)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, diffID, fakeHeaderID, fakeBidId, fakeEnd)
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

		Describe("BidUsr", func() {
			bidUsrMetadata := vdbStorage.ValueMetadata{
				Name: storage.BidUsr,
				Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
				Type: vdbStorage.Address,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "usr",
				Value:            FakeAddress,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_usr",
				Repository:       &repo,
				Metadata:         bidUsrMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("BidGal", func() {
			bidGalMetadata := vdbStorage.ValueMetadata{
				Name: storage.BidGal,
				Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
				Type: vdbStorage.Address,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "gal",
				Value:            FakeAddress,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_gal",
				Repository:       &repo,
				Metadata:         bidGalMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("BidTab", func() {
			fakeTabValue := strconv.Itoa(rand.Int())
			bidTabMetadata := vdbStorage.ValueMetadata{
				Name: storage.BidTab,
				Keys: map[vdbStorage.Key]string{constants.BidId: fakeBidId},
				Type: vdbStorage.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:     string(constants.BidId),
				ValueFieldName:   "tab",
				Value:            fakeTabValue,
				Key:              fakeBidId,
				IsAMapping:       true,
				StorageTableName: "maker.flip_bid_tab",
				Repository:       &repo,
				Metadata:         bidTabMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})
})
