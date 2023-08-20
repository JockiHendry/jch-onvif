package action

import (
	"flag"
	"fmt"
	"jch-onvif/internal/util"
)

var DeleteUsersHandler = Handler{
	ActionKey: "delete-user",
	Action: func(args Arguments) error {
		username := flag.Arg(0)
		if username == "" {
			return fmt.Errorf("invalid user name %s", username)
		}
		return deleteUser(args.NetInterface, args.ServiceUrl, username)
	},
}

var deleteUsersPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/DeleteUsers	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:DeleteUsers>	
			<tds:Username>%s</tds:Username>	
		</tds:DeleteUsers>
	</soapenv:Body>
</soapenv:Envelope>
`

func deleteUser(interfaceName string, serviceUrl string, username string) error {
	fmt.Println("Performing DeleteUsers action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(deleteUsersPayloadTemplate, uuid, username))
	if err != nil {
		return err
	}
	return nil
}
