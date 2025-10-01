package service

import (
	"log"
	"time"

	"github.com/beevik/ntp"
)

func GetCurrentTime(host string) (time.Time, error) {
	response, err := ntp.Query(host)
	if err != nil {
		log.Printf("Failed to get response from %s: %v", host, err)
		return time.Time{}, err
	}
	currentTime := time.Now().Add(response.ClockOffset)
	return currentTime, nil
}
