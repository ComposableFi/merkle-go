package mmr

import (
	"reflect"
	"testing"
)

func TestPosHeightInTree(t *testing.T) {
	tests := map[string]struct {
		input uint64
		want  uint32
	}{
		"Pos 0": {input: 0, want: 0},
		"Pos 1": {input: 1, want: 0},
		"Pos 2": {input: 2, want: 1},
		"Pos 3": {input: 3, want: 0},
		"Pos 4": {input: 4, want: 0},
		"Pos 6": {input: 6, want: 2},
		"Pos 7": {input: 7, want: 0},
	}

	for name, test := range tests {
		got := PosHeightInTree(test.input)
		if got != test.want {
			t.Errorf("%s: want %v  got %v", name, test.want, got)
		}
	}
}

func TestGetPeaks(t *testing.T) {
	tests := map[string]struct {
		input uint64
		want  []uint64
	}{
		"Pos 0":  {input: 0, want: []uint64{0}},
		"Pos 1":  {input: 1, want: []uint64{0}},
		"Pos 2":  {input: 2, want: []uint64{0}},
		"Pos 3":  {input: 3, want: []uint64{2}},
		"Pos 4":  {input: 4, want: []uint64{2, 3}},
		"Pos 5":  {input: 5, want: []uint64{2, 3}},
		"Pos 6":  {input: 6, want: []uint64{2, 5}},
		"Pos 7":  {input: 7, want: []uint64{6}},
		"Pos 19": {input: 19, want: []uint64{14, 17, 18}},
	}

	for name, test := range tests {
		got := GetPeaks(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: want %v  got %v", name, test.want, got)
		}
	}
}

func TestLeafIndexToPos(t *testing.T) {
	tests := map[string]struct {
		input uint64
		want  uint64
	}{
		"index 0": {input: 0, want: 0},
		"index 1": {input: 1, want: 1},
		"index 2": {input: 2, want: 3},
	}

	for name, test := range tests {
		got := LeafIndexToPos(test.input)
		if got != test.want {
			t.Errorf("%s: want %v  got %v", name, test.want, got)
		}
	}
}

func TestLeafIndexToMMRSize(t *testing.T) {
	tests := map[string]struct {
		input uint64
		want  uint64
	}{
		"index 0": {input: 0, want: 1},
		"index 1": {input: 1, want: 3},
		"index 2": {input: 2, want: 4},
	}

	for name, test := range tests {
		got := LeafIndexToMMRSize(test.input)
		if got != test.want {
			t.Errorf("%s: want %v  got %v", name, test.want, got)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := map[string]struct {
		input []int
		want  []int
	}{
		"3 items": {input: []int{1, 2, 3}, want: []int{3, 2, 1}},
		"1 item":  {input: []int{0}, want: []int{0}},
		"no item": {input: []int{}, want: []int{}},
	}

	for name, test := range tests {
		reverse(test.input)
		if !reflect.DeepEqual(test.input, test.want) {
			t.Errorf("%s: want %v  got %v", name, test.want, test.input)
		}
	}
}
