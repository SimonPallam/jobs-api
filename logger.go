package main

import "log"

var debugLevel = 0

func debug(format string, v ...interface{}) {
	if debugLevel == 1 {
		println("\n    DEBUG :")
		log.Printf("    "+format, v...)
		println()
	}
	return
}
