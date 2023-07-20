package local

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"rpc/internal/flags"
	"rpc/pkg/utils"
)

func ExecuteCommand(flags *flags.Flags) (int, error) {
	resultCode := utils.Success
	localConfiguration := NewLocalConfiguration(flags)
	switch flags.Command {
	case utils.CommandAMTInfo:
		resultCode = localConfiguration.DisplayAMTInfo()
		break
	case utils.CommandVersion:
		resultCode = localConfiguration.DisplayVersion()
		break
	}
	// activate - ClientControlMode
	// activate - AdminControlMode
	//   for both, need to get LocalSystemAccount username password for wsman connection
	// deactivate
	// maintenance - addwifisettings
	// local.
	// figure out if need to use wsman calls or can do PTHI stuff
	//var password string = config.Password
	//var username string = "admin"
	//if flags.UseCCM || flags.UseACM {
	//	rpsPayload := rps.NewPayload()
	//	lsa, err := rpsPayload.AMT.GetLocalSystemAccount()
	//	if err != nil {
	//		log.Error(err)
	//		// TODO: not a valid return code
	//		return -1, err
	//	}
	//	password = lsa.Password
	//	username = lsa.Username
	//}
	//client := wsman.NewClient("http://"+utils.LMSAddress+":"+utils.LMSPort+"/wsman", username, password, true)
	//localConfiguration := NewLocalConfiguration(flags)
	//if flags.UseCCM {
	//	localConfiguration.ActivateCCM()
	//} else {
	//	localConfiguration.Configure8021xWiFi()
	//}

	return resultCode, nil
}

func (local *LocalConfiguration) setupWsmanClient(username string, password string) bool {
	//var password string = local.config.Password
	//var username string = "admin"
	//// this system username password is only for local activation
	//// otherwise the password should be the AMT password from flags.
	//if local.flags.UseCCM || local.flags.UseACM {
	//	amtCommand := amt.NewAMTCommand()
	//	lsa, err := amtCommand.GetLocalSystemAccount()
	//	if err != nil {
	//		log.Error(err)
	//		return false
	//	}
	//	password = lsa.Password
	//	username = lsa.Username
	//}
	local.client = wsman.NewClient("http://"+utils.LMSAddress+":"+utils.LMSPort+"/wsman", username, password, true)
	return true
}
