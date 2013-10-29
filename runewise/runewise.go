/*
Copyright 2013 Zack Pierce.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.
*/
/*
Package stralgo/runewise implements various string algorithms
with an emphasis on similarity metrics, implemented with
support for multi-byte runes.

This bytewise approach is suited for accurate comparisons
when the encoded data truly makes use of UTF8.
*/
package runewise

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// HammingDistance calculates the Hamming distance between
// two strings containing equal numbers of runes.
//
// The Hamming distance is the total number of runes
// that differ at the same index in the resolved array
// of runes from each string.
// The higher the result, the more different the strings.

// See: http://en.wikipedia.org/wiki/Hamming_distance
//
// Returns an error if the string rune counts are not equal.
func HammingDistance(a, b string) (uint, error) {
	aRunes, aLen := runeSlice(a)
	bRunes, bLen := runeSlice(b)

	if aLen != bLen {
		return 0, errors.New("Hamming distance is undefined between strings of unequal length.")
	}
	var d uint
	for i := 0; i < aLen; i++ {
		if aRunes[i] != bRunes[i] {
			d++
		}
	}
	return d, nil
}

// DiceCoefficent calculates the simiarlity of two
// strings per the Sorensen-Dice coefficient, runewise.
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
// Returns an error if both of the input strings contain
// less than two runes.
func DiceCoefficient(a, b string) (float64, error) {
	aRunes, aLen := runeSlice(a)
	bRunes, bLen := runeSlice(b)
	aLimit := aLen - 1
	bLimit := bLen - 1
	if aLimit < 1 && bLimit < 1 {
		return 0, errors.New("At least one of the input strings must contain 2 or more runes for the bigram-based DiceCoefficient to be calculated.")
	}
	aSet := make(map[runeBigram]bool, aLimit)
	totalBigrams := 0.0
	var bigram runeBigram
	for i := 0; i < aLimit; i++ {
		bigram = runeBigram{rA: aRunes[i], rB: aRunes[i+1]}
		if !aSet[bigram] {
			totalBigrams++
			aSet[bigram] = true
		}
	}

	bSet := make(map[runeBigram]bool, bLimit)
	sharedBigrams := 0.0
	for i := 0; i < bLimit; i++ {
		bigram = runeBigram{rA: bRunes[i], rB: bRunes[i+1]}
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

// WhiteSimilarity calculates the similarity of two
// strings through a variation on the Sorensen-Dice
// Coefficient algorithm.
//
// The resulting value is scaled between 0 and 1.0,
// and a higher value means a higher similarity.
//
// WhiteSimilarity differs from DiceCoefficient in that
// it disregards bigrams that include whitespace, and
// applies an upper-case filter, and accounts for bigram
// frequency.
//
// See: http://www.catalysoft.com/articles/strikeamatch.html
//
// Returns an error if neither of the input strings
// contains at least one rune bigram without whitespace.
func WhiteSimilarity(a, b string) (float64, error) {
	aPairs, aLen := wordLetterPairs(strings.ToUpper(a))
	bPairs, bLen := wordLetterPairs(strings.ToUpper(b))
	union := aLen + bLen
	if union == 0 {
		return 0.0, errors.New("At least one of the input strings must contain two or more non-whitespace rune bigrams in order to calculate the White Similarity.")
	}
	intersection := 0.0
	for _, aBigram := range aPairs {
		for j, bBigram := range bPairs {
			if aBigram == bBigram {
				intersection++
				bPairs[j] = runeBigram{}
				break
			}
		}
	}
	return 2 * intersection / float64(union), nil
}

func wordLetterPairs(s string) ([]runeBigram, int) {
	runes, n := runeSlice(s)
	limit := n - 1
	if limit < 1 {
		return make([]runeBigram, 0), 0
	}
	bigrams := make([]runeBigram, limit)
	var a rune
	var b rune
	numPairs := 0
	for i := 0; i < limit; i++ {
		a = runes[i]
		b = runes[i+1]
		if unicode.IsSpace(b) {
			i++
			continue
		}
		if unicode.IsSpace(a) {
			continue
		}
		bigrams[numPairs] = runeBigram{rA: a, rB: b}
		numPairs++
	}
	bigrams = bigrams[0:numPairs]
	return bigrams, numPairs
}

// LevenshteinDistance calculates the magnitude of
// difference between two strings using the
// Levenshtein Distance metric.
//
// This edit distance is the minimum number of single-rune
// edits (insertions, deletions, or substitutions) needed
// to transform one string into the other.
//
// The larger the result, the more different the strings.
//
// See: http://en.wikipedia.org/wiki/Levenshtein_distance
func LevenshteinDistance(a, b string) (int, error) {
	aRunes, aLen := runeSlice(a)
	bRunes, bLen := runeSlice(b)
	if aLen == 0 {
		return bLen, nil
	}
	if bLen == 0 {
		return aLen, nil
	}
	if aLen == bLen && a == b {
		return 0, nil
	}

	rowLen := bLen + 1
	prevRow := make([]int, rowLen, rowLen)
	currRow := make([]int, rowLen, rowLen)
	for h := 0; h < rowLen; h++ {
		prevRow[h] = h
	}
	cost := 0
	for i := 0; i < aLen; i++ {
		currRow[0] = i + 1
		for j := 0; j < bLen; j++ {
			if aRunes[i] == bRunes[j] {
				cost = 0
			} else {
				cost = 1
			}
			currRow[j+1] = min(
				currRow[j]+1,
				prevRow[j+1]+1,
				prevRow[j]+cost)
		}
		prevRow, currRow = currRow, prevRow
	}
	return prevRow[bLen], nil
}

func min(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		return c
	}
	return m
}

type runeBigram struct {
	rA, rB rune
}

func (r runeBigram) String() string {
	return fmt.Sprintf("{%q, %q}", r.rA, r.rB)
}

func runeSlice(s string) ([]rune, int) {
	n := 0
	runes := make([]rune, len(s))
	for _, r := range s {
		runes[n] = r
		n++
	}
	runes = runes[0:n]
	return runes, n
}
