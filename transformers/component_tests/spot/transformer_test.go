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

package spot

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		err               error
		ilkID             int64
		storageKeysLookup = storage.NewKeysLookup(spot.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = spot.SpotStorageRepository{}
		contractAddress   = "a57d4123c8a80ac410e924df9d5e47765ffd1375"
		transformer       = storage.Transformer{
			HashedAddress:     vdbStorage.HexToKeccak256Hash(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		headerID int64
		header   = fakes.FakeHeader
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		ilk := "0x434f4c352d410000000000000000000000000000000000000000000000000000"
		ilkID, err = shared.GetOrCreateIlk(ilk, db)
		Expect(err).NotTo(HaveOccurred())
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		header.Id = headerID
	})

	It("reads in a Spot Vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("00000000000000000000000057aa8b02f5d3e28371fedcf672c8668869f9aac7")
		spotVatDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)
		err := transformer.Execute(spotVatDiff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.spot_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, spotVatDiff.ID, headerID, "0x57aA8B02F5D3E28371FEdCf672C8668869f9AAC7")
	})

	It("reads in a Spot Par storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		spotParDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)
		err := transformer.Execute(spotParDiff)
		Expect(err).NotTo(HaveOccurred())

		var parResult test_helpers.VariableRes
		err = db.Get(&parResult, `SELECT diff_id, header_id, par AS value FROM maker.spot_par`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(parResult, spotParDiff.ID, headerID, "1000000000000000000000000000")
	})

	It("reads in a Spot Live storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		spotLiveDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(spotLiveDiff)
		Expect(err).NotTo(HaveOccurred())
		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT diff_id, header_id, live AS value FROM maker.spot_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, spotLiveDiff.ID, headerID, "1")
	})

	It("reads in a Spot Ilk Pip storage diff row and persists it", func() {
		key := common.HexToHash("1730ac98111482efebd8acadb14d7fa301298e0d95bf3c34c3378ef524670bc6")
		value := common.HexToHash("000000000000000000000000a53e6efb4cbed841eace02220498860905e94998")
		spotIlkPipDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(spotIlkPipDiff)
		Expect(err).NotTo(HaveOccurred())

		var ilkPipResult test_helpers.MappingRes
		err = db.Get(&ilkPipResult, `SELECT diff_id, header_id, ilk_id AS key, pip AS value FROM maker.spot_ilk_pip`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkPipResult, spotIlkPipDiff.ID, headerID, strconv.FormatInt(ilkID, 10), "0xA53e6EFB4cBeD841Eace02220498860905E94998")
	})

	It("reads in a Spot Ilk Mat storage diff row and persists it", func() {
		key := common.HexToHash("1730ac98111482efebd8acadb14d7fa301298e0d95bf3c34c3378ef524670bc7")
		value := common.HexToHash("000000000000000000000000000000000000000006765c793fa10079d0000000")
		spotIlkMatDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(spotIlkMatDiff)
		Expect(err).NotTo(HaveOccurred())

		var ilkRhoResult test_helpers.MappingRes
		err = db.Get(&ilkRhoResult, `SELECT diff_id, header_id, ilk_id AS key, mat AS value FROM maker.spot_ilk_mat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkRhoResult, spotIlkMatDiff.ID, headerID, strconv.FormatInt(ilkID, 10), "2000000000000000000000000000")
	})
})
