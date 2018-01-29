/*
Generate and validate Swedish personal identity numbers.

The library was based on the details that can be found here: https://en.wikipedia.org/wiki/Personal_identity_number_(Sweden)
*/
package pin

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// validDate validates a birthday date.
var validDate *regexp.Regexp

// validPIN validates a personal identity number.
var validPIN *regexp.Regexp

// init initializes the regular expressions.
func init() {
	// Initialize regular expressions
	validDate = regexp.MustCompile(`^[0-9]{8}$`)
	validPIN = regexp.MustCompile(`^[0-9]{8}[-+][0-9]{4}$`)
}

// rndDate returns a random date between 1970 and the current date.
func rndDate() string {
	min := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	sec := rnd.Int63n(min+max) + min
	date := time.Unix(sec, 0)
	return date.Format("20060102")
}

// rndNumber generates a random number between min and max.
func rndNumber(min, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return rnd.Intn(max-min) + min
}

// Generate generates a randomized personal identity number.
func Generate() (string, error) {
	return GenerateFromDate(rndDate())
}

// GenerateFromDate generates a personal identity number for the given date.
func GenerateFromDate(date string) (string, error) {
	// Make sure that we get a valid date.
	if !validDate.Match([]byte(date)) {
		return "", fmt.Errorf("%s is not a valid date, expects format YYYYMMDD", date)
	}

	// Generate a birth number.
	// This should be a number between 100 and 999.
	bn := strconv.Itoa(rndNumber(100, 999))

	// Get the control number.
	// Exclude the first two digits of the date, 1984 becomes 84.
	cn := getControlNumber(strings.Join([]string{date[2:], bn}, ""))

	// Return birthday, birth number and control number.
	return fmt.Sprintf("%s-%s%s", date, bn, string(cn)), nil
}

// getControlNumber calculates the control number in a personal identity number.
//
// The control number is calculated by using the Lugn algorithm.
// More information about it can be found at https://en.wikipedia.org/wiki/Luhn_algorithm
func getControlNumber(s string) byte {
	// Store the sum of all added numbers here.
	sum := 0

	// Iterate over s and multiply each digit with either 1 or 2.
	for i, c := range s {
		// Convert the character to an integer.
		n, _ := strconv.Atoi(string(c))

		// Even numbers are multiplied with 2, odd with 1.
		var t int
		if i%2 == 0 {
			t = n * 2
		} else {
			t = n * 1
		}

		// If t is greater or equal to 10, split it up in to two digits.
		if t >= 10 {
			sum += 1
			sum += t - 10
		} else {
			sum += t
		}
	}

	// Calculate the control number.
	// We first convert the control number to a string.
	// We know that the number always will be a digit between 0 and 9,
	// so it's safe to cast the result to a byte when we return it.
	return byte(strconv.Itoa((10 - (sum % 10)) % 10)[0])
}

// IsFemale returns true if the personal identity number belongs to a female.
func IsFemale(pin string) (bool, error) {
	// Make sure that the supplied pin is valid.
	if _, err := IsValid(pin); err != nil {
		return false, fmt.Errorf("%s is not a valid personal identity number", pin)
	}

	// IsMale should return an error, otherwise we got a male pin.
	if _, err := IsMale(pin); err == nil {
		return false, fmt.Errorf("%s is not a valid female personal identity number", pin)
	}

	// We've got a woman.
	return true, nil
}

// IsMale returns true if the personal identity number belongs to a male.
func IsMale(pin string) (bool, error) {
	// Make sure that the supplied pin is valid.
	if _, err := IsValid(pin); err != nil {
		return false, fmt.Errorf("%s is not a valid personal identity number", pin)
	}

	// Extract the last but one digit from the personal id number.
	c := pin[len(pin)-2]

	// Convert the digit to an integer.
	n, _ := strconv.Atoi(string(c))

	// Odd numbers indicates that the personal identity number belongs to a male.
	if n%2 == 0 {
		return false, fmt.Errorf("%s is not a valid male personal identity number", pin)
	}

	// It's a man.
	return true, nil
}

// IsValid checks whether the supplied personal identity number is valid.
func IsValid(pin string) (bool, error) {
	// Make sure that the input matches the validPIN regular expression.
	if !validPIN.Match([]byte(pin)) {
		return false, fmt.Errorf("%s is not a valid personal identity number", pin)
	}

	// This variable will contain a cleaned version of the personal identity number.
	// We begin by stripping the first two digits of the value.
	c := pin[2:]

	// The separator will always be located at the position 7 in the string.
	// So we will concatenate a new string and skip position 7 to remove it.
	c = c[:6] + c[7:]

	// The control number is calculated based on the first 9 digits of the
	// personal identity number.
	cn := getControlNumber(c[:9])

	// Compare the calculated control number with the supplied control number.
	if cn != c[9] {
		return false, fmt.Errorf("%s is not a valid personal identity number", pin)
	}

	// The control digit is valid, return true.
	return true, nil
}
