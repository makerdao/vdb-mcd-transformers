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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
)

var _ = Describe("Vow storage repository test", func() {
	var (
		fakeBlockNumber int
		fakeBlockHash   string
		fakeAddress     string
		fakeUint256     string
		db              *postgres.DB
		err             error
		repo            vow.VowStorageRepository
	)

	BeforeEach(func() {
		fakeBlockNumber = 123
		fakeBlockHash = "expected_block_hash"
		fakeAddress = fakes.FakeAddress.Hex()
		fakeUint256 = "12345"
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vow.VowStorageRepository{}
		repo.SetDB(db)
	})

	It("persists a vow vat", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value from maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("does not duplicate vow vat", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_vat`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow flapper", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.FlapperMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, flapper AS value from maker.vow_flapper`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("does not duplicate vow flapper", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.FlapperMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.FlapperMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_flapper`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow flopper", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.FlopperMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, flopper AS value from maker.vow_flopper`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("does not duplicate vow flopper", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.FlopperMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.FlopperMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_flopper`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	Describe("vow sin mapping", func() {
		It("writes row", func() {
			timestamp := "1538558052"
			fakeKeys := map[utils.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := utils.GetStorageValueMetadata(vow.SinMapping, fakeKeys, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, vowSinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, era AS key, tab AS value FROM maker.vow_sin_mapping`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, timestamp, fakeUint256)
		})

		It("does not duplicate row", func() {
			timestamp := "1538558052"
			fakeKeys := map[utils.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := utils.GetStorageValueMetadata(vow.SinMapping, fakeKeys, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vowSinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vowSinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_sin_mapping`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing timestamp", func() {
			malformedVowSinMappingMetadata := utils.GetStorageValueMetadata(vow.SinMapping, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedVowSinMappingMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Timestamp}))
		})
	})

	It("persists a vow Sin integer", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.SinIntegerMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, sin AS value from maker.vow_sin_integer`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Sin integer", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.SinIntegerMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.SinIntegerMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_sin_integer`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Ash", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.AshMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, ash AS value from maker.vow_ash`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Ash", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.AshMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.AshMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_ash`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Wait", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.WaitMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, wait AS value from maker.vow_wait`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Wait", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.WaitMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.WaitMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_wait`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Dump", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.DumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, dump AS value from maker.vow_dump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Dump", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.DumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.DumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_dump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Sump", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.SumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, sump AS value from maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Sump", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.SumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.SumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_sump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Bump", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.BumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, bump AS value from maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Bump", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.BumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.BumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_bump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Hump", func() {
		err = repo.Create(fakeBlockNumber, fakeBlockHash, vow.HumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, hump AS value from maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vow Hump", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.HumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vow.HumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_hump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
