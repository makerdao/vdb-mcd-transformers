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
func CatContractAddress() string     { return getEnvironmentString("contract.address.cat") }
func FlapperContractAddress() string { return getEnvironmentString("contract.address.mcd_flap") }
func FlipperContractAddress() string { return getEnvironmentString("contract.address.eth_flip") }
func FlopperContractAddress() string { return getEnvironmentString("contract.address.mcd_flop") }
func JugContractAddress() string     { return getEnvironmentString("contract.address.jug") }
func PepContractAddress() string     { return getEnvironmentString("contract.address.pep") }
func PipContractAddress() string     { return getEnvironmentString("contract.address.pip") }
func PitContractAddress() string     { return getEnvironmentString("contract.address.pit") }
func RepContractAddress() string     { return getEnvironmentString("contract.address.rep") }
func VatContractAddress() string     { return getEnvironmentString("contract.address.vat") }
func OldVatContractAddress() string  { return getEnvironmentString("contract.address.old_vat") }
func VowContractAddress() string     { return getEnvironmentString("contract.address.vow") }

func CatABI() string { return getEnvironmentString("contract.abi.cat") }

func FlapperABI() string    { return getEnvironmentString("contract.abi.mcd_flap") }
func FlipperABI() string    { return getEnvironmentString("contract.abi.eth_flip") }
func FlopperABI() string    { return getEnvironmentString("contract.abi.mcd_flop") }
func JugABI() string        { return getEnvironmentString("contract.abi.jug") }
func MedianizerABI() string { return getEnvironmentString("contract.abi.medianizer") }
func PitABI() string        { return getEnvironmentString("contract.abi.pit") }
func VatABI() string        { return getEnvironmentString("contract.abi.vat") }
func OldVatABI() string     { return getEnvironmentString("contract.abi.old_vat") }
func VowABI() string        { return getEnvironmentString("contract.abi.vow") }

func CatDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.cat") }

func FlapperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.mcd_flap") }
func FlipperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.eth_flip") }
func FlopperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.mcd_flop") }
func JugDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.jug") }
func PepDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.pep") }
func PipDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.pip") }
func PitDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.pit") }
func RepDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.rep") }
func VatDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.vat") }
func OldVatDeploymentBlock() int64  { return getEnvironmentInt64("contract.deployment-block.old_vat") }
func VowDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.vow") }
func MedianizerDeploymentBlock() int64 {
	return getEnvironmentInt64("contract.deployment-block.medianizer")
}
