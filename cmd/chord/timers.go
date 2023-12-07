package main

import (
	"time"
)

func stabilizeTimer(ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func checkPredecessorTimer(ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func fixFingersTimer(ms int) {
	for !finished {
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}
