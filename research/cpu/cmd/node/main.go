package main

import (
	"github.com/kotfalya/hulk/research/cpu/ledger"
	"github.com/kotfalya/hulk/research/cpu/types"
)

type Node struct {
	addr  types.Addr
	pk    types.PK
	block ledger.Block
}

type MemoryBlock struct {
	addr types.Addr
	rp   types.Permission
	wp   types.Permission
}

type StorageBlock struct {
	addr types.Addr
	rp   types.Permission
	wp   types.Permission
}

func main() {

}

//
//func searchRange(nums []int, target int) []int {
//	r := []int{-1, -1}
//	if len(nums) == 0 || nums[0] > target || nums[len(nums)-1] < target {
//		return r
//	}
//
//	ls := leftSearch(nums, target, 0)
//	if ls == -1 {
//		return r
//	}
//
//	rs := rightSearch(nums[ls+1:], target, ls+1)
//	if rs == -1 {
//		rs = ls
//	}
//
//	return []int{ls, rs}
//
//}
//
//func leftSearch(nums []int, target, pos int) int {
//	if nums[0] == target {
//		return pos
//	} else if len(nums) == 1 {
//		return -1
//	}
//
//	mid := len(nums) / 2
//	if nums[mid-1] < target {
//		return leftSearch(nums[mid:], target, pos+mid)
//	} else {
//		return leftSearch(nums[:mid], target, pos)
//	}
//
//}
//
//func rightSearch(nums []int, target, pos int) int {
//	if len(nums) == 0 {
//		return -1
//	}
//	lastIndex := len(nums) - 1
//
//	if nums[lastIndex] == target {
//		return pos + lastIndex
//	} else if len(nums) == 1 {
//		return -1
//	}
//
//	mid := len(nums) / 2
//
//	if nums[mid] > target {
//		return rightSearch(nums[:mid], target, pos)
//	} else {
//		return rightSearch(nums[mid:], target, pos+mid)
//	}
//}
//
//
/////////////////////////////
//
//type Window struct {
//	runes     []rune
//	runePos   map[rune]int
//	sampleLen int
//}
//
//func NewWindow(sampleLen int) *Window {
//	return &Window{
//		runePos:   make(map[rune]int),
//		sampleLen: sampleLen,
//	}
//}
//
//func (w *Window) isFull() bool {
//	return w.sampleLen == len(w.runes)
//}
//
//func (w *Window) dist() int {
//	if len(w.runes) < 2 {
//		return len(w.runes)
//	}
//
//	firstRune := w.runes[0]
//	lastRune := w.runes[len(w.runes)-1]
//	return w.runePos[lastRune] - w.runePos[firstRune] + 1
//}
//
//func (w *Window) isMinimal() bool {
//	return w.sampleLen == w.dist()
//}
//
//func (w *Window) update(i int, r rune) bool {
//	// first item
//	if len(w.runePos) == 0 {
//		w.runePos[r] = i
//		w.runes = []rune{r}
//		return true
//	}
//
//	// first instance of rune
//	if _, ok := w.runePos[r]; !ok {
//		w.runePos[r] = i
//		w.runes = append(w.runes, r)
//		return true
//	}
//	//recalculate existing rune position
//
//	// runes is not full
//	if !w.isFull() {
//		if w.runes[0] == r {
//			w.runePos[r] = i
//			w.runes = append(w.runes[1:], r)
//		}
//		return true
//	}
//	// runes is full
//
//	// rune is not a first item, no point to recalculate
//	if w.runes[0] != r {
//		// stop working on this runes, create new one
//		return false
//	}
//
//	secondRune := w.runes[1]
//
//	// if distance between first and second rune greater then distance between second rune and current position
//	// we have to recalculate position
//	if w.runePos[secondRune]-w.runePos[r] > i-w.runePos[secondRune] {
//		w.runePos[r] = i
//		w.runes = append(w.runes[1:], r)
//		return true
//	}
//
//	// stop working on this runes, create new one
//	return false
//}
//
//func minWindow(s string, t string) string {
//	window := NewWindow(len(t))
//	minWindow := window
//	for i, r := range []rune(s) {
//		if strings.ContainsRune(t, r) &&
//			!window.update(i, r) {
//			if window.isMinimal() {
//				break
//			}
//
//			if minWindow.dist() > window.dist() {
//				minWindow = window
//			}
//
//			window = NewWindow(len(t))
//			window.update(i, r)
//		}
//	}
//
//	if window.isFull() && minWindow.dist() > window.dist() {
//		minWindow = window
//	}
//
//	if minWindow.isFull() {
//		firstRunePos := minWindow.runePos[minWindow.runes[0]]
//		return s[firstRunePos : firstRunePos+minWindow.dist()]
//	} else {
//		return ""
//	}
//}
//
/////////////////////////////
//
//type Line struct {
//	words  []string
//	length int
//}
//
//func (l *Line) AddWord(maxWidth int, word string) bool {
//	if l.length+len(l.words)+len(word) > maxWidth {
//		return false
//	}
//
//	l.length += len(word)
//	l.words = append(l.words, word)
//
//	return true
//}
//
//func (l *Line) Print(maxWidth int, last bool) string {
//	if len(l.words) == 1 {
//		return l.words[0] + strings.Repeat(" ", maxWidth-l.length)
//	}
//
//	if last {
//		return strings.Join(l.words, " ") + strings.Repeat(" ", maxWidth-l.length-len(l.words)+1)
//	}
//
//	extraSpaceDefault := (maxWidth - l.length) / (len(l.words) - 1)
//	extraSpaceMod := (maxWidth - l.length) % (len(l.words) - 1)
//
//	res := ""
//	for i := range l.words {
//		res += l.words[i]
//		if i == len(l.words)-1 {
//			break
//		}
//
//		if i+1 <= extraSpaceMod {
//			res += strings.Repeat(" ", extraSpaceDefault+1)
//		} else {
//			res += strings.Repeat(" ", extraSpaceDefault)
//		}
//	}
//
//	return res
//
//}
//
//func fullJustify(words []string, maxWidth int) []string {
//	lines := []*Line{new(Line)}
//	lineCount := 1
//
//	for i := range words {
//		if !lines[lineCount-1].AddWord(maxWidth, words[i]) {
//			lines = append(lines, new(Line))
//			lineCount++
//			lines[lineCount-1].AddWord(maxWidth, words[i])
//		}
//	}
//
//	var res []string
//	for i := range lines {
//		res = append(res, lines[i].Print(maxWidth, i == len(lines)-1))
//	}
//
//	return res
//}
