package keyringservice_test

// WE Need to figure out how to run this test in the github runner

// import (
// 	"testing"
// 	"time"

// 	"github.com/motain/of-catalog/internal/services/keyringservice"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/zalando/go-keyring"
// )

// type mockKeyring struct {
// 	data map[string]string
// }

// func (m *mockKeyring) Set(service, user, secret string) error {
// 	m.data[service+user] = secret
// 	return nil
// }

// func (m *mockKeyring) Get(service, user string) (string, error) {
// 	secret, ok := m.data[service+user]
// 	if !ok {
// 		return "", keyring.ErrNotFound
// 	}
// 	return secret, nil
// }

// func (m *mockKeyring) Delete(service, user string) error {
// 	delete(m.data, service+user)
// 	return nil
// }

// func TestKeyringService_Set(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		service string
// 		user    string
// 		secret  string
// 		wantErr bool
// 	}{
// 		{"Set secret", "service1", "user1", "secret1", false},
// 		// Need to think about how to test this, I would need to refactor the code to make it testable
// 		// {"Set secret timeout", "service2", "user2", "secret2", true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ks := keyringservice.NewKeyringService()
// 			if tt.wantErr {
// 				time.Sleep(4 * time.Second)
// 			}
// 			err := ks.Set(tt.service, tt.user, tt.secret)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestKeyringService_Get(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		service   string
// 		user      string
// 		want      string
// 		wantErr   bool
// 		setupFunc func(*mockKeyring)
// 	}{
// 		{
// 			"Get secret",
// 			"service1",
// 			"user1",
// 			"secret1",
// 			false,
// 			func(m *mockKeyring) {
// 				m.Set("service1", "user1", "secret1")
// 			},
// 		},
// 		// {
// 		// 	"Get secret not found",
// 		// 	"service2",
// 		// 	"user2",
// 		// 	"",
// 		// 	true,
// 		// 	func(m *mockKeyring) {},
// 		// },
// 		{
// 			"Get secret timeout",
// 			"service3",
// 			"user3",
// 			"",
// 			true,
// 			func(m *mockKeyring) {
// 				time.Sleep(4 * time.Second)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mock := &mockKeyring{data: make(map[string]string)}
// 			ks := keyringservice.NewKeyringService()
// 			tt.setupFunc(mock)
// 			got, err := ks.Get(tt.service, tt.user)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want, got)
// 			}
// 		})
// 	}
// }

// func TestKeyringService_Delete(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		service   string
// 		user      string
// 		wantErr   bool
// 		setupFunc func(*mockKeyring)
// 	}{
// 		{
// 			"Delete secret",
// 			"service1",
// 			"user1",
// 			false,
// 			func(m *mockKeyring) {
// 				m.Set("service1", "user1", "secret1")
// 			},
// 		},
// 		{
// 			"Delete secret not found",
// 			"service2",
// 			"user2",
// 			true,
// 			func(m *mockKeyring) {},
// 		},
// 		{
// 			"Delete secret timeout",
// 			"service3",
// 			"user3",
// 			true,
// 			func(m *mockKeyring) {
// 				time.Sleep(4 * time.Second)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mock := &mockKeyring{data: make(map[string]string)}
// 			ks := keyringservice.NewKeyringService()
// 			tt.setupFunc(mock)
// 			err := ks.Delete(tt.service, tt.user)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
