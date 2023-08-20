package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetHostnameHandler = Handler{
	ActionKey: "get-hostname",
	Action: func(args Arguments) error {
		return getHostname(args.NetInterface, args.ServiceUrl)
	},
}

var getHostnamePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetHostname	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetHostname />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetHostnameResponse struct {
	XMLName  xml.Name            `xml:"Envelope"`
	Header   OnvifResponseHeader `xml:"Header"`
	FromDHCP string              `xml:"Body>GetHostnameResponse>HostnameInformation>FromDHCP"`
	Name     string              `xml:"Body>GetHostnameResponse>HostnameInformation>Name"`
}

func getHostname(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetHostname action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getHostnamePayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getHostnameResponse GetHostnameResponse
	err = xml.Unmarshal(content, &getHostnameResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-10s: %s\n", "FromDHCP", getHostnameResponse.FromDHCP)
	fmt.Printf("%-10s: %s\n", "Hostname", getHostnameResponse.Name)
	return nil
}
