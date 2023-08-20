package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetServiceCapabilitiesHandler = Handler{
	ActionKey: "get-service-capabilities",
	Action: func(args Arguments) error {
		return getServiceCapabilities(args.NetInterface, args.ServiceUrl)
	},
}

var getServiceCapabilitiesPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetServiceCapabilities	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetServiceCapabilities />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetServiceCapabilitiesResponse struct {
	XMLName  xml.Name            `xml:"Envelope"`
	Header   OnvifResponseHeader `xml:"Header"`
	Network  NetworkCapability   `xml:"Body>GetServiceCapabilitiesResponse>Capabilities>Network"`
	Security SecurityCapability  `xml:"Body>GetServiceCapabilitiesResponse>Capabilities>Security"`
	System   SystemCapability    `xml:"Body>GetServiceCapabilitiesResponse>Capabilities>System"`
}

type NetworkCapability struct {
	NTP                string `xml:"NTP,attr"`
	HostnameFromDHCP   string `xml:"HostnameFromDHCP,attr"`
	Dot11Configuration string `xml:"Dot11Configuration,attr"`
	DynDNS             string `xml:"DynDNS,attr"`
	IPVersion6         string `xml:"IPVersion6,attr"`
	ZeroConfiguration  string `xml:"ZeroConfiguration,attr"`
	IPFilter           string `xml:"IPFilter,attr"`
}

type SecurityCapability struct {
	SupportedEAPMethods  string `xml:"SupportedEAPMethods,attr"`
	RELToken             string `xml:"RELToken,attr"`
	HttpDigest           string `xml:"HttpDigest,attr"`
	UsernameToken        string `xml:"UsernameToken,attr"`
	KerberosToken        string `xml:"KerberosToken,attr"`
	SAMLToken            string `xml:"SAMLToken,attr"`
	X509Token            string `xml:"X.509Token,attr"`
	RemoteUserHandling   string `xml:"RemoteUserHandling,attr"`
	Dot1X                string `xml:"Dot1X,attr"`
	AccessPolicyConfig   string `xml:"AccessPolicyConfig,attr"`
	OnboardKeyGeneration string `xml:"OnboardKeyGeneration,attr"`
	TLS12                string `xml:"TLS1.2,attr"`
	TLS11                string `xml:"TLS1.1,attr"`
	TLS10                string `xml:"TLS1.0,attr"`
}

type SystemCapability struct {
	HttpSupportInformation string `xml:"HttpSupportInformation,attr"`
	HttpSystemLogging      string `xml:"HttpSystemLogging,attr"`
	HttpSystemBackup       string `xml:"HttpSystemBackup,attr"`
	HttpFirmwareUpgrade    string `xml:"HttpFirmwareUpgrade,attr"`
	FirmwareUpgrade        string `xml:"FirmwareUpgrade,attr"`
	SystemLogging          string `xml:"SystemLogging,attr"`
	SystemBackup           string `xml:"SystemBackup,attr"`
	RemoteDiscovery        string `xml:"RemoteDiscovery,attr"`
	DiscoveryBye           string `xml:"DiscoveryBye,attr"`
	DiscoveryResolve       string `xml:"DiscoveryResolve,attr"`
}

func getServiceCapabilities(interfaceName string, serviceUrl string) error {
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getServiceCapabilitiesPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getServiceCapabilitiesResponse GetServiceCapabilitiesResponse
	err = xml.Unmarshal(content, &getServiceCapabilitiesResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	networkCapabilities := getServiceCapabilitiesResponse.Network
	fmt.Println("Network Capabilities")
	fmt.Println("====================")
	fmt.Printf("%-22s: %s\n", "NTP", networkCapabilities.NTP)
	fmt.Printf("%-22s: %s\n", "HostnameFromDHCP", networkCapabilities.HostnameFromDHCP)
	fmt.Printf("%-22s: %s\n", "Dot11Configuration", networkCapabilities.Dot11Configuration)
	fmt.Printf("%-22s: %s\n", "DynDNS", networkCapabilities.DynDNS)
	fmt.Printf("%-22s: %s\n", "IPVersion6", networkCapabilities.IPVersion6)
	fmt.Printf("%-22s: %s\n", "ZeroConfiguration", networkCapabilities.ZeroConfiguration)
	fmt.Printf("%-22s: %s\n", "IPFilter", networkCapabilities.IPFilter)
	fmt.Println()
	securityCapabilities := getServiceCapabilitiesResponse.Security
	fmt.Println("Security Capabilities")
	fmt.Println("=====================")
	fmt.Printf("%-22s: %s\n", "SupportedEAPMethods", securityCapabilities.SupportedEAPMethods)
	fmt.Printf("%-22s: %s\n", "RELToken", securityCapabilities.RELToken)
	fmt.Printf("%-22s: %s\n", "HttpDigest", securityCapabilities.HttpDigest)
	fmt.Printf("%-22s: %s\n", "UsernameToken", securityCapabilities.UsernameToken)
	fmt.Printf("%-22s: %s\n", "KerberosToken", securityCapabilities.KerberosToken)
	fmt.Printf("%-22s: %s\n", "SAMLToken", securityCapabilities.SAMLToken)
	fmt.Printf("%-22s: %s\n", "X.509Token", securityCapabilities.X509Token)
	fmt.Printf("%-22s: %s\n", "RemoteUserHandling", securityCapabilities.RemoteUserHandling)
	fmt.Printf("%-22s: %s\n", "Dot1X", securityCapabilities.Dot1X)
	fmt.Printf("%-22s: %s\n", "AccessPolicyConfig", securityCapabilities.AccessPolicyConfig)
	fmt.Printf("%-22s: %s\n", "OnboardKeyGeneration", securityCapabilities.OnboardKeyGeneration)
	fmt.Printf("%-22s: %s\n", "TLS1.2", securityCapabilities.TLS12)
	fmt.Printf("%-22s: %s\n", "TLS1.1", securityCapabilities.TLS11)
	fmt.Printf("%-22s: %s\n", "TLS1.0", securityCapabilities.TLS10)
	fmt.Println()
	systemCapabilities := getServiceCapabilitiesResponse.System
	fmt.Println("System Capabilities")
	fmt.Println("===================")
	fmt.Printf("%-22s: %s\n", "HttpSupportInfo", systemCapabilities.HttpSupportInformation)
	fmt.Printf("%-22s: %s\n", "HttpSystemLogging", systemCapabilities.HttpSystemLogging)
	fmt.Printf("%-22s: %s\n", "HttpSystemBackup", systemCapabilities.HttpSystemBackup)
	fmt.Printf("%-22s: %s\n", "HttpFirmwareUpgrade", systemCapabilities.HttpFirmwareUpgrade)
	fmt.Printf("%-22s: %s\n", "FirmwareUpgrade", systemCapabilities.FirmwareUpgrade)
	fmt.Printf("%-22s: %s\n", "SystemLogging", systemCapabilities.SystemLogging)
	fmt.Printf("%-22s: %s\n", "SystemBackup", systemCapabilities.SystemBackup)
	fmt.Printf("%-22s: %s\n", "RemoteDiscovery", systemCapabilities.RemoteDiscovery)
	fmt.Printf("%-22s: %s\n", "DiscoveryBye", systemCapabilities.DiscoveryBye)
	fmt.Printf("%-22s: %s\n", "DiscoveryResolve", systemCapabilities.DiscoveryResolve)
	return nil
}
