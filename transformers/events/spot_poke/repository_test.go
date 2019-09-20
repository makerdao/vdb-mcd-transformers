//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package spot_poke_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_poke"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("Spot Poke repository", func() {

	Describe("Create", func() {
		var (
			spotPokeRepository spot_poke.SpotPokeRepository
			db                 *postgres.DB
			headerID, logID    int64
			model              spot_poke.SpotPokeModel
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			spotPokeRepository = spot_poke.SpotPokeRepository{}
			spotPokeRepository.SetDB(db)
			headerID = test_data.CreateTestHeader(db)
			spotPokeLog := test_data.CreateTestLog(headerID, db)
			logID = spotPokeLog.ID
			model = test_data.SpotPokeModel
			model.HeaderID = headerID
			model.LogID = logID
		})

		It("persists a spot poke record", func() {
			insertErr := spotPokeRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			ilkID, err := shared.GetOrCreateIlk(test_data.SpotPokeModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())

			var dbSpotPoke spot_poke.SpotPokeModel
			err = db.Get(&dbSpotPoke, `SELECT header_id, log_id, ilk_id, value, spot, log_id FROM maker.spot_poke WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbSpotPoke.HeaderID).To(Equal(headerID))
			Expect(dbSpotPoke.LogID).To(Equal(logID))
			Expect(dbSpotPoke.Ilk).To(Equal(strconv.FormatInt(ilkID, 10)))
			Expect(dbSpotPoke.Value).To(Equal(test_data.SpotPokeModel.Value))
			Expect(dbSpotPoke.Spot).To(Equal(test_data.SpotPokeModel.Spot))
		})

		It("marks log as transformed", func() {
			insertErr := spotPokeRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			modelWithDifferentLogID := test_data.SpotPokeModel
			modelWithDifferentLogID.HeaderID = headerID
			spotPokeLogTwo := test_data.CreateTestLog(headerID, db)
			modelWithDifferentLogID.LogID = spotPokeLogTwo.ID

			insertOneErr := spotPokeRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := spotPokeRepository.Create([]interface{}{modelWithDifferentLogID})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("handles events with the same header_id + log_id combo by upserting", func() {
			insertOneErr := spotPokeRepository.Create([]interface{}{model})
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := spotPokeRepository.Create([]interface{}{model})
			Expect(insertTwoErr).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			insertErr := spotPokeRepository.Create([]interface{}{model})
			Expect(insertErr).NotTo(HaveOccurred())

			_, deleteErr := db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(deleteErr).NotTo(HaveOccurred())

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.spot_poke`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			insertErr := spotPokeRepository.Create([]interface{}{test_data.WrongModel{}})

			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))
		})

		It("rolls back the transaction if insertion fails", func() {
			badSpotPokeModel := spot_poke.SpotPokeModel{}
			badSpotPokeModel.HeaderID = headerID

			insertErr := spotPokeRepository.Create([]interface{}{model, badSpotPokeModel})
			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("invalid input syntax for type numeric"))

			var spotPokeCount int
			getOneErr := db.Get(&spotPokeCount, `SELECT count(*) FROM maker.spot_poke`)
			Expect(getOneErr).NotTo(HaveOccurred())
			Expect(spotPokeCount).To(Equal(0))

			var ilkCount int
			getTwoErr := db.Get(&ilkCount, `SELECT count(*) FROM maker.ilks`)
			Expect(getTwoErr).NotTo(HaveOccurred())
			Expect(ilkCount).To(Equal(0))
		})

		It("rolls back the transaction if the given model is of the wrong type", func() {
			insertErr := spotPokeRepository.Create([]interface{}{model, test_data.WrongModel{}})
			Expect(insertErr).To(HaveOccurred())
			Expect(insertErr.Error()).To(ContainSubstring("model of type"))

			var count int
			getErr := db.QueryRow(`SELECT count(*) FROM maker.spot_poke`).Scan(&count)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
})
