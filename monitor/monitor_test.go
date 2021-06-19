package monitor

import (
	"github.com/Stef2k16/goalive/notification"
	"log"
	"sync"
	"testing"
	"time"
)

func TestStatusSummarySucceeded(t *testing.T) {
	timestamp := time.Date(2021, time.Month(6), 11, 20, 36, 14, 0, time.UTC)
	prevStatus := map[string]status{
		"https://example.com": {
			timestamp:          timestamp,
			success:            true,
			code:               200,
			body:               "Health OK",
			url:                "https://example.com",
			notificationFailed: false,
		},
	}
	m := Monitor{
		logger:             log.Default(),
		NotificationClient: notification.NewMockClient(),
		urls:               []string{"https://example.com"},
		pollingInterval:    0,
		prevStatus:         prevStatus,
		prevStatusMutex:    &sync.Mutex{},
	}
	summary := m.statusSummary()
	expectedSummary := "https://example.com polled at Fri, 11 Jun 2021 20:36:14 UTC:\n"
	expectedSummary += "\tRequest SUCCEEDED with Status 200\n\tBody: Health OK\n\n"
	if !(summary == expectedSummary) {
		t.Errorf(`m.statusSummary()
			Got: %s
			Expected: %s`, summary, expectedSummary)
	}
}

func TestStatusSummaryFailed(t *testing.T) {
	timestamp := time.Date(2021, time.Month(6), 11, 20, 36, 14, 0, time.UTC)
	prevStatus := map[string]status{
		"https://example.com": {
			timestamp:          timestamp,
			success:            false,
			code:               404,
			body:               "Not found",
			url:                "https://example.com",
			notificationFailed: false,
		},
	}
	m := Monitor{
		logger:             log.Default(),
		NotificationClient: notification.NewMockClient(),
		urls:               []string{"https://example.com"},
		pollingInterval:    0,
		prevStatus:         prevStatus,
		prevStatusMutex:    &sync.Mutex{},
	}
	summary := m.statusSummary()
	expectedSummary := "https://example.com polled at Fri, 11 Jun 2021 20:36:14 UTC:\n"
	expectedSummary += "\tRequest FAILED with Status 404\n\tBody: Not found\n\n"
	if !(summary == expectedSummary) {
		t.Errorf(`m.statusSummary()
			Got: %s
			Expected: %s`, summary, expectedSummary)
	}
}

func TestStatusSummaryNotificationFailed(t *testing.T) {
	timestamp := time.Date(2021, time.Month(6), 11, 20, 36, 14, 0, time.UTC)
	prevStatus := map[string]status{
		"https://example.com": {
			timestamp:          timestamp,
			success:            false,
			code:               0,
			body:               "Sending notification failed",
			url:                "https://example.com",
			notificationFailed: true,
		},
	}
	m := Monitor{
		logger:             log.Default(),
		NotificationClient: notification.NewMockClient(),
		urls:               []string{"https://example.com"},
		pollingInterval:    0,
		prevStatus:         prevStatus,
		prevStatusMutex:    &sync.Mutex{},
	}
	summary := m.statusSummary()
	expectedSummary := "https://example.com polled at Fri, 11 Jun 2021 20:36:14 UTC:\n"
	expectedSummary += "\tRequest FAILED with Status 0\n\tBody: Sending notification failed\n\n"
	if !(summary == expectedSummary) {
		t.Errorf(`m.statusSummary()
			Got: %s
			Expected: %s`, summary, expectedSummary)
	}
}

func TestRespondToStatusSuccess(t *testing.T) {
	url := "https://httpstat.us/200"
	prevStatus := map[string]status{
		url: {
			timestamp:          time.Time{},
			success:            true,
			code:               200,
			body:               "Health OK",
			url:                url,
			notificationFailed: false,
		},
	}
	m := Monitor{
		logger:             log.Default(),
		NotificationClient: notification.NewMockClient(),
		urls:               []string{url},
		pollingInterval:    10,
		prevStatus:         prevStatus,
		prevStatusMutex:    &sync.Mutex{},
	}

	m.respondToStatus(url)
	st := m.prevStatus[url]
	st.timestamp = time.Time{}
	expected := status{
		timestamp:          time.Time{},
		success:            true,
		code:               200,
		body:               "",
		url:                url,
		notificationFailed: false,
	}
	if st != expected {
		t.Errorf(`m.respondToStatus(%s)
			Got: %v
			Expected: %v`, url, st, expected)
	}
}

func TestRespondToStatusSuccessWithPreviousFailure(t *testing.T) {
	url := "https://httpstat.us/200"
	prevStatus := map[string]status{
		url: {
			timestamp:          time.Time{},
			success:            false,
			code:               404,
			body:               "Health not OK",
			url:                url,
			notificationFailed: false,
		},
	}
	m := Monitor{
		logger:             log.Default(),
		NotificationClient: notification.NewMockClient(),
		urls:               []string{url},
		pollingInterval:    10,
		prevStatus:         prevStatus,
		prevStatusMutex:    &sync.Mutex{},
	}

	m.respondToStatus(url)
	st := m.prevStatus[url]
	st.timestamp = time.Time{}
	expected := status{
		timestamp:          time.Time{},
		success:            true,
		code:               200,
		body:               "",
		url:                url,
		notificationFailed: false,
	}
	if st != expected {
		t.Errorf(`m.respondToStatus(%s)
			Got: %v
			Expected: %v`, url, st, expected)
	}
}

func TestRespondToStatusFailureWithPreviousSuccess(t *testing.T) {
	url := "https://httpstat.us/404"
	prevStatus := map[string]status{
		url: {
			timestamp:          time.Time{},
			success:            true,
			code:               200,
			body:               "Health OK",
			url:                url,
			notificationFailed: false,
		},
	}
	m := Monitor{
		logger:             log.Default(),
		NotificationClient: notification.NewMockClient(),
		urls:               []string{url},
		pollingInterval:    10,
		prevStatus:         prevStatus,
		prevStatusMutex:    &sync.Mutex{},
	}

	m.respondToStatus(url)
	st := m.prevStatus[url]
	st.timestamp = time.Time{}
	expected := status{
		timestamp:          time.Time{},
		success:            false,
		code:               404,
		body:               "",
		url:                url,
		notificationFailed: false,
	}
	if st != expected {
		t.Errorf(`m.respondToStatus(%s)
			Got: %v
			Expected: %v`, url, st, expected)
	}
}
