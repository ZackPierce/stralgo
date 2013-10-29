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
// The higher the result, the more different the strings.

// See: http://en.wikipedia.org/wiki/Hamming_distance
//
// Returns an error if the string lengths are not equal.
//
// Note that this algorithm implementation operates upon
// individual bytes, and does not account for multibyte
// unicode runes.
func HammingDistance(a, b string) (uint, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen != bLen {
		return 0, errors.New("Hamming distance is undefined between strings of unequal length.")
	}
	var d uint
	for i := 0; i < aLen; i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d, nil
}

// DiceCoefficent calculates the simiarlity of two
// strings per the Sorensen-Dice coefficient, bytewise.
//
// The resulting value is scaled between 0 and 1.0,
// and a higher value means a higher similarity.
//
// This algorithm is also known as the Sorensen Index, and
// is very close to the White Similarity metric, with the key
// distinctions that DiceCoefficient does not differentiate
// between whitespace and other characters and also does not
// account for bigram frequency count differences between
// the compared strings.
//
// See: http://en.wikipedia.org/wiki/Sorensen-Dice_coefficient
//
// Note that this algorithm implementation operates upon
// individual bytes and does not acocunt for multibyte
// unicode runes.
func DiceCoefficient(a, b string) (float64, error) {
	aLimit := len(a) - 1
	bLimit := len(b) - 1
	if aLimit < 1 && bLimit < 1 {
		return 0, errors.New("At least one of the input strings must have a length of 2 or greater for the bigram-based DiceCoefficient to be calculated.")
	}
	aSet := make(map[string]bool, aLimit)
	totalBigrams := 0.0
	var bigram string
	for i := 0; i < aLimit; i++ {
		bigram = a[i : i+2]
		if !aSet[bigram] {
			totalBigrams++
			aSet[bigram] = true
		}
	}

	bSet := make(map[string]bool, bLimit)
	sharedBigrams := 0.0
	for i := 0; i < bLimit; i++ {
		bigram = b[i : i+2]
		if !bSet[bigram] {
			totalBigrams++
			bSet[bigram] = true
			if aSet[bigram] {
				sharedBigrams++
			}
		}
	}
	return 2 * sharedBigrams / totalBigrams, nil

}
