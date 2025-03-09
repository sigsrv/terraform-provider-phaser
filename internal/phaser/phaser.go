// Copyright (c) EcmaXp.
// SPDX-License-Identifier: MPL-2.0

package phaser

import (
	"fmt"
	"slices"
)

func GetNextPhaseSequential(phases []string, currentPhase string) (string, error) {
	return getNextPhase(phases, currentPhase, false)
}

func getNextPhase(phases []string, currentPhase string, loop bool) (string, error) {
	index := slices.Index(phases, currentPhase)
	if index == -1 {
		return "", fmt.Errorf("phase %q not found in phases=%v", currentPhase, phases)
	}

	index++

	if loop {
		index %= len(phases)
	} else if index >= len(phases) {
		return currentPhase, nil
	}

	return phases[index], nil
}
