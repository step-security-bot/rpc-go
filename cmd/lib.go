package main

import "C"
import (
	log "github.com/sirupsen/logrus"
	"rpc/internal/amt"
	"rpc/pkg/utils"
	"strings"
)

//export checkAccess
func checkAccess() int {
	amtCommand := amt.NewAMTCommand()
	result, err := amtCommand.Initialize()
	if !result || err != nil {
		log.Error("Unable to launch application. " +
			"Please ensure that Intel ME is present, " +
			"the MEI driver is installed, " +
			"and that this application is run with administrator or root privileges.")
		return utils.ReturnCode_BASIC_FAIL
	}
	return utils.ReturnCode_SUCCESS
}

//export rpcExec
func rpcExec(Input *C.char, Output **C.char) int {
	status := checkAccess()
	if status != utils.ReturnCode_SUCCESS {
		return status
	}

	//create argument array from input string
	args := strings.Fields(C.GoString(Input))
	args = append([]string{"rpc"}, args...)
	runRPC(args)
	*Output = C.CString("test output")
	return utils.ReturnCode_SUCCESS
}
