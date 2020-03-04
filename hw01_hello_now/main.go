package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalln(err)
	}
	localTime := time.Now()
	fmt.Printf("current time: %s\n", localTime.String())
	fmt.Printf("exact time: %s\n", ntpTime.String())
}
