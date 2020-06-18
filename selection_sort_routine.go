package main

/*
    Parallel Sorting Demo
    Copyright (C) 2020 Robert Sheridan

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// SelectionSortRoutine - sorts by finding the smallest unsorted element and moving it into place
type SelectionSortRoutine struct {
	data                 []int32
	dataSize             int32
	swapChannel          chan SwapEvent
	comparisonChannel    chan ComparisonEvent
	knownToBeSortedCount int32
}

// NewSelectionSortRoutine factory
func NewSelectionSortRoutine(startSlice []int32) *SelectionSortRoutine {
	ssr := new(SelectionSortRoutine)
	ssr.dataSize = int32(len(startSlice))
	ssr.data = make([]int32, ssr.dataSize)
	_ = copy(ssr.data, startSlice)
	cc := make(chan ComparisonEvent, 1000)
	ssr.comparisonChannel = cc
	sc := make(chan SwapEvent, 1000)
	ssr.swapChannel = sc
	ssr.knownToBeSortedCount = 0
	return ssr
}

func (ssr SelectionSortRoutine) getComparisonChannel() chan ComparisonEvent {
	return ssr.comparisonChannel
}

func (ssr SelectionSortRoutine) getSwapChannel() chan SwapEvent {
	return ssr.swapChannel
}

func (ssr SelectionSortRoutine) run() {
	var top int32 = int32(0)
	var bottom int32 = int32(len(ssr.data) - 1)
	for top = int32(0); top < bottom; top = top + 1 {
		var indexOfLowest = top
		var scanPos int32
		for scanPos = bottom; scanPos > top; scanPos = scanPos - 1 {
			if compareElementsAt(ssr.data, scanPos, indexOfLowest, ssr.knownToBeSortedCount, ssr.comparisonChannel) {
				indexOfLowest = scanPos
			}
		}
		swapElementsAt(ssr.data, top, indexOfLowest, ssr.knownToBeSortedCount, ssr.swapChannel)
		ssr.knownToBeSortedCount = top
	}
	sortingRoutineComplete(ssr.comparisonChannel, ssr.swapChannel)
}
