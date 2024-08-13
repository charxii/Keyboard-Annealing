package keyboard

import (
	"testing"
)

func CompareSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestLayoutStringConversion(t *testing.T) {
	row := "qwertyuiopasdfghjkl;'zxcvbnm,./"
	col := "qazwsxedcrfvtgbyhnujmik,ol.p;/'"

	colOut := RowLayoutToCol(row)
	rowOut := ColLayoutToRow(col)

	if row != rowOut {
		t.Errorf("ColLayoutToRow() = %s but want %s", row, rowOut)
	}

	if col != colOut {
		t.Errorf("RowLayoutToCol() = %s but want %s", col, colOut)
	}
}

func TestLayoutGroupConversion(t *testing.T) {
	row := "qwertyuiopasdfghjkl;'zxcvbnm,./"
	col := "qazwsxedcrfvtgbyhnujmik,ol.p;/'"

	expected := []string{
		"qaz",
		"wsx",
		"edc",
		"rfvtgb",
		"yhnujm",
		"ik,",
		"ol.",
		"p;/'",
	}

	groups := RowLayoutToGroups(row)

	if !CompareSlices(groups, expected) {
		t.Errorf("RowLayoutToGroups() = %s but want %s", groups, expected)
	}

	groups = ColLayoutToGroups(col)

	if !CompareSlices(groups, expected) {
		t.Errorf("ColLayoutToGroups() = %s but want %s", groups, expected)
	}

}

func TestLeftRightKeyboard(t *testing.T) {
	kb := NewKeyboard("qwertyuiopasdfghjkl;'zxcvbnm,./")

	left := "qazwsxedcrfvtgb"
	right := "yhnujmik,ol.p;/'"

	for _, r := range left {
		if !kb.OnLeft(r) {
			t.Errorf("OnLeft(%s) = false but want true", string(r))
		}
	}

	for _, r := range right {
		if !kb.OnRight(r) {
			t.Errorf("OnRight(%s) = false but want true", string(r))
		}
	}
}

func TestHomerowOptimization(t *testing.T) {
	kb := NewKeyboard("qwertyuiopasdfghjkl;'zxcvbnm,./")
	expected := "qxcbvjmk.;asetrhniop/zwdgfyu,l'"

	cf, err := CharFreqFromFolder("../CharFreqData/mt-quotes")
	if err != nil {
		t.Error(err)
	}
	kb = OptimizeHomerow(kb, cf, false, false)

	if kb.Layout != expected {
		expected_kb := NewKeyboard(expected)
		t.Errorf("OptimizeHomerow() =\n%s \n -- should be -- \n%s", kb.GetKeyboardString(), expected_kb.GetKeyboardString())
	}
}

func TestGroupCount(t *testing.T) {
	kb := NewKeyboard("qwertyuiopasdfghjkl;'zxcvbnm,./")
	group_cnt := []int{3, 3, 3, 6, 6, 3, 3, 4}

	for i, group := range kb.Groups {
		if len(group) != group_cnt[i] {
			t.Errorf("Group %d: Expected %d, got %d", i, group_cnt[i], len(group))
		}
	}
}

func TestCharFreq(t *testing.T) {
	cf, err := CharFreqFromFolder("../CharFreqData/mt-quotes")

	if err != nil {
		t.Error(err)
	}

	expected := map[rune]int{
		'e':  145118,
		't':  107577,
		'o':  96634,
		'a':  87401,
		'n':  81596,
		'i':  72571,
		's':  68735,
		'h':  64001,
		'r':  63181,
		'l':  48900,
		'd':  43002,
		'u':  37947,
		'y':  30286,
		'm':  30021,
		'w':  27418,
		'c':  25115,
		'g':  24867,
		'f':  24042,
		'.':  19755,
		'p':  17712,
		'b':  16676,
		',':  16331,
		'v':  13018,
		'k':  11520,
		'\'': 9027,
		'j':  1579,
		'x':  1354,
		'z':  757,
		'q':  689,
		';':  625,
		' ':  0,
	}

	for c, val := range expected {
		if val != cf.Chars[c] {
			t.Errorf("CharFreq() = %d but want %d", cf.Chars[c], val)
		}
	}
}

func TestGroupsToCol(t *testing.T) {
	col := "qazwsxedcrfvtgbyhnujmik,ol.p;/'"
	groups := []string{
		"qaz",
		"wsx",
		"edc",
		"rfvtgb",
		"yhnujm",
		"ik,",
		"ol.",
		"p;/'",
	}
	colOut := GroupsToCol(groups)

	if col != colOut {
		t.Errorf("GroupToCol() = %s but want %s", GroupsToCol(groups), col)
	}
}
