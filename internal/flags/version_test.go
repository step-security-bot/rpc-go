package flags

import (
	"github.com/stretchr/testify/assert"
	"rpc/pkg/utils"
	"testing"
)

// TODO: move this functionality testing
func TestHandleVersionCommand(t *testing.T) {
	f := NewFlags([]string{
		"rpc",
		"version",
	})

	result := f.handleVersionCommand()
	assert.EqualValues(t, result, utils.Success)
	assert.EqualValues(t, f.Local, true)
	// Test non-JSON output
	//expectedOutput := "RPC\nVersion " + utils.ProjectVersion + "\nProtocol " + utils.ProtocolVersion + "\n"
	//actualOutput, result := f.handleVersionCommand()
	//if result != utils.Success {
	//	t.Errorf("Expected success, but got %d", result)
	//}
	//
	//if actualOutput != expectedOutput {
	//	t.Errorf("Unexpected output.\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	//}
	//
	//// Test JSON output
	//f.JsonOutput = true
	//dataStruct := map[string]string{
	//	"app":      "RPC",
	//	"version":  utils.ProjectVersion,
	//	"protocol": utils.ProtocolVersion,
	//}
	//outBytes, _ := json.MarshalIndent(dataStruct, "", "  ")
	//expectedOutput = string(outBytes)
	//actualOutput, result = f.handleVersionCommand()
	//if result != utils.Success {
	//	t.Errorf("Expected success, but got %d", result)
	//}
	//
	//if actualOutput != expectedOutput {
	//	t.Errorf("Unexpected output.\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	//}
}
