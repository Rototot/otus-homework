package main

import (
	"fmt"
	"io"
	"time"

	"github.com/beevik/ntp"
)

type timeGetter = func() (time.Time, error)

func WriteCurrentTime(w io.Writer) error {
	var formattedTime string

	currentTime, err := getTime(getCurrentTime)
	if err != nil {
		return err
	}
	formattedTime += formatTime(currentTime, "current time: %s")

	exactTime, err := getTime(getExactTime)
	if err != nil {
		return err
	}
	formattedTime += formatTime(exactTime, "\nexact time: %s")

	_, err = w.Write([]byte(formattedTime))

	return err
}

func formatTime(t time.Time, layout string) string {
	return fmt.Sprintf(layout, t.Round(0).String())
}

func getTime(adapter timeGetter) (time.Time, error) {
	receivedTime, err := adapter()

	if err != nil {
		return time.Time{}, fmt.Errorf("cannot get time. Reason: %s", err.Error())
	}

	return receivedTime, nil
}

func getCurrentTime() (time.Time, error) {
	return time.Now(), nil
}

func getExactTime() (time.Time, error) {
	return ntp.Time("0.beevik-ntp.pool.ntp.org")
}
