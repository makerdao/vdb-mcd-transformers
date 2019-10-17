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
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                *postgres.DB
		storageKeysLookup = cat.StorageKeysLookup{StorageRepository: &mcdStorage.MakerStorageRepository{}}
		repository        = cat.CatStorageRepository{}
		contractAddress   = "81f7aa9c1570de564eb511b3a1e57dae558c65b5"
		transformer       = storage.Transformer{
			HashedAddress: utils.HexToKeccak256Hash(contractAddress),
			Mappings:      &storageKeysLookup,
			Repository:    &repository,
		}
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
	})

	It("reads in a Cat Live storage diff row and persists it", func() {
		blockNumber := 10980005
		catLineRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("82755877daf87d6eb8a228ee757450890f4b9b7bef96b749d8831fcfa466e6d7"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT block_number, block_hash, live AS value FROM maker.cat_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, blockNumber, "0x82755877daf87d6eb8a228ee757450890f4b9b7bef96b749d8831fcfa466e6d7", "1")
	})

	It("reads in a Cat Vat storage diff row and persists it", func() {
		blockNumber := 10980005
		catLineRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("82755877daf87d6eb8a228ee757450890f4b9b7bef96b749d8831fcfa466e6d7"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue:  common.HexToHash("000000000000000000000000acdd1ee0f74954ed8f0ac581b081b7b86bd6aad9"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value FROM maker.cat_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, "0x82755877daf87d6eb8a228ee757450890f4b9b7bef96b749d8831fcfa466e6d7", "0xaCdd1ee0F74954Ed8F0aC581b081B7b86bD6aad9")
	})

	It("reads in a Cat Vow storage diff row and persists it", func() {
		blockNumber := 10980005
		catLineRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("82755877daf87d6eb8a228ee757450890f4b9b7bef96b749d8831fcfa466e6d7"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004"),
			StorageValue:  common.HexToHash("00000000000000000000000021444ac712ccd21ce82af24ea1aec64cf07361d2"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableRes
		err = db.Get(&vowResult, `SELECT block_number, block_hash, vow AS value FROM maker.cat_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, blockNumber, "0x82755877daf87d6eb8a228ee757450890f4b9b7bef96b749d8831fcfa466e6d7", "0x21444AC712cCD21ce82AF24eA1aEc64Cf07361D2")
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
			blockNumber := 10980036
			catLineRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				BlockHeight:   blockNumber,
				BlockHash:     common.HexToHash("88aa77e1bdd6f06c304b8f674a10689c8b96e48deccb2e358597198e8a96a3ef"),
				StorageKey:    common.HexToHash("ddedd75666d350fcd985cb35e3b9f2d4f288318d97268199e03d4405df947015"),
				StorageValue:  common.HexToHash("000000000000000000000000b88d2655aba486a06e638707fbebd858d430ac6e"),
			}
			err := transformer.Execute(catLineRow)
			Expect(err).NotTo(HaveOccurred())

			var ilkFlipResult test_helpers.MappingRes
			err = db.Get(&ilkFlipResult, `SELECT block_number, block_hash, ilk_id AS key, flip AS value FROM maker.cat_ilk_flip`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkFlipResult, blockNumber, "0x88aa77e1bdd6f06c304b8f674a10689c8b96e48deccb2e358597198e8a96a3ef", strconv.FormatInt(ilkID, 10), "0xB88d2655abA486A06e638707FBEbD858D430AC6E")
		})

		It("reads in a Cat Ilk Chop storage diff row and persists it", func() {
			blockNumber := 10980036
			catLineRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				BlockHeight:   blockNumber,
				BlockHash:     common.HexToHash("88aa77e1bdd6f06c304b8f674a10689c8b96e48deccb2e358597198e8a96a3ef"),
				StorageKey:    common.HexToHash("ddedd75666d350fcd985cb35e3b9f2d4f288318d97268199e03d4405df947016"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000"),
			}
			err := transformer.Execute(catLineRow)
			Expect(err).NotTo(HaveOccurred())

			var ilkChopResult test_helpers.MappingRes
			err = db.Get(&ilkChopResult, `SELECT block_number, block_hash, ilk_id AS key, chop AS value FROM maker.cat_ilk_chop`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkChopResult, blockNumber, "0x88aa77e1bdd6f06c304b8f674a10689c8b96e48deccb2e358597198e8a96a3ef", strconv.FormatInt(ilkID, 10), "1000000000000000000000000000")
		})

		It("reads in a Cat Ilk Lump storage diff row and persists it", func() {
			blockNumber := 10980036
			catLineRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				BlockHeight:   blockNumber,
				BlockHash:     common.HexToHash("88aa77e1bdd6f06c304b8f674a10689c8b96e48deccb2e358597198e8a96a3ef"),
				StorageKey:    common.HexToHash("ddedd75666d350fcd985cb35e3b9f2d4f288318d97268199e03d4405df947017"),
				StorageValue:  common.HexToHash("000000000000000000000006d79f82328ea3da61e066ebb2f88a000000000000"),
			}
			err := transformer.Execute(catLineRow)
			Expect(err).NotTo(HaveOccurred())

			var ilkLumpResult test_helpers.MappingRes
			err = db.Get(&ilkLumpResult, `SELECT block_number, block_hash, ilk_id AS key, lump AS value FROM maker.cat_ilk_lump`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkLumpResult, blockNumber, "0x88aa77e1bdd6f06c304b8f674a10689c8b96e48deccb2e358597198e8a96a3ef", strconv.FormatInt(ilkID, 10), "10000000000000000000000000000000000000000000000000")
		})
	})
})
