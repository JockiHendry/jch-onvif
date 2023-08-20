package action

import (
	"flag"
	"fmt"
	"jch-onvif/internal/util"
	"net"
)

var SetNetworkInterfacesHandler = Handler{
	ActionKey: "set-network-ipv4",
	Action: func(args Arguments) error {
		networkToken := flag.Arg(0)
		if networkToken == "" {
			return fmt.Errorf("invalid network token %s", networkToken)
		}
		ip := flag.Arg(1)
		if net.ParseIP(ip) == nil {
			return fmt.Errorf("invalid IP address %s", ip)
		}
		prefixLength := "24"
		if flag.Arg(2) != "" {
			prefixLength = flag.Arg(2)
		}
		useDHCP := "false"
		if flag.Arg(3) != "" {
			useDHCP = flag.Arg(3)
		}
		return setNetworkInterface(args.NetInterface, args.ServiceUrl, networkToken, ip, prefixLength, useDHCP)
	},
}

var setNetworkInterfacesPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SetNetworkInterfaces	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SetNetworkInterfaces>
			<tds:InterfaceToken>%s</tds:InterfaceToken>
			<tds:NetworkInterface>
				<tt:Enabled>true</tt:Enabled>
				<tt:IPv4>
					<tt:Enabled>true</tt:Enabled>
					<tt:Manual>
						<tt:Address>%s</tt:Address>
						<tt:PrefixLength>%s</tt:PrefixLength>
					</tt:Manual>
					<tt:DHCP>%s</tt:DHCP>
				</tt:IPv4>
			</tds:NetworkInterface>
		</tds:SetNetworkInterfaces>
	</soapenv:Body>
</soapenv:Envelope>
`

func setNetworkInterface(interfaceName string, serviceUrl string, interfaceToken string, ipv4 string, ipv4Prefix string, useDHCP string) error {
	fmt.Println("Performing SetNetworkInterface action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(setNetworkInterfacesPayloadTemplate, uuid, interfaceToken, ipv4, ipv4Prefix, useDHCP))
	if err != nil {
		return err
	}
	return nil
}
