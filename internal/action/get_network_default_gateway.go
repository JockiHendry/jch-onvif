package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetNetworkDefaultGatewayHandler = Handler{
	ActionKey: "get-network-default-gateway",
	Action: func(args Arguments) error {
		return getNetworkDefaultGateway(args.NetInterface, args.ServiceUrl)
	},
}

var getDefaultGatewayPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetNetworkDefaultGateway	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetNetworkDefaultGateway />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetNetworkDefaultGatewayResponse struct {
	XMLName     xml.Name            `xml:"Envelope"`
	Header      OnvifResponseHeader `xml:"Header"`
	IPv4Address string              `xml:"Body>GetNetworkDefaultGatewayResponse>NetworkGateway>IPv4Address"`
	IPv6Address string              `xml:"Body>GetNetworkDefaultGatewayResponse>NetworkGateway>IPv6Address"`
}

func getNetworkDefaultGateway(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetNetworkDefaultGateway action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getDefaultGatewayPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getNetworkDefaultGatewayResponse GetNetworkDefaultGatewayResponse
	err = xml.Unmarshal(content, &getNetworkDefaultGatewayResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "IPv4 Gateway", getNetworkDefaultGatewayResponse.IPv4Address)
	fmt.Printf("%-15s: %s\n", "IPv6 Gateway", getNetworkDefaultGatewayResponse.IPv6Address)
	return nil
}
