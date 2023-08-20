package action

import (
	"fmt"
	"jch-onvif/internal/util"
)

var SetNTPHandler = Handler{
	ActionKey: "set-ntp",
	Action: func(args Arguments) error {
		return setNTP(args.NetInterface, args.ServiceUrl, args.Args[0])
	},
}

var setNTPPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SetNTP	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SetNTP>
			<tds:FromDHCP>false</tds:FromDHCP>
			<tds:NTPManual>
				<tt:Type>IPv4</tt:Type>
				<tt:IPv4Address>%s</tt:IPv4Address>
			</tds:NTPManual>
		</tds:SetNTP>
	</soapenv:Body>
</soapenv:Envelope>
`

func setNTP(interfaceName string, serviceUrl string, ntpServerAddress string) error {
	fmt.Println("Performing SetNTP action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(setNTPPayloadTemplate, uuid, ntpServerAddress))
	if err != nil {
		return err
	}
	return nil
}
