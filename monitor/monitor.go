// Package monitor implements functionality to poll endpoints and notify about the retrieved status.
package monitor

import (
	"fmt"
	"galive/notification"
	"log"
	"os"
	"sync"
	"time"
)

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