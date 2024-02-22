package local

import (
	"rpc/internal/flags"
	"rpc/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	f := &flags.Flags{}

	t.Run("expect error for unhandled Subcommand", func(t *testing.T) {
		lps := setupService(&flags.Flags{})
		err := lps.Configure()
		assert.Equal(t, utils.IncorrectCommandLineParameters, err)
	})
	t.Run("expect error for SubCommandAddWifiSettings", func(t *testing.T) {
		f.SubCommand = utils.SubCommandAddWifiSettings
		errEnableWiFi = errTestError
		lps := setupService(f)
		err := lps.Configure()
		assert.Error(t, err)
		errEnableWiFi = nil
	})
	t.Run("expect success for SubCommandAddWifiSettings", func(t *testing.T) {
		f.SubCommand = utils.SubCommandAddWifiSettings
		lps := setupService(f)
		err := lps.Configure()
		assert.NoError(t, err)
	})
	t.Run("expect error for SubCommandEnableWifiPort", func(t *testing.T) {
		f.SubCommand = utils.SubCommandEnableWifiPort
		errEnableWiFi = errTestError
		lps := setupService(f)
		err := lps.Configure()
		assert.Error(t, err)
		errEnableWiFi = nil
	})
	t.Run("expect success for SubCommandEnableWifiPort", func(t *testing.T) {
		f.SubCommand = utils.SubCommandEnableWifiPort
		lps := setupService(f)
		err := lps.Configure()
		assert.NoError(t, err)
	})
	t.Run("expect error for SetMebx", func(t *testing.T) {
		f.SubCommand = utils.SubCommandSetMEBx
		lps := setupService(f)
		mockSetupAndConfigurationErr = errTestError
		err := lps.Configure()
		assert.Error(t, err)
		mockSetupAndConfigurationErr = nil
	})
	t.Run("expect success for SetMebx", func(t *testing.T) {
		f.SubCommand = utils.SubCommandSetMEBx
		lps := setupService(f)
		mockControlMode = 2
		err := lps.Configure()
		assert.NoError(t, err)
	})
	t.Run("expect error for Syncclock", func(t *testing.T) {
		f.SubCommand = utils.SubCommandSyncClock
		lps := setupService(f)
		mockGetLowAccuracyTimeSynchErr = errTestError
		err := lps.Configure()
		assert.Error(t, err)
		mockGetLowAccuracyTimeSynchErr = nil
	})
	t.Run("expect success for Syncclock", func(t *testing.T) {
		f.SubCommand = utils.SubCommandSyncClock
		lps := setupService(f)
		mockControlMode = 2
		err := lps.Configure()
		assert.NoError(t, err)
	})
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid http URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "Valid https URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "Missing scheme",
			url:     "://example.com",
			wantErr: true,
		},
		{
			name:    "Missing host",
			url:     "http://",
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			url:     "ht!tp://[::1]/",
			wantErr: true,
		},
		{
			name:    "Relative URL without scheme and host",
			url:     "/path/to/resource",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &flags.Flags{}
			mockAMT := new(MockAMT)
			mockWsman := new(MockWSMAN)
			service := NewProvisioningService(f)
			service.amtCommand = mockAMT
			service.interfacedWsmanMessage = mockWsman
			err := service.ValidateURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err, "ValidateURL() should return an error")
			} else {
				assert.NoError(t, err, "ValidateURL() should not return an error")
			}
		})
	}
}
