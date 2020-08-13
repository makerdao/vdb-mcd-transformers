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
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow storage repository test", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		diffID, fakeHeaderID int64
		fakeAddress          = "0x" + fakes.RandomString(40)
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
		repo                 vow.StorageRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = vow.StorageRepository{}
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
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := shared.GetOrCreateAddress(fakeUserAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())
			Expect(result.AddressID).To(Equal(contractAddressID))
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

	Describe("Vat", func() {
		metadata := types.ValueMetadata{Name: storage.Vat}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Vat,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowVatTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Flapper", func() {
		metadata := types.ValueMetadata{Name: vow.Flapper}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Flapper,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowFlapperTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Flopper", func() {
		metadata := types.ValueMetadata{Name: vow.Flopper}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Flopper,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowFlopperTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("vow sin mapping", func() {
		It("writes row", func() {
			timestamp := "1538558052"
			fakeKeys := map[types.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := types.GetValueMetadata(vow.SinMapping, fakeKeys, types.Uint256)

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
			fakeKeys := map[types.Key]string{constants.Timestamp: timestamp}
			vowSinMetadata := types.GetValueMetadata(vow.SinMapping, fakeKeys, types.Uint256)
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
			malformedVowSinMappingMetadata := types.GetValueMetadata(vow.SinMapping, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedVowSinMappingMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Timestamp}))
		})
	})

	Describe("Sin integer", func() {
		metadata := types.ValueMetadata{Name: vow.SinInteger}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.SinInteger,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowSinIntegerTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Ash", func() {
		metadata := types.ValueMetadata{Name: vow.Ash}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Ash,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowAshTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Wait", func() {
		metadata := types.ValueMetadata{Name: vow.Wait}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Wait,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowWaitTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Dump", func() {
		metadata := types.ValueMetadata{Name: vow.Dump}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Dump,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowDumpTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Sump", func() {
		metadata := types.ValueMetadata{Name: vow.Sump}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Sump,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowSumpTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Bump", func() {
		metadata := types.ValueMetadata{Name: vow.Bump}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Bump,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowBumpTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Hump", func() {
		metadata := types.ValueMetadata{Name: vow.Hump}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: vow.Hump,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowHumpTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("Live", func() {
		metadata := types.ValueMetadata{Name: storage.Live}
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: storage.Live,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.VowLiveTable,
			Repository:     &repo,
			Metadata:       metadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})
})
