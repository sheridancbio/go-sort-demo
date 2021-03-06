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
	"math/rand"
	"strconv"
	"time"
)

func makeRandomizedDataArray(size int32) []int32 {
	data := make([]int32, 0, size)
	var pos int32
	for pos = 0; pos < size; pos = pos + 1 {
		data = append(data, pos)
	}
	randomSource := rand.NewSource(time.Now().UnixNano())
	for pos = 0; pos < size; pos = pos + 1 {
		var pos2 int32 = rand.New(randomSource).Int31n(size)
		temp := data[pos]
		data[pos] = data[pos2]
		data[pos2] = temp
	}
	return data
}

func printDataArray(data []int32) {
	for pos := 0; pos < len(data); pos = pos + 1 {
		fmt.Println(data[pos])
	}
}

func arrayIsSorted(data []int32) bool {
	for pos := 0; pos < len(data)-1; pos = pos + 1 {
		if data[pos] > data[pos+1] {
			fmt.Println("Incorrect order: " + strconv.Itoa(int(data[pos])) + " was positioned before " + strconv.Itoa(int(data[pos+1])))
			return false
		}
	}
	return true
}

func reportFinalSortResults(data []int32, name string) {
	if arrayIsSorted(data) {
		fmt.Println(name + " correctly sorted.")
	} else {
		fmt.Println(name + " not correctly sorted - some elements are out of order")
		printDataArray(data)
	}
}

func main() {
	var startSlice []int32 = makeRandomizedDataArray(1000)
	// create supervisory channels and start processing
	var masterSupervisorChannel chan string = make(chan string)
	var compareSupervisorChannel chan int = make(chan int)
	go monitorSupervisorChannel(compareSupervisorChannel, masterSupervisorChannel, ALL_COMPARISONS_COMPLETE_MESSAGE)
	var swapSupervisorChannel chan int = make(chan int)
	go monitorSupervisorChannel(swapSupervisorChannel, masterSupervisorChannel, ALL_SWAPS_COMPLETE_MESSAGE)
	// create algorithm routines
	var bsr *BubbleSortRoutine = NewBubbleSortRoutine(startSlice)
	var ssr *SelectionSortRoutine = NewSelectionSortRoutine(startSlice)
	var isr *InsertionSortRoutine = NewInsertionSortRoutine(startSlice)
	var shsr *ShellSortRoutine = NewShellSortRoutine(startSlice)
	var qsr *QuickSortRoutine = NewQuickSortRoutine(startSlice)
	// start up algorithms and channel processors
	startSupervisionOfSort(compareSupervisorChannel, swapSupervisorChannel, ALGORITHM_BUBBLE_SORT)
	go processComparisonChannel(bsr.getComparisonChannel(), ALGORITHM_BUBBLE_SORT, compareSupervisorChannel)
	go processSwapChannel(bsr.getSwapChannel(), ALGORITHM_BUBBLE_SORT, swapSupervisorChannel)
	startSupervisionOfSort(compareSupervisorChannel, swapSupervisorChannel, ALGORITHM_SELECTION_SORT)
	go processComparisonChannel(ssr.getComparisonChannel(), ALGORITHM_SELECTION_SORT, compareSupervisorChannel)
	go processSwapChannel(ssr.getSwapChannel(), ALGORITHM_SELECTION_SORT, swapSupervisorChannel)
	startSupervisionOfSort(compareSupervisorChannel, swapSupervisorChannel, ALGORITHM_INSERTION_SORT)
	go processComparisonChannel(isr.getComparisonChannel(), ALGORITHM_INSERTION_SORT, compareSupervisorChannel)
	go processSwapChannel(isr.getSwapChannel(), ALGORITHM_INSERTION_SORT, swapSupervisorChannel)
	startSupervisionOfSort(compareSupervisorChannel, swapSupervisorChannel, ALGORITHM_SHELL_SORT)
	go processComparisonChannel(shsr.getComparisonChannel(), ALGORITHM_SHELL_SORT, compareSupervisorChannel)
	go processSwapChannel(shsr.getSwapChannel(), ALGORITHM_SHELL_SORT, swapSupervisorChannel)
	startSupervisionOfSort(compareSupervisorChannel, swapSupervisorChannel, ALGORITHM_QUICK_SORT)
	go processComparisonChannel(qsr.getComparisonChannel(), ALGORITHM_QUICK_SORT, compareSupervisorChannel)
	go processSwapChannel(qsr.getSwapChannel(), ALGORITHM_QUICK_SORT, swapSupervisorChannel)
	// start sorting algorithms
	fmt.Println("beginning sorting routines")
	go bsr.run()
	go ssr.run()
	go isr.run()
	go shsr.run()
	go qsr.run()
	waitForEverythingComplete(masterSupervisorChannel)
	reportFinalSortResults(bsr.data, "bubble sort")
	reportFinalSortResults(ssr.data, "selection sort")
	reportFinalSortResults(isr.data, "insertion sort")
	reportFinalSortResults(ssr.data, "shell sort")
	reportFinalSortResults(qsr.data, "quick sort")
	fmt.Println("program complete")
}
