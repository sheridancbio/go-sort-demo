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

// InsertionSortRoutine - sorts by adding/moving one element at a time into the correct position in a sorted list
type InsertionSortRoutine struct {
	data                 []int32
	dataSize             int32
	swapChannel          chan SwapEvent
	comparisonChannel    chan ComparisonEvent
	knownToBeSortedCount int32
}

// NewInsertionSortRoutine factory
func NewInsertionSortRoutine(startSlice []int32) *InsertionSortRoutine {
	isr := new(InsertionSortRoutine)
	isr.dataSize = int32(len(startSlice))
	isr.data = make([]int32, isr.dataSize)
	_ = copy(isr.data, startSlice)
	cc := make(chan ComparisonEvent, 1000)
	isr.comparisonChannel = cc
	sc := make(chan SwapEvent, 1000)
	isr.swapChannel = sc
	isr.knownToBeSortedCount = 0
	return isr
}

func (isr InsertionSortRoutine) getComparisonChannel() chan ComparisonEvent {
	return isr.comparisonChannel
}

func (isr InsertionSortRoutine) getSwapChannel() chan SwapEvent {
	return isr.swapChannel
}

func (isr InsertionSortRoutine) run() {
	var top int32 = int32(0)
	var bottom int32 = top
	for bottom < int32(len(isr.data)-1) {
		var scanPos int32
		for scanPos = bottom + 1; scanPos > top; scanPos = scanPos - 1 {
			if compareElementsAt(isr.data, scanPos, scanPos-1, isr.knownToBeSortedCount, isr.comparisonChannel) {
				swapElementsAt(isr.data, scanPos, scanPos-1, isr.knownToBeSortedCount, isr.swapChannel)
			}
		}
		bottom = bottom + 1
		isr.knownToBeSortedCount = bottom
	}
	sortingRoutineComplete(isr.comparisonChannel, isr.swapChannel)
}
