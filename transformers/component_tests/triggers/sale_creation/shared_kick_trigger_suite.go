package sale_creation

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func SharedSaleCreationTriggerTests(tableName, contractAddress string, kickModel *event.InsertionModel) {
	Describe("Updating sale table", func() {
		var (
			blockOne,
			blockTwo int
			headerOne,
			headerTwo core.Header
			rawTimestampOne,
			rawTimestampTwo int64
			saleID              string
			logID               int64
			db                  = test_config.NewTestDB(test_config.NewTestNode())
			getTimeCreatedQuery = fmt.Sprintf(`SELECT created FROM maker.%s ORDER BY block_number`, tableName)
			insertEmptyRowQuery = fmt.Sprintf(`INSERT INTO maker.%s (block_number, sale_id, address_id, updated) VALUES ($1, $2, $3, $4)`, tableName)
			//deleteRowQuery      = fmt.Sprintf(`DELETE FROM maker.%s_kick WHERE header_id = $1`, tableName)
		)

		BeforeEach(func() {
			test_config.CleanTestDB(db)
			blockOne = rand.Int()
			blockTwo = blockOne + 1
			rawTimestampOne = int64(rand.Int31())
			rawTimestampTwo = rawTimestampOne + 1
			headerOne = CreateHeader(rawTimestampOne, blockOne, db)
			headerTwo = CreateHeader(rawTimestampTwo, blockTwo, db)
			logID = test_data.CreateTestLog(headerOne.Id, db).ID
			saleID = strconv.Itoa(rand.Int())
			Expect(saleID).NotTo(Equal(10))
			kickModel.ColumnValues = updateKickKeys(*kickModel, headerOne.Id, logID, contractAddress, saleID, db)
		})

		It("updates time created of all records for a sale", func() {
			_, setupErr := db.Exec(insertEmptyRowQuery, headerTwo.BlockNumber, saleID,
				kickModel.ColumnValues[event.AddressFK], FormatTimestamp(rawTimestampTwo))
			Expect(setupErr).NotTo(HaveOccurred())
			expectedTimeCreated := sql.NullString{Valid: true, String: FormatTimestamp(rawTimestampOne)}

			kickErr := event.PersistModels([]event.InsertionModel{*kickModel}, db)
			Expect(kickErr).NotTo(HaveOccurred())

			var creationTimes []sql.NullString
			queryErr := db.Select(&creationTimes, getTimeCreatedQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(creationTimes)).To(Equal(1))
			Expect(creationTimes[0]).To(Equal(expectedTimeCreated))
		})

		/*It("does not update records from a different contract", func() {
			randomAddressID, addressErr := repository.GetOrCreateAddress(db, test_data.RandomString(40))
			Expect(addressErr).NotTo(HaveOccurred())
			_, setupErr := db.Exec(insertEmptyRowQuery, headerTwo.BlockNumber,
				kickModel.ColumnValues[constants.BidIDColumn], randomAddressID, FormatTimestamp(rawTimestampTwo))
			Expect(setupErr).NotTo(HaveOccurred())

			kickErr := event.PersistModels([]event.InsertionModel{*kickModel}, db)
			Expect(kickErr).NotTo(HaveOccurred())

			var creationTimes []sql.NullString
			queryErr := db.Select(&creationTimes, getTimeCreatedQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(creationTimes)).To(Equal(1))
			Expect(creationTimes[0].Valid).To(BeFalse())
		})

		It("does not update records with a different bid_id", func() {
			randomBidID := strconv.Itoa(rand.Int())
			_, setupErr := db.Exec(insertEmptyRowQuery, headerTwo.BlockNumber, randomBidID,
				kickModel.ColumnValues[event.AddressFK], FormatTimestamp(rawTimestampTwo))
			Expect(setupErr).NotTo(HaveOccurred())

			kickErr := event.PersistModels([]event.InsertionModel{*kickModel}, db)
			Expect(kickErr).NotTo(HaveOccurred())

			var creationTimes []sql.NullString
			queryErr := db.Select(&creationTimes, getTimeCreatedQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(creationTimes)).To(Equal(1))
			Expect(creationTimes[0].Valid).To(BeFalse())
		})

		It("sets created to null when record is deleted", func() {
			_, setupErr := db.Exec(insertEmptyRowQuery, headerOne.BlockNumber, saleID,
				kickModel.ColumnValues[event.AddressFK], FormatTimestamp(rawTimestampOne))
			Expect(setupErr).NotTo(HaveOccurred())
			kickErr := event.PersistModels([]event.InsertionModel{*kickModel}, db)
			Expect(kickErr).NotTo(HaveOccurred())

			_, err := db.Exec(deleteRowQuery, kickModel.ColumnValues[event.HeaderFK])
			Expect(err).NotTo(HaveOccurred())

			var creationTimes []sql.NullString
			queryErr := db.Select(&creationTimes, getTimeCreatedQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(creationTimes[0].Valid).To(BeFalse())
		}) */
	})
}

func updateKickKeys(kickModel event.InsertionModel, headerID, logID int64, contractAddress, saleID string, db *postgres.DB) event.ColumnValues {
	addressID, addressErr := repository.GetOrCreateAddress(db, contractAddress)
	Expect(addressErr).NotTo(HaveOccurred())

	kickModel.ColumnValues[event.HeaderFK] = headerID
	kickModel.ColumnValues[event.LogFK] = logID
	kickModel.ColumnValues[event.AddressFK] = addressID
	kickModel.ColumnValues[constants.SaleIDColumn] = saleID
	return kickModel.ColumnValues
}
