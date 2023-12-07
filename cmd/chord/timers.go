package main

import (
	"time"
	"log"
)

func stabilizeTimer(ms int) {
	for !finished {
		log.Println("stab")
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func checkPredecessorTimer(ms int) {
	for !finished {
		log.Println("check pred")
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func fixFingersTimer(ms int) {
	for !finished {
		log.Println("fix fingers")
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}
