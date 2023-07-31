package local

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
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

	t.Run("returns ActivationFailed when GeneralSettings.Get() fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondServerError(w)
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns ActivationFailed when xml.unmarshal GeneralSettings.Get() fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondBadXML(t, w)
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.ActivationFailed, resultCode)
	})

	t.Run("returns ActivationFailed when HostBasedSetupService.Setup fails", func(t *testing.T) {
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

	t.Run("returns ActivationFailed when HostBasedSetupService.Setup bad xml", func(t *testing.T) {
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

	t.Run("returns ActivationFailed when HostBasedSetupService.Setup return value is not success", func(t *testing.T) {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if calls == 1 {
				respondGeneralSettings(t, w)
			} else if calls == 2 {
				respondHostBasedSetup(t, w, 1)
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
				respondHostBasedSetup(t, w, 0)
			}
		})
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.ActivateCCM()
		assert.Equal(t, utils.Success, resultCode)
	})
}

func TestActivateACM(t *testing.T) {
	f := &flags.Flags{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondServerError(w)
	})
	lps := setupWithWsmanClient(f, handler)
	resultCode := lps.ActivateACM()
	assert.Equal(t, utils.ActivationFailed, resultCode)
}
