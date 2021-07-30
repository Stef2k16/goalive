// Package monitor implements functionality to poll endpoints and send notifications about the retrieved status.
package monitor

import (
	"fmt"
	"github.com/Stef2k16/goalive/config"
	"github.com/Stef2k16/goalive/notification"
	"log"
	"os"
	"sync"
	"time"
)

// Monitor provides the functionality to monitor http endpoints.
type Monitor struct {
	logger             *log.Logger
	NotificationClient notification.Client
	urls               []string
	pollingInterval    time.Duration
	prevStatus         map[string]status
	prevStatusMutex    *sync.Mutex
}

// New sets up a new monitor to check services.
func New(conf config.Config) (*Monitor, error) {
	file, err := os.OpenFile(conf.LogFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	lg := log.New(file, "", log.Ldate|log.Ltime)

	c, err := notification.GetClient(conf.Notification)
	if err != nil {
		return nil, err
	}

	m := &Monitor{
		logger:             lg,
		NotificationClient: c,
		urls:               conf.URL,
		pollingInterval:    time.Duration(conf.PollingInterval) * time.Second,
		prevStatus:         make(map[string]status),
		prevStatusMutex:    &sync.Mutex{},
	}
	c.AddStatusHandler(m.statusSummary)
	if err := c.Start(); err != nil {
		return nil, err
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

// StatusSummary creates a summary of the currently cached status for all monitored urls.
func (m *Monitor) statusSummary() string {
	var summary string
	m.prevStatusMutex.Lock()
	for _, st := range m.prevStatus {
		part1 := fmt.Sprintf("%s polled at %s:\n", st.url, st.timestamp.Format(time.RFC1123))
		readableStatus := "FAILED"
		if st.success {
			readableStatus = "SUCCEEDED"
		}
		part2 := fmt.Sprintf("\tRequest %s with Status %d\n\tBody: %s\n\n", readableStatus, st.code, st.body)
		summary += part1 + part2
	}
	m.prevStatusMutex.Unlock()
	return summary
}

// log writes a new entry to the Monitor's log.
func (m *Monitor) log(message string) {
	m.logger.Println(message)
}

// respondToStatus sends notifications over the monitors channel considering the state of the current and the previous status.
func (m *Monitor) respondToStatus(url string) {
	st := getStatus(url)
	message := st.String()
	m.prevStatusMutex.Lock()
	prevStatus, ok := m.prevStatus[url]
	m.prevStatusMutex.Unlock()
	prevStatusFailed := ok && !prevStatus.success
	prevStatusSucceeded := ok && prevStatus.success
	prevStatusNotificationFailed := ok && prevStatus.notificationFailed

	if st.success && (prevStatusFailed || prevStatusNotificationFailed) {
		err := m.NotificationClient.SendNotification(message)
		if err != nil {
			st.notificationFailed = true
			m.log(fmt.Sprintf("Sending SUCCESS notification for %s failed.\n\tError: %v", st.url, err))
		}

	} else if !st.success && (!ok || prevStatusSucceeded || prevStatusNotificationFailed) {
		m.log(message)
		err := m.NotificationClient.SendNotification(message)
		if err != nil {
			st.notificationFailed = true
			m.log(fmt.Sprintf("Sending FAILURE notification for %s failed.\n\tError: %v", st.url, err))
		}
	}
	m.prevStatusMutex.Lock()
	m.prevStatus[url] = st
	m.prevStatusMutex.Unlock()
}
