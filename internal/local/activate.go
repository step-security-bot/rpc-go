package local

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"rpc/pkg/utils"
	"strings"

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
	certObject, fingerPrint, err := service.GetProvisioningCertObj()
	log.Info(certObject, fingerPrint)
	if err != nil {
		log.Error(err)
		return utils.ActivationFailed
	}
	generalSettings, err := service.GetGeneralSettings()
	if err != nil {
		log.Error(err)
		return utils.ActivationFailed
	}
	log.Info(generalSettings)
	getHostBasedSetupResponse, err := service.GetHostBasedSetupService()
	if err != nil {
		log.Error(err)
		return utils.ActivationFailed
	}
	log.Info(getHostBasedSetupResponse)
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
	tempstring := string(response)
	log.Info(tempstring)
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
	certs []*x509.Certificate
	keys  []interface{}
}

type CertificateObject struct {
	pem     string
	subject string
	issuer  string
}

type ProvisioningCertObj struct {
	certChain []string
	privateKey crypto.PrivateKey
}

func dumpPfx(pfxobj CertsAndKeys) (ProvisioningCertObj, string, error) {
	var provisioningCertificateObj ProvisioningCertObj
	var interObj []CertificateObject
	var leaf CertificateObject
	var root CertificateObject
	var fingerprint string

	if len(pfxobj.certs) > 0 {
		for i, cert := range pfxobj.certs {
			pemBlock := &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			}
			pem := string(pem.EncodeToMemory(pemBlock))
			pem = strings.Replace(pem, "-----BEGIN CERTIFICATE-----", "", -1)
		    pem = strings.Replace(pem, "-----END CERTIFICATE-----", "", -1)

			if i == 0 {
				leaf = CertificateObject{pem: pem, subject: cert.Subject.String(), issuer: cert.Issuer.String()}
			} else if cert.Subject.String() == cert.Issuer.String() {
				root = CertificateObject{pem: pem, subject: cert.Subject.String(), issuer: cert.Issuer.String()}
				der := cert.Raw
				hash := sha256.Sum256(der)
				fingerprint = hex.EncodeToString(hash[:])
			} else {
				inter := CertificateObject{pem: pem, subject: cert.Subject.String(), issuer: cert.Issuer.String()}
				interObj = append(interObj, inter)
			}
		}
	} else {
		return ProvisioningCertObj{}, "", errors.New("no certificates found")
	}

	provisioningCertificateObj.certChain = append(provisioningCertificateObj.certChain, leaf.pem)
	for _, inter := range interObj {
		provisioningCertificateObj.certChain = append(provisioningCertificateObj.certChain, inter.pem)
	}
	provisioningCertificateObj.certChain = append(provisioningCertificateObj.certChain, root.pem)

	if len(pfxobj.keys) > 0 {
		provisioningCertificateObj.privateKey = pfxobj.keys[0]
	}

	return provisioningCertificateObj, fingerprint, nil

}

func convertPfxToObject(pfxb64 string, passphrase string) (CertsAndKeys, error) {
	var pfxOut = CertsAndKeys{certs: []*x509.Certificate{}, keys: []interface{}{}}
	pfx, err := base64.StdEncoding.DecodeString(pfxb64)
	if err != nil {
		return pfxOut, err
	}

	privateKey, certificate, extraCerts, err := pkcs12.DecodeChain(pfx, passphrase)
	if err != nil {
		return pfxOut, errors.New("Decrypting provisioning certificate failed")
	}

	pfxOut.certs = append(pfxOut.certs, certificate)
	pfxOut.certs = append(pfxOut.certs, extraCerts...)
	pfxOut.keys = append(pfxOut.keys, privateKey)

	return pfxOut, nil
}

func (service *ProvisioningService) GetProvisioningCertObj() (ProvisioningCertObj, string, error) {
	config := service.config.ACMSettings
	certsAndKeys, err := convertPfxToObject(config.ProvisioningCert, config.ProvisioningCertPwd)
	if err != nil {
		log.Error("Failed to convert the certificate pfx to an object", err)
	}
	result, fingerprint, err := dumpPfx(certsAndKeys)
	if err != nil {
		log.Error("Failed to convert the certificate pfx to an object", err)
	} 
	return result, fingerprint, nil
}

func (service *ProvisioningService) CompareCertHashes() () {
	result, err := service.amtCommand.GetCertificateHashes()
	if err != nil {
		log.Error(err)
	}
	certs := make(map[string]interface{})
	for _, v := range result {
		certs[v.Name] = v
	}
}