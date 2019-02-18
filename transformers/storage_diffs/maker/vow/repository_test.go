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

package vow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	. "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/vow"
)

var _ = Describe("Vow storage repository test", func() {
	var (
		blockNumber int
		blockHash   string
		db          *postgres.DB
		err         error
		repo        vow.VowStorageRepository
	)

	BeforeEach(func() {
		blockNumber = 123
		blockHash = "expected_block_hash"
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vow.VowStorageRepository{}
		repo.SetDB(db)
	})

	It("persists a vow vat", func() {
		expectedVat := "123"

		err = repo.Create(blockNumber, blockHash, vow.VatMetadata, expectedVat)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value from maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedVat)
	})

	It("persists a vow cow", func() {
		expectedCow := "123"

		err = repo.Create(blockNumber, blockHash, vow.CowMetadata, expectedCow)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, cow AS value from maker.vow_cow`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedCow)
	})

	It("persists a vow row", func() {
		expectedRow := "123"

		err = repo.Create(blockNumber, blockHash, vow.RowMetadata, expectedRow)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, row AS value from maker.vow_row`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedRow)
	})

	It("persists a vow Sin", func() {
		expectedSow := "123"

		err = repo.Create(blockNumber, blockHash, vow.SinMetadata, expectedSow)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, sin AS value from maker.vow_sin`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedSow)
	})

	It("persists a vow woe", func() {
		expectedWoe := "123"

		err = repo.Create(blockNumber, blockHash, vow.WoeMetadata, expectedWoe)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, woe AS value from maker.vow_woe`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedWoe)
	})

	It("persists a vow Ash", func() {
		expectedAsh := "123"

		err = repo.Create(blockNumber, blockHash, vow.AshMetadata, expectedAsh)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, ash AS value from maker.vow_ash`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedAsh)
	})

	It("persists a vow Wait", func() {
		expectedWait := "123"

		err = repo.Create(blockNumber, blockHash, vow.WaitMetadata, expectedWait)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, wait AS value from maker.vow_wait`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedWait)
	})

	It("persists a vow Bump", func() {
		expectedBump := "123"

		err = repo.Create(blockNumber, blockHash, vow.BumpMetadata, expectedBump)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, bump AS value from maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedBump)
	})

	It("persists a vow Sump", func() {
		expectedSump := "123"

		err = repo.Create(blockNumber, blockHash, vow.SumpMetadata, expectedSump)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, sump AS value from maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedSump)
	})

	It("persists a vow Hump", func() {
		expectedHump := "123"

		err = repo.Create(blockNumber, blockHash, vow.HumpMetadata, expectedHump)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, hump AS value from maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, blockNumber, blockHash, expectedHump)
	})
})
