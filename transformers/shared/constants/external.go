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

package constants

import (
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var initialized = false

func initConfig() {
	if initialized {
		return
	}

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Sprintf("Could not find environment file: %v", err))
	}
	initialized = true
}

func getEnvironmentString(key string) string {
	initConfig()
	value := viper.GetString(key)
	if value == "" {
		logrus.Fatalf("No environment configuration variable set for key: \"%v\"", key)
	}
	return value
}

/* Returns all contract config names from transformer configuration:
[exporter.vow_file]
	path = "transformers/events/vow_file/initializer"
	type = "eth_event"
	repository = "github.com/makerdao/vdb-mcd-transformers"
	migrations = "db/migrations"
	contracts = ["MCD_VOW"]   <----
	rank = "0"
*/
func GetTransformerContractNames(transformerLabel string) []string {
	initConfig()
	configKey := "exporter." + transformerLabel + ".contracts"
	contracts := viper.GetStringSlice(configKey)
	if len(contracts) == 0 {
		logrus.Fatalf("No contracts configured for transformer: \"%v\"", transformerLabel)
	}
	return contracts
}

// GetABIFromContractsWithMatchingABI gets the ABI for multiple contracts from config
// Makes sure the ABI matches for all, since a single transformer may run against many contracts.
func GetABIFromContractsWithMatchingABI(contractNames []string) string {
	initConfig()
	if len(contractNames) < 1 {
		logrus.Fatalf("No contracts to get ABI for")
	}
	contractABI := getContractABI(contractNames[0])
	parsedABI, parseErr := eth.ParseAbi(contractABI)
	if parseErr != nil {
		panic(fmt.Sprintf("unable to parse ABI for %s", contractNames[0]))
	}
	for _, contractName := range contractNames[1:] {
		nextABI := getContractABI(contractName)
		nextParsedABI, nextParseErr := eth.ParseAbi(nextABI)
		if nextParseErr != nil {
			panic(fmt.Sprintf("unable to parse ABI for %s", contractName))
		}
		if !parsedABIsAreEqual(parsedABI, nextParsedABI) {
			panic(fmt.Sprintf("ABIs don't match for contracts: %s", contractNames))
		}
	}
	return contractABI
}

// GetFirstABI returns the ABI from the first contract in a collection in config
func GetFirstABI(contractNames []string) string {
	initConfig()
	if len(contractNames) < 1 {
		logrus.Fatalf("No contracts to get ABI for")
	}
	return getContractABI(contractNames[0])
}

func getContractABI(contractName string) string {
	configKey := "contract." + contractName + ".abi"
	contractABI := viper.GetString(configKey)
	if contractABI == "" {
		logrus.Fatalf("No ABI configured for contract: \"%v\"", contractName)
	}
	return contractABI
}

// Get the minimum deployment block for multiple contracts from config
func GetMinDeploymentBlock(contractNames []string) int64 {
	if len(contractNames) < 1 {
		logrus.Fatalf("No contracts supplied")
	}
	initConfig()
	minBlock := int64(math.MaxInt64)
	for _, c := range contractNames {
		deployed := getDeploymentBlock(c)
		if deployed < minBlock {
			minBlock = deployed
		}
	}
	return minBlock
}

func getDeploymentBlock(contractName string) int64 {
	configKey := "contract." + contractName + ".deployed"
	value := viper.GetInt64(configKey)
	if value == -1 {
		logrus.Infof("No deployment block configured for contract \"%v\", defaulting to 0.", contractName)
		return 0
	}
	return value
}

// Get the addresses for multiple contracts from config
func GetContractAddresses(contractNames []string) (addresses []string) {
	if len(contractNames) < 1 {
		logrus.Fatalf("No contracts supplied")
	}
	initConfig()
	for _, contractName := range contractNames {
		addresses = append(addresses, GetContractAddress(contractName))
	}
	return
}

func GetContractAddress(contract string) string {
	return getEnvironmentString("contract." + contract + ".address")
}

func parsedABIsAreEqual(a, b abi.ABI) bool {
	if a.Constructor.String() != b.Constructor.String() {
		return false
	}
OuterMethods:
	for ak, av := range a.Methods {
		for bk, bv := range b.Methods {
			if ak == bk {
				if av.String() == bv.String() {
					continue OuterMethods
				}
			}
		}
		return false
	}
OuterEvents:
	for ak, av := range a.Events {
		for bk, bv := range b.Events {
			if ak == bk && av.String() == bv.String() {
				continue OuterEvents
			}
		}
		return false
	}
	return true
}
