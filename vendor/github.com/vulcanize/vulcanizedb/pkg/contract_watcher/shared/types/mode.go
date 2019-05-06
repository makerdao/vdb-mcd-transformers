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

package types

import "fmt"

type Mode int

const (
	HeaderSync Mode = iota
	FullSync
)

func (mode Mode) IsValid() bool {
	return mode >= HeaderSync && mode <= FullSync
}

func (mode Mode) String() string {
	switch mode {
	case HeaderSync:
		return "header"
	case FullSync:
		return "full"
	default:
		return "unknown"
	}
}

func (mode Mode) MarshalText() ([]byte, error) {
	switch mode {
	case HeaderSync:
		return []byte("header"), nil
	case FullSync:
		return []byte("full"), nil
	default:
		return nil, fmt.Errorf("contract watcher: unknown mode %d, want HeaderSync or FullSync", mode)
	}
}

func (mode *Mode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "header":
		*mode = HeaderSync
	case "full":
		*mode = FullSync
	default:
		return fmt.Errorf(`contract watcher: unknown mode %q, want "header" or "full"`, text)
	}
	return nil
}
