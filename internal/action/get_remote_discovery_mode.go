package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetRemoteDiscoveryModeHandler = Handler{
	ActionKey: "get-remote-discovery-mode",
	Action: func(args Arguments) error {
		return getRemoteDiscoveryMode(args.NetInterface, args.ServiceUrl)
	},
}

var getRemoteDiscoveryModePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetRemoteDiscoveryMode	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetRemoteDiscoveryMode />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetRemoteDiscoveryModeResponse struct {
	XMLName       xml.Name            `xml:"Envelope"`
	Header        OnvifResponseHeader `xml:"Header"`
	DiscoveryMode string              `xml:"Body>GetRemoteDiscoveryModeResponse>RemoteDiscoveryMode"`
}

func getRemoteDiscoveryMode(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetRemoteDiscoveryMode action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getRemoteDiscoveryModePayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getRemoteDiscoveryModeResponse GetRemoteDiscoveryModeResponse
	err = xml.Unmarshal(content, &getRemoteDiscoveryModeResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("Discovery Mode: %s\n", getRemoteDiscoveryModeResponse.DiscoveryMode)
	return nil
}
