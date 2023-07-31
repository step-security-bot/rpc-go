package local

import (
	"encoding/xml"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/amt/setupandconfiguration"
	log "github.com/sirupsen/logrus"
	"rpc/pkg/utils"
)

func (service *ProvisioningService) Deactivate() int {

	controlMode, err := service.amtCommand.GetControlMode()
	if err != nil {
		log.Error(err)
		return utils.AMTConnectionFailed
	}
	if controlMode == 1 {
		return service.DeactivateCCM()
	} else if controlMode == 2 {
		return service.DeactivateACM()
	}
	log.Error("Deactivation failed. Device control mode: " + utils.InterpretControlMode(controlMode))
	return utils.UnableToDeactivate
}

func (service *ProvisioningService) DeactivateACM() int {
	service.setupWsmanClient("admin", service.flags.Password)
	msg := service.amtMessages.SetupAndConfigurationService.Unprovision(1)
	response, err := service.client.Post(msg)
	if err != nil {
		log.Error("Status: Unable to deactivate ", err)
		return utils.UnableToDeactivate
	}
	var setupResponse setupandconfiguration.UnprovisionResponse
	err = xml.Unmarshal([]byte(response), &setupResponse)
	if err != nil {
		log.Error("Status: Failed to deactivate ", err)
		return utils.DeactivationFailed
	}
	if setupResponse.Body.Unprovision_OUTPUT.ReturnValue != 0 {
		log.Error("Status: Failed to deactivate. ReturnValue: ", setupResponse.Body.Unprovision_OUTPUT.ReturnValue)
		return utils.DeactivationFailed
	}
	log.Info("Status: Device deactivated in ACM.")
	return utils.Success
}

func (service *ProvisioningService) DeactivateCCM() int {
	status, err := service.amtCommand.Unprovision()
	if err != nil || status != 0 {
		log.Error("Status: Failed to deactivate ", err)
		return utils.DeactivationFailed
	}
	log.Info("Status: Device deactivated.")
	return utils.Success
}
