// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package trigger_test

import (
	"database/sql"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updating flip table", func() {
	var (
		blockOne,
		blockTwo int
		headerOne,
		headerTwo core.Header
		rawTimestampOne,
		rawTimestampTwo int64
		bidID               string
		logID               int64
		flipKickModel       event.InsertionModel
		db                  = test_config.NewTestDB(test_config.NewTestNode())
		getTimeCreatedQuery = `SELECT created FROM maker.flip ORDER BY block_number`
		insertEmptyRowQuery = `INSERT INTO maker.flip (block_number, bid_id, address_id, updated) VALUES ($1, $2, $3, $4)`
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
		bidID = strconv.Itoa(rand.Int())
		flipKickModel = createFlipKickModel(headerOne.Id, logID, test_data.EthFlipAddress(), bidID, db)
	})

	It("updates time created of all records for a bid", func() {
		_, setupErr := db.Exec(insertEmptyRowQuery, headerTwo.BlockNumber, bidID,
			flipKickModel.ColumnValues[event.AddressFK], FormatTimestamp(rawTimestampTwo))
		Expect(setupErr).NotTo(HaveOccurred())
		expectedTimeCreated := sql.NullString{Valid: true, String: FormatTimestamp(rawTimestampOne)}

		kickErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
		Expect(kickErr).NotTo(HaveOccurred())

		var flipStates []flipState
		queryErr := db.Select(&flipStates, getTimeCreatedQuery)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(len(flipStates)).To(Equal(1))
		Expect(flipStates[0].Created).To(Equal(expectedTimeCreated))
	})

	It("does not update records from a different contract", func() {
		randomAddressID, addressErr := shared.GetOrCreateAddress(test_data.RandomString(40), db)
		Expect(addressErr).NotTo(HaveOccurred())
		_, setupErr := db.Exec(insertEmptyRowQuery, headerTwo.BlockNumber,
			flipKickModel.ColumnValues[constants.BidIDColumn], randomAddressID, FormatTimestamp(rawTimestampTwo))
		Expect(setupErr).NotTo(HaveOccurred())

		kickErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
		Expect(kickErr).NotTo(HaveOccurred())

		var flipStates []flipState
		queryErr := db.Select(&flipStates, getTimeCreatedQuery)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(len(flipStates)).To(Equal(1))
		Expect(flipStates[0].Created.Valid).To(BeFalse())
	})

	It("does not update records with a different bid_id", func() {
		randomBidID := strconv.Itoa(rand.Int())
		_, setupErr := db.Exec(insertEmptyRowQuery, headerTwo.BlockNumber, randomBidID,
			flipKickModel.ColumnValues[event.AddressFK], FormatTimestamp(rawTimestampTwo))
		Expect(setupErr).NotTo(HaveOccurred())

		kickErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
		Expect(kickErr).NotTo(HaveOccurred())

		var flipStates []flipState
		queryErr := db.Select(&flipStates, getTimeCreatedQuery)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(len(flipStates)).To(Equal(1))
		Expect(flipStates[0].Created.Valid).To(BeFalse())
	})
})

func createFlipKickModel(headerID, logID int64, contractAddress, bidID string, db *postgres.DB) event.InsertionModel {
	addressID, addressErr := shared.GetOrCreateAddress(contractAddress, db)
	Expect(addressErr).NotTo(HaveOccurred())

	vatInit := test_data.FlipKickModel()
	vatInit.ColumnValues[event.HeaderFK] = headerID
	vatInit.ColumnValues[event.LogFK] = logID
	vatInit.ColumnValues[event.AddressFK] = addressID
	vatInit.ColumnValues[constants.BidIDColumn] = bidID
	return vatInit
}

type flipState struct {
	BlockNumber string `db:"block_number"`
	AddressID   string `db:"address_id"`
	BidID       string `db:"bid_id"`
	Guy         string
	Tic         string
	End         string
	Lot         string
	Bid         string
	Gal         string
	Tab         string
	Created     sql.NullString
	Updated     string
}
