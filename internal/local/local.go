package local

import (
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
	if flags.Command == utils.CommandAMTInfo {
		resultCode = localConfiguration.DisplayAMTInfo()
	}
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
