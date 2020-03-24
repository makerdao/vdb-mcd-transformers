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

package shared_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shared transformer utils", func() {
	Describe("VerifyLog", func() {
		It("returns err if log is missing topics", func() {
			log := types.Log{Data: fakes.FakeHash.Bytes()}

			err := shared.VerifyLog(log, 1, true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrLogMissingTopics(1, 0)))
		})

		It("returns error if log has fewer than required number of topics", func() {
			log := types.Log{
				Data: fakes.FakeHash.Bytes(),
				Topics: []common.Hash{
					fakes.FakeHash,
				},
			}

			err := shared.VerifyLog(log, 2, true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrLogMissingTopics(2, 1)))
		})

		It("returns err if log is missing required data", func() {
			log := types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			}

			err := shared.VerifyLog(log, 4, true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared.ErrLogMissingData))
		})

		It("does not return error for missing data if data not required", func() {
			log := types.Log{
				Topics: []common.Hash{{}, {}, {}, {}},
			}

			err := shared.VerifyLog(log, 4, false)

			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return error for valid log", func() {
			log := types.Log{
				Topics: []common.Hash{{}, {}},
				Data:   fakes.FakeHash.Bytes(),
			}

			err := shared.VerifyLog(log, 2, true)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("getLogNoteData", func() {
		var (
			db = test_config.NewTestDB(test_config.NewTestNode())
		)
		BeforeEach(func() {
			test_config.CleanTestDB(db)
		})

		It("gets event log data when there is one indexed value", func() {
			var expectedIDs []int64
			accountIDs, accountErr := shared.GetLogNoteData(2, test_data.MedianDropLogWithOneAccount.Log.Data, db)
			Expect(accountErr).NotTo(HaveOccurred())

			expectedID, addressErr := shared.GetOrCreateAddress("0xacB48fD097f1E0B24d3853BeAd826E5E9278B700", db)
			expectedIDs = append(expectedIDs, expectedID, 0, 0, 0)
			Expect(addressErr).NotTo(HaveOccurred())
			Expect(accountIDs).To(Equal(expectedIDs))
		})

		It("gets event log data when there are four indexed values", func() {
			var expectedIDs []int64
			accountIDs, accountErr := shared.GetLogNoteData(2, test_data.MedianDropLogWithFourAccounts.Log.Data, db)
			Expect(accountErr).NotTo(HaveOccurred())

			expectedID1, addressErr := shared.GetOrCreateAddress("0xc45E7858EEf1318337A803Ede8C5A9bE12E2B40f", db)
			Expect(addressErr).NotTo(HaveOccurred())

			expectedID2, addressErr2 := shared.GetOrCreateAddress("0xEf6B95815E215635bd77851f1Fc42e8750873024", db)
			Expect(addressErr2).NotTo(HaveOccurred())

			expectedID3, addressErr3 := shared.GetOrCreateAddress("0x8efccC4eCb27F7f233A7fF4e74E86c5E979d1c43", db)
			Expect(addressErr3).NotTo(HaveOccurred())

			expectedID4, addressErr4 := shared.GetOrCreateAddress("0xc2D2D553A39cc08e7e294427edE2C38A89c0066A", db)
			Expect(addressErr4).NotTo(HaveOccurred())

			expectedIDs = append(expectedIDs, expectedID1, expectedID2, expectedID3, expectedID4)

			Expect(accountIDs).To(Equal(expectedIDs))
		})
	})
})
