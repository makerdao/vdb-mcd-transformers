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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/vulcanize/vulcanizedb/pkg/eth"
)

func getEventTopicZero(solidityEventSignature string) string {
	eventSignature := []byte(solidityEventSignature)
	hash := crypto.Keccak256Hash(eventSignature)
	return hash.Hex()
}

func getLogNoteTopicZero(solidityFunctionSignature string) string {
	rawSignature := getEventTopicZero(solidityFunctionSignature)
	return "0x" + rawSignature[2:10] + "00000000000000000000000000000000000000000000000000000000"
}

func getSolidityFunctionSignature(abi, name string) string {
	parsedAbi, _ := eth.ParseAbi(abi)

	if method, ok := parsedAbi.Methods[name]; ok {
		return method.Sig()
	} else if event, ok := parsedAbi.Events[name]; ok {
		return getEventSignature(event)
	}
	panic("Error: could not get Solidity method signature for: " + name)
}

func getEventSignature(event abi.Event) string {
	types := make([]string, len(event.Inputs))
	for i, input := range event.Inputs {
		types[i] = input.Type.String()
		i++
	}

	return fmt.Sprintf("%v(%v)", event.Name, strings.Join(types, ","))
}

type ContractMethod struct {
	Name   string
	Inputs []MethodInput
}

type MethodInput struct {
	Type string
}

func getOverloadedFunctionSignature(rawAbi, name string, paramTypes []string) string {
	result, err := findSignatureInAbi(rawAbi, name, paramTypes)
	if err != nil {
		panic(err)
	}
	return result
}

func findSignatureInAbi(rawAbi, name string, paramTypes []string) (string, error) {
	contractMethods := make([]ContractMethod, 0)
	err := json.Unmarshal([]byte(rawAbi), &contractMethods)
	if err != nil {
		return "", errors.New("unable to parse ABI")
	}
	signature := fmt.Sprintf("%v(%v)", name, strings.Join(paramTypes, ","))
	if containsMatchingMethod(contractMethods, name, paramTypes) == false {
		return "", errors.New("method " + signature + " does not exist in ABI")
	}
	return signature, nil
}

func containsMatchingMethod(methods []ContractMethod, name string, paramTypes []string) bool {
	for _, method := range methods {
		if method.Name == name && hasMatchingParams(method, paramTypes) {
			return true
		}
	}
	return false
}

func hasMatchingParams(method ContractMethod, expectedParamTypes []string) bool {
	params := method.Inputs
	actualParamTypes := make([]string, len(params))
	for i, param := range params {
		actualParamTypes[i] = param.Type
	}
	return areEqual(expectedParamTypes, actualParamTypes)
}

func areEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}
