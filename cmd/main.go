/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package main

import (
	"errors"
	"os"
	"rpc"
	"rpc/internal/amt"
	"rpc/internal/client"
	"rpc/internal/rps"
	"rpc/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func runRPC(args []string) {
	// process cli flags/env vars
	flags, err := handleFlags(args)
	if err != nil {
		log.Error(err.Error())
		return
	}

	startMessage, err := rps.PrepareInitialMessage(flags)
	if err != nil {
		log.Error(err.Error())
		return
	}

	rpc, err := client.NewExecutor(*flags)
	if err != nil {
		log.Error(err.Error())
		return
	}
	rpc.MakeItSo(startMessage)
}

func handleFlags(args []string) (*rpc.Flags, error) {
	//process flags
	flags := rpc.NewFlags(args)
	flagParsed, result := flags.ParseFlags()
	if !result {
		return nil, errors.New("Failed parsing flag " + flagParsed)
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
	return flags, nil
}

func main() {
	// ensure we are admin/sudo
	status := checkAdminAccess()
	if status != utils.ReturnCode_SUCCESS {
		return
	}
	runRPC(os.Args)
}

func checkAdminAccess() int {
	amt := amt.NewAMTCommand()
	result, err := amt.Initialize()
	if !result || err != nil {
		log.Error(err) //Print the errors
		return utils.ReturnCode_BASIC_FAIL
	}
	return utils.ReturnCode_SUCCESS
}
