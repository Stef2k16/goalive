package monitor

import (
	"testing"
	"time"
)

func TestStatusToStringSuccess(t *testing.T) {
	timestamp := time.Date(2021, time.Month(6), 11, 20, 36, 14, 0, time.UTC)
	st := status{
		timestamp:          timestamp,
		success:            true,
		code:               200,
		body:               "Health OK",
		url:                "www.example.com",
		notificationFailed: false,
	}
	stString := st.String()
	expected := "Request for www.example.com SUCCEEDED at Fri, 11 Jun 2021 20:36:14 UTC\n\tStatus: 200\n\tBody: Health OK"
	if !(stString == expected) {
		t.Errorf("st.String(): \nGot\n%s\nExpected:\n%s", stString, expected)
	}
}

func TestStatusToStringFailure(t *testing.T) {
	timestamp := time.Date(2021, time.Month(6), 11, 20, 36, 14, 0, time.UTC)
	st := status{
		timestamp:          timestamp,
		success:            false,
		code:               500,
		body:               "Healthcheck failed",
		url:                "www.example.com",
		notificationFailed: false,
	}
	stString := st.String()
	expected := "Request for www.example.com FAILED at Fri, 11 Jun 2021 20:36:14 UTC\n\tStatus: 500\n\tBody: Healthcheck failed"
	if !(stString == expected) {
		t.Errorf("st.String(): \nGot\n%s\nExpected:\n%s", stString, expected)
	}
}
