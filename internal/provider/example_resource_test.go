// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccExampleResource(t *testing.T) {
	os.Setenv("TF_ACC_TERRAFORM_PATH", "/Users/sgoods/Hashicorp/terraform/terraform")
	resource.UnitTest(t, resource.TestCase{
		// Write-only attributes are only available in 1.11 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `ephemeral "scaffolding_example" "test-ephemeral" {}

				resource "scaffolding_example" "foo" {
					configurable_attribute = "test"
  					write_only_attribute = ephemeral.scaffolding_example.test-ephemeral.id
				}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue("scaffolding_example.foo", tfjsonpath.New("write_only_attribute"), knownvalue.Null()),
						plancheck.ExpectUnknownValue("scaffolding_example.foo", tfjsonpath.New("write_only_state_attribute")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("scaffolding_example.foo", tfjsonpath.New("write_only_attribute"), knownvalue.Null()),
					statecheck.ExpectKnownValue("scaffolding_example.foo", tfjsonpath.New("write_only_state_attribute"), knownvalue.NotNull()),
				},
			},
		},
	})
}
