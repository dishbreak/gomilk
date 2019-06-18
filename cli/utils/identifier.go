package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func ResolveIdentifier(identifers []string) ([]int, error) {
	result := make([]int, 0)

	for _, identifier := range identifers {
		if parts := strings.Split(identifier, "-"); len(parts) == 2 {
			var start, end int
			start, err := strconv.Atoi(parts[0])
			if err != nil {
				return result, fmt.Errorf("Error finding start of range '%s': %s", identifier, err)
			}

			end, err = strconv.Atoi(parts[1])
			if err != nil {
				return result, fmt.Errorf("Error finding end of range '%s': %s", identifier, err)
			}

			if start > end {
				return result, fmt.Errorf("Error in range '%s': start must be less than end", identifier)
			}

			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else if val, err := strconv.Atoi(identifier); err == nil {
			result = append(result, val)
		} else {
			return result, fmt.Errorf("failed to parse identifier: %s. must be int or range", identifier)
		}
	}

	sort.Ints(result)

	return uniq(result), nil

}

func uniq(input []int) []int {
	output := make([]int, 0)
	for idx, value := range input {
		if idx == 0 || input[idx-1] != value {
			output = append(output, value)
		}
	}
	return output
}
