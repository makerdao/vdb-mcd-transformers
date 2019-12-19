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
		fakeAddress          = "0x" + fakes.RandomString(20)
		fakeUint256          = strconv.Itoa(rand.Int())
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
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS KEY, pip AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkPipTable))
				err = db.Get(&result, query)
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
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkPipTable))
				getCountErr := db.Get(&count, query)
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
				Schema:        constants.MakerSchema,
				TableName:     constants.SpotIlkPipTable,
			})
		})

		Describe("Mat", func() {
			It("writes a row", func() {
				ilkMatMetadata := storage.GetValueMetadata(spot.IlkMat, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkMatMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS KEY, mat AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkMatTable))
				err = db.Get(&result, query)
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
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.SpotIlkMatTable))
				getCountErr := db.Get(&count, query)
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
