package local_test

import (
	"rpc/internal/config"
	"rpc/internal/flags"
	"rpc/internal/local"
	"testing"
)

func TestActivationCCM(t *testing.T) {
	f := &flags.Flags{}
	f.LocalConfig = config.Config{
		Password: "P@ssw0rd",
	}
	localConfiguration := local.NewLocalConfiguration(f)
	localConfiguration.ActivateCCM()

	//localConfig := config.Config{
	//	Password: "P@ssw0rd",
	//}
	//// gets required information for us
	//rpsPayload := rps.NewPayload()
	//lsa, err := rpsPayload.AMT.GetLocalSystemAccount()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// payload.Username = lsa.Username
	////localConfig.Password = lsa.Password
	//
	//
	//client := wsman.NewClient("http://"+utils.LMSAddress+":"+utils.LMSPort+"/wsman", lsa.Username, lsa.Password, true)
	//localConnection := local.NewLocalConfiguration(localConfig, client)
	//localConnection.ActivateCCM()
}
