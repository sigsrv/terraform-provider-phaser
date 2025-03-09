// Copyright (c) EcmaXp.
// SPDX-License-Identifier: MPL-2.0

package phaser

import (
	"testing"
)

func TestGetNextPhaseSequential(t *testing.T) {
	t.Helper()

	for _, tc := range []struct {
		Name         string
		CurrentPhase string
		Phases       []string
		ExpectPhases []string
	}{
		{
			Name:         "single_phase",
			Phases:       []string{"running"},
			ExpectPhases: []string{"running"},
		},
		{
			Name:         "single_phases_with_multiple_running",
			Phases:       []string{"running"},
			ExpectPhases: []string{"running", "running"},
		},
		{
			Name:         "two_phases",
			Phases:       []string{"prepare", "ready"},
			ExpectPhases: []string{"ready", "ready"},
		},
		{
			Name:         "three_phases",
			Phases:       []string{"prepare", "ready", "running"},
			ExpectPhases: []string{"ready", "running", "running"},
		},
		{
			Name:         "three_phases_with_multiple_running",
			Phases:       []string{"prepare", "ready", "running"},
			ExpectPhases: []string{"ready", "running", "running", "running"},
		},
		{
			Name:         "four_phases_with_skip_prepare",
			CurrentPhase: "ready",
			Phases:       []string{"prepare", "ready", "running", "completed"},
			ExpectPhases: []string{"running", "completed", "completed", "completed"},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			phase := tc.Phases[0]
			if tc.CurrentPhase != "" {
				phase = tc.CurrentPhase
			}

			for _, expectPhase := range tc.ExpectPhases {
				nextPhase, err := GetNextPhaseSequential(tc.Phases, phase)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if expectPhase != nextPhase {
					t.Errorf("expected %s, got %s", expectPhase, nextPhase)
				}

				phase = nextPhase
			}
		})
	}
}
