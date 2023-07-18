package flags

import (
	"fmt"
	"rpc/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func (f *Flags) handleDeactivateCommand() int {
	f.amtDeactivateCommand.BoolVar(&f.Local, "local", false, "Execute command to AMT directly without cloud interaction")
	if len(f.commandLineArgs) == 2 {
		f.amtDeactivateCommand.PrintDefaults()
		return utils.IncorrectCommandLineParameters
	}
	if err := f.amtDeactivateCommand.Parse(f.commandLineArgs[2:]); err != nil {
		return utils.IncorrectCommandLineParameters
	}
	if f.Local && f.URL != "" {
		fmt.Println("provide either a 'url' or a 'local', but not both")
		return utils.InvalidParameters
	}

	if !f.Local {
		if errCode := f.handleRemoteDeactivation(); errCode != utils.Success {
			return errCode
		}
	} else {
		if errCode := f.handleLocalDeactivation(); errCode != utils.Success {
			return errCode
		}
		log.Info("Status: Device deactivated.")
	}
	return utils.Success
}

func (f *Flags) handleRemoteDeactivation() int {
	if f.URL == "" {
		fmt.Println("-u flag is required and cannot be empty")
		f.amtDeactivateCommand.Usage()
		return utils.MissingOrIncorrectURL
	}
	if f.Password == "" {
		if _, errCode := f.readPasswordFromUser(); errCode != 0 {
			return utils.MissingOrIncorrectPassword
		}
	}
	// TODO: don't put the command together here cause prompt for PW at one place is better
	f.Command = "deactivate --password " + f.Password
	if f.Force {
		f.Command = f.Command + " -f"
	}
	return utils.Success
}
func (f *Flags) handleLocalDeactivation() int {
	controlMode, err := f.amtCommand.GetControlMode()
	if err != nil {
		fmt.Println("Device local deactivation failed.")
		return utils.DeactivationFailed
	}
	if controlMode != 1 {
		fmt.Println("Device is in " + utils.InterpretControlMode(controlMode) + ". Local deactivation is only supported for client control mode.")
		return utils.UnableToDeactivate
	}
	status, err := f.amtCommand.Unprovision()
	if err != nil || status != 0 {
		fmt.Println("Device local deactivation failed.")
		return utils.DeactivationFailed
	}
	return utils.Success
}
