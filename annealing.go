package main

import (
	"fmt"
	kb "kbannealing/keyboard"
	m "kbannealing/metrics"
	"math"
	"math/rand"
	"strings"
	"unicode"
)

func SimulatedAnnealing(metricFn m.Metric, initTemp float64, cf *kb.CharFreq, lowerIsBetter bool, lockSymbols bool) *kb.Keyboard {
	temp := initTemp
	coolRate := 0.9995

	var isBetter func(int, int) bool
	if lowerIsBetter {
		isBetter = func(a, b int) bool { return a < b }
	} else {
		isBetter = func(a, b int) bool { return a > b }
	}

	// bestKb := kb.NewKeyboard("hdlvzcpui.ntrsxwfeoa;mkqj/gby',")
	bestKb := kb.NewKeyboard("qwertyuiopasdfghjkl;'zxcvbnm,./")
	bestScore := metricFn(bestKb, cf)
	lockedIndexes := []int{}

	if lockSymbols {
		for i, r := range bestKb.Layout {
			if !unicode.IsLetter(r) {
				lockedIndexes = append(lockedIndexes, i)
			}
		}
	}

	progress := ""
	epochs := int(math.Ceil(math.Log(1/initTemp) / math.Log(coolRate)))

	for temp > 1.0 {
		swaps := int(math.Max(rand.Float64()*3+1, temp/500))
		curKb := kb.MutateKeyboard(bestKb, swaps, lockedIndexes)
		score := metricFn(curKb, cf)

		if isBetter(score, bestScore) {
			bestScore = score
			bestKb = curKb
		} else {
			var delta int
			if lowerIsBetter {
				delta = score - bestScore
			} else {
				delta = bestScore - score
			}

			if rand.Float64() < math.Exp(float64(-delta)/temp) {
				bestScore = score
				bestKb = curKb
			}
		}
		epochs -= 1
		temp *= coolRate
		fmt.Print("\r" + strings.Repeat(" ", len(progress)))
		progress = fmt.Sprintf("\rScore: %d, Temp: %.2f, Epochs Left: %d", bestScore, temp, epochs)
		fmt.Print(progress)

	}
	fmt.Println("")

	return bestKb
}
