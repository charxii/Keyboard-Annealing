package metrics

import (
	kbd "kbannealing/keyboard"
	"testing"
)

// Make sure that optimizing the homerow of the keyboard does not change any metrics scores
func TestHomerowOptimized(t *testing.T) {
	kb := kbd.NewKeyboard("',hqt;yezsainkmwpufv/o.ldjcgrxb")

	cf, err := kbd.CharFreqFromFolder("../CharFreqData/mt-quotes")
	if err != nil {
		t.Error(err)
	}

	optimizedKb := kbd.OptimizeHomerow(kb, cf, false, false)

	unoptimizedMetric := AllMetrics(kb, cf)
	optimizedMetric := AllMetrics(optimizedKb, cf)

	for key, val := range unoptimizedMetric {
		if key != "3roll" && val != optimizedMetric[key] {
			t.Errorf("(3roll False) Key %s: Expected %d, got %d", key, val, optimizedMetric[key])
		}
	}

	threeRollOptimized := kbd.OptimizeHomerow(kb, cf, true, true)
	optimizedMetric = AllMetrics(threeRollOptimized, cf)

	for key, val := range unoptimizedMetric {
		if val != optimizedMetric[key] {
			t.Errorf("(3roll True) Key %s: Expected %d, got %d", key, val, optimizedMetric[key])
		}
	}
}

func TestStringProduct(t *testing.T) {
	for s := range stringProduct("ab", "cd", "ef") {
		t.Log(s)
		if len(s) != 3 {
			t.Errorf("Expected length 3, got %d", len(s))
		}
	}
}
