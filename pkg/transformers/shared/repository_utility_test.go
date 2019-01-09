// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/rand"
)

var _ = Describe("Repository utilities", func() {
	Describe("MissingHeaders", func() {
		var (
			db                       *postgres.DB
			headerRepository         datastore.HeaderRepository
			startingBlockNumber      int64
			endingBlockNumber        int64
			eventSpecificBlockNumber int64
			blockNumbers             []int64
			headerIDs                []int64
			notCheckedSQL            string
			err                      error
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			headerRepository = repositories.NewHeaderRepository(db)

			columnNames, err := shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			notCheckedSQL = shared.CreateNotCheckedSQL(columnNames)

			startingBlockNumber = rand.Int63()
			eventSpecificBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2
			outOfRangeBlockNumber := endingBlockNumber + 1

			blockNumbers = []int64{startingBlockNumber, eventSpecificBlockNumber, endingBlockNumber, outOfRangeBlockNumber}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("only treats headers as checked if the event specific logs have been checked", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber, db, notCheckedSQL)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n + 10))
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(err).NotTo(HaveOccurred())
			nodeOneMissingHeaders, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber, db, notCheckedSQL)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(3))
			Expect(nodeOneMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))

			nodeTwoMissingHeaders, err := shared.MissingHeaders(startingBlockNumber, endingBlockNumber+10, dbTwo, notCheckedSQL)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(3))
			Expect(nodeTwoMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber+10), Equal(eventSpecificBlockNumber+10), Equal(endingBlockNumber+10)))
			Expect(nodeTwoMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber+10), Equal(eventSpecificBlockNumber+10), Equal(endingBlockNumber+10)))
		})
	})

	Describe("GetCheckedColumnNames", func() {
		It("gets the column names from checked_headers", func() {
			db := test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			expectedColumnNames := getExpectedColumnNames()
			actualColumnNames, err := shared.GetCheckedColumnNames(db)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualColumnNames).To(Equal(expectedColumnNames))
		})
	})

	Describe("CreateNotCheckedSQL", func() {
		It("generates a correct SQL string for one column", func() {
			columns := []string{"columnA"}
			expected := "NOT (columnA)"
			actual := shared.CreateNotCheckedSQL(columns)
			Expect(actual).To(Equal(expected))
		})

		It("generates a correct SQL string for several columns", func() {
			columns := []string{"columnA", "columnB"}
			expected := "NOT (columnA AND columnB)"
			actual := shared.CreateNotCheckedSQL(columns)
			Expect(actual).To(Equal(expected))
		})

		It("defaults to FALSE when there are no columns", func() {
			expected := "FALSE"
			actual := shared.CreateNotCheckedSQL([]string{})
			Expect(actual).To(Equal(expected))
		})
	})
})

func getExpectedColumnNames() []string {
	return []string{
		"price_feeds_checked",
		"flip_kick_checked",
		"frob_checked",
		"tend_checked",
		"bite_checked",
		"dent_checked",
		"pit_file_debt_ceiling_checked",
		"pit_file_ilk_checked",
		"vat_init_checked",
		"drip_file_ilk_checked",
		"drip_file_repo_checked",
		"drip_file_vow_checked",
		"deal_checked",
		"drip_drip_checked",
		"cat_file_chop_lump_checked",
		"cat_file_flip_checked",
		"cat_file_pit_vow_checked",
		"flop_kick_checked",
		"vat_move_checked",
		"vat_fold_checked",
		"vat_heal_checked",
		"vat_toll_checked",
		"vat_tune_checked",
		"vat_grab_checked",
		"vat_flux_checked",
		"vat_slip_checked",
		"vow_flog_checked",
		"flap_kick_checked",
	}
}
