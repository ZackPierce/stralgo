package stralgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HammingDistance(t *testing.T) {
	d, err := HammingDistance("toned", "roses")
	assert.Nil(t, err, "HammingDistance should not return an error for two strings of equivalent length.")
	assert.Equal(t, 3, d, "HammingDistance between 'toned' and 'roses' should be 3")

	d, err = HammingDistance("", "")
	assert.Nil(t, err, "HammingDistance should not return an error for two empty strings")
	assert.Equal(t, 0, d, "HammingDistance between two empty strings should be 0")

	d, err = HammingDistance("spam", "spam")
	assert.Nil(t, err, "HammingDistance should not return an error for two equivalent non-empty strings.")
	assert.Equal(t, 0, d, "HammingDistance between 'spam' and 'spam' should be 0, but was %i", d)

	d, err = HammingDistance("green eggs", "ham")
	assert.NotNil(t, err, "HammingDistance between 'green eggs' and 'ham' should produce an error due to unequal lengths")
}
