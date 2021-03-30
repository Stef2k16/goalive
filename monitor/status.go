package monitor

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// status holds the results and metadata of one request to an endpoint.
type status struct {
	timestamp          time.Time
	success            bool
	code               int
	body               string
	url                string
	notificationFailed bool
}

func (st *status) String() string {
	if st.success {
		return fmt.Sprintf("Request for %v SUCCEEDED at %s\n\tStatus: %d\n\tBody: %s",
			st.url, st.timestamp.Format(time.RFC1123), st.code, st.body)
	}
	return fmt.Sprintf("Request for %s FAILED at %s\n\tStatus: %d\n\tBody: %s",
		st.url, st.timestamp.Format(time.RFC1123), st.code, st.body)
}

// getStatus sends a get request to request the status of a service.
func getStatus(url string) status {
	resp, err := http.Get(url)
	if err != nil {
		return status{
			timestamp:          time.Now(),
			success:            false,
			code:               0,
			body:               err.Error(),
			url:                url,
			notificationFailed: false,
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return status{
			timestamp:          time.Now(),
			success:            false,
			code:               0,
			body:               err.Error(),
			url:                url,
			notificationFailed: false,
		}
	}
	successfulResponse := resp.StatusCode >= 200 && resp.StatusCode < 300
	st := status{
		timestamp:          time.Now(),
		success:            successfulResponse,
		code:               resp.StatusCode,
		body:               string(body),
		url:                url,
		notificationFailed: false,
	}
	return st
}
