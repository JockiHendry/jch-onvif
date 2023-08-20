package action

import (
	"flag"
	"fmt"
	"jch-onvif/internal/util"
)

var SetDiscoveryModeHandler = Handler{
	ActionKey: "set-discovery-mode",
	Action: func(args Arguments) error {
		discoveryMode := flag.Arg(0)
		if (discoveryMode != "Discoverable") && (discoveryMode != "NonDiscoverable") {
			return fmt.Errorf("invalid discovery mode %s", discoveryMode)
		}
		return setDiscoveryMode(args.NetInterface, args.ServiceUrl, discoveryMode)
	},
}

var setDiscoveryModeTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SetDiscoveryMode	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SetDiscoveryMode>
			<tds:DiscoveryMode>%s</tds:DiscoveryMode>	
		</tds:SetDiscoveryMode>
	</soapenv:Body>
</soapenv:Envelope>
`

func setDiscoveryMode(interfaceName string, serviceUrl string, discoveryMode string) error {
	fmt.Println("Performing SetDiscoveryMode action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(setDiscoveryModeTemplate, uuid, discoveryMode))
	if err != nil {
		return err
	}
	return nil
}
