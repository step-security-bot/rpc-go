/*********************************************************************
 * Copyright (c) Intel Corporation 2022
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package client

import (
	"os"
	"os/signal"
	"rpc"
	"rpc/internal/lm"
	"rpc/internal/rps"
	"rpc/pkg/utils"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type Executor struct {
	server          rps.AMTActivationServer
	localManagement lm.LocalMananger
	isLME           bool
	payload         rps.Payload
	data            chan []byte
	errors          chan error
	status          chan bool
}

func NewExecutor(flags rpc.Flags) (Executor, error) {
	// these are closed in the close function for each lm implementation
	lmDataChannel := make(chan []byte)
	lmErrorChannel := make(chan error)

	client := Executor{
		server:          rps.NewAMTActivationServer(&flags),
		localManagement: lm.NewLMSConnection(utils.LMSAddress, utils.LMSPort, lmDataChannel, lmErrorChannel),
		data:            lmDataChannel,
		errors:          lmErrorChannel,
	}

	// TEST CONNECTION TO SEE IF LMS EXISTS
	err := client.localManagement.Connect()

	if err != nil {
		// client.localManagement.Close()
		log.Trace("LMS not running.  Using LME Connection\n")
		client.status = make(chan bool)
		client.localManagement = lm.NewLMEConnection(lmDataChannel, lmErrorChannel, client.status)
		client.isLME = true
		client.localManagement.Initialize()
	} else {
		log.Trace("Using existing LMS\n")
		client.localManagement.Close()
	}

	err = client.server.Connect(flags.SkipCertCheck)
	if err != nil {
		log.Error("error connecting to RPS")
		// TODO: should the connection be closed?
		// client.localManagement.Close()
	}
	return client, err
}

func (e Executor) MakeItSo(messageRequest rps.Message) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	rpsDataChannel := e.server.Listen()

	log.Debug("sending activation request to RPS")
	err := e.server.Send(messageRequest)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer e.localManagement.Close()
	defer close(e.data)
	defer close(e.errors)
	if e.status != nil {
		defer close(e.status)
	}

	for {
		select {
		case dataFromServer := <-rpsDataChannel:
			shallIReturn := e.HandleDataFromRPS(dataFromServer)
			if shallIReturn { //quits the loop -- we're either done or reached a point where we need to stop
				return
			}
		case <-interrupt:
			e.HandleInterrupt()
			return
		}
	}

}

func (e Executor) HandleInterrupt() {
	log.Info("interrupt")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	// err := e.localManagement.Close()
	// if err != nil {
	// 	log.Error("Connection close failed", err)
	// 	return
	// }

	err := e.server.Close()
	if err != nil {
		log.Error("Connection close failed", err)
		return
	}
}

func (e Executor) HandleDataFromRPS(dataFromServer []byte) bool {
	msgPayload := e.server.ProcessMessage(dataFromServer)
	if msgPayload == nil {
		return true
	} else if string(msgPayload) == "heartbeat" {
		return false
	}

	// send channel open
	err := e.localManagement.Connect()
	go e.localManagement.Listen()

	if err != nil {
		log.Error(err)
		return true
	}
	if e.isLME {
		// wait for channel open confirmation
		<-e.status
		log.Trace("Channel open confirmation received")
	} else {
		//with LMS we open/close websocket on every request, so setup close for when we're done handling LMS data
		defer e.localManagement.Close()
	}

	// send our data to LMX
	err = e.localManagement.Send(msgPayload)
	if err != nil {
		log.Error(err)
		return true
	}

	for {
		select {
		case dataFromLM := <-e.data:
			//log.Trace("before HandleDataFromLM: ", dataFromLM)
			e.HandleDataFromLM(dataFromLM)
			//log.Trace("after HandleDataFromLM: ", dataFromLM)
			if e.isLME {
				<-e.status
			}
			return false
		case errFromLMS := <-e.errors:
			if errFromLMS != nil {
				log.Error("error from LMS")
				return true
			}
		}
	}
}

func (e Executor) HandleDataFromLM(data []byte) {
	if len(data) > 0 {
		log.Debug("received data from LMX:")
		data = e.checkForJumbling(data)
		err := e.server.Send(e.payload.CreateMessageResponse(data))
		if err != nil {
			log.Error(err)
		}
	}
}

var jumbleRsps = []string{
	//"<g:AMT_GeneralSettings>",
	//"<g:IPS_HostBasedSetupService>",
	//"<a:Body><g:PullResponse><g:Items><h:AMT_EthernetPortSettings>",
	//"<a:Body><g:AMT_EthernetPortSettings>",
	//"CIM_WiFiEndpointSettings</c:ResourceURI></a:Header><a:Body><g:EnumerateResponse><g:EnumerationContext>",
	//"CIM_WiFiEndpointSettings</c:ResourceURI></a:Header><a:Body><g:PullResponse>",
	//"<a:Body><g:AMT_WiFiPortConfigurationService>",
	//"<g:AMT_RedirectionService>",
	//"<g:IPS_OptInService>",
	//"<g:CIM_KVMRedirectionSAP>",
	//"<g:AMT_RedirectionService>",
	//"<g:IPS_OptInService>",
	//"<g:AddTrustedRootCertificate_OUTPUT>",
	//"<g:AddMpServer_OUTPUT>",
	//"<a:Body><g:PullResponse><g:Items><h:AMT_ManagementPresenceRemoteSAP>",
	//"<a:Body><g:AddRemoteAccessPolicyRule_OUTPUT>",
	//"AMT_RemoteAccessPolicyAppliesToMPS</c:ResourceURI></a:Header><a:Body><g:EnumerateResponse>",
	//"AMT_RemoteAccessPolicyAppliesToMPS</c:ResourceURI></a:Header><a:Body><g:PullResponse>",
	//"<g:AMT_EnvironmentDetectionSettingData>",
	//"<g:AMT_EnvironmentDetectionSettingData>",
	// ---BEGIN: TLS state flow
	"AMT_PublicKeyCertificate</c:ResourceURI></a:Header><a:Body><g:EnumerateResponse>",
	"AMT_PublicKeyCertificate</c:ResourceURI></a:Header><a:Body><g:PullResponse>",
	"AMT_PublicKeyManagementService</c:ResourceURI></a:Header><a:Body><g:AddTrustedRootCertificate_OUTPUT><g:CreatedCertificate>",
	"AMT_PublicKeyManagementService</c:ResourceURI></a:Header><a:Body><g:GenerateKeyPair_OUTPUT>",
	"AMT_PublicPrivateKeyPair</c:ResourceURI></a:Header><a:Body><g:EnumerateResponse>",
	"AMT_PublicPrivateKeyPair</c:ResourceURI></a:Header><a:Body><g:PullResponse>",
	"AMT_PublicKeyManagementService</c:ResourceURI></a:Header><a:Body><g:AddCertificate_OUTPUT><g:CreatedCertificate>",
	"AMT_TLSCredentialContext</c:ResourceURI></a:Header><a:Body><g:ResourceCreated>",
	"AMT_TimeSynchronizationService</c:ResourceURI></a:Header><a:Body><g:GetLowAccuracyTimeSynch_OUTPUT>",
	"AMT_TimeSynchronizationService</c:ResourceURI></a:Header><a:Body><g:SetHighAccuracyTimeSynch_OUTPUT>",
	"AMT_TLSSettingData</c:ResourceURI></a:Header><a:Body><g:EnumerateResponse>",
	"AMT_TLSSettingData</c:ResourceURI></a:Header><a:Body><g:PullResponse>",
	"AMT_TLSSettingData</c:ResourceURI></a:Header><a:Body><g:AMT_TLSSettingData>",
	"</g:AcceptNonSecureConnections><g:ElementName>Intel(r) AMT 802.3 TLS Settings</g:ElementName>",
	"</g:AcceptNonSecureConnections><g:ElementName>Intel(r) AMT LMS TLS Settings</g:ElementName>",
	//"AMT_SetupAndConfigurationService</c:ResourceURI></a:Header><a:Body><g:CommitChanges_OUTPUT>",
	// ---END: TLS state flow
}
var jumbleCounts = make([]int, len(jumbleRsps))

func (e Executor) checkForJumbling(data []byte) []byte {

	var shouldJumble = false
	for i, v := range jumbleRsps {
		if jumbleCounts[i] < 3 && strings.Contains(string(data), v) {
			jumbleCounts[i]++
			shouldJumble = jumbleCounts[i] < 3
			log.Debug("jumbling", v)
			break
		}
	}

	if shouldJumble {
		jumbles := strings.Split(string(data), "\r\n")
		_, jumbles = jumbles[0], jumbles[1:]
		var returnData = []byte(strings.Join(jumbles, ""))
		log.Trace(string(returnData))
		return returnData
	} else {
		return data
	}
}
