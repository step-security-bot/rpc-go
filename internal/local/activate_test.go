package local

import (
	"encoding/xml"
	"errors"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/ips/hostbasedsetup"
	"github.com/stretchr/testify/assert"
	"net/http"
	amt2 "rpc/internal/amt"
	"rpc/internal/flags"
	"rpc/pkg/utils"
	"testing"
)

func TestActivation(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondServerError(w)
	})
	lps := setupWithTestServer(&flags.Flags{}, handler)
	lps.flags.Command = utils.CommandActivate
	lps.flags.LocalConfig.Password = "P@ssw0rd"

	t.Run("returns AMTConnectionFailed when GetControlMode fails", func(t *testing.T) {
		mockControlModeErr = errors.New("yep it failed")
		resultCode := lps.Activate()
		assert.Equal(t, utils.AMTConnectionFailed, resultCode)
		mockControlModeErr = nil
	})

	t.Run("returns UnableToActivate when already activated", func(t *testing.T) {
		mockControlMode = 1
		resultCode := lps.Activate()
		assert.Equal(t, utils.UnableToActivate, resultCode)
		mockControlMode = 0
	})

	t.Run("returns AMTConnectionFailed when GetLocalSystemAccount fails", func(t *testing.T) {
		mockLocalSystemAccountErr = errors.New("yep it failed")
		resultCode := lps.Activate()
		assert.Equal(t, utils.AMTConnectionFailed, resultCode)
		mockLocalSystemAccountErr = nil
	})

	t.Run("returns ActivationFailed when UseACM and responses are not mocked", func(t *testing.T) {
		lps.flags.UseACM = true
		resultCode := lps.Activate()
		assert.Equal(t, utils.ActivationFailed, resultCode)
		lps.flags.UseACM = false
	})

	t.Run("returns ActivationFailed when UseCCM and responses are not mocked", func(t *testing.T) {
		lps.flags.UseCCM = true
		resultCode := lps.Activate()
		assert.Equal(t, utils.ActivationFailed, resultCode)
		lps.flags.UseCCM = false
	})
}

func TestActivateCCM(t *testing.T) {
	f := &flags.Flags{}

	t.Run("returns ActivationFailed on GeneralSettings.Get() server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondServerError(w)
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns ActivationFailed on GeneralSettings.Get() xml.unmarshal error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondBadXML(t, w)
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns ActivationFailed on HostBasedSetupService.Setup server error", func(t *testing.T) {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if calls == 1 {
				respondGeneralSettings(t, w)
			} else if calls == 2 {
				respondServerError(w)
			}
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns ActivationFailed on HostBasedSetupService.Setup xml.unmarshal error", func(t *testing.T) {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if calls == 1 {
				respondGeneralSettings(t, w)
			} else if calls == 2 {
				respondBadXML(t, w)
			}
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns ActivationFailed on HostBasedSetupService.Setup ReturnValue is not success (0)", func(t *testing.T) {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if calls == 1 {
				respondGeneralSettings(t, w)
			} else if calls == 2 {
				mockHostBasedSetupResponse.Body.Setup_OUTPUT.ReturnValue = 1
				respondHostBasedSetup(t, w)
				mockHostBasedSetupResponse.Body.Setup_OUTPUT.ReturnValue = 0
			}
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns Success on happy path", func(t *testing.T) {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if calls == 1 {
				respondGeneralSettings(t, w)
			} else if calls == 2 {
				respondHostBasedSetup(t, w)
			}
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.Success, resultCode)
	})
}

func TestGetHostBasedSetupService(t *testing.T) {
	f := &flags.Flags{}

	t.Run("returns error on server error response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondServerError(w)
		})
		lps := setupWithWsmanClient(f, handler)
		_, err := lps.GetHostBasedSetupService()
		assert.NotNil(t, err)
	})

	t.Run("returns error on xml.unmarshal error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondBadXML(t, w)
		})
		lps := setupWithWsmanClient(f, handler)
		_, err := lps.GetHostBasedSetupService()
		assert.NotNil(t, err)
	})

	t.Run("returns valid response on happy path", func(t *testing.T) {
		expected := "test_name"
		mockHostBasedSetupResponse.Body.IPS_HostBasedSetupService.SystemName = expected
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondHostBasedSetup(t, w)
		})
		lps := setupWithWsmanClient(f, handler)
		rsp, err := lps.GetHostBasedSetupService()
		assert.Nil(t, err)
		assert.Equal(t, expected, rsp.Body.IPS_HostBasedSetupService.SystemName)
		mockHostBasedSetupResponse.Body.IPS_HostBasedSetupService.SystemName = ""
	})
}

func TestGetGeneralSettings(t *testing.T) {
	f := &flags.Flags{}

	t.Run("returns error on server error response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondServerError(w)
		})
		lps := setupWithWsmanClient(f, handler)
		_, err := lps.GetGeneralSettings()
		assert.NotNil(t, err)
	})

	t.Run("returns error on xml.unmarshal error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondBadXML(t, w)
		})
		lps := setupWithWsmanClient(f, handler)
		_, err := lps.GetGeneralSettings()
		assert.NotNil(t, err)
	})

	t.Run("returns valid response on happy path", func(t *testing.T) {
		expected := "test_name"
		mockGenerlSettingsResponse.Body.AMTGeneralSettings.HostName = expected
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondGeneralSettings(t, w)
		})
		lps := setupWithWsmanClient(f, handler)
		rsp, err := lps.GetGeneralSettings()
		assert.Nil(t, err)
		assert.Equal(t, expected, rsp.Body.AMTGeneralSettings.HostName)
		mockGenerlSettingsResponse.Body.AMTGeneralSettings.HostName = ""
	})
}

func TestActivateACM(t *testing.T) {
	f := &flags.Flags{}
	f.LocalConfig.ACMSettings.AMTPassword = "P@ssw0rd"
	f.LocalConfig.ACMSettings.MEBxPassword = "P@ssw0rd"
	f.LocalConfig.ACMSettings.ProvisioningCert = "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"
	f.LocalConfig.ACMSettings.ProvisioningCertPwd = "P@ssw0rd"
	mockCertHashes = []amt2.CertHashEntry{
		{
			Hash:      "cb3ccbb76031e5e0138f8dd39a23f9de47ffc35e43c1144cea27d46a5ab1cb5f",
			Name:      "",
			Algorithm: "",
			IsActive:  true,
			IsDefault: true,
		},
	}
	calls := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		var err error
		calls++
		if calls == 1 {
			rsp := general.Response{}
			xmlString, err := xml.Marshal(rsp)
			assert.Nil(t, err)
			_, err = w.Write(xmlString)
		} else if calls == 2 {
			rsp := hostbasedsetup.Response{}
			xmlString, err := xml.Marshal(rsp)
			assert.Nil(t, err)
			_, err = w.Write(xmlString)
		}
		assert.Nil(t, err)
	})
	lps := setupWithWsmanClient(f, handler)
	resultCode := lps.ActivateACM()
	assert.Equal(t, utils.ActivationFailed, resultCode)
}
