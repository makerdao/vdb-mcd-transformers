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

package spot_test

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spot storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 spot.SpotStorageRepository
		fakeAddress          = "0x12345"
		fakeUint256          = "12345"
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = spot.SpotStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Ilk", func() {
		Describe("Pip", func() {
			It("writes a row", func() {
				ilkPipMetadata := storage.GetValueMetadata(spot.IlkPip, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Address)

				err := repo.Create(diffID, fakeHeaderID, ilkPipMetadata, fakeAddress)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT diff_id, header_id, ilk_id AS key, pip AS VALUE FROM maker.spot_ilk_pip`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkPipMetadata := storage.GetValueMetadata(spot.IlkPip, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Address)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkPipMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkPipMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_ilk_pip`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkPipMetadata := storage.GetValueMetadata(spot.IlkPip, nil, storage.Address)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkPipMetadata, fakeAddress)
				Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      storage.GetValueMetadata(spot.IlkPip, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Address),
				PropertyName:  "Pip",
				PropertyValue: fakeAddress,
				TableName:     "maker.spot_ilk_pip",
			})
		})

		Describe("Mat", func() {
			It("writes a row", func() {
				ilkMatMetadata := storage.GetValueMetadata(spot.IlkMat, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT diff_id, header_id, ilk_id AS KEY, mat AS VALUE FROM maker.spot_ilk_mat`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkMatMetadata := storage.GetValueMetadata(spot.IlkMat, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_ilk_mat`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkMatMetadata := storage.GetValueMetadata(spot.IlkMat, nil, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkMatMetadata, fakeUint256)
				Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      storage.GetValueMetadata(spot.IlkMat, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256),
				PropertyName:  "Mat",
				PropertyValue: strconv.Itoa(rand.Int()),
				TableName:     "maker.spot_ilk_mat",
			})
		})
	})

	It("persists a spot vat", func() {
		err := repo.Create(diffID, fakeHeaderID, spot.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, vat AS value FROM maker.spot_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeAddress)
	})

	It("does not duplicate spot vat", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, spot.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, spot.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_vat`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a spot par", func() {
		err := repo.Create(diffID, fakeHeaderID, spot.ParMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, par AS value FROM maker.spot_par`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate spot par", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, spot.ParMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, spot.ParMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_par`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a spot live", func() {
		err := repo.Create(diffID, fakeHeaderID, spot.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT diff_id, header_id, live AS value FROM maker.spot_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate spot live", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, spot.LiveMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, spot.LiveMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_live`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
