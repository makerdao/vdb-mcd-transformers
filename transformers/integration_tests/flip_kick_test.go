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
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlipKick Event Transformer", func() {
	Context("MCD_FLIP_BAT_A_1.0.0", func() {
		expected := flipKickModel{
			BidId: "473",
			Lot:   "455405056000000000000",
			Bid:   "0",
			Tab:   "84750000000000000000297127152220292821786803075",
			Usr:   "0xfFBFD66b36D9B58e9dA63817E89F40CF5c74C789",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10431926, test_data.FlipBatV100Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_BAT_A_1.0.9", func() {
		expected := flipKickModel{
			BidId: "6",
			Lot:   "27500000000000000000000",
			Bid:   "0",
			Tab:   "7006052974872711082626367362638391740684916432363",
			Usr:   "0x7D5BeC38D8DA48846EADF410a601f96Fee047FC5",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10742760, test_data.FlipBatV109Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_BAT_A_1.1.0", func() {
		expected := flipKickModel{
			BidId: "1",
			Lot:   "2479962706275246500000",
			Bid:   "0",
			Tab:   "581950000000000000001041312137300424782413684325",
			Usr:   "0x8B7b68B93Cb709976F4feFdC05408039e9927246",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10782907, test_data.FlipBatV110Address(), constants.FlipV110ABI(), expected)
	})

	Context("MCD_FLIP_ETH_A_1.0.0", func() {
		expected := flipKickModel{
			BidId: "116",
			Lot:   "50000000000000000000",
			Bid:   "0",
			Tab:   "5046619216084543990261356563876808629308883826941",
			Usr:   "0x0A051CD913dFD1820dbf87a9bf62B04A129F88A5",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(8997383, test_data.FlipEthAV100Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_ETH_A_1.0.9", func() {
		expected := flipKickModel{
			BidId: "55",
			Lot:   "1235952758123546544",
			Bid:   "0",
			Tab:   "395500000000000000000858799990562985454201380269",
			Usr:   "0xA0be0B36Cfd52C3D8c6ABF08f61306fa52c9Ec20",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10767181, test_data.FlipEthAV109Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_ETH_A_1.1.0", func() {
		expected := flipKickModel{
			BidId: "4",
			Lot:   "128149858188440510103",
			Bid:   "0",
			Tab:   "45397083100174537229400150129540694835373946249282",
			Usr:   "0x62a2851197a313901990D277e32cEdEc0905ECF5",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10780886, test_data.FlipEthAV110Address(), constants.FlipV110ABI(), expected)
	})

	Context("MCD_FLIP_KNC_A_1.0.8", func() {
		expected := flipKickModel{
			BidId: "5",
			Lot:   "25000000000000000000",
			Bid:   "0",
			Tab:   "22643013812771213593715910213152135973231480498",
			Usr:   "0xAd0B5a59BeD3acc083Fb7A875EBc99d7bDD0a34E",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10542247, test_data.FlipKncAV108Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_KNC_A_1.0.9", func() {
		expected := flipKickModel{
			BidId: "1",
			Lot:   "25000000000000000000",
			Bid:   "0",
			Tab:   "24001291519747595527495780820793681538737102619",
			Usr:   "0x26E3C7bE1c16bD0B7Db44D10b2367614f1c6aCe3",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10550327, test_data.FlipKncAV109Address(), constants.FlipV100ABI(), expected)
	})

	// TODO: add MCD_FLIP_KNC_A_1.1.0 test when available

	Context("MCD_FLIP_MANA_A_1.0.9", func() {
		expected := flipKickModel{
			BidId: "2",
			Lot:   "10391000000000000000000",
			Bid:   "0",
			Tab:   "566274443571663195468127549713779774444398081796",
			Usr:   "0xB84f7D2eF7Baa0189120B2e829518216160636D5",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10688829, test_data.FlipManaAV109Address(), constants.FlipV100ABI(), expected)
	})

	// TODO: add MCD_FLIP_MANA_A_1.1.0 test when available

	// Note: no auctions for stablecoin collateral types (SAI, TUSD, USDC)

	Context("MCD_FLIP_WBTC_A_1.0.6", func() {
		expected := flipKickModel{
			BidId: "3",
			Lot:   "9963900000000000",
			Bid:   "0",
			Tab:   "68512815152442944620090016739680878351903013931",
			Usr:   "0xac8c465EAFC13F2933DA6D63375C574477Ab0085",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10397960, test_data.FlipWbtcAV106Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_WBTC_A_1.0.9", func() {
		expected := flipKickModel{
			BidId: "1",
			Lot:   "282645530000000000",
			Bid:   "0",
			Tab:   "2429500000000000000000362179051175015558856619323",
			Usr:   "0xD32eaAeA81f26F0f5aDa251BDeCBe9353fAd8763",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10731373, test_data.FlipWbtcAV109Address(), constants.FlipV100ABI(), expected)
	})

	// TODO: add MCD_FLIP_WBTC_A_1.1.0 test when available

	Context("MCD_FLIP_ZRX_A_1.0.8", func() {
		expected := flipKickModel{
			BidId: "1",
			Lot:   "120000000000000000000",
			Bid:   "0",
			Tab:   "25576507885258353714710545074204251538458988726",
			Usr:   "0xced5A2c5BA030e224420C1749c96bAaC417EdADB",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10353299, test_data.FlipZrxAV108Address(), constants.FlipV100ABI(), expected)
	})

	Context("MCD_FLIP_ZRX_A_1.0.9", func() {
		expected := flipKickModel{
			BidId: "2",
			Lot:   "93000000000000000000",
			Bid:   "0",
			Tab:   "22831002453619890592757908844533369358717753741",
			Usr:   "0x2A4ed588c523Cf7F8003c83419d8cEfd18eE599B",
			Gal:   "0xA950524441892A31ebddF91d3cEEFa04Bf454466",
		}
		flipKickIntegrationTest(10570493, test_data.FlipZrxAV109Address(), constants.FlipV100ABI(), expected)
	})

	// TODO: add MCD_FLIP_ZRX_A_1.1.0 test when available
})

func flipKickIntegrationTest(blockNumber int64, contractAddressHex, contractABI string, expectedResult flipKickModel) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		flipKickConfig := event.TransformerConfig{
			ContractAbi:         contractABI,
			ContractAddresses:   []string{contractAddressHex},
			EndingBlockNumber:   blockNumber,
			StartingBlockNumber: blockNumber,
			Topic:               constants.FlipKickSignature(),
			TransformerName:     constants.FlipKickTable,
		}

		header, err := persistHeader(db, blockNumber, blockChain)
		Expect(err).NotTo(HaveOccurred())

		initializer := event.ConfiguredTransformer{
			Config:      flipKickConfig,
			Transformer: flip_kick.Transformer{},
		}
		transformer := initializer.NewTransformer(db)

		logFetcher := fetcher.NewLogFetcher(blockChain)
		logs, err := logFetcher.FetchLogs(
			[]common.Address{common.HexToAddress(flipKickConfig.ContractAddresses[0])},
			[]common.Hash{common.HexToHash(flipKickConfig.Topic)},
			header)
		Expect(err).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		execErr := transformer.Execute(eventLogs)
		Expect(execErr).NotTo(HaveOccurred())

		var dbResult flipKickModel
		err = db.Get(&dbResult, `SELECT bid_id, lot, bid, tab, usr, gal, address_id FROM maker.flip_kick`)
		Expect(err).NotTo(HaveOccurred())

		flipContractAddressId, addrErr := shared.GetOrCreateAddress(contractAddressHex, db)
		Expect(addrErr).NotTo(HaveOccurred())
		expectedResult.AddressId = flipContractAddressId

		Expect(dbResult).To(Equal(expectedResult))
	})
}

type flipKickModel struct {
	BidId     string `db:"bid_id"`
	Lot       string
	Bid       string
	Tab       string
	Usr       string
	Gal       string
	AddressId int64 `db:"address_id"`
}
