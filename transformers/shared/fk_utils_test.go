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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shared repository", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())
	const hexIlk = "0x464b450000000000000000000000000000000000000000000000000000000000"

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	Describe("GetOrCreateIlk", func() {
		It("returns ID for same ilk with or without hex prefix", func() {
			ilkIDOne, insertErrOne := shared.GetOrCreateIlk(hexIlk, db)
			Expect(insertErrOne).NotTo(HaveOccurred())

			ilkIDTwo, insertErrTwo := shared.GetOrCreateIlk(hexIlk[2:], db)
			Expect(insertErrTwo).NotTo(HaveOccurred())

			Expect(ilkIDOne).To(Equal(ilkIDTwo))
		})
	})

	Describe("GetOrCreateIlkInTransaction", func() {
		It("returns ID for same ilk with or without hex prefix", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			ilkIDOne, insertErrOne := shared.GetOrCreateIlkInTransaction(hexIlk, tx)
			Expect(insertErrOne).NotTo(HaveOccurred())

			ilkIDTwo, insertErrTwo := shared.GetOrCreateIlkInTransaction(hexIlk[2:], tx)
			Expect(insertErrTwo).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			Expect(ilkIDOne).NotTo(BeZero())
			Expect(ilkIDOne).To(Equal(ilkIDTwo))
		})
	})
})
