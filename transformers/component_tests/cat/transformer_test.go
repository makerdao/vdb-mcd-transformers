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

package cat

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
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
		contractAddress   = test_data.CatAddress()
		repository        = cat.CatStorageRepository{ContractAddress: contractAddress}
		storageKeysLookup = storage.NewKeysLookup(cat.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
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
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		header.Id = headerID
	})

	It("reads in a Cat Live storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		catLiveDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(catLiveDiff)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT diff_id, header_id, live AS value FROM maker.cat_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, catLiveDiff.ID, headerID, "1")
	})

	It("reads in a Cat Vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("000000000000000000000000acdd1ee0f74954ed8f0ac581b081b7b86bd6aad9")
		catVatDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(catVatDiff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.cat_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, catVatDiff.ID, headerID, "0xaCdd1ee0F74954Ed8F0aC581b081B7b86bD6aad9")
	})

	It("reads in a Cat Vow storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("00000000000000000000000021444ac712ccd21ce82af24ea1aec64cf07361d2")
		catVowDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(catVowDiff)
		Expect(err).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableRes
		err = db.Get(&vowResult, `SELECT diff_id, header_id, vow AS value FROM maker.cat_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, catVowDiff.ID, headerID, "0x21444AC712cCD21ce82AF24eA1aEc64Cf07361D2")
	})

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			catAddressID, catAddressErr := shared.GetOrCreateAddress(test_data.CatAddress(), db)
			Expect(catAddressErr).NotTo(HaveOccurred())

			userAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = catAddressID
			denyModel.ColumnValues[constants.UsrColumn] = userAddressID
			insertErr := event.PersistModels([]event.InsertionModel{denyModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("d017f9e94318721f8f69ea22f089555ab823c8cda66a25f15651a0f72ffd3d98")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var wardsResult test_helpers.WardsMappingRes
			err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
			Expect(err).NotTo(HaveOccurred())
			Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(catAddressID, 10)))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, headerID, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	Describe("ilk", func() {
		var (
			ilk    string
			ilkID  int64
			ilkErr error
		)

		BeforeEach(func() {
			ilk = "0x4554482d41000000000000000000000000000000000000000000000000000000"
			ilkID, ilkErr = shared.GetOrCreateIlk(ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
		})

		It("reads in a Cat Ilk Flip storage diff row and persists it", func() {
			key := common.HexToHash("ddedd75666d350fcd985cb35e3b9f2d4f288318d97268199e03d4405df947015")
			value := common.HexToHash("000000000000000000000000b88d2655aba486a06e638707fbebd858d430ac6e")
			catIlkFlipDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			err := transformer.Execute(catIlkFlipDiff)
			Expect(err).NotTo(HaveOccurred())

			var ilkFlipResult test_helpers.MappingRes
			err = db.Get(&ilkFlipResult, `SELECT diff_id, header_id, ilk_id AS key, flip AS value FROM maker.cat_ilk_flip`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkFlipResult, catIlkFlipDiff.ID, headerID, strconv.FormatInt(ilkID, 10), "0xB88d2655abA486A06e638707FBEbD858D430AC6E")
		})

		It("reads in a Cat Ilk Chop storage diff row and persists it", func() {
			key := common.HexToHash("ddedd75666d350fcd985cb35e3b9f2d4f288318d97268199e03d4405df947016")
			value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
			catIlkChopDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			err := transformer.Execute(catIlkChopDiff)
			Expect(err).NotTo(HaveOccurred())

			var ilkChopResult test_helpers.MappingRes
			err = db.Get(&ilkChopResult, `SELECT diff_id, header_id, ilk_id AS key, chop AS value FROM maker.cat_ilk_chop`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkChopResult, catIlkChopDiff.ID, headerID, strconv.FormatInt(ilkID, 10), "1000000000000000000000000000")
		})

		It("reads in a Cat Ilk Lump storage diff row and persists it", func() {
			key := common.HexToHash("ddedd75666d350fcd985cb35e3b9f2d4f288318d97268199e03d4405df947017")
			value := common.HexToHash("000000000000000000000006d79f82328ea3da61e066ebb2f88a000000000000")
			catIlkLumpDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			err := transformer.Execute(catIlkLumpDiff)
			Expect(err).NotTo(HaveOccurred())

			var ilkLumpResult test_helpers.MappingRes
			err = db.Get(&ilkLumpResult, `SELECT diff_id, header_id, ilk_id AS key, lump AS value FROM maker.cat_ilk_lump`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkLumpResult, catIlkLumpDiff.ID, headerID, strconv.FormatInt(ilkID, 10), "10000000000000000000000000000000000000000000000000")
		})
	})
})
