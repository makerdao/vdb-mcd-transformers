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
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jug storage repository", func() {
	var (
		db           *postgres.DB
		repo         jug.JugStorageRepository
		fakeAddress  = "0x12345"
		fakeUint256  = "12345"
		fakeHeaderID int64
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = jug.JugStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	Describe("Ilk", func() {
		Describe("Rho", func() {
			It("writes a row", func() {
				ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeHeaderID, ilkRhoMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT header_id, ilk_id AS key, rho AS VALUE FROM maker.jug_ilk_rho`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeHeaderID, ilkRhoMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeHeaderID, ilkRhoMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_ilk_rho`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, nil, utils.Uint256)

				err := repo.Create(fakeHeaderID, malformedIlkRhoMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
				PropertyName:  "Rho",
				PropertyValue: strconv.Itoa(rand.Int()),
			})
		})

		Describe("Duty", func() {
			It("writes a row", func() {
				ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeHeaderID, ilkDutyMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT header_id, ilk_id AS KEY, duty AS VALUE FROM maker.jug_ilk_duty`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeHeaderID, ilkDutyMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeHeaderID, ilkDutyMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_ilk_duty`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, nil, utils.Uint256)

				err := repo.Create(fakeHeaderID, malformedIlkDutyMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
				PropertyName:  "Duty",
				PropertyValue: strconv.Itoa(rand.Int()),
			})
		})
	})

	It("persists a jug vat", func() {
		err := repo.Create(fakeHeaderID, jug.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT header_id, vat AS value FROM maker.jug_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate jug vat", func() {
		insertOneErr := repo.Create(fakeHeaderID, jug.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, jug.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_vat`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a jug vow", func() {
		err := repo.Create(fakeHeaderID, jug.VowMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT header_id, vow AS value FROM maker.jug_vow`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate jug vow", func() {
		insertOneErr := repo.Create(fakeHeaderID, jug.VowMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, jug.VowMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_vow`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a jug base", func() {
		err := repo.Create(fakeHeaderID, jug.BaseMetadata, fakeUint256)
		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT header_id, base AS value FROM maker.jug_base`)
		Expect(err).NotTo(HaveOccurred())

		AssertVariable(result, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate jug base", func() {
		insertOneErr := repo.Create(fakeHeaderID, jug.BaseMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, jug.BaseMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_base`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
