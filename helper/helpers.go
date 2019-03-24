package helper

import (
	"errors"
	"math"
)

func BitsToNumbers(str string) int {
	sum := 0
	length := len(str)
	for _, v := range str {
		sum += (int(v) - '0') * int(math.Pow(2, float64(length-1)))
		length -= 1
	}
	return sum
}

func BitsToComplementNumber(str string) int {
	sum := 0
	length := len(str)
	for i := 1; i < length; i++ {
		sum += (int(str[i]) - '0') * int(math.Pow(2, float64(length-1)))
		length -= 1
	}
	if str[0] == '1' {
		sum = (int(math.Pow(2, float64(len(str)-1))) - sum) * -1
	}
	return sum
}

func Bits2Ascii(str string) (string, error) {
	result := ""
	for i := 0; i < len(str); i += 6 {
		value := BitsToAscii[str[i:i+6]]
		if value == "" {
			return "", errors.New("bits can't change to ascii")
		}
		result += value
	}
	return result, nil
}
