package flags

import (
	"rpc/pkg/utils"
)

func (f *Flags) handleVersionCommand() int {
	if err := f.versionCommand.Parse(f.commandLineArgs[2:]); err != nil {
		return utils.IncorrectCommandLineParameters
	}
	// runs locally
	f.Local = true
	return utils.Success
}
