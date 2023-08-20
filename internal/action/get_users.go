package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetUsersHandler = Handler{
	ActionKey: "get-users",
	Action: func(args Arguments) error {
		return getUsers(args.NetInterface, args.ServiceUrl)
	},
}

var getUsersPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetUsers	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetUsers />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetUsersResponse struct {
	XMLName xml.Name            `xml:"Envelope"`
	Header  OnvifResponseHeader `xml:"Header"`
	Users   []User              `xml:"Body>GetUsersResponse>User"`
}

type User struct {
	Username  string `xml:"Username"`
	UserLevel string `xml:"UserLevel"`
}

func getUsers(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetUsers action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getUsersPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getUsersResponse GetUsersResponse
	err = xml.Unmarshal(content, &getUsersResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	for _, user := range getUsersResponse.Users {
		fmt.Printf("%s => %s\n", user.Username, user.UserLevel)
	}
	return nil
}
