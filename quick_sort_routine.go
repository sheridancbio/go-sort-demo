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

// QuickSortRoutine - sorts by picking a pivot element and partitioning each sublist into a larger and a smaller partition. Recur.
type QuickSortRoutine struct {
	data                 []int32
	dataSize             int32
	swapChannel          chan SwapEvent
	comparisonChannel    chan ComparisonEvent
	knownToBeSortedCount int32
}

// NewQuickSortRoutine factory
func NewQuickSortRoutine(startSlice []int32) *QuickSortRoutine {
	qsr := new(QuickSortRoutine)
	qsr.dataSize = int32(len(startSlice))
	qsr.data = make([]int32, qsr.dataSize)
	_ = copy(qsr.data, startSlice)
	cc := make(chan ComparisonEvent, 1000)
	qsr.comparisonChannel = cc
	sc := make(chan SwapEvent, 1000)
	qsr.swapChannel = sc
	qsr.knownToBeSortedCount = 0
	return qsr
}

func (qsr QuickSortRoutine) getComparisonChannel() chan ComparisonEvent {
	return qsr.comparisonChannel
}

func (qsr QuickSortRoutine) getSwapChannel() chan SwapEvent {
	return qsr.swapChannel
}

type sortRange struct {
	top    int32
	bottom int32
}

func (qsr QuickSortRoutine) selectPivot(top int32) int32 {
	if compareElementsAt(qsr.data, top, top+1, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
		// e0 < e1
		if compareElementsAt(qsr.data, top+1, top+2, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
			// e0 < e1 < e2
			return top + 1
		}
		// e0 < e1 && e2 < e1
		if compareElementsAt(qsr.data, top, top+2, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
			// e0 < e2 < e1
			return top + 2
		}
		// e2 < e0 < e1
		return top
	}
	// e1 < e0
	if compareElementsAt(qsr.data, top+1, top+2, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
		// e1 < e0 && e1 < e2
		if compareElementsAt(qsr.data, top, top+2, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
			// e1 < e0 < e2
			return top
		}
		// e1 < e2 < e0
		return top + 2
	}
	// e2 < e1 < e0
	return top + 1
}

func (qsr QuickSortRoutine) insertionSort(rangeToSort sortRange) {
	var bottom int32 = rangeToSort.top
	for bottom < rangeToSort.bottom {
		var scanPos int32
		for scanPos = bottom + 1; scanPos > rangeToSort.top; scanPos = scanPos - 1 {
			if compareElementsAt(qsr.data, scanPos, scanPos-1, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
				swapElementsAt(qsr.data, scanPos, scanPos-1, qsr.knownToBeSortedCount, qsr.swapChannel)
			}
		}
		bottom = bottom + 1
		qsr.knownToBeSortedCount = qsr.knownToBeSortedCount + 1
	}
}

/* Quick Sort
 * partition the list around a selected pivot element into two sublists -- one with elements
 * larger than the pivot and one with elements smaller than the pivot
 * recursively apply quicksort to each sublist
 * when a sublist is 5 elements or less, use insertion sort instead
 * select a pivot by considering the first three elements in the list and choosing the
 * middle-sized element
 */
func (qsr QuickSortRoutine) run() {
	var rangesToSort []sortRange = make([]sortRange, 0)
	rangesToSort = append(rangesToSort, sortRange{0, int32(len(qsr.data) - 1)})
	for len(rangesToSort) > 0 {
		// pop the next range to sort
		var rangeToSort = rangesToSort[0]
		rangesToSort = rangesToSort[1:]
		if rangeToSort.bottom-rangeToSort.top < 6 {
			qsr.insertionSort(rangeToSort)
		} else {
			var pivotPos int32 = qsr.selectPivot(rangeToSort.top)
			if pivotPos != rangeToSort.top {
				swapElementsAt(qsr.data, pivotPos, rangeToSort.top, qsr.knownToBeSortedCount, qsr.swapChannel)
				pivotPos = rangeToSort.top
			}
			var scanFromTop int32 = rangeToSort.top + 1
			var scanFromBottom int32 = rangeToSort.bottom
			var anySwapWasMade bool = false
			for scanFromTop < scanFromBottom {
				for scanFromTop < scanFromBottom && compareElementsAt(qsr.data, scanFromTop, pivotPos, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
					scanFromTop = scanFromTop + 1
				}
				if scanFromTop < scanFromBottom && anySwapWasMade {
					// we know the element at scanFromBottom is >= pivot element if a swap has occurred in this range - no comparison needed
					scanFromBottom = scanFromBottom - 1
				}
				for scanFromTop < scanFromBottom && compareElementsAt(qsr.data, pivotPos, scanFromBottom, qsr.knownToBeSortedCount, qsr.comparisonChannel) {
					scanFromBottom = scanFromBottom - 1
				}
				if scanFromTop < scanFromBottom {
					// both incorrectly positioned elements found, so swap them
					swapElementsAt(qsr.data, scanFromTop, scanFromBottom, qsr.knownToBeSortedCount, qsr.swapChannel)
					anySwapWasMade = true
					scanFromTop = scanFromTop + 1
				}
			}
			// partition completed - exchange pivot element with the final smaller element
			// we know from the selection of pivot approach that at least one element smaller and
			// one element larger than the pivot exists in rangeToBeSorted
			// so at the end of partitioning scanFromTop will have moved at least one step past rangeToSort.top
			swapElementsAt(qsr.data, pivotPos, scanFromTop-1, qsr.knownToBeSortedCount, qsr.swapChannel)
			qsr.knownToBeSortedCount = qsr.knownToBeSortedCount + 1                                 // pivot element is in its final position
			rangesToSort = append([]sortRange{{scanFromTop, rangeToSort.bottom}}, rangesToSort...)  // queue larger sublist
			rangesToSort = append([]sortRange{{rangeToSort.top, scanFromTop - 2}}, rangesToSort...) // queue smaller sublist
		}
	}
	sortingRoutineComplete(qsr.comparisonChannel, qsr.swapChannel)
}
