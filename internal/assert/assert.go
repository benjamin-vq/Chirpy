package assert

import "log"

func That(condition bool, message string) {
	if !condition {
		log.Panicf(message)
	}
}
