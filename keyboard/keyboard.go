package keyboard

import (
	"fmt"
	utils "kbannealing/utils"
	"math/rand"
	"slices"
	"strings"
	"unicode"
)

type Keyboard struct {
	Layout    string
	Left      string
	Right     string
	GroupId   map[rune]int
	Groups    []string
	colLayout string
}

// layout is a string of 31 characters representing a keyboard row by row
// e.g. The standard qwerty layout is: "qwertyuiopasdfghjkl;'zxcvbnm,./"
func NewKeyboard(layout string) *Keyboard {
	if len(layout) != 31 {
		panicMsg := fmt.Sprint("Invalid layout length ", layout, " length is ", len(layout))
		panic(panicMsg)
	}

	colLayout := RowLayoutToCol(layout)
	groups := ColLayoutToGroups(colLayout)
	groupId := map[rune]int{}

	for i, group := range groups {
		for _, c := range group {
			groupId[c] = i
		}
	}

	left := colLayout[0:15]
	right := colLayout[15:31]

	return &Keyboard{layout, left, right, groupId, groups, colLayout}
}

func (k *Keyboard) GetGroup(r rune) int {
	return k.GroupId[r]
}

func (k *Keyboard) OnLeft(r rune) bool {
	idx := strings.Index(k.colLayout, string(r))
	if idx == -1 {
		return false
	}
	return idx < 15
}

func (k *Keyboard) OnRight(r rune) bool {
	idx := strings.Index(k.colLayout, string(r))
	if idx == -1 {
		return false
	}
	return idx >= 15
}

// Returns a formatted string in the shape of a keyboard. To get a one-line string, that is stored in k.Layout
func (k *Keyboard) GetKeyboardString() string {
	row1 := k.Layout[0:10]
	row2 := k.Layout[10:21]
	row3 := k.Layout[21:31]

	tmp := strings.Split(row1, "")
	row1 = strings.Join(tmp, "  ")

	tmp = strings.Split(row2, "")
	row2 = strings.Join(tmp, "  ")

	tmp = strings.Split(row3, "")
	row3 = strings.Join(tmp, "  ")

	return row1 + "\n" + row2 + "\n" + row3
}

func (k *Keyboard) PrintKeyboard() {
	fmt.Println(k.GetKeyboardString())
}

// Creates a derivative keyboard by swappinng random characters in the layout.
// Locked indexes will not be swapped. Index refers to the position of a character in k.Layout
func MutateKeyboard(k *Keyboard, swaps int, lockedIndexes []int) *Keyboard {
	chars := []rune(k.Layout)
	unlocked := []int{}

	for i := 0; i < len(chars); i++ {
		if !utils.Contains(lockedIndexes, i) {
			unlocked = append(unlocked, i)
		}
	}

	for i := 0; i < swaps; i++ {
		a := rand.Intn(len(unlocked))
		b := rand.Intn(len(unlocked))
		chars[unlocked[a]], chars[unlocked[b]] = chars[unlocked[b]], chars[unlocked[a]]
	}
	return NewKeyboard(string(chars))
}

func OptimizeHomerow(k *Keyboard, cf *CharFreq, lockSymbols bool, lockColumns bool) *Keyboard {
	priorities := [][]int{
		{3, 1, 2},
		{3, 1, 2},
		{3, 1, 2},
		{5, 1, 3, 6, 2, 4},
		{6, 2, 4, 5, 1, 3},
		{3, 1, 2},
		{3, 1, 2},
		{3, 1, 2, 4},
	}

	newGroups := make([]string, len(k.Groups))

	for i, priority := range priorities {
		chars := []rune(k.Groups[i])
		unlocked := make([]rune, 0, len(chars))
		unlockedIdx := make([]int, 0, len(chars))
		unlockedPriority := make([]int, 0, len(chars))

		for j, c := range chars {
			if lockSymbols && !unicode.IsLetter(c) {
				continue
			}
			unlocked = append(unlocked, c)
			unlockedIdx = append(unlockedIdx, j)
			unlockedPriority = append(unlockedPriority, priority[j])
		}

		// highest priority index (i.e index where priority is 1 should be first)
		priorityMap := make(map[int]int)
		for j, idx := range unlockedIdx {
			priorityMap[idx] = unlockedPriority[j]
		}

		slices.SortFunc(unlockedIdx, func(a, b int) int {
			return priorityMap[a] - priorityMap[b]
		})

		// sort by frequency, so most frequent characters are matched with highest priority
		slices.SortFunc(unlocked, func(a, b rune) int {
			return cf.Chars[b] - cf.Chars[a]
		})

		for j, c := range unlocked {
			chars[unlockedIdx[j]] = c
		}

		newGroups[i] = string(chars)
	}

	// technically you can swap some of the rows without changing the score:
	// groups [0, 3) are swappable. groups [5, 7) are swappable.
	// as long as we aren't optimizing for 3 rolls, we can do this.
	if !lockColumns {
		slices.SortFunc(newGroups[0:3], func(a, b string) int {
			return cf.Chars[rune(a[1])] - cf.Chars[rune(b[1])]
		})

		slices.SortFunc(newGroups[5:7], func(a, b string) int {
			return cf.Chars[rune(b[1])] - cf.Chars[rune(a[1])]
		})
	}

	return NewKeyboard(GroupsToRow(newGroups))
}

func ColLayoutToRow(colLayout string) string {
	// except for the last character
	columns := make([]string, 0, 10)
	for i := 0; i < 30; i += 3 {
		columns = append(columns, colLayout[i:i+3])
	}

	row1 := make([]byte, 0, 10)
	row2 := make([]byte, 0, 11)
	row3 := make([]byte, 0, 10)

	for i := 0; i < 10; i++ {
		row1 = append(row1, columns[i][0])
		row2 = append(row2, columns[i][1])
		row3 = append(row3, columns[i][2])
	}
	row2 = append(row2, colLayout[len(colLayout)-1])

	return string(row1) + string(row2) + string(row3)
}

func RowLayoutToCol(rowLayout string) string {
	row1 := rowLayout[0:10]
	row2 := rowLayout[10:21]
	row3 := rowLayout[21:31]

	columns := make([]byte, 31)

	rowPtr := 0
	for i := 0; i < 28; i += 3 {
		columns[i] = row1[rowPtr]
		columns[i+1] = row2[rowPtr]
		columns[i+2] = row3[rowPtr]
		rowPtr++
	}
	columns[30] = row2[len(row2)-1]

	return string(columns)
}

func ColLayoutToGroups(colLayout string) []string {
	return []string{
		colLayout[0:3],
		colLayout[3:6],
		colLayout[6:9],
		colLayout[9:15],
		colLayout[15:21],
		colLayout[21:24],
		colLayout[24:27],
		colLayout[27:31],
	}
}

func RowLayoutToGroups(rowLayout string) []string {
	return ColLayoutToGroups(RowLayoutToCol(rowLayout))
}

func GroupsToCol(groups []string) string {
	return strings.Join(groups, "")
}

func GroupsToRow(groups []string) string {
	return ColLayoutToRow(GroupsToCol(groups))
}
