package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetNetworkInterfacesHander = Handler{
	ActionKey: "get-network-interfaces",
	Action: func(args Arguments) error {
		return getNetworkInterfaces(args.NetInterface, args.ServiceUrl)
	},
}

var getNetworkInterfacesPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetNetworkInterfaces	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetNetworkInterfaces />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetNetworkInterfacesResponse struct {
	XMLName           xml.Name            `xml:"Envelope"`
	Header            OnvifResponseHeader `xml:"Header"`
	NetworkInterfaces NetworkInterfaces   `xml:"Body>GetNetworkInterfacesResponse>NetworkInterfaces"`
}

type NetworkInterfaces struct {
	Token   string               `xml:"token,attr"`
	Enabled string               `xml:"Enabled"`
	Info    NetworkInterfaceInfo `xml:"Info"`
	IPv4    IPNetworkInterface   `xml:"IPv4"`
	IPv6    IPNetworkInterface   `xml:"IPv6"`
}

type NetworkInterfaceInfo struct {
	Name      string `xml:"Name"`
	HwAddress string `xml:"HwAddress"`
	MTU       string `xml:"MTU"`
}

type IPNetworkInterface struct {
	Enabled string         `xml:"Enabled"`
	Manual  IPManualConfig `xml:"Config>Manual"`
	DHCP    string         `xml:"Config>DHCP"`
}

type IPManualConfig struct {
	Address      string `xml:"Address"`
	PrefixLength string `xml:"PrefixLength"`
}

func getNetworkInterfaces(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetNetworkInterfaces action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getNetworkInterfacesPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getNetworkInterfacesResponse GetNetworkInterfacesResponse
	err = xml.Unmarshal(content, &getNetworkInterfacesResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	networkInterface := getNetworkInterfacesResponse.NetworkInterfaces
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "Token", networkInterface.Token)
	fmt.Printf("%-15s: %s\n", "Enabled", networkInterface.Enabled)
	fmt.Printf("%-15s: %s\n", "Name", networkInterface.Info.Name)
	fmt.Printf("%-15s: %s\n", "HwAddress", networkInterface.Info.HwAddress)
	fmt.Printf("%-15s: %s\n", "MTU", networkInterface.Info.MTU)
	fmt.Println()
	if networkInterface.IPv4.Enabled == "true" {
		fmt.Println("IPv4")
		fmt.Println("==========")
		fmt.Printf("%-15s: %s\n", "Enabled", networkInterface.IPv4.Enabled)
		fmt.Printf("%-15s: %s\n", "Address", networkInterface.IPv4.Manual.Address)
		fmt.Printf("%-15s: %s\n", "Prefix Length", networkInterface.IPv4.Manual.PrefixLength)
		fmt.Printf("%-15s: %s\n", "DHCP", networkInterface.IPv4.DHCP)
	}
	if networkInterface.IPv6.Enabled == "true" {
		fmt.Println("IPv6")
		fmt.Println("==========")
		fmt.Printf("%-15s: %s\n", "Enabled", networkInterface.IPv6.Enabled)
		fmt.Printf("%-15s: %s\n", "Address", networkInterface.IPv6.Manual.Address)
		fmt.Printf("%-15s: %s\n", "Prefix Length", networkInterface.IPv6.Manual.PrefixLength)
		fmt.Printf("%-15s: %s\n", "DHCP", networkInterface.IPv6.DHCP)
	}
	return nil
}
