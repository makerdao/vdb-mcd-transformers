package shared_behaviors

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type BidTriggerTestInput struct {
	Repository      vdbStorage.Repository
	Metadata        types.ValueMetadata
	ContractAddress string
	TriggerTable    string
	FieldTable      string
	ColumnName      event.ColumnName
	PackedValueType types.ValueType
}

func UpdateBidSnapshotTriggerTests(input BidTriggerTestInput) {
	Describe(fmt.Sprintf(`updating %s trigger table`, input.TriggerTable), func() {
		var (
			bidID int
			headerOne,
			headerTwo,
			headerThree core.Header
			timestampOne,
			timestampTwo,
			timestampThree string
			diffID,
			addressID int64
			repo             = input.Repository
			db               = test_config.NewTestDB(test_config.NewTestNode())
			hashOne          = common.BytesToHash([]byte{1, 2, 3, 4, 5})
			hashTwo          = common.BytesToHash([]byte{5, 4, 3, 2, 1})
			hashThree        = common.BytesToHash([]byte{6, 7, 8, 9, 0})
			getFieldQuery    = fmt.Sprintf(`SELECT "%s" FROM maker.%s ORDER BY block_number`, input.ColumnName, input.TriggerTable)
			insertFieldQuery = fmt.Sprintf(`INSERT INTO maker.%s (address_id, block_number, bid_id, "%s", updated) VALUES ($1, $2, $3, $4, $5)`, input.TriggerTable, input.ColumnName)
			deleteRowQuery   = fmt.Sprintf(`DELETE FROM maker.%s WHERE header_id = $1`, input.FieldTable)
		)

		BeforeEach(func() {
			test_config.CleanTestDB(db)
			repo.SetDB(db)
			blockOne := rand.Int()
			blockTwo := blockOne + 1
			blockThree := blockTwo + 1
			rawTimestampOne := int64(rand.Int31())
			rawTimestampTwo := rawTimestampOne + 1
			rawTimestampThree := rawTimestampTwo + 1
			timestampOne = FormatTimestamp(rawTimestampOne)
			timestampTwo = FormatTimestamp(rawTimestampTwo)
			timestampThree = FormatTimestamp(rawTimestampTwo)
			headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, db)
			headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, db)
			headerThree = CreateHeaderWithHash(hashThree.String(), rawTimestampThree, blockThree, db)
			diffID = CreateFakeDiffRecord(db)
			var parseErr error
			bidID, parseErr = strconv.Atoi(input.Metadata.Keys[constants.BidId])
			Expect(parseErr).NotTo(HaveOccurred())
			var addressErr error
			addressID, addressErr = shared.GetOrCreateAddress(input.ContractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
		})

		Describe("updating history with new value", func() {
			It(fmt.Sprintf("updates field in subsequent blocks until next %s diff", input.FieldTable), func() {
				_, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				_, setupErr := db.Exec(insertFieldQuery, addressID, headerTwo.BlockNumber, bidID, initialColumnVal, timestampTwo)
				Expect(setupErr).NotTo(HaveOccurred())

				newRepoVal, newColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				err := repo.Create(diffID, headerOne.Id, input.Metadata, newRepoVal)
				Expect(err).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(valueHistory)).To(Equal(2))
				Expect(valueHistory[0].String).To(Equal(newColumnVal))
				Expect(valueHistory[1].String).To(Equal(newColumnVal))
			})

			It("ignores rows from blocks after the next time the field is updated", func() {
				initialRepoVal, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErr := repo.Create(diffID, headerTwo.Id, input.Metadata, initialRepoVal)
				Expect(setupErr).NotTo(HaveOccurred())

				newRepoVal, newColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				err := repo.Create(diffID, headerOne.Id, input.Metadata, newRepoVal)
				Expect(err).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(valueHistory)).To(Equal(2))
				Expect(valueHistory[0].String).To(Equal(newColumnVal))
				Expect(valueHistory[1].String).To(Equal(initialColumnVal))
			})

			It("ignores rows from earlier blocks", func() {
				_, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				_, setupErr := db.Exec(insertFieldQuery, addressID, headerOne.BlockNumber, bidID, initialColumnVal, timestampOne)
				Expect(setupErr).NotTo(HaveOccurred())

				newRepoVal, newColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				err := repo.Create(diffID, headerTwo.Id, input.Metadata, newRepoVal)
				Expect(err).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(valueHistory)).To(Equal(2))
				Expect(valueHistory[0].String).To(Equal(initialColumnVal))
				Expect(valueHistory[1].String).To(Equal(newColumnVal))
			})

			It("ignores rows from different address", func() {
				_, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				differentAddressID, addressErr := shared.GetOrCreateAddress(test_data.RandomString(40), db)
				Expect(addressErr).NotTo(HaveOccurred())
				_, setupErr := db.Exec(insertFieldQuery, differentAddressID, headerTwo.BlockNumber, bidID, initialColumnVal, timestampTwo)
				Expect(setupErr).NotTo(HaveOccurred())

				newRepoVal, newColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				err := repo.Create(diffID, headerOne.Id, input.Metadata, newRepoVal)
				Expect(err).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(valueHistory)).To(Equal(2))
				Expect(valueHistory[0].String).To(Equal(newColumnVal))
				Expect(valueHistory[1].String).To(Equal(initialColumnVal))
			})

			It("ignores rows with different bidID", func() {
				_, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				differentBidID := rand.Int()
				_, setupErr := db.Exec(insertFieldQuery, addressID, headerTwo.BlockNumber, differentBidID, initialColumnVal, timestampTwo)
				Expect(setupErr).NotTo(HaveOccurred())

				newRepoVal, newColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				err := repo.Create(diffID, headerOne.Id, input.Metadata, newRepoVal)
				Expect(err).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(valueHistory)).To(Equal(2))
				Expect(valueHistory[0].String).To(Equal(newColumnVal))
				Expect(valueHistory[1].String).To(Equal(initialColumnVal))
			})
		})

		Describe("when diff is deleted", func() {
			It("updates field to previous value in subsequent rows", func() {
				initialRepoVal, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialRepoVal)
				Expect(setupErrOne).NotTo(HaveOccurred())

				subsequentRepoVal, subsequentColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErrTwo := repo.Create(diffID, headerTwo.Id, input.Metadata, subsequentRepoVal)
				Expect(setupErrTwo).NotTo(HaveOccurred())
				_, setupErrThree := db.Exec(insertFieldQuery, addressID, headerThree.BlockNumber, bidID, subsequentColumnVal, timestampThree)
				Expect(setupErrThree).NotTo(HaveOccurred())

				_, deleteErr := db.Exec(deleteRowQuery, headerTwo.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(valueHistory[0].String).To(Equal(initialColumnVal))
				Expect(valueHistory[1].String).To(Equal(initialColumnVal))
			})

			It("sets field in subsequent rows to null if no previous diff exists", func() {
				initialRepoVal, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialRepoVal)
				Expect(setupErrOne).NotTo(HaveOccurred())
				_, setupErrTwo := db.Exec(insertFieldQuery, addressID, headerTwo.BlockNumber, bidID, initialColumnVal, timestampTwo)
				Expect(setupErrTwo).NotTo(HaveOccurred())

				_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(valueHistory[0].Valid).To(BeFalse())
			})

			It("deletes bid state associated with diff if identical to previous state", func() {
				initialRepoVal, initialColumnVal := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialRepoVal)
				Expect(setupErrOne).NotTo(HaveOccurred())

				subsequentRepoVal, _ := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErrTwo := repo.Create(diffID, headerTwo.Id, input.Metadata, subsequentRepoVal)
				Expect(setupErrTwo).NotTo(HaveOccurred())
				_, setupErrThree := db.Exec(insertFieldQuery, addressID, headerThree.BlockNumber, bidID, initialColumnVal, timestampThree)
				Expect(setupErrThree).NotTo(HaveOccurred())

				_, deleteErr := db.Exec(deleteRowQuery, headerTwo.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(valueHistory)).To(Equal(2))
				Expect(valueHistory[0].String).To(Equal(initialColumnVal))
				Expect(valueHistory[1].String).To(Equal(initialColumnVal))
			})

			It("deletes bid state associated with diff if it's the earliest state in the table", func() {
				initialRepoVal, _ := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
				setupErr := repo.Create(diffID, headerOne.Id, input.Metadata, initialRepoVal)
				Expect(setupErr).NotTo(HaveOccurred())

				_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var valueHistory []sql.NullString
				queryErr := db.Select(&valueHistory, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(valueHistory).To(BeEmpty())
			})
		})
	})
}

func InsertFlipBidSnapshotTriggerTests(input BidTriggerTestInput) {
	Describe("inserting a new field", func() {
		var (
			bidID                int
			headerOne, headerTwo core.Header
			diffID, addressID    int64
			repo                 = input.Repository
			db                   = test_config.NewTestDB(test_config.NewTestNode())
			hashOne              = common.BytesToHash([]byte{1, 2, 3, 4, 5})
			hashTwo              = common.BytesToHash([]byte{5, 4, 3, 2, 1})
			getFlipStateQuery    = `SELECT address_id, block_number, bid_id, guy, tic, "end", lot, bid, usr, gal, tab, updated FROM maker.flip ORDER BY block_number`
		)

		BeforeEach(func() {
			test_config.CleanTestDB(db)
			repo.SetDB(db)
			blockOne := rand.Int()
			blockTwo := blockOne + 1
			rawTimestampOne := int64(rand.Int31())
			rawTimestampTwo := rawTimestampOne + 1
			headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, db)
			headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, db)
			diffID = CreateFakeDiffRecord(db)
			var parseErr error
			bidID, parseErr = strconv.Atoi(input.Metadata.Keys[constants.BidId])
			Expect(parseErr).NotTo(HaveOccurred())
			var addressErr error
			addressID, addressErr = shared.GetOrCreateAddress(input.ContractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
		})

		It("inserts a row for new bid + blockNumber combination", func() {
			initialBidValues := test_helpers.GetFlipStorageValues(0, test_helpers.FakeIlk.Hex, bidID)
			test_helpers.InsertValues(db, input.Repository, headerOne, initialBidValues,
				test_helpers.GetFlipMetadatas(strconv.Itoa(bidID)))

			updatedRepoVal, _ := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
			err := input.Repository.Create(diffID, headerTwo.Id, input.Metadata, updatedRepoVal)
			Expect(err).NotTo(HaveOccurred())

			var bidSnapshots []flipBidSnapshot
			queryErr := db.Select(&bidSnapshots, getFlipStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(bidSnapshots)).To(Equal(2))
			initialBid := flipBidSnapshotFromValues(bidID, headerOne.BlockNumber, addressID, headerOne.Timestamp,
				initialBidValues)
			initialBidValues[input.Metadata.Name] = updatedRepoVal
			expectedBid := flipBidSnapshotFromValues(bidID, headerTwo.BlockNumber, addressID, headerTwo.Timestamp,
				initialBidValues)
			assertFlipSnapshot(bidSnapshots[0], initialBid, headerOne.BlockNumber)
			assertSingleField(bidSnapshots[0], initialBid, input.ColumnName)
			assertFlipSnapshot(bidSnapshots[1], expectedBid, headerTwo.BlockNumber)
			assertSingleField(bidSnapshots[1], expectedBid, input.ColumnName)
		})

		It("updates record if bid record already exists for block", func() {
			initialBidValues := test_helpers.GetFlipStorageValues(0, test_helpers.FakeIlk.Hex, bidID)
			test_helpers.InsertValues(db, repo, headerOne, initialBidValues,
				test_helpers.GetFlipMetadatas(strconv.Itoa(bidID)))

			updatedRepoVal, _ := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
			err := input.Repository.Create(diffID, headerOne.Id, input.Metadata, updatedRepoVal)
			Expect(err).NotTo(HaveOccurred())

			var bidSnapshots []flipBidSnapshot
			queryErr := db.Select(&bidSnapshots, getFlipStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(bidSnapshots)).To(Equal(1))
			initialBidValues[input.Metadata.Name] = updatedRepoVal
			expectedBid := flipBidSnapshotFromValues(bidID, headerTwo.BlockNumber, addressID, headerOne.Timestamp,
				initialBidValues)
			assertFlipSnapshot(bidSnapshots[0], expectedBid, headerOne.BlockNumber)
			assertSingleField(bidSnapshots[0], expectedBid, input.ColumnName)
		})
	})
}

func InsertBidSnapshotTriggerTests(input BidTriggerTestInput) {
	Describe("inserting a new field", func() {
		var (
			bidID int
			headerOne,
			headerTwo core.Header
			diffID,
			addressID int64
			repo             = input.Repository
			db               = test_config.NewTestDB(test_config.NewTestNode())
			hashOne          = common.BytesToHash([]byte{1, 2, 3, 4, 5})
			hashTwo          = common.BytesToHash([]byte{5, 4, 3, 2, 1})
			getBidStateQuery = fmt.Sprintf(`SELECT address_id, block_number, bid_id, guy, tic, "end", lot, bid, updated FROM maker.%s ORDER BY block_number`, input.TriggerTable)
		)

		BeforeEach(func() {
			test_config.CleanTestDB(db)
			repo.SetDB(db)
			blockOne := rand.Int()
			blockTwo := blockOne + 1
			rawTimestampOne := int64(rand.Int31())
			rawTimestampTwo := rawTimestampOne + 1
			headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, db)
			headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, db)
			diffID = CreateFakeDiffRecord(db)
			var parseErr error
			bidID, parseErr = strconv.Atoi(input.Metadata.Keys[constants.BidId])
			Expect(parseErr).NotTo(HaveOccurred())
			var addressErr error
			addressID, addressErr = shared.GetOrCreateAddress(input.ContractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
		})

		It("inserts a row for new bid + blockNumber combination", func() {
			initialBidValues := test_helpers.GetCommonBidStorageValues(0, bidID)
			test_helpers.InsertValues(db, input.Repository, headerOne, initialBidValues,
				test_helpers.GetCommonBidMetadatas(strconv.Itoa(bidID)))

			updatedRepoVal, _ := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
			err := input.Repository.Create(diffID, headerTwo.Id, input.Metadata, updatedRepoVal)
			Expect(err).NotTo(HaveOccurred())

			var bidSnapshots []bidSnapshot
			queryErr := db.Select(&bidSnapshots, getBidStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(bidSnapshots)).To(Equal(2))
			initialBid := bidSnapshotFromValues(bidID, headerOne.BlockNumber, addressID, headerOne.Timestamp,
				initialBidValues)
			initialBidValues[input.Metadata.Name] = updatedRepoVal
			expectedBid := bidSnapshotFromValues(bidID, headerTwo.BlockNumber, addressID, headerTwo.Timestamp,
				initialBidValues)
			assertBidSnapshot(bidSnapshots[0], initialBid, headerOne.BlockNumber)
			assertSingleField(bidSnapshots[0], initialBid, input.ColumnName)
			assertBidSnapshot(bidSnapshots[1], expectedBid, headerTwo.BlockNumber)
			assertSingleField(bidSnapshots[1], expectedBid, input.ColumnName)
		})

		It("updates record if bid record already exists for block", func() {
			initialBidValues := test_helpers.GetCommonBidStorageValues(0, bidID)
			test_helpers.InsertValues(db, repo, headerOne, initialBidValues,
				test_helpers.GetCommonBidMetadatas(strconv.Itoa(bidID)))

			updatedRepoVal, _ := randomBidStorageValue(input.Metadata.Type, input.PackedValueType)
			err := input.Repository.Create(diffID, headerOne.Id, input.Metadata, updatedRepoVal)
			Expect(err).NotTo(HaveOccurred())

			var bidSnapshots []bidSnapshot
			queryErr := db.Select(&bidSnapshots, getBidStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(bidSnapshots)).To(Equal(1))
			initialBidValues[input.Metadata.Name] = updatedRepoVal
			expectedBid := bidSnapshotFromValues(bidID, headerTwo.BlockNumber, addressID, headerOne.Timestamp,
				initialBidValues)
			assertBidSnapshot(bidSnapshots[0], expectedBid, headerOne.BlockNumber)
			assertSingleField(bidSnapshots[0], expectedBid, input.ColumnName)
		})
	})
}

func randomBidStorageValue(valueType types.ValueType, packedValueType types.ValueType) (interface{}, string) {
	var repoVal interface{}
	var columnVal string
	var err error
	switch valueType {
	case types.Address:
		val := "0x" + test_data.RandomString(40)
		repoVal = val
		columnVal = val
	case types.Uint256, types.Uint48:
		val := strconv.Itoa(rand.Int())
		repoVal = val
		columnVal = val
	case types.PackedSlot:
		if packedValueType == types.Uint48 {
			columnVal = strconv.Itoa(rand.Int())
		} else {
			columnVal = "0x" + test_data.RandomString(40)
		}
		repoVal = map[int]string{0: columnVal}
	default:
		err = errors.New("ValueType not implemented")
	}
	Expect(err).NotTo(HaveOccurred())
	return repoVal, columnVal
}

type bidSnapshot struct {
	AddressID   string `db:"address_id"`
	BlockNumber string `db:"block_number"`
	BidID       string `db:"bid_id"`
	Guy         sql.NullString
	Tic         sql.NullString
	End         sql.NullString
	Lot         sql.NullString
	Bid         sql.NullString
	Updated     string
}

type flipBidSnapshot struct {
	bidSnapshot
	Usr sql.NullString
	Gal sql.NullString
	Tab sql.NullString
}

func bidSnapshotFromValues(bidID int, blockNumber, addressID int64, updated string, bidValues map[string]interface{}) bidSnapshot {
	parsedUpdated, parseErr := strconv.ParseInt(updated, 10, 64)
	Expect(parseErr).NotTo(HaveOccurred())
	packedValues := bidValues[mcdStorage.Packed].(map[int]string)

	return bidSnapshot{
		AddressID:   strconv.FormatInt(addressID, 10),
		BlockNumber: strconv.FormatInt(blockNumber, 10),
		BidID:       strconv.Itoa(bidID),
		Guy:         test_helpers.GetValidNullString(packedValues[0]),
		Tic:         test_helpers.GetValidNullString(packedValues[1]),
		End:         test_helpers.GetValidNullString(packedValues[2]),
		Lot:         test_helpers.GetValidNullString(bidValues[mcdStorage.BidLot].(string)),
		Bid:         test_helpers.GetValidNullString(bidValues[mcdStorage.BidBid].(string)),
		Updated:     FormatTimestamp(parsedUpdated),
	}
}

func flipBidSnapshotFromValues(bidID int, blockNumber, addressID int64, updated string, bidValues map[string]interface{}) flipBidSnapshot {
	return flipBidSnapshot{
		bidSnapshot: bidSnapshotFromValues(bidID, blockNumber, addressID, updated, bidValues),
		Usr:         test_helpers.GetValidNullString(bidValues[mcdStorage.BidUsr].(string)),
		Gal:         test_helpers.GetValidNullString(bidValues[mcdStorage.BidGal].(string)),
		Tab:         test_helpers.GetValidNullString(bidValues[mcdStorage.BidTab].(string)),
	}
}

func assertBidSnapshot(actualBid, expectedBid bidSnapshot, expectedBlockNumber int64) {
	Expect(actualBid.AddressID).To(Equal(expectedBid.AddressID))
	Expect(actualBid.BlockNumber).To(Equal(strconv.FormatInt(expectedBlockNumber, 10)))
	Expect(actualBid.BidID).To(Equal(expectedBid.BidID))
	Expect(actualBid.Lot).To(Equal(expectedBid.Lot))
	Expect(actualBid.Bid).To(Equal(expectedBid.Bid))
	Expect(actualBid.Updated).To(Equal(expectedBid.Updated))
}

func assertFlipSnapshot(actualFlip, expectedFlip flipBidSnapshot, expectedBlockNumber int64) {
	assertBidSnapshot(actualFlip.bidSnapshot, expectedFlip.bidSnapshot, expectedBlockNumber)
	Expect(actualFlip.Usr).To(Equal(expectedFlip.Usr))
	Expect(actualFlip.Gal).To(Equal(expectedFlip.Gal))
	Expect(actualFlip.Tab).To(Equal(expectedFlip.Tab))
}

func assertSingleField(actualFlip, expectedFlip interface{}, fieldName event.ColumnName) {
	Expect(getBidProperty(actualFlip, string(fieldName))).To(Equal(getBidProperty(expectedFlip, string(fieldName))))
}

func getBidProperty(bid interface{}, fieldName string) string {
	r := reflect.ValueOf(bid)
	property := reflect.Indirect(r).FieldByName(fieldName)
	return property.String()
}
