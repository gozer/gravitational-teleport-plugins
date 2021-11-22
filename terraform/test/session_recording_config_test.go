/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func (s *TerraformSuite) TestSessionRecordingConfig() {
	res := "teleport_session_recording_config"

	create := s.terraformConfig + `
		resource "` + res + `" "test" {
			metadata {
				labels = {
					"example" = "yes"
				}
			}
							
			spec {
				proxy_checks_host_keys = true
			}			
		}
	`

	update := s.terraformConfig + `
		resource "` + res + `" "test" {
			metadata {
				labels = {
					"example" = "no"
				}
			}
			
			spec {
				proxy_checks_host_keys = false
			}
		}
	`
	name := res + ".test"

	resource.Test(s.T(), resource.TestCase{
		ProviderFactories:         s.terraformProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: create,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "kind", "session_recording_config"),
					resource.TestCheckResourceAttr(name, "spec.0.proxy_checks_host_keys", "true"),
				),
			},
			{
				Config:   create, // Check that there is no state drift
				PlanOnly: true,
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "kind", "session_recording_config"),
					resource.TestCheckResourceAttr(name, "spec.0.proxy_checks_host_keys", "false"),
				),
			},
			{
				Config:   update, // Check that there is no state drift
				PlanOnly: true,
			},
		},
	})
}