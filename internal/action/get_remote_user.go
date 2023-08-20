package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetRemoteUserHandler = Handler{
	ActionKey: "get-remote-user",
	Action: func(args Arguments) error {
		return getRemoteUser(args.NetInterface, args.ServiceUrl)
	},
}

var getRemoteUserPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetRemoteUser	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetRemoteUser />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetRemoteUserResponse struct {
	XMLName            xml.Name            `xml:"Envelope"`
	Header             OnvifResponseHeader `xml:"Header"`
	Username           string              `xml:"Body>GetRemoteUserResponse>Username"`
	UseDerivedPassword string              `xml:"Body>GetRemoteUserResponse>UseDerivedPassword"`
}

func getRemoteUser(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetRemoteUser action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getRemoteUserPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getRemoteUserResponse GetRemoteUserResponse
	err = xml.Unmarshal(content, &getRemoteUserResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "Username", getRemoteUserResponse.Username)
	fmt.Printf("%-15s: %s\n", "Derived Pwd", getRemoteUserResponse.UseDerivedPassword)
	return nil
}
