package flap_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
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

var _ = Describe("Flap storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		flapContractAddress  = test_data.FlapV100Address()
		repository           = &flap.StorageRepository{ContractAddress: flapContractAddress}
		blockNumber          int64
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repository.SetDB(db)
		blockNumber = rand.Int63()
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}
		flapCreate := func() {
			repository.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
		}

		Expect(flapCreate).Should(Panic())
	})

	It("rolls back the record and address insertions if there's a failure", func() {
		var begMetadata = types.ValueMetadata{Name: storage.Beg}
		err := repository.Create(diffID, fakeHeaderID, begMetadata, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

		var addressCount int
		countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
		Expect(countErr).NotTo(HaveOccurred())
		Expect(addressCount).To(Equal(0))
	})

	It("rolls back the records and address insertions for records with bid_ids if there's a failure", func() {
		var bidBidMetadata = types.ValueMetadata{
			Name: storage.BidBid,
			Keys: map[types.Key]string{constants.BidId: strconv.Itoa(rand.Int())},
			Type: types.Uint256,
		}

		err := repository.Create(diffID, fakeHeaderID, bidBidMetadata, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

		var addressCount int
		countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
		Expect(countErr).NotTo(HaveOccurred())
		Expect(addressCount).To(Equal(0))
	})

	It("gets or creates the address record", func() {
		var begMetadata = types.ValueMetadata{Name: storage.Beg}
		createErr := repository.Create(diffID, fakeHeaderID, begMetadata, strconv.Itoa(rand.Int()))
		Expect(createErr).NotTo(HaveOccurred())

		addressId, addressErr := shared.GetOrCreateAddress(repository.ContractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		var ttlContractAddressId int64
		query := "SELECT address_id FROM maker.flap_beg LIMIT 1"
		getErr := db.Get(&ttlContractAddressId, query)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(ttlContractAddressId).To(Equal(addressId))
	})

	Describe("Wards mapping", func() {
		var fakeUint256 = strconv.Itoa(rand.Intn(1000000))

		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			setupErr := repository.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			err := db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repository.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := shared.GetOrCreateAddress(fakeUserAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())
			AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(userAddressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)
			insertOneErr := repository.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repository.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing user", func() {
			malformedWardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{}, types.Uint256)

			err := repository.Create(diffID, fakeHeaderID, malformedWardsMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.User}))
		})
	})

	Describe("vat", func() {
		var vatMetadata = types.ValueMetadata{Name: storage.Vat}
		var fakeAddress = FakeAddress

		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Vat,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.FlapVatTable,
			Repository:     repository,
			Metadata:       vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("gem", func() {
		var gemMetadata = types.ValueMetadata{Name: storage.Gem}
		var fakeAddress = FakeAddress
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Gem,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.FlapGemTable,
			Repository:     repository,
			Metadata:       gemMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("beg", func() {
		var begMetadata = types.ValueMetadata{Name: storage.Beg}
		var fakeBeg = strconv.Itoa(rand.Int())
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Beg,
			Schema:         constants.MakerSchema,
			TableName:      constants.FlapBegTable,
			Repository:     repository,
			Metadata:       begMetadata,
			Value:          fakeBeg,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			err := repository.Create(diffID, fakeHeaderID, begMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("ttl and tau", func() {
		packedNames := make(map[int]string)
		packedNames[0] = storage.Ttl
		packedNames[1] = storage.Tau
		var ttlAndTauMetadata = types.ValueMetadata{
			Name:        storage.Packed,
			PackedNames: packedNames,
		}

		var fakeTtl = strconv.Itoa(rand.Intn(100))
		var fakeTau = strconv.Itoa(rand.Intn(100))
		values := make(map[int]string)
		values[0] = fakeTtl
		values[1] = fakeTau

		It("persists ttl and tau records", func() {
			err := repository.Create(diffID, fakeHeaderID, ttlAndTauMetadata, values)
			Expect(err).NotTo(HaveOccurred())

			var ttlResult VariableRes
			ttlQuery := fmt.Sprintf(`SELECT diff_id, header_id, ttl AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlapTtlTable))
			err = db.Get(&ttlResult, ttlQuery)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(ttlResult, diffID, fakeHeaderID, fakeTtl)

			var tauResult VariableRes
			tauQuery := fmt.Sprintf(`SELECT diff_id, header_id, tau AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlapTauTable))
			err = db.Get(&tauResult, tauQuery)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(tauResult, diffID, fakeHeaderID, fakeTau)
		})

		It("panics if the packed name is not recognized", func() {
			packedNames := make(map[int]string)
			packedNames[0] = "notRecognized"

			var badMetadata = types.ValueMetadata{
				Name:        storage.Packed,
				PackedNames: packedNames,
			}

			createFunc := func() {
				repository.Create(diffID, fakeHeaderID, badMetadata, values)
			}
			Expect(createFunc).To(Panic())
		})

		It("returns an error if inserting fails", func() {
			badValues := make(map[int]string)
			badValues[0] = ""
			err := repository.Create(diffID, fakeHeaderID, ttlAndTauMetadata, badValues)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: invalid input syntax"))
		})
	})

	Describe("kicks", func() {
		var kicksMetadata = types.ValueMetadata{Name: storage.Kicks}
		var fakeKicks = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Kicks,
			Schema:         constants.MakerSchema,
			TableName:      constants.FlapKicksTable,
			Repository:     repository,
			Metadata:       kicksMetadata,
			Value:          fakeKicks,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			err := repository.Create(diffID, fakeHeaderID, kicksMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("live", func() {
		var liveMetadata = types.ValueMetadata{Name: storage.Live}
		var fakeLive = strconv.Itoa(rand.Intn(100))
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Live,
			Schema:         constants.MakerSchema,
			TableName:      constants.FlapLiveTable,
			Repository:     repository,
			Metadata:       liveMetadata,
			Value:          fakeLive,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

		It("returns an error if inserting fails", func() {
			err := repository.Create(diffID, fakeHeaderID, liveMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))
		})
	})

	Describe("bids", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("for mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := types.ValueMetadata{
				Name: storage.BidBid,
				Keys: map[types.Key]string{},
				Type: types.Uint256,
			}
			err := repository.Create(diffID, fakeHeaderID, badMetadata, "")
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("bid_bid", func() {
			var bidBidMetadata = types.ValueMetadata{
				Name: storage.BidBid,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Uint256,
			}

			var fakeBidValue = strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "bid",
				Value:          fakeBidValue,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlapBidBidTable,
				Repository:     repository,
				Metadata:       bidBidMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repository,
				Metadata:        bidBidMetadata,
				ContractAddress: repository.ContractAddress,
				TriggerTable:    constants.FlapTable,
				FieldTable:      constants.FlapBidBidTable,
				ColumnName:      constants.BidColumn,
			}
			shared_behaviors.InsertBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)

			It("triggers an update to the flap table", func() {
				err := repository.Create(diffID, fakeHeaderID, bidBidMetadata, fakeBidValue)
				Expect(err).NotTo(HaveOccurred())

				var flap FlapRes
				queryErr := db.Get(&flap, `SELECT block_number, bid_id, bid FROM maker.flap`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flap.BlockNumber).To(Equal(blockNumber))
				Expect(flap.BidId).To(Equal(fakeBidId))
				Expect(flap.Bid).To(Equal(fakeBidValue))
			})
		})

		Describe("bid_lot", func() {
			var bidLotMetadata = types.ValueMetadata{
				Name: storage.BidLot,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Uint256,
			}

			var fakeLotValue = strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "lot",
				Value:          fakeLotValue,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlapBidLotTable,
				Repository:     repository,
				Metadata:       bidLotMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repository,
				Metadata:        bidLotMetadata,
				ContractAddress: repository.ContractAddress,
				TriggerTable:    constants.FlapTable,
				FieldTable:      constants.FlapBidLotTable,
				ColumnName:      constants.LotColumn,
			}
			shared_behaviors.InsertBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)

			It("triggers an update to the flap table", func() {
				err := repository.Create(diffID, fakeHeaderID, bidLotMetadata, fakeLotValue)
				Expect(err).NotTo(HaveOccurred())

				var flap FlapRes
				queryErr := db.Get(&flap, `SELECT block_number, bid_id, lot FROM maker.flap`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(flap.BlockNumber).To(Equal(blockNumber))
				Expect(flap.BidId).To(Equal(fakeBidId))
				Expect(flap.Lot).To(Equal(fakeLotValue))
			})
		})

		Describe("bid_guy, bid_tic and bid_end packed storage", func() {
			packedNames := make(map[int]string)
			packedNames[0] = storage.BidGuy
			packedNames[1] = storage.BidTic
			packedNames[2] = storage.BidEnd
			var bidGuyTicEndMetadata = types.ValueMetadata{
				Name:        storage.Packed,
				Keys:        map[types.Key]string{constants.BidId: fakeBidId},
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

					err := repository.Create(diffID, fakeHeaderID, bidGuyTicEndMetadata, values)
					Expect(err).NotTo(HaveOccurred())
				})

				It("persists bid guy record", func() {
					var guyResult MappingRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, bid_id AS key, guy AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlapBidGuyTable))
					selectErr := db.Get(&guyResult, query)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, diffID, fakeHeaderID, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, bid_id AS key, tic AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlapBidTicTable))
					selectErr := db.Get(&ticResult, query)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, diffID, fakeHeaderID, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, bid_id AS key, "end" AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlapBidEndTable))
					selectErr := db.Get(&endResult, query)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, diffID, fakeHeaderID, fakeBidId, fakeEnd)
				})

				It("triggers an update to the flap table with the latest guy, tic, and end values", func() {
					var flap FlapRes
					queryErr := db.Get(&flap, `SELECT block_number, bid_id, guy, tic, "end" FROM maker.flap`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(flap.BlockNumber).To(Equal(blockNumber))
					Expect(flap.BidId).To(Equal(fakeBidId))
					Expect(flap.Guy).To(Equal(fakeGuy))
					Expect(flap.Tic).To(Equal(fakeTic))
					Expect(flap.End).To(Equal(fakeEnd))
				})
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[1] = ""
				err := repository.Create(diffID, fakeHeaderID, bidGuyTicEndMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("pq: invalid input syntax"))
			})
		})

		var bidGuyMetadata = types.ValueMetadata{
			Name:        storage.Packed,
			Keys:        map[types.Key]string{constants.BidId: fakeBidId},
			Type:        types.PackedSlot,
			PackedNames: map[int]string{0: storage.BidGuy},
		}
		guyTriggerInput := shared_behaviors.BidTriggerTestInput{
			Repository:      repository,
			Metadata:        bidGuyMetadata,
			ContractAddress: repository.ContractAddress,
			TriggerTable:    constants.FlapTable,
			FieldTable:      constants.FlapBidGuyTable,
			ColumnName:      constants.GuyColumn,
			PackedValueType: types.Address,
		}
		shared_behaviors.InsertBidSnapshotTriggerTests(guyTriggerInput)
		shared_behaviors.UpdateBidSnapshotTriggerTests(guyTriggerInput)

		var bidTicMetadata = types.ValueMetadata{
			Name:        storage.Packed,
			Keys:        map[types.Key]string{constants.BidId: fakeBidId},
			Type:        types.PackedSlot,
			PackedNames: map[int]string{0: storage.BidTic},
		}
		ticTriggerInput := shared_behaviors.BidTriggerTestInput{
			Repository:      repository,
			Metadata:        bidTicMetadata,
			ContractAddress: repository.ContractAddress,
			TriggerTable:    constants.FlapTable,
			FieldTable:      constants.FlapBidTicTable,
			ColumnName:      constants.TicColumn,
			PackedValueType: types.Uint48,
		}
		shared_behaviors.InsertBidSnapshotTriggerTests(ticTriggerInput)
		shared_behaviors.UpdateBidSnapshotTriggerTests(ticTriggerInput)

		var bidEndMetadata = types.ValueMetadata{
			Name:        storage.Packed,
			Keys:        map[types.Key]string{constants.BidId: fakeBidId},
			Type:        types.PackedSlot,
			PackedNames: map[int]string{0: storage.BidEnd},
		}
		endTriggerInput := shared_behaviors.BidTriggerTestInput{
			Repository:      repository,
			Metadata:        bidEndMetadata,
			ContractAddress: repository.ContractAddress,
			TriggerTable:    constants.FlapTable,
			FieldTable:      constants.FlapBidEndTable,
			ColumnName:      constants.EndColumn,
			PackedValueType: types.Uint48,
		}
		shared_behaviors.InsertBidSnapshotTriggerTests(endTriggerInput)
		shared_behaviors.UpdateBidSnapshotTriggerTests(endTriggerInput)
	})
})
