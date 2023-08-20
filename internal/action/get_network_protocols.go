package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetNetworkProtocolsHandler = Handler{
	ActionKey: "get-network-protocols",
	Action: func(args Arguments) error {
		return getNetworkProtocols(args.NetInterface, args.ServiceUrl)
	},
}

var getNetworkProtocolsPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetNetworkProtocols	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetNetworkProtocols />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetNetworkProtocolsResponse struct {
	XMLName          xml.Name            `xml:"Envelope"`
	Header           OnvifResponseHeader `xml:"Header"`
	NetworkProtocols []NetworkProtocol   `xml:"Body>GetNetworkProtocolsResponse>NetworkProtocols"`
}

type NetworkProtocol struct {
	Name    string `xml:"Name"`
	Enabled string `xml:"Enabled"`
	Port    string `xml:"Port"`
}

func getNetworkProtocols(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetNetworkProtocols action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getNetworkProtocolsPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getNetworkProtocolsResponse GetNetworkProtocolsResponse
	err = xml.Unmarshal(content, &getNetworkProtocolsResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	for _, networkProtocol := range getNetworkProtocolsResponse.NetworkProtocols {
		fmt.Printf("%-15s: %s\n", "Name", networkProtocol.Name)
		fmt.Printf("%-15s: %s\n", "Enabled", networkProtocol.Enabled)
		fmt.Printf("%-15s: %s\n", "Port", networkProtocol.Port)
		fmt.Println()
	}
	return nil
}
