package ptcollection

func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr))
}
func partition(arr []int, start, end int) int {
	temp := arr[start]
	for start < end {
		for start < end && arr[end] < temp {
			end--
		}
		arr[start] = arr[end]
		for start < end && arr[start] <= temp {
			start++
		}
		arr[end] = arr[start]
	}
	arr[start] = temp
	return start
}
func quickSort(arr []int, start, end int) {
	if start < end {
		base := partition(arr, start, end)
		quickSort(arr, start, base-1)
		quickSort(arr, base+1, end)
	}
}
