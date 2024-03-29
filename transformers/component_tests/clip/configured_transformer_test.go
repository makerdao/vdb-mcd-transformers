package clip

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db              = test_config.NewTestDB(test_config.NewTestNode())
		contractAddress = common.HexToAddress(test_data.ClipLinkAV130Address())
		transformer     storage.Transformer
		header          = fakes.FakeHeader
	)

	BeforeEach(func() {
		storageKeysLookup := storage.NewKeysLookup(clip.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress.Hex()))
		repository := clip.StorageRepository{ContractAddress: contractAddress.Hex()}
		transformer = storage.Transformer{
			Address:           contractAddress,
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		header.Id, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Clip Dog storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		value := common.HexToHash("000000000000000000000000135954d155898d42c90d2a57824c690e0c7bef1b")
		clipDogDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(clipDogDiff)
		Expect(err).NotTo(HaveOccurred())

		diffAddressID, diffAddressErr := repository.GetOrCreateAddress(db, "0x135954d155898d42c90d2a57824c690e0c7bef1b")
		Expect(diffAddressErr).NotTo(HaveOccurred())

		var dogResult test_helpers.VariableResWithAddress
		err = db.Get(&dogResult, `SELECT diff_id, header_id, address_id, dog AS value FROM maker.clip_dog`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(dogResult, clipDogDiff.ID, header.Id, contractAddressID, strconv.FormatInt(diffAddressID, 10))
	})

	It("reads in a Clip Vow storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("000000000000000000000000a950524441892a31ebddf91d3ceefa04bf454466")
		clipVowDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(clipVowDiff)
		Expect(err).NotTo(HaveOccurred())

		diffAddressID, diffAddressErr := repository.GetOrCreateAddress(db, "0xa950524441892a31ebddf91d3ceefa04bf454466")
		Expect(diffAddressErr).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableResWithAddress
		err = db.Get(&vowResult, `SELECT diff_id, header_id, address_id, vow AS value FROM maker.clip_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(vowResult, clipVowDiff.ID, header.Id, contractAddressID, strconv.FormatInt(diffAddressID, 10))
	})

	It("reads in a Clip Spotter storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("00000000000000000000000065c79fcb50ca1594b025960e539ed7a9a6d434a3")
		clipSpotterDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(clipSpotterDiff)
		Expect(err).NotTo(HaveOccurred())

		diffAddressID, diffAddressErr := repository.GetOrCreateAddress(db, "0x65c79fcb50ca1594b025960e539ed7a9a6d434a3")
		Expect(diffAddressErr).NotTo(HaveOccurred())

		var spotterResult test_helpers.VariableResWithAddress
		err = db.Get(&spotterResult, `SELECT diff_id, header_id, address_id, spotter AS value FROM maker.clip_spotter`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(spotterResult, clipSpotterDiff.ID, header.Id, contractAddressID, strconv.FormatInt(diffAddressID, 10))
	})

	It("reads in a Clip Calc storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("0000000000000000000000007b1696677107e48b152e9bf400293e98b7d86eb1")
		clipCalcDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(clipCalcDiff)
		Expect(err).NotTo(HaveOccurred())

		diffAddressID, diffAddressErr := repository.GetOrCreateAddress(db, "0x7b1696677107e48b152e9bf400293e98b7d86eb1")
		Expect(diffAddressErr).NotTo(HaveOccurred())

		var calcResult test_helpers.VariableResWithAddress
		err = db.Get(&calcResult, `SELECT diff_id, header_id, address_id, calc AS value FROM maker.clip_calc`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(calcResult, clipCalcDiff.ID, header.Id, contractAddressID, strconv.FormatInt(diffAddressID, 10))
	})

	It("reads in a Clip Buf storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")
		value := common.HexToHash("0000000000000000000000000000000000000000043355B53628A6B594000000")
		clipBufDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipBufDiff)
		Expect(err).NotTo(HaveOccurred())

		var bufResult test_helpers.VariableRes
		err = db.Get(&bufResult, `SELECT diff_id, header_id, buf AS value FROM maker.clip_buf`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(bufResult, clipBufDiff.ID, header.Id, "1300000000000000000000000000")
	})

	It("reads in a Clip Tail storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006")
		value := common.HexToHash("00000000000000000000000000000000000000000000000000000000000020D0")
		clipTailDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipTailDiff)
		Expect(err).NotTo(HaveOccurred())

		var tailResult test_helpers.VariableRes
		err = db.Get(&tailResult, `SELECT diff_id, header_id, tail AS value FROM maker.clip_tail`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(tailResult, clipTailDiff.ID, header.Id, "8400")
	})

	It("reads in a Clip Cusp storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007")
		value := common.HexToHash("0000000000000000000000000000000000000000014ADF4B7320334B90000000")
		clipCuspDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipCuspDiff)
		Expect(err).NotTo(HaveOccurred())

		var cuspResult test_helpers.VariableRes
		err = db.Get(&cuspResult, `SELECT diff_id, header_id, cusp AS value FROM maker.clip_cusp`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(cuspResult, clipCuspDiff.ID, header.Id, "400000000000000000000000000")
	})

	It("reads in a Clip Chip and Tip packed storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008")
		value := common.HexToHash("00000000348c771b1de11359f9ee9b8d0c9380000000000000038d7ea4c68000")
		clipChipAndTipDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipChipAndTipDiff)
		Expect(err).NotTo(HaveOccurred())

		var chipResult test_helpers.VariableRes
		err = db.Get(&chipResult, `SELECT diff_id, header_id, chip AS value FROM maker.clip_chip`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(chipResult, clipChipAndTipDiff.ID, header.Id, "1000000000000000")

		var tipResult test_helpers.VariableRes
		err = db.Get(&tipResult, `SELECT diff_id, header_id, tip AS value FROM maker.clip_tip`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(tipResult, clipChipAndTipDiff.ID, header.Id, "300000000000000000000000000000000000000000000000")

	})

	It("reads in a Clip Chost storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009")
		value := common.HexToHash("000000000000000000000003DDAAC3295D6441C938631C35C22F400000000000")
		clipChostDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipChostDiff)
		Expect(err).NotTo(HaveOccurred())

		var chostResult test_helpers.VariableRes
		err = db.Get(&chostResult, `SELECT diff_id, header_id, chost AS value FROM maker.clip_chost`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(chostResult, clipChostDiff.ID, header.Id, "5650000000000000000000000000000000000000000000000")
	})

	It("reads in a Clip Kicks storage diff row and persists it", func() {
		key := common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000026")
		clipKicksDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipKicksDiff)
		Expect(err).NotTo(HaveOccurred())

		var kicksResult test_helpers.VariableRes
		err = db.Get(&kicksResult, `SELECT diff_id, header_id, kicks AS value FROM maker.clip_kicks`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(kicksResult, clipKicksDiff.ID, header.Id, "38")
	})

	It("reads the Clip Active storage diff row and persists it", func() {
		key := common.HexToHash("000000000000000000000000000000000000000000000000000000000000000b")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		clipActiveDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		err := transformer.Execute(clipActiveDiff)
		Expect(err).NotTo(HaveOccurred())

		var activeResult test_helpers.VariableRes
		err = db.Get(&activeResult, `SELECT diff_id, header_id, active AS value FROM maker.clip_active`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(activeResult, clipActiveDiff.ID, header.Id, "1")
	})

	It("reads the Clip Active storage diff row and persists the incremented value", func() {
		key := common.HexToHash("0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db9")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000032")
		clipActiveDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		activeErr := transformer.Execute(clipActiveDiff)
		Expect(activeErr).NotTo(HaveOccurred())

		var activeResult test_helpers.VariableRes
		dbErr := db.Get(&activeResult, `SELECT diff_id, header_id, sale_id AS value FROM maker.clip_active_sales`)
		Expect(dbErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(activeResult, clipActiveDiff.ID, header.Id, "50")
	})

	It("reads the Clip Active storage diff row and persists multiple values", func() {
		key := common.HexToHash("0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db9")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000032")
		clipActiveDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		activeErr := transformer.Execute(clipActiveDiff)
		Expect(activeErr).NotTo(HaveOccurred())

		keyTwo := common.HexToHash("0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db9")
		valueTwo := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000033")
		clipActiveDiffTwo := test_helpers.CreateDiffRecord(db, header, contractAddress, keyTwo, valueTwo)

		activeErrTwo := transformer.Execute(clipActiveDiffTwo)
		Expect(activeErrTwo).NotTo(HaveOccurred())

		var activeResult []test_helpers.VariableRes
		dbErr := db.Select(&activeResult, `SELECT diff_id, header_id, sale_id AS value FROM maker.clip_active_sales`)
		Expect(dbErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(activeResult[0], clipActiveDiff.ID, header.Id, "50")
		test_helpers.AssertVariable(activeResult[1], clipActiveDiff.ID+1, header.Id, "51")
	})

	//Data from storage diff on chain https://etherscan.io/tx/0x79f7ac4251a177bc2f86709319a65ff22c4eeb8c4eabe5ece4acabfb0a113b2a
	Describe("Sales", func() {
		BeforeEach(func() {
			clipKickLog := test_data.CreateTestLog(header.Id, db)
			clipKickModel := test_data.ClipKickModel()

			msgSenderAddressID, err := repository.GetOrCreateAddress(db, test_data.ClipLinkAV130Address())
			Expect(err).NotTo(HaveOccurred())

			clipKickModel.ColumnValues[event.HeaderFK] = header.Id
			clipKickModel.ColumnValues[event.LogFK] = clipKickLog.ID
			clipKickModel.ColumnValues[event.AddressFK] = msgSenderAddressID
			clipKickModel.ColumnValues[constants.SaleIDColumn] = "50"

			insertErr := event.PersistModels([]event.InsertionModel{clipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())
		})

		It("reads in a Sales Pos storage diff row and persists it", func() {
			key := common.HexToHash("74c83704300c65b1de76b9ee7537f3f330650a1d59eb262898de510c0c350be2")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
			clipSalesPosDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

			err := transformer.Execute(clipSalesPosDiff)
			Expect(err).NotTo(HaveOccurred())

			var salesPosResult test_helpers.VariableRes
			err = db.Get(&salesPosResult, `SELECT diff_id, header_id, pos AS value FROM maker.clip_sale_pos`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertVariable(salesPosResult, clipSalesPosDiff.ID, header.Id, "0")
		})

		It("reads in a Sales Tab storage diff row and persists it", func() {
			key := common.HexToHash("74c83704300c65b1de76b9ee7537f3f330650a1d59eb262898de510c0c350be3")
			value := common.HexToHash("00000000000000000000000d25f9b0930678426f13816b87849e36612b7adb51")
			clipSalesTabDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

			err := transformer.Execute(clipSalesTabDiff)
			Expect(err).NotTo(HaveOccurred())

			var salesTabResult test_helpers.VariableRes
			err = db.Get(&salesTabResult, `SELECT diff_id, header_id, tab AS value FROM maker.clip_sale_tab`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertVariable(salesTabResult, clipSalesTabDiff.ID, header.Id, "19216322211169085842425265706619325766990804212561")
		})

		It("reads in a Sales Lot storage diff row and persists it", func() {
			key := common.HexToHash("74c83704300c65b1de76b9ee7537f3f330650a1d59eb262898de510c0c350be4")
			value := common.HexToHash("00000000000000000000000000000000000000000000003183f290e991427b71")
			clipSalesLotDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

			err := transformer.Execute(clipSalesLotDiff)
			Expect(err).NotTo(HaveOccurred())

			var salesLotResult test_helpers.VariableRes
			err = db.Get(&salesLotResult, `SELECT diff_id, header_id, lot AS value FROM maker.clip_sale_lot`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertVariable(salesLotResult, clipSalesLotDiff.ID, header.Id, "913398280707939400561")
		})

		It("reads in a Sales usr and tic packed storage diff row and persists it", func() {
			key := common.HexToHash("74c83704300c65b1de76b9ee7537f3f330650a1d59eb262898de510c0c350be5")
			value := common.HexToHash("000000000000000061330b62b116df5da53d75e3670bdf13905c87008b7d0ad5")
			clipSalesUsrTicDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

			err := transformer.Execute(clipSalesUsrTicDiff)
			Expect(err).NotTo(HaveOccurred())

			var salesUsrResult test_helpers.VariableRes
			err = db.Get(&salesUsrResult, `SELECT diff_id, header_id, usr AS value FROM maker.clip_sale_usr`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertVariable(salesUsrResult, clipSalesUsrTicDiff.ID, header.Id, "0xB116dF5DA53d75e3670bDF13905c87008B7D0ad5")

			var salesTicResult test_helpers.VariableRes
			err = db.Get(&salesTicResult, `SELECT diff_id, header_id, tic AS value FROM maker.clip_sale_tic`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertVariable(salesTicResult, clipSalesUsrTicDiff.ID, header.Id, "1630735202")
		})

		It("reads in a Sales Top storage diff row and persists it", func() {
			key := common.HexToHash("74c83704300c65b1de76b9ee7537f3f330650a1d59eb262898de510c0c350be6")
			value := common.HexToHash("00000000000000000000000000000000000000008077b09bcd374e2b30940000")
			clipSalesTopDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

			err := transformer.Execute(clipSalesTopDiff)
			Expect(err).NotTo(HaveOccurred())

			var salesTopResult test_helpers.VariableRes
			err = db.Get(&salesTopResult, `SELECT diff_id, header_id, top AS value FROM maker.clip_sale_top`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertVariable(salesTopResult, clipSalesTopDiff.ID, header.Id, "39758777440200000000000000000")
		})
	})
})
