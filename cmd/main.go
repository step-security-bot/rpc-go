/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package main

import "C"
import (
	"os"
	"rpc"
	"rpc/internal/amt"
	"rpc/internal/client"
	"rpc/internal/rps"
	"rpc/pkg/utils"
	"strings"

	log "github.com/sirupsen/logrus"
)

func checkAccess() (int, error) {
	amtCommand := amt.NewAMTCommand()
	result, err := amtCommand.Initialize()
	if !result || err != nil {
		return utils.ReturnCode_ACCESS, err
	}
	return utils.ReturnCode_SUCCESS, nil
}

func runRPC(args []string) (int, error) {
	// process cli flags/env vars
	flags, keepGoing := handleFlags(args)
	if keepGoing == false {
		return utils.ReturnCode_SUCCESS, nil
	}

	startMessage, err := rps.PrepareInitialMessage(flags)
	if err != nil {
		return utils.ReturnCode_GENERAL_ERR, err
	}

	executor, err := client.NewExecutor(*flags)
	if err != nil {
		return utils.ReturnCode_GENERAL_ERR, err
	}

	executor.MakeItSo(startMessage)
	return utils.ReturnCode_SUCCESS, nil
}

func handleFlags(args []string) (*rpc.Flags, bool) {
	//process flags
	flags := rpc.NewFlags(args)
	_, result := flags.ParseFlags()
	if !result {
		return nil, false
	}
	if flags.Verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		lvl, err := log.ParseLevel(flags.LogLevel)
		if err != nil {
			log.Warn(err)
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(lvl)
		}

	}
	if flags.SyncClock {
		log.Info("Syncing the clock")
	}
	if flags.JsonOutput {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
	return flags, true
}

const AccessErrMsg = "Failed to execute due to access issues. " +
	"Please ensure that Intel ME is present, " +
	"the MEI driver is installed, " +
	"and the runtime has administrator or root privileges."

func main() {
	status, err := checkAccess()
	if status != utils.ReturnCode_SUCCESS {
		if err != nil {
			log.Error(err.Error())
		}
		log.Error(AccessErrMsg)
		return
	}
	status, err = runRPC(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Completed with status code: ", status)
}

//export rpcExec
func rpcExec(Input *C.char, Output **C.char) int {
	status, err := checkAccess()
	if status != utils.ReturnCode_SUCCESS {
		if err != nil {
			log.Error(err.Error())
		}
		*Output = C.CString(AccessErrMsg)
		return status
	}

	//create argument array from input string
	inputString := C.GoString(Input)
	args := strings.Fields(inputString)
	args = append([]string{"rpc"}, args...)
	status, err = runRPC(args)
	if status != utils.ReturnCode_SUCCESS {
		if err != nil {
			log.Error(err.Error())
		}
		*Output = C.CString("rpcExec failed: " + inputString)
	}
	return status
}
