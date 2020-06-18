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

func sortingRoutineComplete(cc *chan ComparisonEvent, sc *chan SwapEvent) {
	*cc <- sortingCompleteComparisonEvent
	*sc <- sortingCompleteSwapEvent
}

func compareElementsAt(data []int32, i int32, j int32, ktbsc int32, c *chan ComparisonEvent) bool {
	var e ComparisonEvent = ComparisonEvent{[2]int32{i, j}, [2]int32{data[i], data[j]}, data[i] < data[j], ktbsc}
	*c <- e
	return e.firstWasLower
}

func swapElementsAt(data []int32, i int32, j int32, ktbsc int32, c *chan SwapEvent) {
	var e SwapEvent = SwapEvent{[2]int32{i, j}, [2]int32{data[i], data[j]}, ktbsc}
	*c <- e
	var t int32 = data[i]
	data[i] = data[j]
	data[j] = t
}
