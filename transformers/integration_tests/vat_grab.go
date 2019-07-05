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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_grab"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdConstants "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var _ = XDescribe("Vat Grab Transformer", func() {
	vatGrabConfig := transformer.EventTransformerConfig{
		TransformerName:   mcdConstants.VatGrabLabel,
		ContractAddresses: []string{mcdConstants.VatContractAddress()},
		ContractAbi:       mcdConstants.VatABI(),
		Topic:             mcdConstants.VatGrabSignature(),
	}

	// TODO: Replace block number once there's a grab event on the updated Vat
	It("transforms VatGrab log events", func() {
		blockNumber := int64(8958230)
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

		tr := shared.LogNoteTransformer{
			Config:     vatGrabConfig,
			Converter:  &vat_grab.VatGrabConverter{},
			Repository: &vat_grab.VatGrabRepository{},
		}.NewLogNoteTransformer(db)

		err = tr.Execute(logs, header)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []vatGrabModel
		err = db.Select(&dbResult, `SELECT urn_id, v, w, dink, dart from maker.vat_grab`)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(1))
		urnID, err := shared.GetOrCreateUrn("0000000000000000000000006a3ae20c315e845b2e398e68effe39139ec6060c",
			"0x5245500000000000000000000000000000000000000000000000000000000000", db)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult[0].Urn).To(Equal(strconv.Itoa(urnID)))
		Expect(dbResult[0].V).To(Equal("0000000000000000000000002f34f22a00ee4b7a8f8bbc4eaee1658774c624e0")) //cat contract address as bytes32
		Expect(dbResult[0].W).To(Equal("0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1"))
		expectedDink := new(big.Int)
		expectedDink.SetString("115792089237316195423570985008687907853269984665640564039455584007913129639936", 10)
		Expect(dbResult[0].Dink).To(Equal(expectedDink.String()))
		expectedDart := new(big.Int)
		expectedDart.SetString("115792089237316195423570985008687907853269984665640564039441803007913129639936", 10)
		Expect(dbResult[0].Dart).To(Equal(expectedDart.String()))
		Expect(dbResult[0].TransactionIndex).To(Equal(uint(0)))
	})
})

type vatGrabModel struct {
	Ilk              string
	Urn              string `db:"urn_id"`
	V                string
	W                string
	Dink             string
	Dart             string
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
