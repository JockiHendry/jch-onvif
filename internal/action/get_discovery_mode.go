package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetDiscoveryModeHandler = Handler{
	ActionKey: "get-discovery-mode",
	Action: func(args Arguments) error {
		return getDiscoveryMode(args.NetInterface, args.ServiceUrl)
	},
}

var getDiscoveryModePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetDiscoveryMode	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetDiscoveryMode />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetDiscoveryModeResponse struct {
	XMLName       xml.Name            `xml:"Envelope"`
	Header        OnvifResponseHeader `xml:"Header"`
	DiscoveryMode string              `xml:"Body>GetDiscoveryModeResponse>DiscoveryMode"`
}

func getDiscoveryMode(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetDiscoveryMode action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getDiscoveryModePayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getDiscoveryModeResponse GetDiscoveryModeResponse
	err = xml.Unmarshal(content, &getDiscoveryModeResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("Discovery Mode: %s\n", getDiscoveryModeResponse.DiscoveryMode)
	return nil
}
