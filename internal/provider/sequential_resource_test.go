// Copyright (c) EcmaXp.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSequentialResource(t *testing.T) {
	resourceAddress := "phaser_sequential.test"
	buildConfig := func(phases []string) string {
		return fmt.Sprintf(
			`resource "phaser_sequential" "test" { phases = ["%s"] }`,
			strings.Join(phases, "\", \""),
		)
	}

	for _, tc := range []struct {
		configs             []string            // optional, default=phases
		preSteps            []resource.TestStep // optional, default=[]
		phases              []string
		expectPhases        []string // optional, default=phases
		expectNonEmptyPlans []bool
	}{
		{
			phases:              []string{"running"},
			expectNonEmptyPlans: []bool{false},
		},
		{
			phases:              []string{"prepare", "ready"},
			expectNonEmptyPlans: []bool{true, false},
		},
		{
			phases:              []string{"prepare", "ready", "running"},
			expectNonEmptyPlans: []bool{true, true, false},
		},
		{
			phases:              []string{"running"},
			expectPhases:        []string{"running", "running"},
			expectNonEmptyPlans: []bool{false, false},
		},
		{
			phases:              []string{"prepare", "ready", "running"},
			expectPhases:        []string{"prepare", "ready", "running", "running", "running"},
			expectNonEmptyPlans: []bool{true, true, false, false, false},
		},
		{
			configs: []string{
				buildConfig([]string{"prepare"}),
				buildConfig([]string{"prepare", "ready"}),
				buildConfig([]string{"ready"}),
				buildConfig([]string{"ready", "running"}),
				buildConfig([]string{"ready", "running"}),
			},
			expectPhases:        []string{"prepare", "ready", "ready", "running", "running"},
			expectNonEmptyPlans: []bool{false, false, false, false, false, true},
		},
		{
			preSteps: []resource.TestStep{
				{
					Config:        buildConfig([]string{"prepare"}),
					ImportState:   true,
					ImportStateId: "prepare",
					ResourceName:  resourceAddress,
				},
			},
			phases:              []string{"prepare", "ready", "running"},
			expectPhases:        []string{"prepare", "ready", "running"},
			expectNonEmptyPlans: []bool{true, true, false},
		},
	} {
		expectPhases := tc.expectPhases
		if expectPhases == nil {
			expectPhases = tc.phases
		}

		configs := tc.configs
		if configs == nil {
			config := buildConfig(tc.phases)
			for range expectPhases {
				configs = append(configs, config)
			}
		}

		var preChecks []statecheck.StateCheck
		if tc.phases != nil {
			preChecks = append(preChecks, buildCheckPhases(resourceAddress, tc.phases))
		}

		steps := tc.preSteps
		for i, phase := range expectPhases {
			checks := append(
				preChecks,
				buildCheckPhase(resourceAddress, phase),
			)

			steps = append(steps, resource.TestStep{
				Config:             configs[i],
				ConfigStateChecks:  checks,
				ExpectNonEmptyPlan: tc.expectNonEmptyPlans[i],
			})
		}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps:                    steps,
		})
	}
}

func buildCheckPhases(resourceAddress string, phases []string) statecheck.StateCheck {
	return statecheck.ExpectKnownValue(
		resourceAddress,
		tfjsonpath.New("phases"),
		knownvalue.ListExact(
			func() (checks []knownvalue.Check) {
				for _, phase := range phases {
					checks = append(checks, knownvalue.StringExact(phase))
				}
				return
			}(),
		),
	)
}

func buildCheckPhase(resourceAddress string, phase string) statecheck.StateCheck {
	return statecheck.ExpectKnownValue(
		resourceAddress,
		tfjsonpath.New("phase"),
		knownvalue.StringExact(phase),
	)
}
