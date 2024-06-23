package assert

import "log"

// That Condition should be true, otherwise the program could be in an invalid state, might as well panic.
func That(condition bool, message string, args ...any) {
	if !condition {
		log.Panicf(message, args)
	}
}

// NoError Error should be nil, otherwise the program could be in an invalid state, might as well panic.
func NoError(err error, message string, args ...any) {
	if err != nil {
		log.Panicf(message, args)
	}
}
