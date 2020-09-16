package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

// This file contains "shortcuts" to some configuration values useful for testing

func Cat100Address() string   { return checksum(constants.GetContractAddress("MCD_CAT_1.0.0")) }
func Cat110Address() string   { return checksum(constants.GetContractAddress("MCD_CAT_1.1.0")) }
func FlapV100Address() string { return checksum(constants.GetContractAddress("MCD_FLAP_1.0.0")) }
func FlapV109Address() string { return checksum(constants.GetContractAddress("MCD_FLAP_1.0.9")) }
func FlipAddresses() []string {
	var addressesResult []string
	flipAddresses := constants.GetContractAddresses([]string{
		"MCD_FLIP_BAT_A_1.0.0",
		"MCD_FLIP_BAT_A_1.0.9",
		"MCD_FLIP_ETH_A_1.0.0",
		"MCD_FLIP_ETH_A_1.0.9",
		"MCD_FLIP_KNC_A_1.0.8",
		"MCD_FLIP_KNC_A_1.0.9",
		"MCD_FLIP_MANA_A_1.0.9",
		"MCD_FLIP_SAI_1.0.0",
		"MCD_FLIP_TUSD_A_1.0.7",
		"MCD_FLIP_TUSD_A_1.0.9",
		"MCD_FLIP_USDC_A_1.0.4",
		"MCD_FLIP_USDC_A_1.0.9",
		"MCD_FLIP_USDC_B_1.0.7",
		"MCD_FLIP_USDC_B_1.0.9",
		"MCD_FLIP_WBTC_A_1.0.6",
		"MCD_FLIP_WBTC_A_1.0.9",
		"MCD_FLIP_ZRX_A_1.0.8",
		"MCD_FLIP_ZRX_A_1.0.9",
	})

	for _, address := range flipAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}

func FlipBatV100Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAT_A_1.0.0"))
}
func FlipBatV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAT_A_1.0.9"))
}
func FlipBatV110Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_BAT_A_1.1.0"))
}
func FlipEthV100Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A_1.0.0"))
}
func FlipEthV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A_1.0.9"))
}
func FlipEthV110Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A_1.1.0"))
}
func FlipKncAV108Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_KNC_A_1.0.8"))
}
func FlipKncAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_KNC_A_1.0.9"))
}
func FlipManaAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_MANA_A_1.0.9"))
}
func FlipTusdAV107Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_TUSD_A_1.0.7"))
}
func FlipTusdAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_TUSD_A_1.0.9"))
}
func FlipUsdcAV104Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_A_1.0.4"))
}
func FlipUsdcAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_A_1.0.9"))
}
func FlipUsdcBV107Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_B_1.0.7"))
}
func FlipUsdcBV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_USDC_B_1.0.9"))
}
func FlipWbtcAV106Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_WBTC_A_1.0.6"))
}
func FlipWbtcAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_WBTC_A_1.0.9"))
}
func FlipZrxAV108Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ZRX_A_1.0.8"))
}
func FlipZrxAV109Address() string {
	return checksum(constants.GetContractAddress("MCD_FLIP_ZRX_A_1.0.9"))
}
func FlopV101Address() string { return checksum(constants.GetContractAddress("MCD_FLOP_1.0.1")) }
func FlopV109Address() string { return checksum(constants.GetContractAddress("MCD_FLOP_1.0.9")) }
func JugAddress() string      { return checksum(constants.GetContractAddress("MCD_JUG")) }
func MedianAddresses() []string {
	var addressesResult []string
	medianAddresses := constants.GetContractAddresses([]string{
		"MEDIAN_BAT",
		"MEDIAN_ETH",
		"MEDIAN_KNC",
		"MEDIAN_MANA",
		"MEDIAN_WBTC",
		"MEDIAN_ZRX",
	})

	for _, address := range medianAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}
func MedianBatAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_BAT")) }
func MedianEthAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_ETH")) }
func MedianKncAddress() string  { return checksum(constants.GetContractAddress("MEDIAN_KNC")) }
func MedianManaAddress() string { return checksum(constants.GetContractAddress("MEDIAN_MANA")) }
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
		"OSM_BAT",
		"OSM_ETH",
		"OSM_KNC",
		"OSM_MANA",
		"OSM_WBTC",
		"OSM_ZRX",
	})

	for _, address := range osmAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}
func OsmBatAddress() string     { return checksum(constants.GetContractAddress("OSM_BAT")) }
func OsmEthAddress() string     { return checksum(constants.GetContractAddress("OSM_ETH")) }
func OsmKncAddress() string     { return checksum(constants.GetContractAddress("OSM_KNC")) }
func OsmManaAddress() string    { return checksum(constants.GetContractAddress("OSM_MANA")) }
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
