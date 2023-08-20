package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetNTPHandler = Handler{
	ActionKey: "get-ntp",
	Action: func(args Arguments) error {
		return getNTP(args.NetInterface, args.ServiceUrl)
	},
}

var getNTPPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetNTP	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetNTP />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetNTPResponse struct {
	XMLName     xml.Name            `xml:"Envelope"`
	Header      OnvifResponseHeader `xml:"Header"`
	FromDHCP    string              `xml:"Body>GetNTPResponse>NTPInformation>FromDHCP"`
	NTPFromDHCP []NTPAddress        `xml:"Body>GetNTPResponse>NTPInformation>NTPFromDHCP"`
	NTPManual   []NTPAddress        `xml:"Body>GetNTPResponse>NTPInformation>NTPManual"`
}

type NTPAddress struct {
	Type        string `xml:"Type"`
	IPv4Address string `xml:"IPv4Address"`
	IPv6Address string `xml:"IPv6Address"`
}

func getNTP(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetNTP action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getNTPPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getNTPResponse GetNTPResponse
	err = xml.Unmarshal(content, &getNTPResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "FromDHCP", getNTPResponse.FromDHCP)
	if getNTPResponse.FromDHCP == "true" {
		for i, addr := range getNTPResponse.NTPFromDHCP {
			fmt.Printf("%-14d : %s %s %s\n", i, addr.Type, addr.IPv4Address, addr.IPv6Address)
		}
	} else {
		for i, addr := range getNTPResponse.NTPManual {
			fmt.Printf("%-14d : %s %s %s\n", i, addr.Type, addr.IPv4Address, addr.IPv6Address)
		}
	}
	return nil
}
