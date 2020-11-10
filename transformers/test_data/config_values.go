package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

// This file contains "shortcuts" to some configuration values useful for testing

func Cat100Address() string   { return checksum(constants.GetContractAddress("MCD_CAT_1_0_0")) }
func Cat110Address() string   { return checksum(constants.GetContractAddress("MCD_CAT_1_1_0")) }
func FlapV100Address() string { return checksum(constants.GetContractAddress("MCD_FLAP_1_0_0")) }
func FlapV109Address() string { return checksum(constants.GetContractAddress("MCD_FLAP_1_0_9")) }
func Flip100Addresses() []string {
	var addressesResult []string
	flipAddresses := constants.GetContractAddresses([]string{
		"MCD_FLIP_BAT_A_1_0_0",
		"MCD_FLIP_BAT_A_1_0_9",
		"MCD_FLIP_ETH_A_1_0_0",
		"MCD_FLIP_ETH_A_1_0_9",
		"MCD_FLIP_KNC_A_1_0_8",
		"MCD_FLIP_KNC_A_1_0_9",
		"MCD_FLIP_MANA_A_1_0_9",
		"MCD_FLIP_SAI_1_0_0",
		"MCD_FLIP_TUSD_A_1_0_7",
		"MCD_FLIP_TUSD_A_1_0_9",
		"MCD_FLIP_USDC_A_1_0_4",
		"MCD_FLIP_USDC_A_1_0_9",
		"MCD_FLIP_USDC_B_1_0_7",
		"MCD_FLIP_USDC_B_1_0_9",
		"MCD_FLIP_WBTC_A_1_0_6",
		"MCD_FLIP_WBTC_A_1_0_9",
		"MCD_FLIP_ZRX_A_1_0_8",
		"MCD_FLIP_ZRX_A_1_0_9",
	})

	for _, address := range flipAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}

func Flip110Addresses() []string {
	var addressesResult []string
	flipAddresses := constants.GetContractAddresses([]string{
		"MCD_FLIP_BAL_A_1_1_14",
		"MCD_FLIP_BAT_A_1_1_0",
		"MCD_FLIP_COMP_A_1_1_2",
		"MCD_FLIP_ETH_A_1_1_0",
		"MCD_FLIP_ETH_B_1_1_3",
		"MCD_FLIP_LINK_A_1_1_2",
		"MCD_FLIP_LRC_A_1_1_2",
		"MCD_FLIP_KNC_A_1_1_0",
		"MCD_FLIP_MANA_A_1_1_0",
		"MCD_FLIP_PAXUSD_A_1_1_1",
		"MCD_FLIP_TUSD_A_1_1_0",
		"MCD_FLIP_USDC_A_1_1_0",
		"MCD_FLIP_USDC_B_1_1_0",
		"MCD_FLIP_USDT_A_1_1_1",
		"MCD_FLIP_WBTC_A_1_1_0",
		"MCD_FLIP_ZRX_A_1_1_0",
	})

	for _, address := range flipAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}

func FlipBalV110Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAL_A_1_1_14"))
}
func FlipBatV100Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAT_A_1_0_0"))
}
func FlipBatV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAT_A_1_0_9"))
}
func FlipBatV110Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAT_A_1_1_0"))
}
func FlipCompV112Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_COMP_A_1_1_2"))
}
func FlipEthAV100Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A_1_0_0"))
}
func FlipEthAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A_1_0_9"))
}
func FlipEthAV110Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A_1_1_0"))
}
func FlipEthBV113Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_B_1_1_3"))
}
func FlipKncAV108Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_KNC_A_1_0_8"))
}
func FlipKncAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_KNC_A_1_0_9"))
}
func FlipLinkV112Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_LINK_A_1_1_2"))
}
func FlipLrcV112Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_LRC_A_1_1_2"))
}
func FlipManaAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_MANA_A_1_0_9"))
}
func FlipPaxusdAV111Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_PAXUSD_A_1_1_1"))
}
func FlipTusdAV107Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_TUSD_A_1_0_7"))
}
func FlipTusdAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_TUSD_A_1_0_9"))
}
func FlipUsdcAV104Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_A_1_0_4"))
}
func FlipUsdcAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_A_1_0_9"))
}
func FlipUsdcBV107Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_B_1_0_7"))
}
func FlipUsdcBV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_B_1_0_9"))
}
func FlipUsdtAV111Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDT_A_1_1_1"))
}
func FlipWbtcAV106Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_WBTC_A_1_0_6"))
}
func FlipWbtcAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_WBTC_A_1_0_9"))
}
func FlipZrxAV108Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ZRX_A_1_0_8"))
}
func FlipZrxAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ZRX_A_1_0_9"))
}
func FlopV101Address() string { return checksum(constants.GetContractAddress("MCD_FLOP_1_0_1")) }
func FlopV109Address() string { return checksum(constants.GetContractAddress("MCD_FLOP_1_0_9")) }
func JugAddress() string      { return checksum(constants.GetContractAddress("MCD_JUG")) }
func MedianAddresses() []string {
	var addressesResult []string
	medianAddresses := constants.GetContractAddresses([]string{
		"MEDIAN_BAL",
		"MEDIAN_BAT",
		"MEDIAN_COMP",
		"MEDIAN_ETH",
		"MEDIAN_KNC",
		"MEDIAN_LINK",
		"MEDIAN_LRC",
		"MEDIAN_MANA",
		"MEDIAN_WBTC",
		"MEDIAN_ZRX",
	})

	for _, address := range medianAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}
func MedianBalAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_BAL")) }
func MedianBatAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_BAT")) }
func MedianCompAddress() string { return checksum(constants.GetContractAddress("MEDIAN_COMP")) }
func MedianEthAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_ETH")) }
func MedianKncAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_KNC")) }
func MedianLinkAddress() string { return checksum(constants.GetContractAddress("MEDIAN_LINK")) }
func MedianLrcAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_LRC")) }
func MedianManaAddress() string { return checksum(constants.GetContractAddress("MEDIAN_MANA")) }
func MedianUsdtAddress() string { return checksum(constants.GetContractAddress("MEDIAN_USDT")) }
func MedianWbtcAddress() string { return checksum(constants.GetContractAddress("MEDIAN_WBTC")) }
func MedianZrxAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_ZRX")) }
func OasisAddresses() []string {
	var addressesResult []string
	oasisAddresses := constants.GetContractAddresses([]string{"OASIS_MATCHING_MARKET_ONE", "OASIS_MATCHING_MARKET_TWO"})

	for _, address := range oasisAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult

}
func OsmAddresses() []string {
	var addressesResult []string
	// Does not include OSM_USDC since that's actually just a DSValue contract right now, not an OSM
	osmAddresses := constants.GetContractAddresses([]string{
		"OSM_BAL",
		"OSM_BAT",
		"OSM_COMP",
		"OSM_ETH",
		"OSM_KNC",
		"OSM_LINK",
		"OSM_LRC",
		"OSM_MANA",
		"OSM_WBTC",
		"OSM_ZRX",
	})

	for _, address := range osmAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}
func OsmBalAddress() string     { return checksum(constants.GetContractAddress("OSM_BAL")) }
func OsmBatAddress() string     { return checksum(constants.GetContractAddress("OSM_BAT")) }
func OsmCompAddress() string    { return checksum(constants.GetContractAddress("OSM_COMP")) }
func OsmEthAddress() string     { return checksum(constants.GetContractAddress("OSM_ETH")) }
func OsmKncAddress() string     { return checksum(constants.GetContractAddress("OSM_KNC")) }
func OsmLinkAddress() string    { return checksum(constants.GetContractAddress("OSM_LINK")) }
func OsmLrcAddress() string     { return checksum(constants.GetContractAddress("OSM_LRC")) }
func OsmManaAddress() string    { return checksum(constants.GetContractAddress("OSM_MANA")) }
func OsmUsdtAddress() string    { return checksum(constants.GetContractAddress("OSM_USDT")) }
func OsmWbtcAddress() string    { return checksum(constants.GetContractAddress("OSM_WBTC")) }
func OsmZrxAddress() string     { return checksum(constants.GetContractAddress("OSM_ZRX")) }
func PotAddress() string        { return checksum(constants.GetContractAddress("MCD_POT")) }
func SpotAddress() string       { return checksum(constants.GetContractAddress("MCD_SPOT")) }
func VatAddress() string        { return checksum(constants.GetContractAddress("MCD_VAT")) }
func VowAddress() string        { return checksum(constants.GetContractAddress("MCD_VOW")) }
func CdpManagerAddress() string { return checksum(constants.GetContractAddress("CDP_MANAGER")) }

func checksum(addressString string) string {
	return common.HexToAddress(addressString).Hex()
}
