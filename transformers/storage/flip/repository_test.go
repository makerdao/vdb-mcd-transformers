package flip_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
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

var _ = Describe("Flip storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 = &flip.StorageRepository{ContractAddress: test_data.FlipEthV110Address()}
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}
		flipCreate := func() {
			_ = repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
		}

		Expect(flipCreate).Should(Panic())
	})

	It("rolls back the record and address insertions if there's a failure", func() {
		var begMetadata = types.ValueMetadata{Name: storage.Beg}
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
			vatMetadata := types.ValueMetadata{Name: storage.Vat}
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: storage.Vat,
				Value:          FakeAddress,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipVatTable,
				Repository:     repo,
				Metadata:       vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Cat", func() {
			It("writes row", func() {
				catMetadata := types.ValueMetadata{Name: storage.Cat}
				insertErr := repo.Create(diffID, fakeHeaderID, catMetadata, FakeAddress)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, cat AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipCatTable))
				getErr := db.Get(&result, query)
				Expect(getErr).NotTo(HaveOccurred())
				addressID, addressErr := shared.GetOrCreateAddress(FakeAddress, db)
				Expect(addressErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
			})

			It("does not duplicate row", func() {
				catMetadata := types.ValueMetadata{Name: storage.Cat}
				insertOneErr := repo.Create(diffID, fakeHeaderID, catMetadata, FakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, catMetadata, FakeAddress)
				Expect(insertTwoErr).NotTo(HaveOccurred())

				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipCatTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Describe("Wards", func() {
			var fakeUint256 = strconv.Itoa(rand.Intn(1000000))

			It("writes a row", func() {
				fakeUserAddress := "0x" + fakes.RandomString(40)
				wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

				setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
				Expect(setupErr).NotTo(HaveOccurred())

				var result MappingResWithAddress
				query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
				err := db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
				Expect(contractAddressErr).NotTo(HaveOccurred())
				userAddressID, userAddressErr := shared.GetOrCreateAddress(fakeUserAddress, db)
				Expect(userAddressErr).NotTo(HaveOccurred())
				AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(userAddressID, 10), fakeUint256)
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

		Describe("Ilk", func() {
			It("writes row", func() {
				ilkMetadata := types.ValueMetadata{Name: storage.Ilk}
				insertErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipIlkTable))
				getErr := db.Get(&result, query)
				Expect(getErr).NotTo(HaveOccurred())
				ilkID, ilkErr := shared.GetOrCreateIlk(FakeIlk, db)
				Expect(ilkErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10))
			})

			It("does not duplicate row", func() {
				ilkMetadata := types.ValueMetadata{Name: storage.Ilk}
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipIlkTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Describe("Beg", func() {
			begMetadata := types.ValueMetadata{Name: storage.Beg}
			fakeBeg := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: storage.Beg,
				Value:          fakeBeg,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipBegTable,
				Repository:     repo,
				Metadata:       begMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Ttl and Tau", func() {
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

			It("persists a ttl record", func() {
				err := repo.Create(diffID, fakeHeaderID, ttlAndTauMetadata, values)
				Expect(err).NotTo(HaveOccurred())

				var ttlResult VariableRes
				ttlQuery := fmt.Sprintf(`SELECT diff_id, header_id, ttl AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipTtlTable))
				err = db.Get(&ttlResult, ttlQuery)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(ttlResult, diffID, fakeHeaderID, fakeTtl)

				var tauResult VariableRes
				tauQuery := fmt.Sprintf(`SELECT diff_id, header_id, tau AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipTauTable))
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
			kicksMetadata := types.ValueMetadata{Name: storage.Kicks}
			fakeKicks := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: storage.Kicks,
				Value:          fakeKicks,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipKicksTable,
				Repository:     repo,
				Metadata:       kicksMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})

	Describe("Bid", func() {
		var fakeBidId = strconv.Itoa(rand.Int())

		It("mappings returns an error if the metadata is missing the bid_id", func() {
			badMetadata := types.ValueMetadata{
				Name: storage.BidBid,
				Keys: map[types.Key]string{},
				Type: types.Uint256,
			}
			err := repo.Create(diffID, fakeHeaderID, badMetadata, "")
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.BidId}))
		})

		Describe("BidBid", func() {
			bidBidMetadata := types.ValueMetadata{
				Name: storage.BidBid,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Uint256,
			}

			fakeBidValue := strconv.Itoa(rand.Int())
			storageInputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "bid",
				Value:          fakeBidValue,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipBidBidTable,
				Repository:     repo,
				Metadata:       bidBidMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&storageInputs)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidBidMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidBidTable,
				ColumnName:      constants.BidColumn,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)
		})

		Describe("BidLot", func() {
			bidLotMetadata := types.ValueMetadata{
				Name: storage.BidLot,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Uint256,
			}

			fakeLotValue := strconv.Itoa(rand.Int())
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "lot",
				Value:          fakeLotValue,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipBidLotTable,
				Repository:     repo,
				Metadata:       bidLotMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidLotMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidLotTable,
				ColumnName:      constants.LotColumn,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)
		})

		Describe("BidGuy, BidTic and BidEnd packed storage", func() {
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
					err := repo.Create(diffID, fakeHeaderID, bidGuyTicEndMetadata, values)
					Expect(err).NotTo(HaveOccurred())
				})

				It("persists bid guy record", func() {
					var guyResult MappingRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, bid_id AS key, guy AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipBidGuyTable))
					selectErr := db.Get(&guyResult, query)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(guyResult, diffID, fakeHeaderID, fakeBidId, fakeGuy)
				})

				It("persists bid tic record", func() {
					var ticResult MappingRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, bid_id AS key, tic AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipBidTicTable))
					selectErr := db.Get(&ticResult, query)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(ticResult, diffID, fakeHeaderID, fakeBidId, fakeTic)
				})

				It("persists bid end record", func() {
					var endResult MappingRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, bid_id AS key, "end" AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.FlipBidEndTable))
					selectErr := db.Get(&endResult, query)
					Expect(selectErr).NotTo(HaveOccurred())
					AssertMapping(endResult, diffID, fakeHeaderID, fakeBidId, fakeEnd)
				})
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[1] = ""
				err := repo.Create(diffID, fakeHeaderID, bidGuyTicEndMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("pq: invalid input syntax"))
			})

			var bidGuyMetadata = types.ValueMetadata{
				Name:        storage.Packed,
				Keys:        map[types.Key]string{constants.BidId: fakeBidId},
				Type:        types.PackedSlot,
				PackedNames: map[int]string{0: storage.BidGuy},
			}
			guyTriggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidGuyMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidGuyTable,
				ColumnName:      constants.GuyColumn,
				PackedValueType: types.Address,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(guyTriggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(guyTriggerInput)

			var bidTicMetadata = types.ValueMetadata{
				Name:        storage.Packed,
				Keys:        map[types.Key]string{constants.BidId: fakeBidId},
				Type:        types.PackedSlot,
				PackedNames: map[int]string{0: storage.BidTic},
			}
			ticTriggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidTicMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidTicTable,
				ColumnName:      constants.TicColumn,
				PackedValueType: types.Uint48,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(ticTriggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(ticTriggerInput)

			var bidEndMetadata = types.ValueMetadata{
				Name:        storage.Packed,
				Keys:        map[types.Key]string{constants.BidId: fakeBidId},
				Type:        types.PackedSlot,
				PackedNames: map[int]string{0: storage.BidEnd},
			}
			endTriggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidEndMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidEndTable,
				ColumnName:      constants.EndColumn,
				PackedValueType: types.Uint48,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(endTriggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(endTriggerInput)
		})

		Describe("BidUsr", func() {
			bidUsrMetadata := types.ValueMetadata{
				Name: storage.BidUsr,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Address,
			}

			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "usr",
				Value:          FakeAddress,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipBidUsrTable,
				Repository:     repo,
				Metadata:       bidUsrMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidUsrMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidUsrTable,
				ColumnName:      constants.UsrColumn,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)
		})

		Describe("BidGal", func() {
			bidGalMetadata := types.ValueMetadata{
				Name: storage.BidGal,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Address,
			}

			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "gal",
				Value:          FakeAddress,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipBidGalTable,
				Repository:     repo,
				Metadata:       bidGalMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidGalMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidGalTable,
				ColumnName:      constants.GalColumn,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)
		})

		Describe("BidTab", func() {
			bidTabMetadata := types.ValueMetadata{
				Name: storage.BidTab,
				Keys: map[types.Key]string{constants.BidId: fakeBidId},
				Type: types.Uint256,
			}

			fakeTabValue := strconv.Itoa(rand.Int())
			storageInput := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.BidId),
				ValueFieldName: "tab",
				Value:          fakeTabValue,
				Key:            fakeBidId,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.FlipBidTabTable,
				Repository:     repo,
				Metadata:       bidTabMetadata,
			}
			shared_behaviors.SharedStorageRepositoryBehaviors(&storageInput)

			triggerInput := shared_behaviors.BidTriggerTestInput{
				Repository:      repo,
				Metadata:        bidTabMetadata,
				ContractAddress: repo.ContractAddress,
				TriggerTable:    constants.FlipTable,
				FieldTable:      constants.FlipBidTabTable,
				ColumnName:      constants.TabColumn,
			}
			shared_behaviors.InsertFlipBidSnapshotTriggerTests(triggerInput)
			shared_behaviors.UpdateBidSnapshotTriggerTests(triggerInput)
		})
	})
})
