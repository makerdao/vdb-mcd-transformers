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

package test_data

import (
	"bytes"
	"encoding/gob"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/gomega"
)

// Returns a deep copy of the given model, so tests aren't getting the same map/slice references
func CopyModel(model event.InsertionModel) event.InsertionModel {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encErr := encoder.Encode(model)
	Expect(encErr).NotTo(HaveOccurred())

	var newModel event.InsertionModel
	decoder := gob.NewDecoder(buf)
	decErr := decoder.Decode(&newModel)
	Expect(decErr).NotTo(HaveOccurred())
	return newModel
}

func AssertDBRecordCount(db *postgres.DB, dbTable string, expectedCount int) {
	var count int
	query := `SELECT count(*) FROM ` + dbTable
	err := db.QueryRow(query).Scan(&count)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(Equal(expectedCount))
}

// Create a header to reference in an event, returning headerID
func CreateTestHeader(db *postgres.DB) int64 {
	headerRepository := repositories.NewHeaderRepository(db)
	headerID, insertHeaderErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
	Expect(insertHeaderErr).NotTo(HaveOccurred())
	return headerID
}

// Create a header sync log to reference in an event, returning inserted header sync log
func CreateTestLog(headerID int64, db *postgres.DB) core.EventLog {
	var blockNumber uint64
	err := db.Get(&blockNumber, `SELECT block_number FROM public.headers WHERE id = $1`, headerID)
	Expect(err).NotTo(HaveOccurred())
	log := types.Log{
		Address:     common.Address{},
		Topics:      nil,
		Data:        nil,
		BlockNumber: blockNumber,
		TxHash:      common.HexToHash("0x" + RandomString(64)),
		TxIndex:     uint(rand.Int31()),
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}
	headerRepo := repositories.NewHeaderRepository(db)
	test_data.CreateMatchingTx(log, headerID, headerRepo)
	eventLogRepository := repositories.NewEventLogRepository(db)
	insertLogsErr := eventLogRepository.CreateEventLogs(headerID, []types.Log{log})
	Expect(insertLogsErr).NotTo(HaveOccurred())

	logCount := getLogCount(db)
	eventLogs, getLogsErr := eventLogRepository.GetUntransformedEventLogs(0, logCount)
	Expect(getLogsErr).NotTo(HaveOccurred())
	for _, EventLog := range eventLogs {
		if EventLog.Log.TxIndex == log.TxIndex {
			return EventLog
		}
	}
	panic("couldn't find inserted test log")
}

func CreateTestLogFromEventLog(headerID int64, log types.Log, db *postgres.DB) core.EventLog {
	return CreateLogs(headerID, []types.Log{log}, db)[0]
}

func CreateLogs(headerID int64, logs []types.Log, db *postgres.DB) []core.EventLog {
	headerRepo := repositories.NewHeaderRepository(db)
	for _, log := range logs {
		test_data.CreateMatchingTx(log, headerID, headerRepo)
	}
	eventLogRepository := repositories.NewEventLogRepository(db)
	insertLogsErr := eventLogRepository.CreateEventLogs(headerID, logs)
	Expect(insertLogsErr).NotTo(HaveOccurred())

	logCount := getLogCount(db)
	eventLogs, getLogsErr := eventLogRepository.GetUntransformedEventLogs(0, logCount)
	Expect(getLogsErr).NotTo(HaveOccurred())
	var results []core.EventLog
	for _, EventLog := range eventLogs {
		for _, log := range logs {
			if EventLog.Log.BlockNumber == log.BlockNumber && EventLog.Log.TxIndex == log.TxIndex && EventLog.Log.Index == log.Index {
				results = append(results, EventLog)
			}
		}
	}
	return results
}

func getLogCount(db *postgres.DB) int {
	var logCount int
	logCountErr := db.Get(&logCount, `SELECT count(*) from public.event_logs`)
	Expect(logCountErr).NotTo(HaveOccurred())

	return logCount
}

func AssignMessageSenderID(log core.EventLog, insertionModel event.InsertionModel, db *postgres.DB) {
	Expect(len(log.Log.Topics)).Should(BeNumerically(">=", 2))
	msgSenderID, msgSenderErr := repository.GetOrCreateAddress(db, log.Log.Topics[1].Hex())
	Expect(msgSenderErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[constants.MsgSenderColumn] = msgSenderID
}

func AssignAddressID(log core.EventLog, insertionModel event.InsertionModel, db *postgres.DB) {
	addressID, addressIDErr := repository.GetOrCreateAddress(db, log.Log.Address.Hex())
	Expect(addressIDErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[event.AddressFK] = addressID
}

func AssignUsrID(log core.EventLog, insertionModel event.InsertionModel, db *postgres.DB) {
	UsrID, UsrIDErr := repository.GetOrCreateAddress(db, log.Log.Topics[1].Hex())
	Expect(UsrIDErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[constants.UsrColumn] = UsrID
}

func AssignUrnID(insertionModel event.InsertionModel, db *postgres.DB) {
	var urnID int64
	urnErr := db.Get(&urnID, `SELECT id FROM maker.urns`)
	Expect(urnErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[constants.UrnColumn] = urnID
}

func AssignIlkID(insertionModel event.InsertionModel, db *postgres.DB) {
	var ilkID int64
	ilkErr := db.Get(&ilkID, `SELECT id FROM maker.ilks`)
	Expect(ilkErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[constants.IlkColumn] = ilkID
}

func AssignClipAddressID(clipAddressHex string, insertionModel event.InsertionModel, db *postgres.DB) {
	clipAddressID, clipAddressErr := repository.GetOrCreateAddress(db, clipAddressHex)
	Expect(clipAddressErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[constants.ClipColumn] = clipAddressID
}

func AssignDataAddressID(dataAddressHex string, insertionModel event.InsertionModel, db *postgres.DB) {
	dataAddressID, dataAddressErr := repository.GetOrCreateAddress(db, dataAddressHex)
	Expect(dataAddressErr).NotTo(HaveOccurred())
	insertionModel.ColumnValues[constants.DataColumn] = dataAddressID
}
