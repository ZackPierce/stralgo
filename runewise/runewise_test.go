﻿package runewise

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_HammingDistance(t *testing.T) {
	d, err := HammingDistance("toned", "roses")
	assert.Nil(t, err, "HammingDistance should not return an error for two strings of equivalent length.")
	assert.Equal(t, uint(3), d, "HammingDistance between 'toned' and 'roses' should be 3")

	d, err = HammingDistance("", "")
	assert.Nil(t, err, "HammingDistance should not return an error for two empty strings")
	assert.Equal(t, uint(0), d, "HammingDistance between two empty strings should be 0")

	d, err = HammingDistance("spam", "spam")
	assert.Nil(t, err, "HammingDistance should not return an error for two equivalent non-empty strings.")
	assert.Equal(t, uint(0), d, "HammingDistance between 'spam' and 'spam' should be 0, but was %i", d)

	d, err = HammingDistance("green eggs", "ham")
	assert.NotNil(t, err, "HammingDistance between 'green eggs' and 'ham' should produce an error due to unequal lengths")
	assert.Equal(t, uint(0), d)

	d, err = HammingDistance("日本語", "日本ゴ")
	assert.Nil(t, err)
	assert.Equal(t, uint(1), d)

	d, err = HammingDistance("日本語", "日本g")
	assert.Nil(t, err)
	assert.Equal(t, uint(1), d)

	d, err = HammingDistance("日本語", "日本go")
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), d)
}

func Test_DiceCoefficient(t *testing.T) {
	c, err := DiceCoefficient("night", "nacht")
	assert.Nil(t, err)
	assert.Equal(t, 1.0/4.0, c)

	c, err = DiceCoefficient("GGGG", "GGGG")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = DiceCoefficient("", "")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = DiceCoefficient("a", "b")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = DiceCoefficient("GG", "GGGG")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c, "Naive Dice coefficient does not account for differences in occurrence-count for bigrams.")

	c, err = DiceCoefficient("日本語", "日本語")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = DiceCoefficient("日本語", "日本ゴ")
	assert.Nil(t, err)
	assert.Equal(t, 2.0/4.0, c)

	c, err = DiceCoefficient("日本語", "日本g")
	assert.Nil(t, err)
	assert.Equal(t, 2.0/4.0, c)

	c, err = DiceCoefficient("日", "本")
	assert.NotNil(t, err, "DiceCoefficient runewise reports an error due to lack of rune-bigrams.")
	assert.Equal(t, 0.0, c)

	// [230 151 165] and [230, 151, 168]
	c, err = DiceCoefficient("日", "旨")
	assert.NotNil(t, err, "DiceCoefficient runewise reports an error due to lack of rune-bigrams.")
	assert.Equal(t, 0.0, c)
}

func Test_WhiteSimilarity(t *testing.T) {
	c, err := WhiteSimilarity("Healed", "Healed")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity("Healed", "Sealed")
	assert.Nil(t, err)
	EqualWithin(t, 0.8, c, 0.01, "Sealed")

	c, err = WhiteSimilarity("Healed", "Healthy")
	assert.Nil(t, err)
	EqualWithin(t, 0.55, c, 0.01, "Healthy")

	c, err = WhiteSimilarity("Healed", "Heard")
	assert.Nil(t, err)
	EqualWithin(t, 0.44, c, 0.01, "Heard")

	c, err = WhiteSimilarity("Healed", "Herded")
	assert.Nil(t, err)
	EqualWithin(t, 0.40, c, 0.01, "Herded")

	c, err = WhiteSimilarity("Healed", "Help")
	assert.Nil(t, err)
	EqualWithin(t, 0.25, c, 0.01, "Help")

	c, err = WhiteSimilarity("Healed", "Sold")
	assert.Nil(t, err)
	EqualWithin(t, 0.0, c, 0.01, "Sold")

	c, err = WhiteSimilarity("Healed ", "HEALed")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity("GGGG", "GGGG")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity("REPUBLIC OF FRANCE", "FRANCE")
	assert.Nil(t, err)
	EqualWithin(t, 0.56, c, 0.01)

	c, err = WhiteSimilarity("FRANCE", "QUEBEC")
	assert.Nil(t, err)
	EqualWithin(t, 0.0, c, 0.001)

	c, err = WhiteSimilarity("FRENCH REPUBLIC", "REPUBLIC OF FRANCE")
	assert.Nil(t, err)
	EqualWithin(t, 0.72, c, 0.01)

	c, err = WhiteSimilarity("FRENCH REPUBLIC", "REPUBLIC OF CUBA")
	assert.Nil(t, err)
	EqualWithin(t, 0.61, c, 0.01)

	c, err = WhiteSimilarity("", "")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = WhiteSimilarity("a", "b")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = WhiteSimilarity("GG", "GGGGG")
	assert.Nil(t, err)
	EqualWithin(t, 0.4, c, 0.01)

	c, err = WhiteSimilarity("GGGGG", "GG")
	assert.Nil(t, err)
	EqualWithin(t, 0.4, c, 0.01)
}

func EqualWithin(t *testing.T, a, b, delta float64, msgAndArgs ...interface{}) bool {
	if math.Abs(a-b) > delta {
		return assert.Fail(t, fmt.Sprintf("Not within delta: Abs(%#v - %#v) > %#v", a, b, delta), msgAndArgs...)
	}

	return true
}