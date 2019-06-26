// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package constants

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var initialized = false

var TTL = int64(10800) // 60 * 60 * 3 == 10800 seconds == 3 hours

func initConfig() {
	if initialized {
		return
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Sprintf("Could not find environment file: %v", err))
	}
	initialized = true
}

func getEnvironmentString(key string) string {
	initConfig()
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("No environment configuration variable set for key: \"%v\"", key))
	}
	return value
}

func getEnvironmentInt64(key string) int64 {
	initConfig()
	value := viper.GetInt64(key)
	if value == -1 {
		panic(fmt.Sprintf("No environment configuration variable set for key: \"%v\"", key))
	}
	return value
}

// Getters for contract addresses from environment files
func CatContractAddress() string        { return getEnvironmentString("contract.address.MCD_CAT") }
func OldFlapperContractAddress() string { return getEnvironmentString("contract.address.MCD_FLAP_OLD") }
func FlapperContractAddress() string    { return getEnvironmentString("contract.address.MCD_FLAP") }
func OldFlipperContractAddress() string { return getEnvironmentString("contract.address.ETH_FLIP_OLD") }
func FlopperContractAddress() string    { return getEnvironmentString("contract.address.MCD_FLOP_OLD") }
func JugContractAddress() string        { return getEnvironmentString("contract.address.MCD_JUG") }
func SpotContractAddress() string       { return getEnvironmentString("contract.address.MCD_SPOT") }
func VatContractAddress() string        { return getEnvironmentString("contract.address.MCD_VAT") }
func VowContractAddress() string        { return getEnvironmentString("contract.address.MCD_VOW") }
func EthFlipContractAddressA() string   { return getEnvironmentString("contract.address.ETH_FLIP_A") }
func EthFlipContractAddressB() string   { return getEnvironmentString("contract.address.ETH_FLIP_B") }
func EthFlipContractAddressC() string   { return getEnvironmentString("contract.address.ETH_FLIP_C") }
func Col1FlipContractAddress() string   { return getEnvironmentString("contract.address.COL1_FLIP") }
func Col2FlipContractAddress() string   { return getEnvironmentString("contract.address.COL2_FLIP") }
func Col3FlipContractAddress() string   { return getEnvironmentString("contract.address.COL3_FLIP") }
func Col4FlipContractAddress() string   { return getEnvironmentString("contract.address.COL4_FLIP") }
func Col5FlipContractAddress() string   { return getEnvironmentString("contract.address.COL5_FLIP") }
func FlipperContractAddresses() []string {
	return []string{EthFlipContractAddressA(), EthFlipContractAddressB(), EthFlipContractAddressC(), Col1FlipContractAddress(), Col2FlipContractAddress(), Col3FlipContractAddress(), Col4FlipContractAddress(), Col5FlipContractAddress()}
}

func CatABI() string        { return getEnvironmentString("contract.abi.MCD_CAT") }
func OldFlapperABI() string { return getEnvironmentString("contract.abi.MCD_FLAP_OLD") }
func FlapperABI() string    { return getEnvironmentString("contract.abi.MCD_FLAP") }
func OldFlipperABI() string { return getEnvironmentString("contract.abi.ETH_FLIP_OLD") }
func FlipperABI() string    { return getEnvironmentString("contract.abi.MCD_FLIP") }
func FlopperABI() string    { return getEnvironmentString("contract.abi.MCD_FLOP_OLD") }
func JugABI() string        { return getEnvironmentString("contract.abi.MCD_JUG") }
func SpotABI() string       { return getEnvironmentString("contract.abi.MCD_SPOT") }
func VatABI() string        { return getEnvironmentString("contract.abi.MCD_VAT") }
func VowABI() string        { return getEnvironmentString("contract.abi.MCD_VOW") }

func CatDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.MCD_CAT") }
func OldFlapperDeploymentBlock() int64 {
	return getEnvironmentInt64("contract.deployment-block.MCD_FLAP_OLD")
}
func FlapperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.MCD_FLAP") }
func OldFlipperDeploymentBlock() int64 {
	return getEnvironmentInt64("contract.deployment-block.ETH_FLIP_OLD")
}
func FlipperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.MCD_FLIP") }
func FlopperDeploymentBlock() int64 {
	return getEnvironmentInt64("contract.deployment-block.MCD_FLOP_OLD")
}
func JugDeploymentBlock() int64  { return getEnvironmentInt64("contract.deployment-block.MCD_JUG") }
func SpotDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.MCD_SPOT") }
func VatDeploymentBlock() int64  { return getEnvironmentInt64("contract.deployment-block.MCD_VAT") }
func VowDeploymentBlock() int64  { return getEnvironmentInt64("contract.deployment-block.MCD_VOW") }
