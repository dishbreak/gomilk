package utils_test

import (
	"testing"

	"github.com/dishbreak/gomilk/cli/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidSingles(t *testing.T) {
	input := []string{
		"1", "4", "5", "4",
	}

	expected := []int{
		1, 4, 5,
	}

	acutal, err := utils.ResolveIdentifier(input)
	if err != nil {
		t.Fatalf("Error resolving: %s", err)
	}

	assert.Equal(t, expected, acutal)
}

func TestValidRanges(t *testing.T) {
	input := []string{
		"2-5", "4", "17-19",
	}

	expected := []int{
		2, 3, 4, 5, 17, 18, 19,
	}

	actual, err := utils.ResolveIdentifier(input)
	if err != nil {
		t.Fatalf("Error resolving: %s", err)
	}

	assert.Equal(t, expected, actual)

}

func TestBadSingle(t *testing.T) {
	input := []string{
		"1", "4", "taco", "4",
	}

	_, err := utils.ResolveIdentifier(input)
	if err == nil {
		t.Fatal("Should have failed to resolve!")
	}
}

func TestBadRangeStart(t *testing.T) {
	input := []string{
		"1", "4", "taco-5", "4",
	}

	_, err := utils.ResolveIdentifier(input)
	if err == nil {
		t.Fatal("Should have failed to resolve!")
	}
}

func TestBadRangeEnd(t *testing.T) {
	input := []string{
		"1", "4", "5-taco", "4",
	}

	_, err := utils.ResolveIdentifier(input)
	if err == nil {
		t.Fatal("Should have failed to resolve!")
	}
}

func TestBadRangeInverted(t *testing.T) {
	input := []string{
		"1", "4", "5-3", "4",
	}

	_, err := utils.ResolveIdentifier(input)
	if err == nil {
		t.Fatal("Should have failed to resolve!")
	}
}
