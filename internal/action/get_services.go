package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetServicesHandler = Handler{
	ActionKey: "get-services",
	Action: func(args Arguments) error {
		return getServices(args.NetInterface, args.ServiceUrl)
	},
}

var getServicesPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetServices	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetServices />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetServicesResponse struct {
	XMLName  xml.Name            `xml:"Envelope"`
	Header   OnvifResponseHeader `xml:"Header"`
	Services []Service           `xml:"Body>GetServicesResponse>Service"`
}

type Service struct {
	Namespace string `xml:"Namespace"`
	XAddr     string `xml:"XAddr"`
}

func getServices(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetServices action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getServicesPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getServicesResponse GetServicesResponse
	err = xml.Unmarshal(content, &getServicesResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	for _, service := range getServicesResponse.Services {
		fmt.Printf("%s => %s\n", service.Namespace, service.XAddr)
	}
	return nil
}
