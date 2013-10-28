/*
Copyright 2013 Zack Pierce.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.
*/
// Package stralgo implements various string algorithms
// with an emphasis on similarity metrics.

package stralgo

import (
	"errors"
)

// HammingDistance calculates the Hamming distance between
// two strings of equal length, bytewise.
//
// The Hamming distance is the total number of indices
// at which the corresponding bytes are different.
// See: http://en.wikipedia.org/wiki/Hamming_distance
//
// Returns an error if the string lengths are not equal.
func HammingDistance(a, b string) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen != bLen {
		return 0, errors.New("Hamming distance is undefined between strings of unequal length.")
	}
	var d int
	for i := 0; i < aLen; i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d, nil
}
// LeeDistance calculates the Lee distance between
// two strings of equal length, bytewise.
//
// See: http://en.wikipedia.org/wiki/Lee_distance
//
// Returns an error if the string lengths are not equal.
func LeeDistance(a, b string, q int) (int, error) {
	if q < 2 {
		return 0, errors.New("Lee distance must have a q-ary alphabet size greater than or equal to 2.")
	}
	aLen := len(a)
	bLen := len(b)
	if aLen != bLen {
		return 0, errors.New("Lee distance is undefined between strings of unequal length.")
	}
	var d int
	for i := 0; i < aLen; i++ {
		var diff int = int(a[i]) - int(b[i])
		if diff < 0 {
			diff = -diff
		}
		qDiff := q - diff
		if diff < qDiff {
			d += diff
		} else {
			d += qDiff
		}
	}
	return d, nil
}
