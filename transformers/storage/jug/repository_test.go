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

package jug_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jug storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		fakeAddress          = "0x12345"
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
		repo                 jug.StorageRepository
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = jug.StorageRepository{ContractAddress: test_data.JugAddress()}
		repo.SetDB(db)

		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Wards", func() {
		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
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
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)
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
			malformedWardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedWardsMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.User}))
		})
	})

	Describe("Ilk", func() {
		Describe("Rho", func() {
			It("writes a row", func() {
				ilkRhoMetadata := types.GetValueMetadata(jug.IlkRho, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkRhoMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS KEY, rho AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugIlkRhoTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkRhoMetadata := types.GetValueMetadata(jug.IlkRho, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkRhoMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkRhoMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugIlkRhoTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkRhoMetadata := types.GetValueMetadata(jug.IlkRho, nil, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkRhoMetadata, fakeUint256)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(jug.IlkRho, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
				PropertyName:  "Rho",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.JugIlkRhoTable,
			})
		})

		Describe("Duty", func() {
			It("writes a row", func() {
				ilkDutyMetadata := types.GetValueMetadata(jug.IlkDuty, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkDutyMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS KEY, duty AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugIlkDutyTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkDutyMetadata := types.GetValueMetadata(jug.IlkDuty, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkDutyMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkDutyMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugIlkDutyTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkDutyMetadata := types.GetValueMetadata(jug.IlkDuty, nil, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkDutyMetadata, fakeUint256)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(jug.IlkDuty, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
				PropertyName:  "Duty",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.JugIlkDutyTable,
			})
		})
	})

	It("persists a jug vat", func() {
		err := repo.Create(diffID, fakeHeaderID, jug.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, vat AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugVatTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate jug vat", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, jug.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, jug.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugVatTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a jug vow", func() {
		err := repo.Create(diffID, fakeHeaderID, jug.VowMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, vow AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugVowTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate jug vow", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, jug.VowMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, jug.VowMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugVowTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a jug base", func() {
		err := repo.Create(diffID, fakeHeaderID, jug.BaseMetadata, fakeUint256)
		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, base AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugBaseTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())

		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate jug base", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, jug.BaseMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, jug.BaseMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.JugBaseTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
