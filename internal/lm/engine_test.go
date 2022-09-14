/*********************************************************************
 * Copyright (c) Intel Corporation 2022
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package lm

import (
	"errors"
	"rpc/pkg/apf"
	"rpc/pkg/pthi"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHECICommands struct{}

var message []byte
var sendBytesWritten uint32 = 12
var sendError error = nil
var initError error = nil
var bufferSize uint32 = 5120

func (c *MockHECICommands) Init(useLME bool) error { return initError }
func (c *MockHECICommands) GetBufferSize() uint32  { return bufferSize } // MaxMessageLength

func (c *MockHECICommands) SendMessage(buffer []byte, done *uint32) (bytesWritten uint32, err error) {
	return sendBytesWritten, sendError
}
func (c *MockHECICommands) ReceiveMessage(buffer []byte, done *uint32) (bytesRead uint32, err error) {
	for i := 0; i < len(message) && i < len(buffer); i++ {
		buffer[i] = message[i]
	}
	return uint32(len(message)), nil
}
func (c *MockHECICommands) Close() {}

var pthiVar pthi.Command

func init() {
	pthiVar = pthi.Command{}
	pthiVar.Heci = &MockHECICommands{}
}
func Test_NewLMEConnection(t *testing.T) {
	lmDataChannel := make(chan []byte)
	lmErrorChannel := make(chan error)
	lmStatusChannel := make(chan bool)
	lme := NewLMEConnection(lmDataChannel, lmErrorChannel, lmStatusChannel)
	assert.Equal(t, lmDataChannel, lme.Session.DataBuffer)
	assert.Equal(t, lmErrorChannel, lme.Session.ErrorBuffer)
	assert.Equal(t, lmStatusChannel, lme.Session.Status)
}
func TestLMEConnection_Initialize(t *testing.T) {
	testError := errors.New("test error")
	tests := []struct {
		name         string
		sendNumBytes uint32
		sendErr      error
		initErr      error
		wantErr      bool
	}{
		{
			name:         "Normal",
			sendNumBytes: uint32(93),
			sendErr:      nil,
			initErr:      nil,
			wantErr:      false,
		},
		{
			name:         "ExpectedFailureOnOpen",
			sendNumBytes: uint32(93),
			sendErr:      nil,
			initErr:      testError,
			wantErr:      true,
		},
		{
			name:         "ExpectedFailureOnExecute",
			sendNumBytes: uint32(93),
			sendErr:      testError,
			initErr:      nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendBytesWritten = tt.sendNumBytes
			sendError = tt.sendErr
			initError = tt.initErr
			lme := &LMEConnection{
				Command:    pthiVar,
				Session:    &apf.LMESession{},
				ourChannel: 1,
			}
			if err := lme.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("LMEConnection.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Send(t *testing.T) {
	sendBytesWritten = 14

	lme := &LMEConnection{
		Command:    pthiVar,
		Session:    &apf.LMESession{},
		ourChannel: 1,
	}
	data := []byte("hello")
	err := lme.Send(data)
	assert.NoError(t, err)
}
func Test_Connect(t *testing.T) {
	sendBytesWritten = 54
	lme := &LMEConnection{
		Command:    pthiVar,
		Session:    &apf.LMESession{},
		ourChannel: 1,
	}
	err := lme.Connect()
	assert.NoError(t, err)
}
func Test_Connect_With_Error(t *testing.T) {
	sendError = errors.New("no such device")
	sendBytesWritten = 54
	lme := &LMEConnection{
		Command:    pthiVar,
		Session:    &apf.LMESession{},
		ourChannel: 1,
	}
	err := lme.Connect()
	assert.Error(t, err)
}
func Test_Listen(t *testing.T) {
	lmDataChannel := make(chan []byte)
	lmErrorChannel := make(chan error)

	lme := &LMEConnection{
		Command: pthiVar,
		Session: &apf.LMESession{
			DataBuffer:  lmDataChannel,
			ErrorBuffer: lmErrorChannel,
			Status:      make(chan bool),
		},
		ourChannel: 1,
	}
	message = []byte{0x94, 0x01}
	defer lme.Close()
	go lme.Listen()
}

func Test_Close(t *testing.T) {
	lme := &LMEConnection{
		Command:    pthiVar,
		Session:    &apf.LMESession{},
		ourChannel: 1,
	}
	err := lme.Close()
	assert.NoError(t, err)
}
func Test_CloseWithChannel(t *testing.T) {
	lmDataChannel := make(chan []byte)
	lmErrorChannel := make(chan error)

	lme := &LMEConnection{
		Command: pthiVar,
		Session: &apf.LMESession{
			DataBuffer:  lmDataChannel,
			ErrorBuffer: lmErrorChannel,
			Status:      make(chan bool),
		},
		ourChannel: 1,
	}
	err := lme.Close()
	assert.NoError(t, err)
}
