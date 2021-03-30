package monitor

import (
	"fmt"
	"galive/notification"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
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

// Monitor provides the functionality to monitor http endpoints.
type Monitor struct {
	logger             *log.Logger
	notificationClient notification.Client
	urls               []string
	pollingInterval    time.Duration
	prevStatus         map[string]status
	prevStatusMutex    *sync.Mutex
}

// New sets up a new monitor to check services.
func New(logfile string, notificationClient notification.Client, urls []string, pollingInterval time.Duration) (*Monitor, error) {
	file, err := os.OpenFile(logfile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	lg := log.New(file, "", log.Ldate|log.Ltime)

	m := &Monitor{
		logger:             lg,
		notificationClient: notificationClient,
		urls:               urls,
		pollingInterval:    pollingInterval,
		prevStatus:         make(map[string]status),
		prevStatusMutex:    &sync.Mutex{},
	}
	return m, err

}

// Start begins the monitoring.
func (m *Monitor) Start() {
	for _, url := range m.urls {
		ticker := time.NewTicker(m.pollingInterval)
		go func(url string) {
			for {
				m.respondToStatus(url)
				<-ticker.C
			}
		}(url)
	}
}

// log writes a new entry to the Monitor's log.
func (m *Monitor) log(message string) {
	m.logger.Println(message)
}

// respondToStatus sends notifications over the monitors channel considering the state of the current and the previous status.
func (m *Monitor) respondToStatus(url string) {
	st := getStatus(url)
	message := st.String()
	prevStatus, ok := m.prevStatus[url]
	prevStatusFailed := ok && !prevStatus.success
	prevStatusSucceeded := ok && prevStatus.success
	prevStatusNotificationFailed := ok && prevStatus.notificationFailed

	if st.success && (prevStatusFailed || prevStatusNotificationFailed) {
		err := m.notificationClient.SendNotification(message)
		if err != nil {
			st.notificationFailed = true
			m.log(fmt.Sprintf("Sending SUCCESS notification for %s failed.\n\tError: %v", st.url, err))
		}

	} else if !st.success && (!ok || prevStatusSucceeded || prevStatusNotificationFailed) {
		m.log(message)
		err := m.notificationClient.SendNotification(message)
		if err != nil {
			st.notificationFailed = true
			m.log(fmt.Sprintf("Sending FAILURE notification for %s failed.\n\tError: %v", st.url, err))
		}
	}
	m.prevStatusMutex.Lock()
	m.prevStatus[url] = st
	m.prevStatusMutex.Unlock()
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
