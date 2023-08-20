package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetDeviceInformationHandler = Handler{
	ActionKey: "get-device-information",
	Action: func(args Arguments) error {
		return getDeviceInformation(args.NetInterface, args.ServiceUrl)
	},
}

var getDeviceInformationPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetDeviceInformation	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetDeviceInformation />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetDeviceInformationResponse struct {
	XMLName         xml.Name            `xml:"Envelope"`
	Header          OnvifResponseHeader `xml:"Header"`
	Manufacturer    string              `xml:"Body>GetDeviceInformationResponse>Manufacturer"`
	Model           string              `xml:"Body>GetDeviceInformationResponse>Model"`
	FirmwareVersion string              `xml:"Body>GetDeviceInformationResponse>FirmwareVersion"`
	SerialNumber    string              `xml:"Body>GetDeviceInformationResponse>SerialNumber"`
	HardwareId      string              `xml:"Body>GetDeviceInformationResponse>HardwareId"`
}

func getDeviceInformation(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetDeviceInformation action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getDeviceInformationPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getDeviceInformationResponse GetDeviceInformationResponse
	err = xml.Unmarshal(content, &getDeviceInformationResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "Manufacturer", getDeviceInformationResponse.Manufacturer)
	fmt.Printf("%-15s: %s\n", "Model", getDeviceInformationResponse.Model)
	fmt.Printf("%-15s: %s\n", "FirmwareVersion", getDeviceInformationResponse.FirmwareVersion)
	fmt.Printf("%-15s: %s\n", "SerialNumber", getDeviceInformationResponse.SerialNumber)
	fmt.Printf("%-15s: %s\n", "HardwareId", getDeviceInformationResponse.HardwareId)
	return nil
}
