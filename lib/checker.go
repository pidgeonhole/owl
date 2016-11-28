package owl

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"strings"
)

const (
	STR = iota
	INT = iota
	FP  = iota
)

/*
Go implementation of the python function of the same name from https://code.google.com/codejam/faq.html#5-9 for handling
floating point answers:

	def IsApproximatelyEqual(x, y, epsilon):
	  """Returns True iff y is within relative or absolute 'epsilon' of x.

	  By default, 'epsilon' is 1e-6.
	  """
	  # Check absolute precision.
	  if -epsilon <= x - y <= epsilon:
	    return True

	  # Is x or y too close to zero?
	  if -epsilon <= x <= epsilon or -epsilon <= y <= epsilon:
	    return False

	  # Check relative precision.
	  return (-epsilon <= (x - y) / x <= epsilon
	       or -epsilon <= (x - y) / y <= epsilon)
*/
func isApproximatelyEqual(x, y, epsilon float64) bool {
	// Check absolute precision.
	if -epsilon <= x-y && x-y <= epsilon {
		return true
	}

	// Is x or y too close to zero?
	if (-epsilon <= x && x <= epsilon) || (-epsilon <= y && y <= epsilon) {
		return false
	}

	// Check relative precision.
	return (-epsilon <= (x-y)/x && (x-y)/x <= epsilon) ||
		(-epsilon <= (x-y)/y && (x-y)/y <= epsilon)
}

func Check(ans, output io.Reader, types []int) (bool, error) {
	ascan := bufio.NewScanner(ans)
	oscan := bufio.NewScanner(output)

	for oscan.Scan() && ascan.Scan() {
		atokens := strings.Split(ascan.Text(), " ")
		otokens := strings.Split(oscan.Text(), " ")

		if len(atokens) != len(types) {
			return false, errors.New("number of tokens in answer does not match types")
		}

		if len(atokens) != len(otokens) {
			return false, nil
		}

		for i, t := range types {
			switch t {
			case STR:
				fallthrough
			case INT:
				if atokens[i] != otokens[i] {
					return false, nil
				}
			case FP:
				afp, err := strconv.ParseFloat(atokens[i], 64)
				if err != nil {
					return false, errors.Wrap(err, "failed to parse floating point token in answer")
				}

				ofp, err := strconv.ParseFloat(otokens[i], 64)
				if err != nil {
					return false, nil
				}

				if !isApproximatelyEqual(afp, ofp, 10E-9) {
					return false, nil
				}
			}
		}
	}

	// check that there is no error in ascan
	if err := ascan.Err(); err != nil {
		return false, errors.Wrap(err, "error while scanning answer")
	}

	// check if there was an error in oscan
	if err := oscan.Err(); err != nil {
		return false, nil
	}

	// check that ascan and oscan have both reached EOF
	if ascan.Scan() || oscan.Scan() {
		return false, nil
	}

	return true, nil
}
