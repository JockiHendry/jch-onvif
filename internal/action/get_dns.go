package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetDNSHandler = Handler{
	ActionKey: "get-dns",
	Action: func(args Arguments) error {
		return getDNS(args.NetInterface, args.ServiceUrl)
	},
}

var getDNSPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetDNS	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetDNS />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetDNSResponse struct {
	XMLName      xml.Name            `xml:"Envelope"`
	Header       OnvifResponseHeader `xml:"Header"`
	FromDHCP     string              `xml:"Body>GetDNSResponse>DNSInformation>FromDHCP"`
	SearchDomain string              `xml:"Body>GetDNSResponse>DNSInformation>SearchDomain"`
	DNSManual    []DNSManual         `xml:"Body>GetDNSResponse>DNSInformation>DNSManual"`
}

type DNSManual struct {
	Type        string `xml:"Type"`
	IPv4Address string `xml:"IPv4Address"`
	IPv6Address string `xml:"IPv6Address"`
}

func getDNS(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetDNS action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getDNSPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getDNSResponse GetDNSResponse
	err = xml.Unmarshal(content, &getDNSResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "FromDHCP", getDNSResponse.FromDHCP)
	fmt.Printf("%-15s: %s\n", "SearchDomain", getDNSResponse.SearchDomain)
	for _, dnsManual := range getDNSResponse.DNSManual {
		if dnsManual.Type == "IPv4" {
			fmt.Printf("%-15s: %s\n", dnsManual.Type, dnsManual.IPv4Address)
		} else {
			fmt.Printf("%-15s: %s\n", dnsManual.Type, dnsManual.IPv6Address)
		}
	}
	return nil
}
