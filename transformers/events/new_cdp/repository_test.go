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

package new_cdp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/new_cdp"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("NewCdp Repository", func() {
	var (
		db               *postgres.DB
		newCdpRepository new_cdp.NewCdpRepository
		headerID, logID  int64
		model            new_cdp.NewCdpModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		newCdpRepository = new_cdp.NewCdpRepository{}
		newCdpRepository.SetDB(db)
		headerID = test_data.CreateTestHeader(db)
		newCdpHeaderSyncLog := test_data.CreateTestLog(headerID, db)
		logID = newCdpHeaderSyncLog.ID
		model = test_data.NewCdpModel
		model.HeaderID = headerID
		model.LogID = logID
	})

	Describe("Create", func() {
		It("persists a new_cdp record", func() {
			err := newCdpRepository.Create([]interface{}{model})
			Expect(err).NotTo(HaveOccurred())

			var count int
			countErr := db.Get(&count, `SELECT count(*) FROM maker.new_cdp`)
			Expect(countErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))

			var dbNewCdp new_cdp.NewCdpModel
			newCdpQuery := `SELECT usr, own, cdp, log_id FROM maker.new_cdp WHERE header_id = $1`
			getErr := db.Get(&dbNewCdp, newCdpQuery, headerID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(dbNewCdp.Usr).To(Equal(test_data.NewCdpModel.Usr))
			Expect(dbNewCdp.Own).To(Equal(test_data.NewCdpModel.Own))
			Expect(dbNewCdp.Cdp).To(Equal(test_data.NewCdpModel.Cdp))
			Expect(dbNewCdp.LogID).To(Equal(logID))
		})

		It("marks log as transformed", func() {
			insertErr := newCdpRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			newCdpHeaderSyncLogTwo := test_data.CreateTestLog(headerID, db)
			modelWithDifferentLogID := test_data.NewCdpModel
			modelWithDifferentLogID.HeaderID = headerID
			modelWithDifferentLogID.LogID = newCdpHeaderSyncLogTwo.ID

			insertOneErr := newCdpRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := newCdpRepository.Create([]interface{}{modelWithDifferentLogID})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("handles events with the same header_id + log_id combo by upserting", func() {
			insertOneErr := newCdpRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := newCdpRepository.Create([]interface{}{model})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			insertErr := newCdpRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.bite`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("removes the log event record if the corresponding log is deleted", func() {
			insertErr := newCdpRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.bite`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			insertErr := newCdpRepository.Create([]interface{}{test_data.WrongModel{}})

			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))
		})

		It("rolls back the transaction if the given model is of the wrong type", func() {
			insertErr := newCdpRepository.Create([]interface{}{model, test_data.WrongModel{}})
			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.bite`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
