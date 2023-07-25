package local

import (
	"encoding/xml"
	"errors"
	"rpc/pkg/utils"
	"encoding/base64"
	"crypto/x509"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/ips/hostbasedsetup"
	log "github.com/sirupsen/logrus"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func (service *ProvisioningService) Activate() int {

	controlMode, err := service.amtCommand.GetControlMode()
	if err != nil {
		log.Error(err)
		return utils.AMTConnectionFailed
	}
	if controlMode != 0 {
		log.Error("Device is already activated")
		return utils.UnableToActivate
	}

	// for local activation, wsman client needs local system account credentials
	lsa, err := service.amtCommand.GetLocalSystemAccount()
	if err != nil {
		log.Error(err)
		return utils.AMTConnectionFailed
	}
	service.setupWsmanClient(lsa.Username, lsa.Password)

	// CCM is the only option supported currently
	// (and is not required on the command line?)
	service.flags.UseACM = true

	resultCode := utils.Success

	if service.flags.UseACM {
		resultCode = service.ActivateACM()
	} else if service.flags.UseCCM {
		resultCode = service.ActivateCCM()
	}

	return resultCode
}

func (service *ProvisioningService) ActivateACM() int {
	generalSettings, err := service.GetGeneralSettings()
	if err != nil {
		log.Error(err)
		return utils.ActivationFailed
	}
	log.Info(generalSettings)
	// getHostBasedSetupResponse, err := service.GetHostBasedSetupService()
	// if err != nil {
	// 	log.Error(err)
	// 	return utils.ActivationFailed
	// }
	// log.Info(getHostBasedSetupResponse)
	certObject, err := service.GetProvisioningCertObj()
	log.Info(certObject)
	return utils.Success
}

func (service *ProvisioningService) ActivateCCM() int {
	generalSettings, err := service.GetGeneralSettings()
	if err != nil {
		log.Error(err)
		return utils.ActivationFailed
	}
	_, err = service.HostBasedSetup(generalSettings.Body.AMTGeneralSettings.DigestRealm, service.config.Password)
	if err != nil {
		log.Error(err)
		return utils.ActivationFailed
	}
	log.Info("Status: Device activated in Client Control Mode")
	return utils.Success
}

func (service *ProvisioningService) GetGeneralSettings() (general.Response, error) {
	message := service.amtMessages.GeneralSettings.Get()
	response, err := service.client.Post(message)
	if err != nil {
		return general.Response{}, err
	}
	var generalSettings general.Response
	err = xml.Unmarshal([]byte(response), &generalSettings)
	if err != nil {
		return general.Response{}, err
	}
	return generalSettings, nil
}

func (service *ProvisioningService) HostBasedSetup(digestRealm string, password string) (int, error) {
	message := service.ipsMessages.HostBasedSetupService.Setup(hostbasedsetup.AdminPassEncryptionTypeHTTPDigestMD5A1, digestRealm, password)
	response, err := service.client.Post(message)
	if err != nil {
		return utils.AMTConnectionFailed, err
	}
	var hostBasedSetupResponse hostbasedsetup.Response
	err = xml.Unmarshal([]byte(response), &hostBasedSetupResponse)
	if err != nil {
		return utils.ActivationFailed, err
	}
	if hostBasedSetupResponse.Body.Setup_OUTPUT.ReturnValue != 0 {
		return utils.ActivationFailed, errors.New("unable to activate CCM, check to make sure the device is not alreacy activated")
	}
	return utils.Success, nil
}

func (service *ProvisioningService) GetHostBasedSetupService() (hostbasedsetup.Response, error) {
	message := service.ipsMessages.HostBasedSetupService.Get()
	response, err := service.client.Post(message)
	log.Info(response)
	if err != nil {
		return hostbasedsetup.Response{}, err
	}
	var getHostBasedSetupResponse hostbasedsetup.Response
	err = xml.Unmarshal([]byte(response), &getHostBasedSetupResponse)
	if err != nil {
		return hostbasedsetup.Response{}, err
	}
	return getHostBasedSetupResponse, nil
}

type CertsAndKeys struct {
	Certs []x509.Certificate
	Keys  []interface{}
}

func (service *ProvisioningService) GetProvisioningCertObj() (CertsAndKeys, error) {
	// Read in cert
	pfxb64 := base64.StdEncoding.EncodeToString([]byte(service.config.ProvisioningCert))

	// Convert the certificate pfx to an object
	pfxOut := CertsAndKeys{
		Certs: []x509.Certificate{},
		Keys:  []interface{}{},
	}

	pfxder, _ := base64.StdEncoding.DecodeString(pfxb64)
	
	privateKey, certificate, extraCerts, err := pkcs12.DecodeChain(pfxder, service.config.ProvisioningCertPwd)
	if err != nil {
		return pfxOut, errors.New("Decrypting provisioning certificate failed")
	}
	pfxOut.Keys = append(pfxOut.Keys, privateKey)
	pfxOut.Certs = append(pfxOut.Certs, *certificate)
	log.Info(extraCerts)
	// pfxOut.Certs = append(pfxOut.Certs, extraCerts...)

	return pfxOut, nil

	// Return the certificate chain pems and private key
	// CertChainPfx, err = DumpPfx(pfxobj)
	// if err != nil {
	// 	log.Fatalf("Failed to dump PFX: %v", err)
	// }
}