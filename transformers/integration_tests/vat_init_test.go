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
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_init"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VatInit EventTransformer", func() {
	vatInitConfig := event.TransformerConfig{
		TransformerName:   constants.VatInitTable,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatInitSignature(),
	}

	It("transforms vat init log events", func() {
		blockNumber := int64(8928180)
		vatInitConfig.StartingBlockNumber = blockNumber
		vatInitConfig.EndingBlockNumber = blockNumber

		test_config.CleanTestDB(db)
		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			event.HexStringsToAddresses(vatInitConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatInitConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		transformer := event.ConfiguredTransformer{
			Config:      vatInitConfig,
			Transformer: vat_init.Transformer{},
		}.NewTransformer(db)

		err = transformer.Execute(eventLogs)
		Expect(err).NotTo(HaveOccurred())

		var ilkID int64
		err = db.Get(&ilkID, `SELECT ilk_id from maker.vat_init`)
		Expect(err).NotTo(HaveOccurred())

		expectedIlkID, ilkErr := shared.GetOrCreateIlk("0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(ilkErr).NotTo(HaveOccurred())

		Expect(ilkID).To(Equal(expectedIlkID))
	})
})
