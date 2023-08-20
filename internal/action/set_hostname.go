package action

import (
	"fmt"
	"jch-onvif/internal/util"
)

var SetHostnameHandler = Handler{
	ActionKey: "set-hostname",
	Action: func(args Arguments) error {
		return setHostname(args.NetInterface, args.ServiceUrl, args.Args[0])
	},
}

var setHostnamePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SetHostname	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SetHostname>
			<tds:Name>%s</tds:Name>
		</tds:SetHostname>
	</soapenv:Body>
</soapenv:Envelope>
`

func setHostname(interfaceName string, serviceUrl string, hostname string) error {
	fmt.Println("Performing SetHostname action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(setHostnamePayloadTemplate, uuid, hostname))
	if err != nil {
		return err
	}
	return nil
}
