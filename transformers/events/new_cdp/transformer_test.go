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

package new_cdp_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/new_cdp"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCdp Transformer", func() {
	var transformer = new_cdp.Transformer{}

	It("converts a log to a Model", func() {
		models, err := transformer.ToModels("", []core.EventLog{test_data.NewCdpEventLog}, nil)
		Expect(err).NotTo(HaveOccurred())

		Expect(models).To(ConsistOf(test_data.NewCdpModel()))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := transformer.ToModels("", []core.EventLog{{}}, nil)
		Expect(err).To(HaveOccurred())
	})
})
