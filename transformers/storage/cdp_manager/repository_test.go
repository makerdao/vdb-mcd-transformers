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

package cdp_manager_test

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CDP Manager storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repository           cdp_manager.StorageRepository
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repository = cdp_manager.StorageRepository{}
		repository.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("panics if the metadata name is not recognized", func() {
		unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}
		repoCreate := func() {
			repository.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
		}

		Expect(repoCreate).Should(Panic())
	})

	Describe("vat", func() {
		var vatMetadata = types.ValueMetadata{Name: cdp_manager.Vat}
		var fakeAddress = FakeAddress

		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: cdp_manager.Vat,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.CdpManagerVatTable,
			Repository:     &repository,
			Metadata:       vatMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("cdpi", func() {
		var (
			cdpiMetadata  = types.ValueMetadata{Name: cdp_manager.Cdpi}
			fakeCdpi      = strconv.Itoa(rand.Int())
			fakeTimestamp int
			header        core.Header
		)

		BeforeEach(func() {
			fakeBlockNumber := rand.Int()
			fakeTimestamp = int(rand.Int31())
			header = fakes.GetFakeHeaderWithTimestamp(int64(fakeTimestamp), int64(fakeBlockNumber))
			headerRepo := repositories.NewHeaderRepository(db)
			var headerErr error
			// TODO: don't shadow fakeHeaderID
			fakeHeaderID, headerErr = headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())
		})

		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: cdp_manager.Cdpi,
			Value:          fakeCdpi,
			Schema:         constants.MakerSchema,
			TableName:      constants.CdpManagerCdpiTable,
			Repository:     &repository,
			Metadata:       cdpiMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

		It("triggers an update to the managed_cdp table", func() {
			fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
			storageDiffRepo := storage.NewDiffRepository(db)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())

			createdTimestamp := time.Unix(int64(fakeTimestamp), 0).UTC().Format(time.RFC3339)
			expectedTimeCreated := sql.NullString{String: createdTimestamp, Valid: true}
			err := repository.Create(diffID, fakeHeaderID, cdpiMetadata, fakeCdpi)
			Expect(err).NotTo(HaveOccurred())

			var cdp test_helpers.ManagedCdp
			queryErr := db.Get(&cdp, `SELECT cdpi, created FROM api.managed_cdp`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(cdp.Id).To(Equal(fakeCdpi))
			Expect(cdp.Created).To(Equal(expectedTimeCreated))
		})
	})

	Describe("cdpi mapping tables", func() {
		fakeCdpi := strconv.Itoa(rand.Int())

		It("returns an error if mapping metadata is missing the key", func() {
			badMetadata := types.ValueMetadata{
				Name: cdp_manager.Urns,
				Keys: map[types.Key]string{},
				Type: types.Address,
			}
			err := repository.Create(diffID, fakeHeaderID, badMetadata, "")
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Cdpi}))
		})

		Describe("urns", func() {
			var fakeUrnsValue = FakeAddress
			var urnsMetadata = types.ValueMetadata{
				Name: cdp_manager.Urns,
				Keys: map[types.Key]string{constants.Cdpi: fakeCdpi},
				Type: types.Address,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Cdpi),
				ValueFieldName: "urn",
				Key:            fakeCdpi,
				Value:          fakeUrnsValue,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerUrnsTable,
				Repository:     &repository,
				Metadata:       urnsMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			It("triggers an update to the managed_cdp table", func() {
				fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
				storageDiffRepo := storage.NewDiffRepository(db)
				var insertDiffErr error
				diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
				Expect(insertDiffErr).NotTo(HaveOccurred())

				err := repository.Create(diffID, fakeHeaderID, urnsMetadata, fakeUrnsValue)
				Expect(err).NotTo(HaveOccurred())

				var cdp test_helpers.ManagedCdp
				queryErr := db.Get(&cdp, `SELECT cdpi, urn_identifier FROM api.managed_cdp`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(cdp.Id).To(Equal(fakeCdpi))
				Expect(cdp.UrnIdentifier).To(Equal(fakeUrnsValue))
			})
		})

		Describe("list_prev", func() {
			var fakePrevValue = strconv.Itoa(rand.Int())
			var prevMetadata = types.ValueMetadata{
				Name: cdp_manager.ListPrev,
				Keys: map[types.Key]string{constants.Cdpi: fakeCdpi},
				Type: types.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Cdpi),
				ValueFieldName: "prev",
				Key:            fakeCdpi,
				Value:          fakePrevValue,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerListPrevTable,
				Repository:     &repository,
				Metadata:       prevMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("list_next", func() {
			var fakeNextValue = strconv.Itoa(rand.Int())
			var nextMetadata = types.ValueMetadata{
				Name: cdp_manager.ListNext,
				Keys: map[types.Key]string{constants.Cdpi: fakeCdpi},
				Type: types.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Cdpi),
				ValueFieldName: "next",
				Key:            fakeCdpi,
				Value:          fakeNextValue,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerListNextTable,
				Repository:     &repository,
				Metadata:       nextMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("owns", func() {
			var fakeOwner = FakeAddress
			var ownsMetadata = types.ValueMetadata{
				Name: cdp_manager.Owns,
				Keys: map[types.Key]string{constants.Cdpi: fakeCdpi},
				Type: types.Address,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Cdpi),
				ValueFieldName: "owner",
				Key:            fakeCdpi,
				Value:          fakeOwner,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerOwnsTable,
				Repository:     &repository,
				Metadata:       ownsMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)

			It("triggers an update to the managed_cdp table", func() {
				fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
				storageDiffRepo := storage.NewDiffRepository(db)
				var insertDiffErr error
				diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
				Expect(insertDiffErr).NotTo(HaveOccurred())

				err := repository.Create(diffID, fakeHeaderID, ownsMetadata, fakeOwner)
				Expect(err).NotTo(HaveOccurred())

				var cdp test_helpers.ManagedCdp
				queryErr := db.Get(&cdp, `SELECT cdpi, usr FROM api.managed_cdp`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(cdp.Id).To(Equal(fakeCdpi))
				Expect(cdp.Usr).To(Equal(fakeOwner))
			})
		})

		Describe("ilks", func() {
			var (
				ilksMetadata = types.ValueMetadata{
					Name: cdp_manager.Ilks,
					Keys: map[types.Key]string{constants.Cdpi: fakeCdpi},
					Type: types.Bytes32,
				}
				fakeIlksValue = test_helpers.FakeIlk.Hex
			)

			BeforeEach(func() {
				fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
				storageDiffRepo := storage.NewDiffRepository(db)
				var insertDiffErr error
				diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
				Expect(insertDiffErr).NotTo(HaveOccurred())
			})

			It("persists a record", func() {
				createErr := repository.Create(diffID, fakeHeaderID, ilksMetadata, fakeIlksValue)
				Expect(createErr).NotTo(HaveOccurred())

				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, cdpi AS key, ilk_id AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CdpManagerIlksTable))
				readErr := db.Get(&result, query)
				Expect(readErr).NotTo(HaveOccurred())

				ilkId, ilkErr := shared.GetOrCreateIlk(fakeIlksValue, db)
				Expect(ilkErr).NotTo(HaveOccurred())

				AssertMapping(result, diffID, fakeHeaderID, fakeCdpi, strconv.FormatInt(ilkId, 10))
			})

			It("doesn't duplicate a record", func() {
				err := repository.Create(diffID, fakeHeaderID, ilksMetadata, fakeIlksValue)
				Expect(err).NotTo(HaveOccurred())

				err = repository.Create(diffID, fakeHeaderID, ilksMetadata, fakeIlksValue)
				Expect(err).NotTo(HaveOccurred())

				var count int
				query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CdpManagerIlksTable))
				err = db.Get(&count, query)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("triggers an update to the managed_cdp table", func() {
				err := repository.Create(diffID, fakeHeaderID, ilksMetadata, fakeIlksValue)
				Expect(err).NotTo(HaveOccurred())

				var cdp test_helpers.ManagedCdp
				queryErr := db.Get(&cdp, `SELECT cdpi, ilk_identifier FROM api.managed_cdp`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(cdp.Id).To(Equal(fakeCdpi))
				Expect(cdp.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
			})
		})
	})

	Describe("owner mapping tables", func() {
		fakeOwner := FakeAddress

		Describe("first", func() {
			var fakeFirstValue = strconv.Itoa(rand.Int())
			var firstMetadata = types.ValueMetadata{
				Name: cdp_manager.First,
				Keys: map[types.Key]string{constants.Owner: fakeOwner},
				Type: types.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Owner),
				ValueFieldName: "first",
				Key:            fakeOwner,
				Value:          fakeFirstValue,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerFirstTable,
				Repository:     &repository,
				Metadata:       firstMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("last", func() {
			var fakeLastValue = strconv.Itoa(rand.Int())
			var lastMetadata = types.ValueMetadata{
				Name: cdp_manager.Last,
				Keys: map[types.Key]string{constants.Owner: fakeOwner},
				Type: types.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Owner),
				ValueFieldName: "last",
				Key:            fakeOwner,
				Value:          fakeLastValue,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerLastTable,
				Repository:     &repository,
				Metadata:       lastMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("count", func() {
			var fakeCountValue = strconv.Itoa(rand.Int())
			var countMetadata = types.ValueMetadata{
				Name: cdp_manager.Count,
				Keys: map[types.Key]string{constants.Owner: fakeOwner},
				Type: types.Uint256,
			}
			inputs := shared_behaviors.StorageBehaviorInputs{
				KeyFieldName:   string(constants.Owner),
				ValueFieldName: "count",
				Key:            fakeOwner,
				Value:          fakeCountValue,
				IsAMapping:     true,
				Schema:         constants.MakerSchema,
				TableName:      constants.CdpManagerCountTable,
				Repository:     &repository,
				Metadata:       countMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})
})
