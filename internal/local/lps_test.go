package local

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	amt2 "rpc/internal/amt"
	"rpc/internal/flags"
	"rpc/pkg/utils"
	"testing"
	"time"
)

// Mock the AMT Hardware
type MockAMT struct{}

func (c MockAMT) Initialize() (int, error) {
	return utils.Success, nil
}
func (c MockAMT) GetVersionDataFromME(key string, amtTimeout time.Duration) (string, error) {
	return "Version", nil
}
func (c MockAMT) GetUUID() (string, error) { return "123-456-789", nil }

var mockControlMode = 0
var mockControlModeErr error = nil

func (c MockAMT) GetControlMode() (int, error)    { return mockControlMode, mockControlModeErr }
func (c MockAMT) GetOSDNSSuffix() (string, error) { return "", nil }
func (c MockAMT) GetDNSSuffix() (string, error)   { return "", nil }
func (c MockAMT) GetCertificateHashes() ([]amt2.CertHashEntry, error) {
	return []amt2.CertHashEntry{}, nil
}
func (c MockAMT) GetRemoteAccessConnectionStatus() (amt2.RemoteAccessStatus, error) {
	return amt2.RemoteAccessStatus{}, nil
}
func (c MockAMT) GetLANInterfaceSettings(useWireless bool) (amt2.InterfaceSettings, error) {
	return amt2.InterfaceSettings{}, nil
}

var mockLocalSystemAccountErr error = nil

func (c MockAMT) GetLocalSystemAccount() (amt2.LocalSystemAccount, error) {
	return amt2.LocalSystemAccount{Username: "Username", Password: "Password"}, mockLocalSystemAccountErr
}

var mockUnprovisionCode = 0
var mockUnprovisionErr error = nil

func (c MockAMT) Unprovision() (int, error) { return mockUnprovisionCode, mockUnprovisionErr }

func setupService(f *flags.Flags) ProvisioningService {
	service := NewProvisioningService(f)
	service.amtCommand = MockAMT{}
	return service
}

func setupWithWsmanClient(f *flags.Flags, handler http.Handler) ProvisioningService {
	server := httptest.NewServer(handler)
	service := setupService(f)
	service.serverURL = server.URL
	service.setupWsmanClient("admin", "password")
	return service
}

func TestExecute(t *testing.T) {
	f := &flags.Flags{}

	t.Run("execute CommandAMTInfo should succeed", func(t *testing.T) {
		f.Command = utils.CommandAMTInfo
		resultCode := ExecuteCommand(f)
		assert.Equal(t, utils.Success, resultCode)
	})

	t.Run("execute CommandVersion should succeed", func(t *testing.T) {
		f.Command = utils.CommandVersion
		resultCode := ExecuteCommand(f)
		assert.Equal(t, utils.Success, resultCode)
	})

	t.Run("execute CommandMaintenance with no SubCommand fails", func(t *testing.T) {
		f.Command = utils.CommandMaintenance
		resultCode := ExecuteCommand(f)
		assert.Equal(t, utils.InvalidParameters, resultCode)
	})

}
