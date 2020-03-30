package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

type TimeGetter = func() (time.Time, error)

func main() {
	err := TimeHandler(os.Stdout)

	if err == nil {
		os.Exit(0)
	}

	log.Fatalln(err)
	//os.Exit(1);
}

func TimeHandler(w io.Writer) error {
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
	var formattedTime string = t.Round(0).String()

	return fmt.Sprintf(layout, formattedTime)
}

func getTime(adapter TimeGetter) (time.Time, error) {
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
	//startTime := time.Now()
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	//latency := time.Now().Sub(startTime)

	if err != nil {
		return time.Time{}, err
	}

	return exactTime, nil
}
