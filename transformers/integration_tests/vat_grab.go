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
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_grab"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat Grab Transformer", func() {
	vatGrabConfig := transformer.EventTransformerConfig{
		TransformerName:   constants.VatGrabLabel,
		ContractAddresses: []string{test_data.VatAddress()},
		ContractAbi:       constants.VatABI(),
		Topic:             constants.VatGrabSignature(),
	}

	It("transforms VatGrab log events", func() {
		blockNumber := int64(14783840)
		vatGrabConfig.StartingBlockNumber = blockNumber
		vatGrabConfig.EndingBlockNumber = blockNumber

		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err := getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())

		db := test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			transformer.HexStringsToAddresses(vatGrabConfig.ContractAddresses),
			[]common.Hash{common.HexToHash(vatGrabConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		headerSyncLogs := test_data.CreateLogs(header.Id, logs, db)

		tr := shared.EventTransformer{
			Config:     vatGrabConfig,
			Converter:  &vat_grab.VatGrabConverter{},
			Repository: &vat_grab.VatGrabRepository{},
		}.NewEventTransformer(db)

		err = tr.Execute(headerSyncLogs)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatGrabModel
		err = db.Select(&dbResult, `SELECT urn_id, v, w, dink, dart from maker.vat_grab`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		urnID, err := shared.GetOrCreateUrn("0xA17B62E20d2d700950cf95ceE5d8cC910850Dd98",
			"0x4554482d41000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Urn).To(Equal(strconv.FormatInt(urnID, 10)))
		Expect(dbResult[0].V).To(Equal("0xAB10DFC4578EE6f9389c3c3F5F010CF9df30ea2B")) //cat contract address as bytes32
		Expect(dbResult[0].W).To(Equal("0x6F29dfc2f7142C6f161abC0E5bE6E79A41269Fa9"))
		expectedDink := new(big.Int)
		expectedDink.SetString("-170000000000000000", 10)
		Expect(dbResult[0].Dink).To(Equal(expectedDink.String()))
		expectedDart := new(big.Int)
		expectedDart.SetString("-20990706494289402873", 10)
		Expect(dbResult[0].Dart).To(Equal(expectedDart.String()))
	})
})

type vatGrabModel struct {
	Ilk  string
	Urn  string `db:"urn_id"`
	V    string
	W    string
	Dink string
	Dart string
}
