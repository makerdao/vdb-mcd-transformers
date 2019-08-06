// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package shared_behaviors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type CreateBehaviorInputs struct {
	CheckedHeaderColumnName string
	Database                *postgres.DB
	LogEventTableName       string
	TestModel               interface{}
	RecheckTestModel        interface{}
	ModelWithDifferentLogID interface{}
	Repository              event.Repository
}

type MarkedHeaderCheckedBehaviorInputs struct {
	CheckedHeaderColumnName string
	Repository              event.Repository
}

func SharedRepositoryCreateBehaviors(inputs *CreateBehaviorInputs) {
	Describe("Create", func() {
		var (
			db            = inputs.Database
			err           error
			repository    = inputs.Repository
			logEventModel = inputs.TestModel
		)

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			err = repository.Create([]interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create([]interface{}{inputs.ModelWithDifferentLogID})
			Expect(err).NotTo(HaveOccurred())
		})

		It("handles events with the same header_id, tx_idx, log_idx combo by upserting", func() {
			err = repository.Create([]interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create([]interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			err = repository.Create([]interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, 0)
			Expect(err).NotTo(HaveOccurred())

			var count int
			query := `SELECT count(*) from ` + inputs.LogEventTableName
			err = db.QueryRow(query).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			err = repository.Create([]interface{}{test_data.WrongModel{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})

		It("rolls back the transaciton if the given model is of the wrong type", func() {
			err = repository.Create([]interface{}{logEventModel, test_data.WrongModel{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))

			var count int
			query := `SELECT count(*) from ` + inputs.LogEventTableName
			err = db.QueryRow(query).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})
}
