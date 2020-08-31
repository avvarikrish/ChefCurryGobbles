package algo

import "log"

// Ksmallest returns the k smallest element in a slice
func Ksmallest(values []int, k int, lo int, hi int) int {
	if k > len(values) {
		log.Fatalf("K cannot be greater than len of values")
	}
	i := partition(values, lo, hi)
	if i == k {
		return values[k-1]
	} else if i > k {
		return Ksmallest(values, k, lo, i-2)
	} else {
		return Ksmallest(values, k, i, hi)
	}
}

func partition(values []int, lo int, hi int) int {
	pivot := values[hi]
	i := lo
	for j := lo; j < hi; j++ {
		if values[j] < pivot {
			if i != j {
				values[i], values[j] = values[j], values[i]
			}
			i++
		}
	}
	values[i], values[hi] = values[hi], values[i]
	return i + 1
}
