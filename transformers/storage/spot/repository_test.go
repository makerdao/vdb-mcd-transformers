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
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	shared2 "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spot storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 spot.StorageRepository
		fakeAddress          = "0x" + fakes.RandomString(20)
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = spot.StorageRepository{ContractAddress: test_data.SpotAddress()}
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
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			err := db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, repo.ContractAddress)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := repository.GetOrCreateAddress(db, fakeUserAddress)
			Expect(userAddressErr).NotTo(HaveOccurred())
			AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(userAddressID, 10), fakeUint256)
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
		Describe("Pip", func() {
			It("writes a row", func() {
				ilkPipMetadata := types.GetValueMetadata(spot.IlkPip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address)

				err := repo.Create(diffID, fakeHeaderID, ilkPipMetadata, fakeAddress)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS KEY, pip AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkPipTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared2.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkPipMetadata := types.GetValueMetadata(spot.IlkPip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkPipMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkPipMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkPipTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkPipMetadata := types.GetValueMetadata(spot.IlkPip, nil, types.Address)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkPipMetadata, fakeAddress)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(spot.IlkPip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address),
				PropertyName:  "Pip",
				PropertyValue: fakeAddress,
				Schema:        constants.MakerSchema,
				TableName:     constants.SpotIlkPipTable,
			})
		})

		Describe("Mat", func() {
			It("writes a row", func() {
				ilkMatMetadata := types.GetValueMetadata(spot.IlkMat, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS KEY, mat AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkMatTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared2.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkMatMetadata := types.GetValueMetadata(spot.IlkMat, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkMatTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkMatMetadata := types.GetValueMetadata(spot.IlkMat, nil, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkMatMetadata, fakeUint256)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(spot.IlkMat, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
				PropertyName:  "Mat",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.SpotIlkMatTable,
			})
		})
	})

	Describe("vat", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: spot.Vat,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.SpotVatTable,
			Repository:     &repo,
			Metadata:       spot.VatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("par", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: spot.Par,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.SpotParTable,
			Repository:     &repo,
			Metadata:       spot.ParMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("live", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: spot.Live,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.SpotLiveTable,
			Repository:     &repo,
			Metadata:       spot.LiveMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})
})
