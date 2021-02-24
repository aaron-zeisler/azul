package testutils

import (
	"github.com/smartystreets/assertions"
)

func ShouldEqualError(actual interface{}, expected ...interface{}) string {
	if expected == nil || expected[0] == nil {
		return assertions.ShouldBeNil(actual)
	}

	actualError, ok := actual.(error)
	if !ok {
		return assertions.ShouldBeError(actual)
	}
	expectedError, ok := expected[0].(error)
	if !ok {
		return assertions.ShouldBeError(expected[0])
	}

	return assertions.ShouldContainSubstring(actualError.Error(), expectedError.Error())
}
