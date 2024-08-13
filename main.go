package main

import (
	"cmp"
	"flag"
	"fmt"
	kbd "kbanalysis/keyboard"
	m "kbanalysis/metrics"
	"slices"
	"strings"
)

func sorted_keys[K cmp.Ordered, T any](dict map[K]T) []K {
	keys := make([]K, 0, len(dict))

	for k := range dict {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	return keys
}

func ProcessStats(keyboards map[string]*kbd.Keyboard, cf *kbd.CharFreq, order []string) {
	bigramSum := 0
	trigramSum := 0

	for _, val := range cf.Bigrams {
		bigramSum += val
	}

	for _, val := range cf.Trigrams {
		trigramSum += val
	}

	// header
	fmt.Printf("%-23s %-10s %-10s %-10s %-10s %-10s\n", "Keyboard", "Alternate", "Roll", "SFB", "3Roll", "A+R+-S")
	fmt.Println(strings.Repeat("-", 65))

	statMap := StatMap{}

	for _, name := range order {
		kb := keyboards[name]
		stats := m.AllMetrics(kb, cf)
		statMap[name] = make(map[string]float64)

		for key, val := range stats {
			if key == "sfb" {
				statMap[name][key] = float64(val) / float64(bigramSum) * 100
			} else {
				statMap[name][key] = float64(val) / float64(trigramSum) * 100
			}
		}
	}

	saveStatsToJSON("stats.json", statMap)

	for _, name := range order {
		metrics := statMap[name]

		fmt.Printf("%-23s ", name)

		alternate := metrics["alternate"]
		roll := metrics["roll"]
		sfb := metrics["sfb"]
		threeRoll := metrics["3roll"]

		fmt.Printf("%-10.2f %-10.2f %-10.2f %-10.2f %-10.2f\n", alternate, roll, sfb, threeRoll, alternate+roll-sfb)
	}
}

func main() {
	annealFlag := flag.Bool("anneal", false, "Use simulated annealing. This will overwrite 000 optimized layouts.")
	textFlag := flag.String("text", "", "Use a wordlist txt file for data")
	folderFlag := flag.String("folder", "CharFreqData/mt-quotes", "Use a folder for data, containing monograms, bigrams, and trigrams.txt")
	lockSymbolsFlag := flag.Bool("symbollock", false, "Lock symbols in the keyboard to qwerty's layout")

	flag.Parse()

	if *textFlag != "" && *folderFlag != "CharFreqData/mt-quotes" {
		fmt.Print("Cannot use both -text and -folder for frequency data.")
		return
	}

	var cf *kbd.CharFreq
	var err error
	if *textFlag != "" {
		cf, err = kbd.NewCharFreq(*textFlag)
	} else {
		cf, err = kbd.CharFreqFromFolder(*folderFlag)
	}

	if err != nil {
		fmt.Printf("Could not load frequency data due to error: %s\n", err)
		return
	}

	combinedMetric := func(kb *kbd.Keyboard, cf *kbd.CharFreq) int {
		return -m.SfbScore(kb, cf) + m.AlternateScore(kb, cf) + m.RollScore(kb, cf)
	}

	layouts, err := loadLayoutFromJSON("layouts.json")
	if err != nil {
		fmt.Print("Could not load layout JSON")
		return
	}

	annealedKeyboards := map[string]*kbd.Keyboard{}
	if *annealFlag {
		fmt.Println("Finding optimal keyboards...")
		startTemp := 1000000.0

		fmt.Println("Optimizing for minimum sfb...")
		annealedKeyboards["000 optimized sfb"] = kbd.OptimizeHomerow(
			SimulatedAnnealing(m.SfbScore, startTemp, cf, true, *lockSymbolsFlag), cf, *lockSymbolsFlag, false)

		fmt.Println("Optimizing for alternate hand use...")
		annealedKeyboards["000 optimized alternate"] = kbd.OptimizeHomerow(
			SimulatedAnnealing(m.AlternateScore, startTemp, cf, false, *lockSymbolsFlag), cf, *lockSymbolsFlag, false)

		fmt.Println("Optimizing for maximum roll...")
		annealedKeyboards["000 optimized roll"] = kbd.OptimizeHomerow(
			SimulatedAnnealing(m.RollScore, startTemp, cf, false, *lockSymbolsFlag), cf, *lockSymbolsFlag, false)

		fmt.Println("Optimizing for 3roll...")
		annealedKeyboards["000 optimized 3roll"] = kbd.OptimizeHomerow(
			SimulatedAnnealing(m.ThreeRollScore, startTemp, cf, false, *lockSymbolsFlag), cf, *lockSymbolsFlag, true)

		fmt.Println("Optimizing for combined metrics... (maximizing altrernate + roll - sfb)")
		annealedKeyboards["000 optimized combined"] = kbd.OptimizeHomerow(
			SimulatedAnnealing(combinedMetric, startTemp, cf, false, *lockSymbolsFlag), cf, *lockSymbolsFlag, false)
	}

	keyboards := map[string]*kbd.Keyboard{}
	for name, layout := range layouts {
		keyboards[name] = kbd.NewKeyboard(layout)
	}
	for name, kb := range annealedKeyboards {
		keyboards[name] = kb
	}

	order := sorted_keys(keyboards)
	ProcessStats(keyboards, cf, order)

	fmt.Println(strings.Repeat("-", 65))

	for name, kb := range keyboards {
		if len(name) >= 3 && name[:3] == "000" {
			fmt.Println(name)
			kb.PrintKeyboard()
			fmt.Println(strings.Repeat("-", 65))
		}
	}

	if *annealFlag {
		for name, kb := range annealedKeyboards {
			layouts[name] = kb.Layout
		}
		saveLayoutToJSON("layouts.json", layouts)
	}
}
