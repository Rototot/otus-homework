package main

import (
	"fmt"
	"io"
	"time"

	"github.com/beevik/ntp"
)

func WriteCurrentTime(w io.Writer) error {
	currentTime := time.Now()

	exactTime, err := getExactTime()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fmt.Sprintf(
		"current time: %s\nexact time: %s\n",
		currentTime.Round(0).String(),
		exactTime.Round(0).String(),
	)))

	return err
}

func getExactTime() (time.Time, error) {
	return ntp.Time("0.beevik-ntp.pool.ntp.org")
}
