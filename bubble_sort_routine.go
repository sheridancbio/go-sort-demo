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

// BubbleSortRoutine - sorts by moving small items to the head of the list iteratively
type BubbleSortRoutine struct {
	data                 []int32
	dataSize             int32
	swapChannel          chan SwapEvent
	comparisonChannel    chan ComparisonEvent
	knownToBeSortedCount int32
}

// NewBubbleSortRoutine factory
func NewBubbleSortRoutine(startSlice []int32) *BubbleSortRoutine {
	bsr := new(BubbleSortRoutine)
	bsr.dataSize = int32(len(startSlice))
	bsr.data = make([]int32, bsr.dataSize)
	_ = copy(bsr.data, startSlice)
	cc := make(chan ComparisonEvent, 1000)
	bsr.comparisonChannel = cc
	sc := make(chan SwapEvent, 1000)
	bsr.swapChannel = sc
	bsr.knownToBeSortedCount = 0
	return bsr
}

func (bsr BubbleSortRoutine) getComparisonChannel() chan ComparisonEvent {
	return bsr.comparisonChannel
}

func (bsr BubbleSortRoutine) getSwapChannel() chan SwapEvent {
	return bsr.swapChannel
}

func (bsr BubbleSortRoutine) run() {
	var top int32 = int32(0)
	var bottom int32 = int32(len(bsr.data) - 1)
	for top = int32(0); top < bottom; top = top + 1 {
		var pos int32
		for pos = bottom - 1; pos >= top; pos = pos - 1 {
			if !compareElementsAt(bsr.data, pos, pos+1, bsr.knownToBeSortedCount, bsr.comparisonChannel) {
				swapElementsAt(bsr.data, pos, pos+1, bsr.knownToBeSortedCount, bsr.swapChannel)
			}
		}
		bsr.knownToBeSortedCount = top
	}
	sortingRoutineComplete(bsr.comparisonChannel, bsr.swapChannel)
}
