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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow storage repository test", func() {
	var (
		diffID, fakeHeaderID int64
		fakeAddress  string
		fakeUint256  string
		db           = test_config.NewTestDB(test_config.NewTestNode())
		err          error
		repo         vow.VowStorageRepository
	)

	BeforeEach(func() {
		fakeAddress = fakes.FakeAddress.Hex()
		fakeUint256 = "12345"
		test_config.CleanTestDB(db)
		repo = vow.VowStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	It("persists a vow vat", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, vat AS value from maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate vow vat", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_vat`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow flapper", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.FlapperMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, flapper AS value from maker.vow_flapper`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate vow flapper", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.FlapperMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.FlapperMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_flapper`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow flopper", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.FlopperMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, flopper AS value from maker.vow_flopper`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate vow flopper", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.FlopperMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.FlopperMetadata, fakeAddress)

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

			err := repo.Create(diffID, fakeHeaderID, vowSinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT diff_id, header_id, era AS key, tab AS value FROM maker.vow_sin_mapping`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, timestamp, fakeUint256)
		})

		It("does not duplicate row", func() {
			timestamp := "1538558052"
			fakeKeys := map[utils.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := utils.GetStorageValueMetadata(vow.SinMapping, fakeKeys, utils.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, vowSinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, vowSinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_sin_mapping`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing timestamp", func() {
			malformedVowSinMappingMetadata := utils.GetStorageValueMetadata(vow.SinMapping, nil, utils.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedVowSinMappingMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Timestamp}))
		})
	})

	It("persists a vow Sin integer", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.SinIntegerMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, sin AS value from maker.vow_sin_integer`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Sin integer", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.SinIntegerMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.SinIntegerMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_sin_integer`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Ash", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.AshMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, ash AS value from maker.vow_ash`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Ash", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.AshMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.AshMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_ash`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Wait", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.WaitMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, wait AS value from maker.vow_wait`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Wait", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.WaitMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.WaitMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_wait`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Dump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.DumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, dump AS value from maker.vow_dump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Dump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.DumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.DumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_dump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Sump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.SumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, sump AS value from maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Sump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.SumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.SumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_sump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Bump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.BumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, bump AS value from maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Bump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.BumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.BumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_bump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Hump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.HumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, hump AS value from maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Hump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.HumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.HumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vow_hump`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
