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
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get Block Heights for New Untransformed Diffs Query", func() {
	const blockHeightsForNewUntransformedStorageDiffs = `SELECT * FROM api.get_block_heights_for_new_untransformed_diffs()`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("has a row for every new storage diff", func() {
		storage_helper.CreateFakeDiffRecord(db)

		var diff []int
		Expect(db.Select(&diff, blockHeightsForNewUntransformedStorageDiffs)).To(Succeed())

		Expect(len(diff)).To(Equal(1))
	})

	It("excludes non-new storage diffs", func() {
		diffID := storage_helper.CreateFakeDiffRecord(db)
		diffRepo := storage.NewDiffRepository(db)
		diffRepo.MarkTransformed(diffID)

		var diff []int
		Expect(db.Select(&diff, blockHeightsForNewUntransformedStorageDiffs)).To(Succeed())

		Expect(diff).To(BeEmpty())
	})

	It("includes the block heights, in ascending order for all new diffs", func() {
		firstHeader := fakes.GetFakeHeader(1)
		secondHeader := fakes.GetFakeHeader(2)
		firstDiff := storage_helper.CreateFakeDiffRecordWithHeader(db, firstHeader)
		secondDiff := storage_helper.CreateFakeDiffRecordWithHeader(db, secondHeader)
		Expect(firstDiff).NotTo(Equal(secondDiff))

		var diff []int
		Expect(db.Select(&diff, blockHeightsForNewUntransformedStorageDiffs)).To(Succeed())

		Expect(len(diff)).To(Equal(2))
		Expect(diff[0]).To(Equal(1))
		Expect(diff[1]).To(Equal(2))
	})
})
