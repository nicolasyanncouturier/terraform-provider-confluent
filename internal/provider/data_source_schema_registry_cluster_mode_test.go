// Copyright 2021 Confluent Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"github.com/walkerus/go-wiremock"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	SchemaRegistryClusterModeDataSourceScenarioName           = "confluent_schema_registry_cluster_mode Data Source Lifecycle"
	testNumberOfSchemaRegistryClusterModeDataSourceAttributes = 5
)

var fullSchemaRegistryClusterModeDataSourceLabel = fmt.Sprintf("data.confluent_schema_registry_cluster_mode.%s", testSchemaResourceLabel)

func TestAccDataSchemaRegistryClusterModeSchema(t *testing.T) {
	mockSchemaTestServerUrl := tc.wiremockUrl
	confluentCloudBaseUrl := ""
	wiremockClient := wiremock.NewClient(mockSchemaTestServerUrl)
	// nolint:errcheck
	defer wiremockClient.Reset()

	// nolint:errcheck
	defer wiremockClient.ResetAllScenarios()

	readCreatedSchemaRegistryClusterModeResponse, _ := ioutil.ReadFile("../testdata/schema_registry_cluster_mode/read_created_schema_registry_cluster_mode.json")
	_ = wiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo(updateSchemaRegistryClusterModePath)).
		InScenario(SchemaRegistryClusterModeDataSourceScenarioName).
		WhenScenarioStateIs(wiremock.ScenarioStateStarted).
		WillReturn(
			string(readCreatedSchemaRegistryClusterModeResponse),
			contentTypeJSONHeader,
			http.StatusOK,
		))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		// https://www.terraform.io/docs/extend/testing/acceptance-tests/teststep.html
		// https://www.terraform.io/docs/extend/best-practices/testing.html#built-in-patterns
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSchemaRegistryClusterModeDataSourceConfig(confluentCloudBaseUrl, mockSchemaTestServerUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(fullSchemaRegistryClusterModeDataSourceLabel),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "id", testStreamGovernanceClusterId),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "schema_registry_cluster.#", "1"),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "schema_registry_cluster.0.%", "1"),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "schema_registry_cluster.0.id", testStreamGovernanceClusterId),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "rest_endpoint", mockSchemaTestServerUrl),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "credentials.#", "1"),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "credentials.0.%", "2"),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "credentials.0.key", testSchemaRegistryKey),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "credentials.0.secret", testSchemaRegistrySecret),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "mode", testSchemaRegistryClusterMode),
					resource.TestCheckResourceAttr(fullSchemaRegistryClusterModeDataSourceLabel, "%", strconv.Itoa(testNumberOfSchemaRegistryClusterModeDataSourceAttributes)),
				),
			},
		},
	})
}

func testAccCheckSchemaRegistryClusterModeDataSourceConfig(confluentCloudBaseUrl, mockServerUrl string) string {
	return fmt.Sprintf(`
	provider "confluent" {
      endpoint = "%s"
    }
	data "confluent_schema_registry_cluster_mode" "%s" {
	  schema_registry_cluster {
        id = "%s"
      }
      rest_endpoint = "%s"
      credentials {
        key = "%s"
        secret = "%s"
	  }
	}
	`, confluentCloudBaseUrl, testSchemaResourceLabel, testStreamGovernanceClusterId, mockServerUrl, testSchemaRegistryKey, testSchemaRegistrySecret)
}
