package flags

import (
	"errors"
	"rpc/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleDeactivateCommandNoFlags(t *testing.T) {
	args := []string{"./rpc", "deactivate"}
	flags := NewFlags(args)
	flags.amtCommand.PTHI = MockPTHICommands{}
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.IncorrectCommandLineParameters)
}
func TestHandleDeactivateInvalidFlag(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-x"}

	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.IncorrectCommandLineParameters)
}

func TestHandleDeactivateCommandNoPasswordPrompt(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-u", "wss://localhost"}
	expected := "deactivate --password password"
	defer userInput(t, "password")()
	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, expected, flags.Command)
}
func TestHandleDeactivateCommandNoPasswordPromptEmpy(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-u", "wss://localhost"}
	defer userInput(t, "")()
	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.MissingOrIncorrectPassword)
}
func TestHandleDeactivateCommandNoURL(t *testing.T) {
	args := []string{"./rpc", "deactivate", "--password", "password"}

	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.MissingOrIncorrectURL)
}
func TestHandleDeactivateCommand(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-u", "wss://localhost", "--password", "password"}
	expected := "deactivate --password password"
	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, expected, flags.Command)
}

func TestHandleDeactivateCommandWithURLAndLocal(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-u", "wss://localhost", "--password", "password", "-local"}
	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.InvalidParameters)
	assert.Equal(t, "wss://localhost", flags.URL)
}
func TestHandleDeactivateCommandWithForce(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-u", "wss://localhost", "--password", "password", "-f"}
	expected := "deactivate --password password -f"
	flags := NewFlags(args)
	success := flags.handleDeactivateCommand()
	assert.EqualValues(t, success, utils.Success)
	assert.Equal(t, "wss://localhost", flags.URL)
	assert.Equal(t, expected, flags.Command)
}

func TestHandleLocalDeactivationWithACM(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-local"}
	flags := NewFlags(args)
	flags.amtCommand.PTHI = MockPTHICommands{}
	mode = 2
	result = 0
	errCode := flags.handleLocalDeactivation()
	assert.Equal(t, errCode, utils.UnableToDeactivate)
	mode = 0
}

func TestHandleLocalDeactivation(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-local"}
	flags := NewFlags(args)
	flags.amtCommand.PTHI = MockPTHICommands{}
	mode = 1
	result = 0
	errCode := flags.handleLocalDeactivation()
	assert.Equal(t, errCode, utils.Success)
	mode = 0
}

func TestHandleDeactivateCommandWithGetControlModeError(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-local"}
	flags := NewFlags(args)
	flags.amtCommand.PTHI = MockPTHICommands{}
	mode = 1
	result = 0
	controlModeErr = errors.New("Failed to get control mode")
	errCode := flags.handleDeactivateCommand()
	assert.Equal(t, errCode, utils.DeactivationFailed)
	mode = 0
}

func TestHandleLocalDeactivationwithUnprovisionError(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-local"}
	flags := NewFlags(args)
	flags.amtCommand.PTHI = MockPTHICommands{}
	mode = 1
	result = -1
	controlModeErr = nil
	errCode := flags.handleLocalDeactivation()
	assert.Equal(t, errCode, utils.DeactivationFailed)
	result = 0
	mode = 0
}

func TestHandleDeactivationWithLocal(t *testing.T) {
	args := []string{"./rpc", "deactivate", "-local"}
	flags := NewFlags(args)
	flags.amtCommand.PTHI = MockPTHICommands{}
	mode = 1
	result = 0
	controlModeErr = nil
	errCode := flags.handleDeactivateCommand()
	assert.Equal(t, errCode, utils.Success)
	mode = 0
}

func TestParseFlagsDeactivate(t *testing.T) {
	args := []string{"./rpc", "deactivate"}
	flags := NewFlags(args)
	result := flags.ParseFlags()
	assert.EqualValues(t, result, utils.IncorrectCommandLineParameters)
	assert.Equal(t, utils.CommandDeactivate, flags.Command)
}
