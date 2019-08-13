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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Shared converter", func() {
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
})
