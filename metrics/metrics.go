package metrics

import (
	kbd "kbannealing/keyboard"
)

type Metric func(*kbd.Keyboard, *kbd.CharFreq) int

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Generates every possible combination, taking one character from each string
// ex. "ab", "cd" -> "ac", "ad", "bc", "bd"
func stringProduct(strings ...string) chan string {
	if len(strings) == 0 {
		out := make(chan string)
		close(out)
		return out
	}

	out := make(chan string)

	go func() {
		defer close(out)

		var generate func(int, []rune)
		generate = func(index int, current []rune) {
			if index == len(strings) {
				out <- string(current)
				return
			}
			for _, c := range strings[index] {
				newCurrent := make([]rune, len(current), len(current)+1)
				copy(newCurrent, current)
				newCurrent = append(newCurrent, c)
				generate(index+1, newCurrent)
			}
		}

		generate(0, []rune{})
	}()

	return out
}

func AlternateScore(kb *kbd.Keyboard, cf *kbd.CharFreq) int {
	score := 0

	for seq := range stringProduct(kb.Left, kb.Right, kb.Left) {
		val, ok := cf.Trigrams[seq]
		if ok {
			score += val
		}
	}

	for seq := range stringProduct(kb.Right, kb.Left, kb.Right) {
		val, ok := cf.Trigrams[seq]
		if ok {
			score += val
		}
	}

	return score
}

// Single Finger Bigrams
func SfbScore(kb *kbd.Keyboard, cf *kbd.CharFreq) int {
	score := 0

	for _, group := range kb.Groups {
		for seq := range stringProduct(group, group) {
			if seq[0] == seq[1] {
				continue
			}
			val, ok := cf.Bigrams[seq]
			if ok {
				score += val
			}
		}
	}

	return score
}

// 2 Roll + Alternate Hand
func RollScore(kb *kbd.Keyboard, cf *kbd.CharFreq) int {
	score := 0
	for seq := range stringProduct(kb.Left, kb.Left, kb.Right) {
		val, ok := cf.Trigrams[seq]
		if ok {
			score += val
		}
		val, ok = cf.Trigrams[reverse(seq)]
		if ok {
			score += val
		}
	}

	for seq := range stringProduct(kb.Right, kb.Right, kb.Left) {
		val, ok := cf.Trigrams[seq]
		if ok {
			score += val
		}
		val, ok = cf.Trigrams[reverse(seq)]
		if ok {
			score += val
		}
	}
	return score
}

// Same side of keyboard, in the same direction. Ex: "lkj" or "jkl"
func ThreeRollScore(kb *kbd.Keyboard, cf *kbd.CharFreq) int {
	score := 0

	// Left groups, indexes 0 - 3
	for i := 0; i < 2; i++ {
		for seq := range stringProduct(kb.Groups[i], kb.Groups[i+1], kb.Groups[i+2]) {
			val, ok := cf.Trigrams[seq]
			if ok {
				score += val
			}
		}
	}

	// Right groups, indexes 4 - 7 in reverse
	for i := 7; i > 5; i-- {
		for seq := range stringProduct(kb.Groups[i], kb.Groups[i-1], kb.Groups[i-2]) {
			val, ok := cf.Trigrams[seq]
			if ok {
				score += val
			}
		}
	}

	return score
}

func AllMetrics(kb *kbd.Keyboard, cf *kbd.CharFreq) map[string]int {
	return map[string]int{
		"alternate": AlternateScore(kb, cf),
		"sfb":       SfbScore(kb, cf),
		"roll":      RollScore(kb, cf),
		"3roll":     ThreeRollScore(kb, cf),
	}
}
