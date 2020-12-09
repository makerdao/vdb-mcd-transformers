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

package queries

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New State Diff Query", func() {
	const transformationStatusQuery = `SELECT * FROM api.storage_transformation_status()`

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})
	It("includes the total number of new storage diffs", func() {
		storage_helper.CreateFakeDiffRecord(db)

		var diff []int
		Expect(db.Select(&diff, transformationStatusQuery)).To(Succeed())

		Expect(len(diff)).To(Equal(1))
	})

	It("excludes non-new storage diffs", func() {
		diffID := storage_helper.CreateFakeDiffRecord(db)
		diffRepo := storage.NewDiffRepository(db)
		diffRepo.MarkTransformed(diffID)

		var diff []int
		Expect(db.Select(&diff, transformationStatusQuery)).To(Succeed())

		Expect(diff).To(BeEmpty())
	})

	It("includes the block numbers, in ascending order", func() {
		firstHeader := fakes.GetFakeHeader(1)
		secondHeader := fakes.GetFakeHeader(2)
		firstDiff := storage_helper.CreateFakeDiffRecordWithHeader(db, firstHeader)
		secondDiff := storage_helper.CreateFakeDiffRecordWithHeader(db, secondHeader)
		Expect(firstDiff).NotTo(Equal(secondDiff))

		var diff []int
		Expect(db.Select(&diff, transformationStatusQuery)).To(Succeed())

		Expect(len(diff)).To(Equal(2))
		Expect(diff[0]).To(Equal(1))
		Expect(diff[1]).To(Equal(2))
	})
})
