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
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow storage repository test", func() {
	var (
		diffID, fakeHeaderID int64
		fakeAddress          string
		fakeUint256          string
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		err                  error
		repo                 vow.VowStorageRepository
	)

	BeforeEach(func() {
		fakeAddress = fakes.FakeAddress.Hex()
		fakeUint256 = strconv.Itoa(rand.Intn(1000000))
		test_config.CleanTestDB(db)
		repo = vow.VowStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Wards mapping", func() {
		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := storage.GetValueMetadata(wards.Wards, map[storage.Key]string{constants.User: fakeUserAddress}, storage.Uint256)

			setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result WardsMappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			err := db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := shared.GetOrCreateAddress(fakeUserAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())
			Expect(result.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))
			AssertMapping(result.MappingRes, diffID, fakeHeaderID, strconv.FormatInt(userAddressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := storage.GetValueMetadata(wards.Wards, map[storage.Key]string{constants.User: fakeUserAddress}, storage.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing user", func() {
			malformedWardsMetadata := storage.GetValueMetadata(wards.Wards, map[storage.Key]string{}, storage.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedWardsMetadata, fakeUint256)
			Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.User}))
		})
	})

	It("persists a vow vat", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, vat AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowVatTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate vow vat", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowVatTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow flapper", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.FlapperMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, flapper AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowFlapperTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate vow flapper", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.FlapperMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.FlapperMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowFlapperTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow flopper", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.FlopperMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, flopper AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowFlopperTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate vow flopper", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.FlopperMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.FlopperMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowFlopperTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	Describe("vow sin mapping", func() {
		It("writes row", func() {
			timestamp := "1538558052"
			fakeKeys := map[storage.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := storage.GetValueMetadata(vow.SinMapping, fakeKeys, storage.Uint256)

			err := repo.Create(diffID, fakeHeaderID, vowSinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, era AS key, tab AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowSinMappingTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, timestamp, fakeUint256)
		})

		It("does not duplicate row", func() {
			timestamp := "1538558052"
			fakeKeys := map[storage.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := storage.GetValueMetadata(vow.SinMapping, fakeKeys, storage.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, vowSinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, vowSinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowSinMappingTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing timestamp", func() {
			malformedVowSinMappingMetadata := storage.GetValueMetadata(vow.SinMapping, nil, storage.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedVowSinMappingMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.Timestamp}))
		})
	})

	It("persists a vow Sin integer", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.SinIntegerMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, sin AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowSinIntegerTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Sin integer", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.SinIntegerMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.SinIntegerMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowSinIntegerTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Ash", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.AshMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, ash AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowAshTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Ash", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.AshMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.AshMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowAshTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Wait", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.WaitMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, wait AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowWaitTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Wait", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.WaitMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.WaitMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowWaitTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Dump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.DumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, dump AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowDumpTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Dump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.DumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.DumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*)  FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowDumpTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Sump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.SumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, sump AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowSumpTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Sump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.SumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.SumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowSumpTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Bump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.BumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, bump AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowBumpTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Bump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.BumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.BumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowBumpTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a vow Hump", func() {
		err = repo.Create(diffID, fakeHeaderID, vow.HumpMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, hump AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowHumpTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vow Hump", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vow.HumpMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vow.HumpMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VowHumpTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
