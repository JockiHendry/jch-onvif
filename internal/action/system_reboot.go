package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var SystemRebootHandler = Handler{
	ActionKey: "reboot",
	Action: func(args Arguments) error {
		return reboot(args.NetInterface, args.ServiceUrl)
	},
}

var systemRebootPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SystemReboot	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SystemReboot />	
	</soapenv:Body>
</soapenv:Envelope>
`

type SystemRebootResponse struct {
	XMLName xml.Name            `xml:"Envelope"`
	Header  OnvifResponseHeader `xml:"Header"`
	Message string              `xml:"Body>SystemRebootResponse>Message"`
}

func reboot(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing SystemReboot action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(systemRebootPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var systemRebootResponse SystemRebootResponse
	err = xml.Unmarshal(content, &systemRebootResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("Message: %s\n", systemRebootResponse.Message)
	return nil
}
