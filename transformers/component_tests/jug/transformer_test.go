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

package jug

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
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
		storageKeysLookup = storage.NewKeysLookup(jug.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = jug.JugStorageRepository{}
		contractAddress   = "25a008bf942ce6d5b362f91ed7ae3e4104286a12"
		transformer       = storage.Transformer{
			HashedAddress:     utils.HexToKeccak256Hash(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		headerID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		ilk := "0x4554480000000000000000000000000000000000000000000000000000000000"
		ilkID, err = shared.GetOrCreateIlk(ilk, db)
		Expect(err).NotTo(HaveOccurred())
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Jug Vat storage diff row and persists it", func() {
		jugVatRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue:  common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(jugVatRow)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT header_id, vat AS value FROM maker.jug_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, headerID, "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Jug Vow storage diff row and persists it", func() {
		jugVowRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue:  common.HexToHash("17560834075da3db54f737db74377e799c865821000000000000000000000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(jugVowRow)
		Expect(err).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableRes
		err = db.Get(&vowResult, `SELECT header_id, vow AS value FROM maker.jug_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, headerID, "0x17560834075da3db54f737db74377e799c865821000000000000000000000000")
	})

	It("reads in a Jug Ilk Duty storage diff row and persists it", func() {
		jugIlkDutyRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97a"),
			StorageValue:  common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(jugIlkDutyRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkDutyResult test_helpers.MappingRes
		err = db.Get(&ilkDutyResult, `SELECT header_id, ilk_id AS key, duty AS value FROM maker.jug_ilk_duty`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkDutyResult, headerID, strconv.FormatInt(ilkID, 10), "1000000000000000000000000000")
	})

	It("reads in a Jug Ilk Rho storage diff row and persists it", func() {
		jugIlkRhoRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97b"),
			StorageValue:  common.HexToHash("000000000000000000000000000000000000000000000000000000005c812808"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(jugIlkRhoRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkRhoResult test_helpers.MappingRes
		err = db.Get(&ilkRhoResult, `SELECT header_id, ilk_id AS key, rho AS value FROM maker.jug_ilk_rho`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkRhoResult, headerID, strconv.FormatInt(ilkID, 10), "1551968264")
	})
})
