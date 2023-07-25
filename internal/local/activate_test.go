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
	f := &flags.Flags{}
	f.Command = utils.CommandActivate
	f.LocalConfig.Password = "P@ssw0rd"

	t.Run("returns AMTConnectionFailed when GetControlMode fails", func(t *testing.T) {
		lps := setupService(f)
		mockControlModeErr = errors.New("yep it failed")
		resultCode := lps.Activate()
		assert.Equal(t, utils.AMTConnectionFailed, resultCode)
		mockControlModeErr = nil
	})

	t.Run("returns UnableToActivate when already activated", func(t *testing.T) {
		lps := setupService(f)
		mockControlMode = 1
		resultCode := lps.Activate()
		assert.Equal(t, utils.UnableToActivate, resultCode)
		mockControlMode = 0
	})

	t.Run("returns AMTConnectionFailed when GetLocalSystemAccount fails", func(t *testing.T) {
		lps := setupService(f)
		mockLocalSystemAccountErr = errors.New("yep it failed")
		resultCode := lps.Activate()
		assert.Equal(t, utils.AMTConnectionFailed, resultCode)
		mockLocalSystemAccountErr = nil
	})

	t.Run("returns AMTFailed when CCM activate responses are not mocked", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.WriteHeader(http.StatusInternalServerError)
		})
		f.UseCCM = false
		lps := setupWithWsmanClient(f, handler)
		resultCode := lps.Activate()
		assert.Equal(t, utils.ActivationFailed, resultCode)
		assert.Equal(t, true, f.UseCCM)
	})

	
}
