package action

import (
	"flag"
	"fmt"
	"jch-onvif/internal/util"
)

var CreateUsersHandler = Handler{
	ActionKey: "create-user",
	Action: func(args Arguments) error {
		username := flag.Arg(0)
		if username == "" {
			return fmt.Errorf("invalid user name %s", username)
		}
		userLevel := flag.Arg(1)
		if userLevel == "" {
			return fmt.Errorf("invalid user level %s", userLevel)
		}
		password := flag.Arg(2)
		if password == "" {
			return fmt.Errorf("invalid password %s", password)
		}
		return createUser(args.NetInterface, args.ServiceUrl, username, userLevel, password)
	},
}

var createUsersPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/CreateUsers	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:CreateUsers>
			<tds:User>
				<tt:Username>%s</tt:Username>
				<tt:UserLevel>%s</tt:UserLevel>
				<tt:Password>%s</tt:Password>
			</tds:User>
		</tds:CreateUsers>
	</soapenv:Body>
</soapenv:Envelope>
`

func createUser(interfaceName string, serviceUrl string, username string, userLevel string, password string) error {
	fmt.Println("Performing CreateUsers action")
	uuid := util.GenerateUUID()
	_, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(createUsersPayloadTemplate, uuid, username, userLevel, password))
	if err != nil {
		return err
	}
	return nil
}
