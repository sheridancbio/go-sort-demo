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

import (
	"math"
)

// ShellSortRoutine - sort list by performing insertion sort on elements separated by distance N, iteratively decreasing N to 1
type ShellSortRoutine struct {
	data                 []int32
	dataSize             int32
	swapChannel          chan SwapEvent
	comparisonChannel    chan ComparisonEvent
	knownToBeSortedCount int32
}

// NewShellSortRoutine factory
func NewShellSortRoutine(startSlice []int32) *ShellSortRoutine {
	ssr := new(ShellSortRoutine)
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

func (ssr ShellSortRoutine) getComparisonChannel() chan ComparisonEvent {
	return ssr.comparisonChannel
}

func (ssr ShellSortRoutine) getSwapChannel() chan SwapEvent {
	return ssr.swapChannel
}

// an insertion sort on all elements in the range separated by an interval
func (ssr ShellSortRoutine) insertionSort(rangeToSort sortRange, interval int32) {
	var bottom int32 = rangeToSort.top
	for bottom <= rangeToSort.bottom-interval {
		var scanPos int32
		for scanPos = bottom + interval; scanPos > rangeToSort.top; scanPos = scanPos - interval {
			if compareElementsAt(ssr.data, scanPos, scanPos-interval, ssr.knownToBeSortedCount, ssr.comparisonChannel) {
				swapElementsAt(ssr.data, scanPos, scanPos-interval, ssr.knownToBeSortedCount, ssr.swapChannel)
			}
		}
		bottom = bottom + interval
		if interval == 1 {
			ssr.knownToBeSortedCount = ssr.knownToBeSortedCount + 1
		}
	}
}

// compute a slice of intervals up to the data size (number of elemetns to be sorted)
func (ssr ShellSortRoutine) findShellGapSizeSeries() []int32 {
	var shellGapSizeSeries = make([]int32, 0)
	const shellGapSizeLimit = math.MaxInt32 / 3
	var lower int32 = 1
	var higher int32 = 1
	for higher <= shellGapSizeLimit && higher < ssr.dataSize {
		shellGapSizeSeries = append(shellGapSizeSeries, higher)
		// add every third fibonacci number
		var index int32 = 0
		for index = 0; index < 3; index = index + 1 {
			var total = lower + higher
			lower = higher
			higher = total
		}
	}
	return shellGapSizeSeries
}

// iterate through interval sizes in decreasing order and call the interval insertion sort on every list partition, starting at each offset in the interval
func (ssr ShellSortRoutine) run() {
	var shellGapSizeSeries []int32 = ssr.findShellGapSizeSeries()
	for intervalIndex := len(shellGapSizeSeries) - 1; intervalIndex >= 0; intervalIndex = intervalIndex - 1 {
		var interval int32 = shellGapSizeSeries[intervalIndex]
		var intervalRangeToBottom int32 = (ssr.dataSize / interval) * interval // probe for bottom by adding the greatest number of whole intervals forward
		var rangeTop int32
		for rangeTop = 0; rangeTop < interval; rangeTop = rangeTop + 1 {
			var rangeBottom int32 = rangeTop + intervalRangeToBottom
			if rangeBottom >= ssr.dataSize {
				rangeBottom = rangeBottom - interval
			}
			ssr.insertionSort(sortRange{rangeTop, rangeBottom}, interval)
		}
	}
	sortingRoutineComplete(ssr.comparisonChannel, ssr.swapChannel)
}
