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
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shared repository", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())
	const hexIlk = "0x464b450000000000000000000000000000000000000000000000000000000000"

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
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

	Describe("GetOrCreateAddress", func() {
		It("creates an address record", func() {
			_, err := shared.GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(err).NotTo(HaveOccurred())

			var address string
			db.Get(&address, `SELECT address from addresses LIMIT 1`)
			Expect(address).To(Equal(fakes.FakeAddress.Hex()))
		})

		It("returns the id for an address that already exists", func() {
			//create the address record
			createAddressId, createErr := shared.GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(createErr).NotTo(HaveOccurred())

			//get the address record
			getAddressId, getErr := shared.GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(createAddressId).To(Equal(getAddressId))

			var addressCount int
			db.Get(&addressCount, `SELECT count(*) from addresses`)
			Expect(addressCount).To(Equal(1))
		})
	})

	Describe("GetOrCreateAddressInTransaction", func() {
		It("creates an address record", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			_, createErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			var address string
			db.Get(&address, `SELECT address from addresses LIMIT 1`)
			Expect(address).To(Equal(fakes.FakeAddress.Hex()))
		})

		It("returns the id for an address that already exists", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			//create the address record
			createAddressId, createErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			//get the address record
			getAddressId, getErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(getErr).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			Expect(createAddressId).To(Equal(getAddressId))

			var addressCount int
			db.Get(&addressCount, `SELECT count(*) from addresses`)
			Expect(addressCount).To(Equal(1))
		})

		It("doesn't persist the address if the transaction is rolled back", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			_, createErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			commitErr := tx.Rollback()
			Expect(commitErr).NotTo(HaveOccurred())

			var addressCount int
			db.Get(&addressCount, `SELECT count(*) from addresses`)
			Expect(addressCount).To(Equal(0))
		})
	})
})
