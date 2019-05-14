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
func CatContractAddress() string     { return getEnvironmentString("contract.address.MCD_CAT") }
func FlapperContractAddress() string { return getEnvironmentString("contract.address.mcd_flap") }
func FlipperContractAddress() string { return getEnvironmentString("contract.address.eth_flip") }
func FlopperContractAddress() string { return getEnvironmentString("contract.address.mcd_flop") }
func JugContractAddress() string     { return getEnvironmentString("contract.address.MCD_JUG") }
func PipEthContractAddress() string  { return getEnvironmentString("contract.address.PIP_ETH") }
func PipCol1ContractAddress() string { return getEnvironmentString("contract.address.PIP_COL1") }
func PipCol2ContractAddress() string { return getEnvironmentString("contract.address.PIP_COL2") }
func PipCol3ContractAddress() string { return getEnvironmentString("contract.address.PIP_COL3") }
func PipCol4ContractAddress() string { return getEnvironmentString("contract.address.PIP_COL4") }
func PipCol5ContractAddress() string { return getEnvironmentString("contract.address.PIP_COL5") }
func VatContractAddress() string     { return getEnvironmentString("contract.address.MCD_VAT") }
func VowContractAddress() string     { return getEnvironmentString("contract.address.MCD_VOW") }

func CatABI() string     { return getEnvironmentString("contract.abi.MCD_CAT") }
func FlapperABI() string { return getEnvironmentString("contract.abi.mcd_flap") }
func FlipperABI() string { return getEnvironmentString("contract.abi.eth_flip") }
func FlopperABI() string { return getEnvironmentString("contract.abi.mcd_flop") }
func JugABI() string     { return getEnvironmentString("contract.abi.MCD_JUG") }
func PipABI() string     { return getEnvironmentString("contract.abi.PIP") }
func VatABI() string     { return getEnvironmentString("contract.abi.MCD_VAT") }
func VowABI() string     { return getEnvironmentString("contract.abi.MCD_VOW") }

func CatDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.MCD_CAT") }
func FlapperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.mcd_flap") }
func FlipperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.eth_flip") }
func FlopperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.mcd_flop") }
func JugDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.MCD_JUG") }
func PipEthDeploymentBlock() int64  { return getEnvironmentInt64("contract.deployment-block.PIP_ETH") }
func PipCol1DeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.PIP_COL1") }
func PipCol2DeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.PIP_COL2") }
func PipCol3DeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.PIP_COL3") }
func PipCol4DeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.PIP_COL4") }
func PipCol5DeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.PIP_COL5") }
func VatDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.MCD_VAT") }
func VowDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.MCD_VOW") }
