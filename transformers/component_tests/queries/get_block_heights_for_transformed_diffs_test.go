// VulcanizeDB
// Copyright © 2020 Vulcanize

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

package queries

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get Block Height for Transformed Diffs Query", func() {
	const blockHeightsForTransformedStorageDiffs = `SELECT * FROM api.get_block_heights_for_transformed_diffs()`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("ignores untransformed diffs", func() {
		test_helpers.CreateFakeDiffRecord(db)

		var diff []int
		Expect(db.Select(&diff, blockHeightsForTransformedStorageDiffs)).To(Succeed())

		Expect(diff).To(BeEmpty())
	})

	It("includes the block heights, in ascending order for all transformed diffs", func() {
		diffRepo := storage.NewDiffRepository(db)
		firstHeader := fakes.GetFakeHeader(1)
		secondHeader := fakes.GetFakeHeader(2)
		firstDiff := test_helpers.CreateFakeDiffRecordWithHeader(db, firstHeader)
		secondDiff := test_helpers.CreateFakeDiffRecordWithHeader(db, secondHeader)
		diffRepo.MarkTransformed(firstDiff)
		diffRepo.MarkTransformed(secondDiff)

		var diff []int
		Expect(db.Select(&diff, blockHeightsForTransformedStorageDiffs)).To(Succeed())

		Expect(len(diff)).To(Equal(2))
		Expect(diff[0]).To(Equal(1))
		Expect(diff[1]).To(Equal(2))
	})
})
