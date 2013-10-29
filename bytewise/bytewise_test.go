package bytewise

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
	
	d, err = HammingDistance("日本語", "日本ゴ")
	assert.Nil(t, err)
	assert.Equal(t, uint(3), d)

	d, err = HammingDistance("日本語", "日本g")
	assert.NotNil(t, err, "HammingDistance bytewise produces an error due to the difference in total bytes in the two strings, even though the rune-counts are equivalent.")
	assert.Equal(t, uint(0), d)

	d, err = HammingDistance("日本語", "日本gon")
	assert.Nil(t, err, "HammingDistance bytewise does not produce an error when comparing strings of equal byte-lengths, even though the rune-counts are different.")
	assert.Equal(t, uint(3), d)
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
	
	c, err = DiceCoefficient("日", "本")
	assert.Nil(t, err, "DiceCoefficient bytewise does not report an error about lack of bigrams for this case, because the runes involved have a width of 2 or greater.")
	assert.Equal(t, 0.0, c)
	
	// [230 151 165] and [230, 151, 168]
	c, err = DiceCoefficient("日", "旨")
	assert.Nil(t, err, "DiceCoefficient bytewise does not report an error about lack of bigrams for this case, because the runes involved have a width of 2 or greater.")
	assert.Equal(t, 2.0/4.0, c)
}
