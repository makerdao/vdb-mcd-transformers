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

package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_fork"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat fork transformer", func() {
	var (
		ilkID         int64
		src           = "0x7939E55BE6A8CB8fceE57F409543E25489C06aaC"
		dst           = "0xd175Dfd88939eE702CB69e96f8bedAa2f93FBFfA"
		vatForkConfig = event.TransformerConfig{
			TransformerName:   constants.VatForkTable,
			ContractAddresses: []string{test_data.VatAddress()},
			ContractAbi:       constants.VatABI(),
			Topic:             constants.VatForkSignature(),
		}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		blockNumber := int64(9003611)
		vatForkConfig.StartingBlockNumber = blockNumber
		vatForkConfig.EndingBlockNumber = blockNumber

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      vatForkConfig,
			Transformer: vat_fork.Transformer{},
		}
		tr := initializer.NewTransformer(db)

		f := fetcher.NewLogFetcher(blockChain)
		logs, logsErr := f.FetchLogs(
			event.HexStringsToAddresses(vatForkConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatForkConfig.Topic)},
			header)
		Expect(logsErr).NotTo(HaveOccurred())
		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		executeErr := tr.Execute(eventLogs)
		Expect(executeErr).NotTo(HaveOccurred())

		var ilkErr error
		ilkID, ilkErr = shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(ilkErr).NotTo(HaveOccurred())
	})

	It("fetches and transforms vat fork event", func() {
		var dbResult vatForkModel
		err := db.Get(&dbResult, `SELECT ilk_id, src, dst, dink, dart FROM maker.vat_fork`)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult.Ilk).To(Equal(ilkID))
		Expect(dbResult.Src).To(Equal(src))
		Expect(dbResult.Dst).To(Equal(dst))
		Expect(dbResult.Dink).To(Equal("3248431462049897973"))
		Expect(dbResult.Dart).To(Equal("218121873079553101113"))
	})

	It("creates urn IDs for src and dst fields in fork payload", func() {
		// necessary so that the API knows about urns that have passed through a fork but not a frob
		var srcUrnID, dstUrnID int64
		srcErr := db.Get(&srcUrnID, `SELECT id FROM maker.urns WHERE identifier = $1 AND ilk_id = $2`, src, ilkID)
		Expect(srcErr).NotTo(HaveOccurred())
		Expect(srcUrnID).NotTo(BeZero())

		dstErr := db.Get(&dstUrnID, `SELECT id FROM maker.urns WHERE identifier = $1 AND ilk_id = $2`, dst, ilkID)
		Expect(dstErr).NotTo(HaveOccurred())
		Expect(dstUrnID).NotTo(BeZero())
	})
})

type vatForkModel struct {
	Ilk  int64 `db:"ilk_id"`
	Src  string
	Dst  string
	Dink string
	Dart string
}
