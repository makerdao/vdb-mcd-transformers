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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

// Returns a deep copy of the given model, so tests aren't getting the same map/slice references
func CopyModel(model shared.InsertionModel) shared.InsertionModel {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encErr := encoder.Encode(model)
	Expect(encErr).NotTo(HaveOccurred())

	var newModel shared.InsertionModel
	decoder := gob.NewDecoder(buf)
	decErr := decoder.Decode(&newModel)
	Expect(decErr).NotTo(HaveOccurred())
	return newModel
}

// TODO rename CopyEventModel => CopyModel
func CopyEventModel(model event.InsertionModel) event.InsertionModel {
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
func CreateTestLog(headerID int64, db *postgres.DB) core.HeaderSyncLog {
	log := types.Log{
		Address:     common.Address{},
		Topics:      nil,
		Data:        nil,
		BlockNumber: 0,
		TxHash:      common.Hash{},
		TxIndex:     uint(rand.Int31()),
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}
	headerSyncLogRepository := repositories.NewHeaderSyncLogRepository(db)
	insertLogsErr := headerSyncLogRepository.CreateHeaderSyncLogs(headerID, []types.Log{log})
	Expect(insertLogsErr).NotTo(HaveOccurred())
	headerSyncLogs, getLogsErr := headerSyncLogRepository.GetUntransformedHeaderSyncLogs()
	Expect(getLogsErr).NotTo(HaveOccurred())
	for _, headerSyncLog := range headerSyncLogs {
		if headerSyncLog.Log.TxIndex == log.TxIndex {
			return headerSyncLog
		}
	}
	panic("couldn't find inserted test log")
}

func CreateLogs(headerID int64, logs []types.Log, db *postgres.DB) []core.HeaderSyncLog {
	headerSyncLogRepository := repositories.NewHeaderSyncLogRepository(db)
	insertLogsErr := headerSyncLogRepository.CreateHeaderSyncLogs(headerID, logs)
	Expect(insertLogsErr).NotTo(HaveOccurred())
	headerSyncLogs, getLogsErr := headerSyncLogRepository.GetUntransformedHeaderSyncLogs()
	Expect(getLogsErr).NotTo(HaveOccurred())
	var results []core.HeaderSyncLog
	for _, headerSyncLog := range headerSyncLogs {
		for _, log := range logs {
			if headerSyncLog.Log.BlockNumber == log.BlockNumber && headerSyncLog.Log.TxIndex == log.TxIndex && headerSyncLog.Log.Index == log.Index {
				results = append(results, headerSyncLog)
			}
		}
	}
	return results
}
