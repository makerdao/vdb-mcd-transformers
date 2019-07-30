package test_data

import "github.com/vulcanize/mcd_transformers/transformers/shared/constants"

// This file contains "shortcuts" to some configuration values useful for testing

func CatAddress() string     { return constants.GetContractAddress("MCD_CAT") }
func FlapperAddress() string { return constants.GetContractAddress("MCD_FLAP") }
func FlipperAddresses() []string {
	return constants.GetContractAddresses([]string{
		"ETH_FLIP_A", "ETH_FLIP_B", "ETH_FLIP_C",
		"MCD_FLIP_REP_A", "MCD_FLIP_ZRX_A", "MCD_FLIP_OMG_A", "MCD_FLIP_BAT_A", "MCD_FLIP_DGD_A", "MCD_FLIP_GNT_A",
	})
}
func FlopperAddress() string { return constants.GetContractAddress("MCD_FLOP") }
func JugAddress() string     { return constants.GetContractAddress("MCD_JUG") }
func SpotAddress() string    { return constants.GetContractAddress("MCD_SPOT") }
func VatAddress() string     { return constants.GetContractAddress("MCD_VAT") }
func VowAddress() string     { return constants.GetContractAddress("MCD_VOW") }

// TODO Can we just nuke these?
func OldFlapperAddress() string { return constants.GetContractAddress("MCD_FLAP_OLD") }
func OldFlipperAddress() string { return constants.GetContractAddress("MCD_FLIP_OLD") }
func OldFlopperAddress() string { return constants.GetContractAddress("MCD_FLOP_OLD") }
