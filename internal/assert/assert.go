package assert

import "log"

// That Condition should be true, otherwise the program could be in an invalid state, might as well panic.
func That(condition bool, message string, args ...any) {
	if !condition {
		log.Panicf(message, args)
	}
}
