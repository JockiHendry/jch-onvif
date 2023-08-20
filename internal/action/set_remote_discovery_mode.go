package action

import (
	"flag"
	"fmt"
	"jch-onvif/internal/util"
)

var SetRemoteDiscoveryModeHandler = Handler{
	ActionKey: "set-remote-discovery-mode",
	Action: func(args Arguments) error {
		discoveryMode := flag.Arg(0)
		if (discoveryMode != "Discoverable") && (discoveryMode != "NonDiscoverable") {
			return fmt.Errorf("invalid discovery mode %s", discoveryMode)
		}
		return setRemoteDiscoveryMode(args.NetInterface, args.ServiceUrl, discoveryMode)
	},
}

var setRemoteDiscoveryModeTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SetRemoteDiscoveryMode	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SetRemoteDiscoveryMode>
			<tds:RemoteDiscoveryMode>%s</tds:RemoteDiscoveryMode>	
		</tds:SetRemoteDiscoveryMode>
	</soapenv:Body>
</soapenv:Envelope>
`

func setRemoteDiscoveryMode(interfaceName string, serviceUrl string, discoveryMode string) error {
	fmt.Println("Performing SetRemoteDiscoveryMode action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(setRemoteDiscoveryModeTemplate, uuid, discoveryMode))
	if err != nil {
		return err
	}
	return nil
}
