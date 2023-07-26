package flags

import (
	"os"
	"rpc/pkg/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleActivateCommandNoFlags(t *testing.T) {
	args := []string{"./rpc", "activate"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.IncorrectCommandLineParameters)
}
func TestHandleActivateCommand(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-profile", "profileName", "-password", "Password"}
	flags := NewFlags(args)
	var AMTTimeoutDuration time.Duration = 120000000000
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, "profileName", flags.Profile)
	assert.Equal(t, utils.CommandActivate, flags.Command)
	assert.Equal(t, "Password", flags.Password)
	assert.Equal(t, "localhost", flags.LMSAddress)
	assert.Equal(t, "16992", flags.LMSPort)
	// 2m default
	assert.Equal(t, AMTTimeoutDuration, flags.AMTTimeoutDuration)
	assert.Equal(t, "", flags.FriendlyName)
}

func TestHandleActivateCommandWithTimeOut(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-profile", "profileName", "-password", "Password", "-t", "2s"}
	flags := NewFlags(args)
	var AMTTimeoutDuration time.Duration = 2000000000
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, "profileName", flags.Profile)
	assert.Equal(t, utils.CommandActivate, flags.Command)
	assert.Equal(t, "Password", flags.Password)
	assert.Equal(t, "localhost", flags.LMSAddress)
	assert.Equal(t, "16992", flags.LMSPort)
	assert.Equal(t, AMTTimeoutDuration, flags.AMTTimeoutDuration)
}
func TestHandleActivateCommandWithLMS(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-profile", "profileName", "-lmsaddress", "1.1.1.1", "-lmsport", "99"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, "profileName", flags.Profile)
	assert.Equal(t, utils.CommandActivate, flags.Command)
	assert.Equal(t, "1.1.1.1", flags.LMSAddress)
	assert.Equal(t, "99", flags.LMSPort)
}
func TestHandleActivateCommandWithFriendlyName(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-profile", "profileName", "-name", "friendlyName"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, "profileName", flags.Profile)
	assert.Equal(t, utils.CommandActivate, flags.Command)
	assert.Equal(t, "friendlyName", flags.FriendlyName)
}
func TestHandleActivateCommandWithENV(t *testing.T) {

	if err := os.Setenv("DNS_SUFFIX", "envdnssuffix.com"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("HOSTNAME", "envhostname"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("PROFILE", "envprofile"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("AMT_PASSWORD", "envpassword"); err != nil {
		t.Error(err)
	}

	args := []string{"./rpc", "activate", "-u", "wss://localhost"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, "envprofile", flags.Profile)
	assert.Equal(t, utils.CommandActivate, flags.Command)
	assert.Equal(t, "envpassword", flags.Password)
	os.Clearenv()
}

func TestHandleActivateCommandNoProfile(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.MissingOrIncorrectProfile)
	assert.Equal(t, "wss://localhost", flags.URL)
}

func TestHandleActivateCommandNoProxy(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-p"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.MissingProxyAddressAndPort)
	assert.Equal(t, "wss://localhost", flags.URL)
}

func TestHandleActivateCommandNoHostname(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-h"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.MissingHostname)
	assert.Equal(t, "wss://localhost", flags.URL)
}

func TestHandleActivateCommandNoDNSSuffix(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-d"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.MissingDNSSuffix)
	assert.Equal(t, "wss://localhost", flags.URL)
}

func TestHandleActivateCommandMissingProfile(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-profile"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.MissingOrIncorrectProfile)
	assert.Equal(t, "wss://localhost", flags.URL)
}

func TestHandleActivateCommandBothURLandLocal(t *testing.T) {
	args := []string{"./rpc", "activate", "-u", "wss://localhost", "-local"}
	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.InvalidParameters)
}

func TestHandleActivateCommandLocal(t *testing.T) {

	tests := map[string]struct {
		cmdLine    string
		wantResult int
	}{
		"should fail with both URL and local": {
			cmdLine:    "./rpc activate -local -u wss://localhost",
			wantResult: utils.InvalidParameters,
		},
		"should fail without acm or ccm specified": {
			cmdLine:    "./rpc activate -local",
			wantResult: utils.InvalidParameters,
		},
		"should fail if both acm and ccm specified": {
			cmdLine:    "./rpc activate -local -acm -ccm",
			wantResult: utils.InvalidParameters,
		},
		"should fail if acm and local config file error": {
			cmdLine:    "./rpc activate -local -acm -config ./nofilehere.txt",
			wantResult: utils.IncorrectCommandLineParameters,
		},
		"should fail if acm and ACM Settings not specified": {
			cmdLine:    "./rpc activate -local -acm",
			wantResult: utils.IncorrectCommandLineParameters,
		},
		"should pass if acm with example config file": {
			cmdLine:    "./rpc activate -local -acm -config ../../config.yaml",
			wantResult: utils.Success,
		},
		"should pass wif acm and ACM Settings specified": {
			cmdLine: "./rpc activate -local -acm " +
				" -amtPassword " + trickyPassword +
				" -mebxPassword " + trickyPassword +
				` -provisioningCert MIIW/gIBAzCCFroGCSqGSIb3DQEHAaCCFqsEghanMIIWozCCBgwGCSqGSIb3DQEHAaCCBf0EggX5MIIF9TCCBfEGCyqGSIb3DQEMCgECoIIE/jCCBPowHAYKKoZIhvc` +
				" -provisioningCertPwd " + trickyPassword,
			wantResult: utils.Success,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args := strings.Fields(tc.cmdLine)
			flags := NewFlags(args)
			gotResult := flags.ParseFlags()
			assert.Equal(t, tc.wantResult, gotResult)
			assert.Equal(t, utils.CommandActivate, flags.Command)
		})
	}

}

// TODO: move to local package
// TODO: refactor no PTHICommands needed
//
//	func TestHandleActivateCommandLocalNoPassword(t *testing.T) {
//		args := []string{"./rpc", "activate", "-local"}
//		flags := NewFlags(args)
//		flags.amtCommand.PTHI = MockPTHICommands{}
//		success := flags.ParseFlags()
//		assert.EqualValues(t, utils.MissingOrIncorrectPassword, success)
//	}
//
//	func TestHandleActivateCommandLocal(t *testing.T) {
//		args := []string{"./rpc", "activate", "-local", "-password", "P@ssw0rd"}
//		flags := NewFlags(args)
//		flags.amtCommand.PTHI = MockPTHICommands{}
//		mode = 0
//		success := flags.ParseFlags()
//		assert.Equal(t, flags.Local, true)
//		assert.EqualValues(t, utils.Success, success)
//	}
//
//	func TestHandleActivateCommandLocalAlreadyActivated(t *testing.T) {
//		args := []string{"./rpc", "activate", "-local", "-password", "P@ssw0rd"}
//		flags := NewFlags(args)
//		flags.amtCommand.PTHI = MockPTHICommands{}
//		mode = 1
//		success := flags.ParseFlags()
//		assert.Equal(t, flags.Local, true)
//		assert.EqualValues(t, utils.UnableToActivate, success)
//		mode = 0
//	}
//
//	func TestHandleActivateCommandLocalControlModeError(t *testing.T) {
//		args := []string{"./rpc", "activate", "-local", "-password", "P@ssw0rd"}
//		flags := NewFlags(args)
//		flags.amtCommand.PTHI = MockPTHICommands{}
//		mode = 0
//		controlModeErr = errors.New("error")
//		success := flags.ParseFlags()
//		assert.Equal(t, flags.Local, true)
//		assert.EqualValues(t, utils.ActivationFailed, success)
//		controlModeErr = nil
//	}
func TestHandleActivateCommandNoURL(t *testing.T) {
	args := []string{"./rpc", "activate", "-profile", "profileName"}

	flags := NewFlags(args)
	success := flags.ParseFlags()
	assert.EqualValues(t, success, utils.MissingOrIncorrectURL)
	assert.Equal(t, "profileName", flags.Profile)
}
