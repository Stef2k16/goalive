package notification

// MockClient provides a mock implementation of an notification provider for testing purpose.
type MockClient struct{}

// NewMockClient returns a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Start opens the mock session.
func (mc *MockClient) Start() error {
	return nil
}

// Stop closes the mock session.
func (mc *MockClient) Stop() error {
	return nil
}

// SendNotification mocks sending of notifications.
func (mc *MockClient) SendNotification(_ string) error {
	return nil
}

// AddStatusHandler mocks adding the status handler.
func (mc *MockClient) AddStatusHandler(_ func() string) {}
