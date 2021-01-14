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

var _ = Describe("Get the Max Block Height of Transformed Diffs Query", func() {
	const maxBlockHeightOfTransformedStorageDiffs = `SELECT * FROM api.get_max_transformed_diff_block()`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("ignores untransformed diffs", func() {
		test_helpers.CreateFakeDiffRecord(db)

		var diff int
		err := db.Get(&diff, maxBlockHeightOfTransformedStorageDiffs)
		Expect(err).Should(HaveOccurred())

		Expect(diff).To(Equal(0))
	})

	It("returns the max block height of transformed diffs", func() {
		diffRepo := storage.NewDiffRepository(db)
		firstHeader := fakes.GetFakeHeader(1)
		secondHeader := fakes.GetFakeHeader(2)
		firstDiff := test_helpers.CreateFakeDiffRecordWithHeader(db, firstHeader)
		secondDiff := test_helpers.CreateFakeDiffRecordWithHeader(db, secondHeader)
		diffRepo.MarkTransformed(firstDiff)
		diffRepo.MarkTransformed(secondDiff)

		var diff int
		Expect(db.Get(&diff, maxBlockHeightOfTransformedStorageDiffs)).To(Succeed())

		Expect(diff).To(Equal(2))
	})
})
