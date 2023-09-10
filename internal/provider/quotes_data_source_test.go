// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuotesDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccQuotesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.theoffice_quotes.test", "quotes.#"),
				),
			},
		},
	})
}

func TestAccQuotesDataSource_filterByEpisode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccQuotesDataSourceConfig_filterByEpisode,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.theoffice_quotes.test", "quotes.#"),
					resource.TestCheckResourceAttr("data.theoffice_quotes.test", "quotes.0.episode", "1"),
				),
			},
		},
	})
}

const testAccQuotesDataSourceConfig = `
data "theoffice_quotes" "test" {
  season = 1
}
`

const testAccQuotesDataSourceConfig_filterByEpisode = `
data "theoffice_quotes" "test" {
  season = 1
  episode = 1
}
`
