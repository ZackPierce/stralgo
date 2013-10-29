/*
Copyright 2013 Zack Pierce.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.
*/
/*
Package stralgo/runewise implements various string algorithms
with an emphasis on similarity metrics, implemented with
support for multi-byte runes.
*/
package runewise

import (
	"errors"
	"fmt"
	"unicode"
)

// HammingDistance calculates the Hamming distance between
// two strings containing equal numbers of runes.
//
// The Hamming distance is the total number of runes
// that differ at the same index within the two strings.
//
// The higher the result, the more different the strings.

// See: http://en.wikipedia.org/wiki/Hamming_distance
//
// Returns an error if the string rune counts are not equal.
func HammingDistance(a, b []rune) (uint, error) {
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
func DiceCoefficient(a, b []rune) (float64, error) {
	aLen := len(a)
	bLen := len(b)
	aLimit := aLen - 1
	bLimit := bLen - 1
	if aLimit < 1 && bLimit < 1 {
		return 0, errors.New("At least one of the input strings must contain 2 or more runes for the bigram-based DiceCoefficient to be calculated.")
	}
	aSet := make(map[runeBigram]bool, aLimit)
	totalBigrams := 0.0
	var bigram runeBigram
	for i := 0; i < aLimit; i++ {
		bigram = runeBigram{rA: a[i], rB: a[i+1]}
		if !aSet[bigram] {
			totalBigrams++
			aSet[bigram] = true
		}
	}

	bSet := make(map[runeBigram]bool, bLimit)
	sharedBigrams := 0.0
	for i := 0; i < bLimit; i++ {
		bigram = runeBigram{rA: b[i], rB: b[i+1]}
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
func WhiteSimilarity(a, b []rune) (float64, error) {
	aPairs, aLen := upperWordLetterPairs(a)
	bPairs, bLen := upperWordLetterPairs(b)
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

func upperWordLetterPairs(runes []rune) ([]runeBigram, int) {
	limit := len(runes) - 1
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
		bigrams[numPairs] = runeBigram{rA: unicode.ToUpper(a), rB: unicode.ToUpper(b)}
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
func LevenshteinDistance(a, b []rune) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen, nil
	}
	if bLen == 0 {
		return aLen, nil
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
			if a[i] == b[j] {
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

// DamerauLevenshteinDistance calculates the magnitude
// of difference between two strings using the Damerau-
// Levenshtein algorithm with adjacent-only transpositions,
// runewise.
//
// This edit distance is the minimum number of single-rune
// edits (insertions, deletions, substitutions, or
// transpositions) to transform one string into the other.
// DamerauLevenshtein differs from Levenshtein primarily
// in that DamerauLevenshtein considers adjacent-rune transpositions.
//
// The larger the result, the more different the strings.
//
// See: http://en.wikipedia.org/wiki/Damerau-Levenshtein_distance
func DamerauLevenshteinDistance(a, b []rune) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen, nil
	} else if bLen == 0 {
		return aLen, nil
	}

	// Swap to ensure a contains the shorter slice
	if aLen > bLen {
		a, aLen, b, bLen = b, bLen, a, aLen
	}
	rowLen := aLen + 1
	tranRow := make([]int, rowLen, rowLen)
	prevRow := make([]int, rowLen, rowLen)
	currRow := make([]int, rowLen, rowLen)
	for h := 0; h < rowLen; h++ {
		prevRow[h] = h
	}
	var prevB rune
	var cost int
	for i := 1; i <= bLen; i++ {
		currB := b[i-1]
		currRow[0] = i

		start := i - bLen - 1
		if start < 1 {
			start = 1
		}
		end := i + bLen + 1
		if end > aLen {
			end = aLen
		}

		var prevA rune
		for j := start; j <= end; j++ {
			currA := a[j-1]
			if currA == currB {
				cost = 0
			} else {
				cost = 1
			}
			entry := min(
				currRow[j-1]+1,
				prevRow[j]+1,
				prevRow[j-1]+cost)
			if currA == prevB && currB == prevA {
				trans := tranRow[j-2] + cost
				if trans < entry {
					entry = trans
				}
			}
			currRow[j] = entry
			prevA = currA
		}
		prevB = currB
		tranRow, prevRow, currRow = prevRow, currRow, tranRow
	}
	return prevRow[aLen], nil
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
