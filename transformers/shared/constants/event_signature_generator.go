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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

func GetEventTopicZero(solidityEventSignature string) string {
	eventSignature := []byte(solidityEventSignature)
	hash := crypto.Keccak256Hash(eventSignature)
	return hash.Hex()
}

// Gets the first ten characters of the topic0 with zero padding; to be used with DSNote modified events.
func GetLogNoteTopicZeroWithZeroPadding(solidityFunctionSignature string) string {
	rawSignature := GetEventTopicZero(solidityFunctionSignature)
	return rawSignature[:10] + "00000000000000000000000000000000000000000000000000000000"
}

// Gets the first ten characters of the topic zero with leading zeros; to be used with the Vat contract's custom Note event modifier.
func GetLogNoteTopicZeroWithLeadingZeros(solidityFunctionSignature string) string {
	rawSignature := GetEventTopicZero(solidityFunctionSignature)
	return "0x00000000000000000000000000000000000000000000000000000000" + rawSignature[2:10]
}

func GetSolidityFunctionSignature(abi, name string) string {
	parsedAbi, _ := geth.ParseAbi(abi)

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
