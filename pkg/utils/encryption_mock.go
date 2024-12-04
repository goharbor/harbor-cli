package utils

import "fmt"

type MockKeyring struct {
	store map[string]map[string]string
}

func NewMockKeyring() *MockKeyring {
	return &MockKeyring{store: make(map[string]map[string]string)}
}

func (m *MockKeyring) Set(service, user, password string) error {
	if m.store[service] == nil {
		m.store[service] = make(map[string]string)
	}
	m.store[service][user] = password
	return nil
}

func (m *MockKeyring) Get(service, user string) (string, error) {
	if val, ok := m.store[service][user]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

func (m *MockKeyring) Delete(service, user string) error {
	if _, ok := m.store[service][user]; ok {
		delete(m.store[service], user)
		return nil
	}
	return fmt.Errorf("key not found")
}
