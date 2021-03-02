package integration_tests

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/fetcher"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rely transformer", func() {
	Context("Cat rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xa9ee75d81d78c36c4163004e6cc7a988eec9433e"
		relyIntegrationTest(8928165, test_data.Cat100Address(), msgSenderAddress, usrAddress)
	})

	Context("Flap v1.0.0 rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xd27a5f3416d8791fc238c148c93630d9e3c882e5"
		relyIntegrationTest(8928163, test_data.FlapV100Address(), msgSenderAddress, usrAddress)
	})

	Context("Flap v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510870, test_data.FlapV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip AAVE v1.2.2 rely events", func() {
		usrAddress := "0xBE8E3E3618F7474F8CB1D074A26AFFEF007E98FB"
		msgSenderAddress := "0xda0fab060e6cc7b1C0AA105d29Bd50D71f036711"
		relyIntegrationTest(11471534, test_data.FlipAaveAV122Address(), msgSenderAddress,
			usrAddress)
	})

	Context("Flip BAL v1.1.14 rely events", func() {
		usrAddress := "0xBE8E3E3618F7474F8CB1D074A26AFFEF007E98FB"
		msgSenderAddress := "0xDa0FaB05039809e63C5D068c897c3e602fA97457"
		relyIntegrationTest(11198257, test_data.FlipBalAV1114Address(), msgSenderAddress,
			usrAddress)
	})

	Context("Flip BAT v1.0.0 rely events", func() {
		usrAddress := "0x9BdDB99625A711bf9bda237044924E34E8570f75"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(9684438, test_data.FlipBatV100Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip BAT v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510871, test_data.FlipBatV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip COMP v1.1.2 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xDa0FaB05039809e63C5D068c897c3e602fA97457"
		relyIntegrationTest(10932692, test_data.FlipCompV112Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip ETH v1.0.0 rely events", func() {
		usrAddress := "0xbaa65281c2FA2baAcb2cb550BA051525A480D3F4"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(8928180, test_data.FlipEthAV100Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip ETH v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510871, test_data.FlipEthAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip ETH_B v1.1.3 rely events", func() {
		usrAddress := "0xc4bE7F74Ee3743bDEd8E0fA218ee5cf06397f472"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(11086830, test_data.FlipEthBV113Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip GUSD_A v1.1.5 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xDa0FaB05039809e63C5D068c897c3e602fA97457"
		relyIntegrationTest(11239632, test_data.FlipGusdAV115Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip KNC-A v1.0.8 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10323433, test_data.FlipKncAV108Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip KNC-A v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510886, test_data.FlipKncAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip LINK v1.1.2 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xDa0FaB05039809e63C5D068c897c3e602fA97457"
		relyIntegrationTest(10932697, test_data.FlipLinkV112Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip LRC v1.1.2 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xDa0FaB05039809e63C5D068c897c3e602fA97457"
		relyIntegrationTest(10932697, test_data.FlipLrcV112Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip MANA-A v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510886, test_data.FlipManaAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip PAXUSD-A v1.1.1 rely events", func() {
		usrAddress := "0xc4bE7F74Ee3743bDEd8E0fA218ee5cf06397f472"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(10821399, test_data.FlipPaxusdAV111Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip RENBTC-A v1.2.1 rely events", func() {
		usrAddress := "0xc4bE7F74Ee3743bDEd8E0fA218ee5cf06397f472"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(11451553, test_data.FlipRenbtcA121Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip TUSD-A v1.0.7 rely events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(10144451, test_data.FlipTusdAV107Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip TUSD-A v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510886, test_data.FlipTusdAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2AAVEETH v1.2.7 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		relyIntegrationTest(11882283, test_data.FlipUniV2AaveEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2DAIETH v1.2.2 rely events", func() {
		usrAddress := "0xc4bE7F74Ee3743bDEd8E0fA218ee5cf06397f472"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(11497506, test_data.FlipUniV2DaiEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2DAIUSDC_A v1.2.5 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		relyIntegrationTest(11791219, test_data.FlipUniV2DaiUsdcAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2DAIUSDT-A v1.2.8 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		relyIntegrationTest(11927179, test_data.FlipUniV2DaiUsdtAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2ETHUSDT_A v1.2.5 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		relyIntegrationTest(11791266, test_data.FlipUniV2EthUsdtAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2USDCETH_A rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab060e6cc7b1C0AA105d29Bd50D71f036711"
		relyIntegrationTest(11700542, test_data.FlipUniV2UsdcEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2UNIETH_A v1.2.6 rely events", func() {
		usrAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		msgSenderAddress := "0x4acdbe9dd0d00b36ec2050e805012b8fc9974f2b"
		relyIntegrationTest(11836943, test_data.FlipUniV2UniEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2LINKETH_A v1.2.6 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		relyIntegrationTest(11837043, test_data.FlipUniV2LinkEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2WBTCDAI_A rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab0700a4389f6e6679abab1692b4601ce9bf"
		relyIntegrationTest(11882322, test_data.FlipUniV2WbtcDaiAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNIV2WBTCETH_A rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab060e6cc7b1C0AA105d29Bd50D71f036711"
		relyIntegrationTest(11700516, test_data.FlipUniV2WbtcEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Flip UNI-A v1.2.1 rely events", func() {
		usrAddress := "0xc4bE7F74Ee3743bDEd8E0fA218ee5cf06397f472"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(11451553, test_data.FlipUniAV121Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip USDC-A v1.0.4 rely events", func() {
		usrAddress := "0x9BdDB99625A711bf9bda237044924E34E8570f75"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(9686502, test_data.FlipUsdcAV104Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip USDC-A v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510885, test_data.FlipUsdcAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip USDC-B v1.0.7 rely events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(10144450, test_data.FlipUsdcBV107Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip USDC-B v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510885, test_data.FlipUsdcBV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip USDT-A v1.1.1 rely events", func() {
		usrAddress := "0xc4bE7F74Ee3743bDEd8E0fA218ee5cf06397f472"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(10821399, test_data.FlipUsdtAV111Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip WBTC v1.0.6 rely events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(9975625, test_data.FlipWbtcAV106Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip WBTC v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510885, test_data.FlipWbtcAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip YFI v1.1.14 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xda0fab05039809e63c5d068c897c3e602fa97457"
		relyIntegrationTest(11198283, test_data.FlipYfiAV1114Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip ZRX-A v1.0.8 rely events", func() {
		usrAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		msgSenderAddress := "0xBAB4FbeA257ABBfe84F4588d4Eedc43656E46Fc5"
		relyIntegrationTest(10323245, test_data.FlipZrxAV108Address(), msgSenderAddress, usrAddress)
	})

	Context("Flip ZRX-A v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510886, test_data.FlipZrxAV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Flop v1.0.1 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(9008136, test_data.FlopV101Address(), msgSenderAddress, usrAddress)
	})

	Context("Flop v1.0.9 rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10510871, test_data.FlopV109Address(), msgSenderAddress, usrAddress)
	})

	Context("Jug rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0x45f0a929889ec8cc2d5b8cd79ab55e3279945cde"
		relyIntegrationTest(8928160, test_data.JugAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median AAVE-A rely events", func() {
		usrAddress := "0xBE8E3E3618F7474F8CB1D074A26AFFEF007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(11472295, test_data.MedianAaveAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median BAL rely events", func() {
		usrAddress := "0xBE8E3E3618F7474F8CB1D074A26AFFEF007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(11199049, test_data.MedianBalAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median BAT rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(8956961, test_data.MedianBatAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median COMP rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225d1f3ea4385627efdc5b4709cab4a21c"
		relyIntegrationTest(10933587, test_data.MedianCompAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median ETH rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(8956896, test_data.MedianEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median KNC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10323303, test_data.MedianKncAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median LINK rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(10933602, test_data.MedianLinkAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median LRC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(10933624, test_data.MedianLrcAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median MANA rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(10511187, test_data.MedianManaAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median UNI-A rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(11424316, test_data.MedianUniAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median USDT rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(10740129, test_data.MedianUsdtAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median WBTC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(8956963, test_data.MedianWbtcAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median YFI rely events", func() {
		usrAddress := "0xBE8E3E3618F7474F8CB1D074A26AFFEF007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(11199055, test_data.MedianYfiAddress(), msgSenderAddress, usrAddress)
	})

	Context("Median ZRX rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10323394, test_data.MedianZrxAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM AAVE-A rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225d1f3ea4385627efdc5b4709cab4a21c"
		relyIntegrationTest(11472399, test_data.OsmAaveAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM BAL rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225d1f3ea4385627efdc5b4709cab4a21c"
		relyIntegrationTest(11198841, test_data.OsmBalAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM BAT rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(9529100, test_data.OsmBatAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM COMP rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(10933766, test_data.OsmCompAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM ETH rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(9529100, test_data.OsmEthAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM KNC rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(10352556, test_data.OsmKncAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM LINK rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225D1F3eA4385627eFDC5B4709Cab4A21c"
		relyIntegrationTest(10933793, test_data.OsmLinkAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM LRC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0x0048d6225d1f3ea4385627efdc5b4709cab4a21c"
		relyIntegrationTest(10933928, test_data.OsmLrcAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM MANA rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10516692, test_data.OsmManaAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM UNI-A rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(11451553, test_data.OsmUniAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM USDT rely events", func() {
		usrAddress := "0x76416A4d5190d071bfed309861527431304aA14f"
		msgSenderAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		relyIntegrationTest(10821399, test_data.OsmUsdtAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM WBTC rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(9975536, test_data.OsmWbtcAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM YFI rely events", func() {
		usrAddress := "0xbe8e3e3618f7474f8cb1d074a26affef007e98fb"
		msgSenderAddress := "0x0048d6225d1f3ea4385627efdc5b4709cab4a21c"
		relyIntegrationTest(11199166, test_data.OsmYfiAddress(), msgSenderAddress, usrAddress)
	})

	Context("OSM ZRX rely events", func() {
		usrAddress := "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
		msgSenderAddress := "0xdDb108893104dE4E1C6d0E47c42237dB4E617ACc"
		relyIntegrationTest(10323394, test_data.OsmZrxAddress(), msgSenderAddress, usrAddress)
	})

	Context("Pot rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0x1a5ee7c64cf874c735968e3a42fa13f1c03427f9"
		relyIntegrationTest(8928160, test_data.PotAddress(), msgSenderAddress, usrAddress)
	})

	Context("Spot rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0xdedd12bcb045c02b2fe11031c2b269bcde457410"
		relyIntegrationTest(8928152, test_data.SpotAddress(), msgSenderAddress, usrAddress)
	})

	Context("Vow rely events", func() {
		usrAddress := "0xbaa65281c2fa2baacb2cb550ba051525a480d3f4"
		msgSenderAddress := "0x68322ca1a9aeb8c1d610b5fc8a8920aa0fba423b"
		relyIntegrationTest(8928163, test_data.VowAddress(), msgSenderAddress, usrAddress)
	})
})

func relyIntegrationTest(blockNumber int64, contractAddressHex, msgSenderAddressHex, usrAddressHex string) {
	It("persists event", func() {
		test_config.CleanTestDB(db)
		logFetcher := fetcher.NewLogFetcher(blockChain)
		relyConfig := event.TransformerConfig{
			TransformerName: constants.RelyTable,
			Topic:           constants.RelySignature(),
		}
		initializer := event.ConfiguredTransformer{
			Config:      relyConfig,
			Transformer: auth.Transformer{TableName: constants.RelyTable},
		}

		header, headerErr := persistHeader(db, blockNumber, blockChain)
		Expect(headerErr).NotTo(HaveOccurred())

		initializer.Config.ContractAddresses = []string{contractAddressHex}
		initializer.Config.StartingBlockNumber = blockNumber
		initializer.Config.EndingBlockNumber = blockNumber

		address := common.HexToAddress(contractAddressHex)
		topics := []common.Hash{common.HexToHash(relyConfig.Topic)}

		logs, logsErr := logFetcher.FetchLogs([]common.Address{address}, topics, header)
		Expect(logsErr).NotTo(HaveOccurred())

		eventLogs := test_data.CreateLogs(header.Id, logs, db)

		relyTransformer := initializer.NewTransformer(db)
		transformErr := relyTransformer.Execute(eventLogs)
		Expect(transformErr).NotTo(HaveOccurred())

		var dbResult []relyModel
		err := db.Select(&dbResult, `SELECT address_id, msg_sender, usr FROM maker.rely ORDER BY log_id`)
		Expect(err).NotTo(HaveOccurred())

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddressHex)
		Expect(contractAddressErr).NotTo(HaveOccurred())
		msgSenderAddressID, msgSenderAddressErr := repository.GetOrCreateAddress(db, msgSenderAddressHex)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())
		usrAddressID, usrAddressErr := repository.GetOrCreateAddress(db, usrAddressHex)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		var matchFound bool
		for _, result := range dbResult {
			if result.AddressID == contractAddressID &&
				result.MsgSender == msgSenderAddressID &&
				result.Usr == usrAddressID {
				matchFound = true
			}
		}

		Expect(matchFound).To(BeTrue(), getRelyFailureMessage(contractAddressHex, blockNumber))
	})
}

type relyModel struct {
	Usr       int64 `db:"usr"`
	MsgSender int64 `db:"msg_sender"`
	AddressID int64 `db:"address_id"`
}

func getRelyFailureMessage(contractAddress string, blockNumber int64) string {
	failureMsgToFmt := "no matching rely event found for contract %s at block %d"
	return fmt.Sprintf(failureMsgToFmt, contractAddress, blockNumber)
}
