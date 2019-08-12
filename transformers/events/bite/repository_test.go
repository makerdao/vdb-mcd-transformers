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

package bite_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("Bite repository", func() {
	Describe("Create", func() {
		var (
			biteRepository  *bite.BiteRepository
			db              *postgres.DB
			headerID, logID int64
			model           bite.BiteModel
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			biteRepository = &bite.BiteRepository{}
			biteRepository.SetDB(db)
			headerID = test_data.CreateTestHeader(db)
			biteLog := test_data.CreateTestLog(headerID, db)
			logID = biteLog.ID
			model = test_data.BiteModel
			model.HeaderID = headerID
			model.LogID = logID
		})

		It("persists a bite record", func() {
			insertErr := biteRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			var dbBite bite.BiteModel
			getErr := db.Get(&dbBite, `SELECT urn_id, ink, art, tab, flip, bite_identifier FROM maker.bite WHERE header_id = $1`, headerID)
			Expect(getErr).NotTo(HaveOccurred())
			urnID, err := shared.GetOrCreateUrn(test_data.BiteModel.Urn, test_data.BiteModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbBite.Urn).To(Equal(strconv.FormatInt(urnID, 10)))
			Expect(dbBite.Ink).To(Equal(test_data.BiteModel.Ink))
			Expect(dbBite.Art).To(Equal(test_data.BiteModel.Art))
			Expect(dbBite.Tab).To(Equal(test_data.BiteModel.Tab))
			Expect(dbBite.Flip).To(Equal(test_data.BiteModel.Flip))
			Expect(dbBite.Id).To(Equal(test_data.BiteModel.Id))
		})

		It("marks log as transformed", func() {
			insertErr := biteRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			biteLogTwo := test_data.CreateTestLog(headerID, db)
			modelWithDifferentLogID := test_data.BiteModel
			modelWithDifferentLogID.HeaderID = headerID
			modelWithDifferentLogID.LogID = biteLogTwo.ID

			insertOneErr := biteRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := biteRepository.Create([]interface{}{modelWithDifferentLogID})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("handles events with the same header_id + log_id combo by upserting", func() {
			insertOneErr := biteRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := biteRepository.Create([]interface{}{model})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			insertErr := biteRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.bite`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("removes the log event record if the corresponding log is deleted", func() {
			insertErr := biteRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.bite`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			insertErr := biteRepository.Create([]interface{}{test_data.WrongModel{}})

			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))
		})

		It("rolls back the transaction if the given model is of the wrong type", func() {
			insertErr := biteRepository.Create([]interface{}{model, test_data.WrongModel{}})
			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.bite`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
