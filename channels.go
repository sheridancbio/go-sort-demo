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
	"fmt"
	"strconv"
)

// ComparisonEvent represents an occurrence of comparing two elements
type ComparisonEvent struct {
	index                [2]int32 // the indexes of compared elements
	value                [2]int32 // the values of compared elements
	firstWasLower        bool     // the result of the comparison (true if first element is "less than" second)
	knownToBeSortedCount int32    // the count of elements currently known to be sorted
}

// SwapEvent represents an occurrence of swapping two elements
type SwapEvent struct {
	index                [2]int32 // the indexes of compared elements
	value                [2]int32 // the values of compared elements
	knownToBeSortedCount int32    // the count of elements currently known to be sorted
}

var sortingCompleteComparisonEvent = ComparisonEvent{
	knownToBeSortedCount: SORTING_COMPLETE_VALUE,
}

var sortingCompleteSwapEvent = SwapEvent{
	knownToBeSortedCount: SORTING_COMPLETE_VALUE,
}

func startSupervisionOfSort(cm *chan int, sm *chan int, algorithm int) {
	*cm <- algorithm
	*sm <- algorithm
}

func monitorSupervisorChannel(m *chan int, msc *chan string, completeMessage string) {
	var alg int
	var runningAlgorithms = map[int]bool{}
	// always begin with starting the first algorithm
	alg = <-*m
	if alg > 0 {
		runningAlgorithms[alg] = true
	}
	for len(runningAlgorithms) > 0 {
		alg = <-*m
		fmt.Println("supervisor got " + strconv.Itoa(alg) + " " + completeMessage)
		if alg > 0 {
			runningAlgorithms[alg] = true
		} else if alg < 0 {
			delete(runningAlgorithms, -alg)
		}
	}
	// no more algorithms are running
	*msc <- completeMessage
}

func processComparisonChannel(c *chan ComparisonEvent, algorithm int, m *chan int) {
	var ce ComparisonEvent
	for true {
		ce = <-*c
		if ce == sortingCompleteComparisonEvent {
			*m <- -algorithm // signal that this channel processing is done
			return
		}
	}
}

func processSwapChannel(c *chan SwapEvent, algorithm int, m *chan int) {
	var se SwapEvent
	for true {
		se = <-*c
		if se == sortingCompleteSwapEvent {
			*m <- -algorithm // signal that this channel processing is done
			return
		}
	}
}

func waitForEverythingComplete(msc *chan string) {
	var comparisonProcessingComplete = false
	var swapProcessingComplete = false
	var m string
	for !comparisonProcessingComplete || !swapProcessingComplete {
		m = <-*msc
		if m == ALL_COMPARISONS_COMPLETE_MESSAGE {
			comparisonProcessingComplete = true
		}
		if m == ALL_SWAPS_COMPLETE_MESSAGE {
			swapProcessingComplete = true
		}
	}
}
