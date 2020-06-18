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

const ALGORITHM_BUBBLE_SORT int = 1
const ALGORITHM_SELECTION_SORT int = 2
const ALGORITHM_INSERTION_SORT int = 3
const ALGORITHM_SHELL_SORT int = 4
const ALGORITHM_TREE_SORT int = 5
const ALGORITHM_HEAP_SORT int = 6
const ALGORITHM_MERGE_SORT int = 7
const ALGORITHM_QUICK_SORT int = 8
const ALGORITHM_RANDOM_SORT int = 9

const SORTING_COMPLETE_VALUE int32 = -1

const ALL_COMPARISONS_COMPLETE_MESSAGE string = "all comparisons complete"
const ALL_SWAPS_COMPLETE_MESSAGE string = "all swaps complete"
