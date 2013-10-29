package stralgo

import (
	"github.com/stretchr/testify/assert"
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
}
