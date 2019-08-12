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

package flap_kick_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("Flap Kick Repository", func() {
	Describe("Create", func() {
		var (
			db                 *postgres.DB
			flapKickRepository flap_kick.FlapKickRepository
			headerID, logID    int64
			model              flap_kick.FlapKickModel
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			flapKickRepository = flap_kick.FlapKickRepository{}
			flapKickRepository.SetDB(db)
			headerID = test_data.CreateTestHeader(db)
			flapKickLogs := test_data.CreateLogs(headerID, []types.Log{test_data.FlapKickHeaderSyncLog.Log}, db)
			Expect(len(flapKickLogs)).To(Equal(1))
			logID = flapKickLogs[0].ID
			model = test_data.FlapKickModel
			model.HeaderID = headerID
			model.LogID = logID
		})

		It("persists a flap kick record", func() {
			insertErr := flapKickRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			test_data.AssertDBRecordCount(db, "maker.flap_kick", 1)
			test_data.AssertDBRecordCount(db, "addresses", 1)

			var addressId string
			addressErr := db.Get(&addressId, `SELECT id FROM addresses`)
			Expect(addressErr).NotTo(HaveOccurred())

			var dbResult flap_kick.FlapKickModel
			getErr := db.Get(&dbResult, `SELECT log_id, bid, bid_id, lot, address_id FROM maker.flap_kick WHERE header_id = $1`, headerID)

			Expect(getErr).NotTo(HaveOccurred())
			Expect(dbResult.Bid).To(Equal(test_data.FlapKickModel.Bid))
			Expect(dbResult.BidId).To(Equal(test_data.FlapKickModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.FlapKickModel.Lot))
			Expect(dbResult.LogID).To(Equal(logID))
			Expect(dbResult.ContractAddress).To(Equal(addressId))
		})

		It("doesn't insert a new address if the flap kick insertion fails", func() {
			badFlapKick := test_data.FlapKickModel
			badFlapKick.Bid = ""
			err := flapKickRepository.Create([]interface{}{model, badFlapKick})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("invalid input syntax for type numeric"))
			test_data.AssertDBRecordCount(db, "maker.flap_kick", 0)
			// address is inserted with the header_sync_log
			// TODO: include address_id on the header_sync_log?
			test_data.AssertDBRecordCount(db, "addresses", 1)
		})

		It("marks log as transformed", func() {
			insertErr := flapKickRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			flapKickLogTwo := test_data.CreateTestLog(headerID, db)
			modelWithDifferentLogID := test_data.FlapKickModel
			modelWithDifferentLogID.HeaderID = headerID
			modelWithDifferentLogID.LogID = flapKickLogTwo.ID

			insertOneErr := flapKickRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := flapKickRepository.Create([]interface{}{modelWithDifferentLogID})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("handles events with the same header_id + log_id combo by upserting", func() {
			insertOneErr := flapKickRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := flapKickRepository.Create([]interface{}{model})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			insertErr := flapKickRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.flap_kick`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			insertErr := flapKickRepository.Create([]interface{}{test_data.WrongModel{}})

			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))
		})

		It("rolls back the transaction if the given model is of the wrong type", func() {
			insertErr := flapKickRepository.Create([]interface{}{model, test_data.WrongModel{}})
			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.flap_kick`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
