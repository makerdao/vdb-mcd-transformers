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

package flip_kick_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("FlipKick Repository", func() {
	Describe("Create", func() {
		var (
			db                 *postgres.DB
			flipKickRepository flip_kick.FlipKickRepository
			headerID, logID    int64
			model              flip_kick.FlipKickModel
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			flipKickRepository = flip_kick.FlipKickRepository{}
			flipKickRepository.SetDB(db)

			headerID = test_data.CreateTestHeader(db)
			insertedLogs := test_data.CreateLogs(headerID, []types.Log{test_data.FlipKickHeaderSyncLog.Log}, db)
			Expect(len(insertedLogs)).To(Equal(1))
			logID = insertedLogs[0].ID

			model = test_data.FlipKickModel
			model.HeaderID = headerID
			model.LogID = logID
		})

		It("persists flip_kick records", func() {
			insertErr := flipKickRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			test_data.AssertDBRecordCount(db, "maker.flip_kick", 1)

			var addressId string
			addressErr := db.Get(&addressId, `SELECT id FROM addresses`)
			Expect(addressErr).NotTo(HaveOccurred())

			dbResult := test_data.FlipKickDBRow{}
			getErr := db.QueryRowx(`SELECT * FROM maker.flip_kick`).StructScan(&dbResult)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(dbResult.HeaderID).To(Equal(headerID))
			Expect(dbResult.LogID).To(Equal(logID))
			Expect(dbResult.BidId).To(Equal(test_data.FlipKickModel.BidId))
			Expect(dbResult.Lot).To(Equal(test_data.FlipKickModel.Lot))
			Expect(dbResult.Bid).To(Equal(test_data.FlipKickModel.Bid))
			Expect(dbResult.Tab).To(Equal(test_data.FlipKickModel.Tab))
			Expect(dbResult.Usr).To(Equal(test_data.FlipKickModel.Usr))
			Expect(dbResult.Gal).To(Equal(test_data.FlipKickModel.Gal))
			Expect(dbResult.ContractAddress).To(Equal(addressId))
		})

		It("doesn't insert a new address if the flip kick insertion fails", func() {
			badFlipKick := test_data.FlipKickModel
			badFlipKick.Bid = ""
			err := flipKickRepository.Create([]interface{}{model, badFlipKick})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("invalid input syntax for type numeric"))

			test_data.AssertDBRecordCount(db, "maker.flip_kick", 0)
			test_data.AssertDBRecordCount(db, "addresses", 1)
		})

		It("marks log as transformed", func() {
			insertErr := flipKickRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			modelWithDifferentLogID := test_data.FlipKickModel
			modelWithDifferentLogID.HeaderID = headerID
			insertedLog := test_data.CreateTestLog(headerID, db)
			modelWithDifferentLogID.LogID = insertedLog.ID

			insertOneErr := flipKickRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := flipKickRepository.Create([]interface{}{modelWithDifferentLogID})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("handles events with the same header_id + log_id combo by upserting", func() {
			insertOneErr := flipKickRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := flipKickRepository.Create([]interface{}{model})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			insertErr := flipKickRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.flip_kick`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			insertErr := flipKickRepository.Create([]interface{}{test_data.WrongModel{}})

			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))
		})

		It("rolls back the transaction if the given model is of the wrong type", func() {
			insertErr := flipKickRepository.Create([]interface{}{model, test_data.WrongModel{}})
			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.flip_kick`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
