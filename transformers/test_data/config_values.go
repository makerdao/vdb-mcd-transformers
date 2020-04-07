package test_data

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
)

// This file contains "shortcuts" to some configuration values useful for testing

func CatAddress() string  { return checksum(constants.GetContractAddress("MCD_CAT")) }
func FlapAddress() string { return checksum(constants.GetContractAddress("MCD_FLAP")) }
func FlipAddresses() []string {
	var addressesResult []string
	flipAddresses := constants.GetContractAddresses([]string{
		"MCD_FLIP_ETH_A", "MCD_FLIP_BAT_A", "MCD_FLIP_SAI",
	})

	for _, address := range flipAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}

func EthFlipAddress() string   { return checksum(constants.GetContractAddress("MCD_FLIP_ETH_A")) }
func EthMedianAddress() string { return checksum(constants.GetContractAddress("MEDIAN_ETH")) }
func EthOsmAddress() string    { return checksum(constants.GetContractAddress("OSM_ETH")) }
func FlopAddress() string      { return checksum(constants.GetContractAddress("MCD_FLOP")) }
func JugAddress() string       { return checksum(constants.GetContractAddress("MCD_JUG")) }
func MedianAddresses() []string {
	var addressesResult []string
	medianAddresses := constants.GetContractAddresses([]string{"MEDIAN_ETH", "MEDIAN_BAT"})

	for _, address := range medianAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}
func OasisAddresses() []string {
	var addressesResult []string
	oasisAddresses := constants.GetContractAddresses([]string{"OASIS_ONE", "OASIS_TWO"})

	for _, address := range oasisAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult

}
func OsmAddresses() []string {
	var addressesResult []string
	osmAddresses := constants.GetContractAddresses([]string{"OSM_ETH", "OSM_BAT"})

	for _, address := range osmAddresses {
		addressesResult = append(addressesResult, checksum(address))
	}
	return addressesResult
}
func PotAddress() string        { return checksum(constants.GetContractAddress("MCD_POT")) }
func SpotAddress() string       { return checksum(constants.GetContractAddress("MCD_SPOT")) }
func VatAddress() string        { return checksum(constants.GetContractAddress("MCD_VAT")) }
func VowAddress() string        { return checksum(constants.GetContractAddress("MCD_VOW")) }
func CdpManagerAddress() string { return checksum(constants.GetContractAddress("CDP_MANAGER")) }

func checksum(addressString string) string {
	return common.HexToAddress(addressString).Hex()
}
