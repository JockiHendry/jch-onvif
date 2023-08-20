package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetSystemDateTimeHandler = Handler{
	ActionKey: "get-system-date-time",
	Action: func(args Arguments) error {
		return getSystemDateTime(args.NetInterface, args.ServiceUrl)
	},
}

var getSystemDateAndTimePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/GetSystemDateAndTime	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:GetSystemDateAndTime />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetSystemDateAndTimeResponse struct {
	XMLName         xml.Name            `xml:"Envelope"`
	Header          OnvifResponseHeader `xml:"Header"`
	DateTimeType    string              `xml:"Body>GetSystemDateAndTimeResponse>SystemDateAndTime>DateTimeType"`
	DaylightSavings string              `xml:"Body>GetSystemDateAndTimeResponse>SystemDateAndTime>DaylightSavings"`
	TimeZone        string              `xml:"Body>GetSystemDateAndTimeResponse>SystemDateAndTime>TimeZone>TZ"`
	UTCDateTime     DateTime            `xml:"Body>GetSystemDateAndTimeResponse>SystemDateAndTime>UTCDateTime"`
	LocalDateTime   DateTime            `xml:"Body>GetSystemDateAndTimeResponse>SystemDateAndTime>LocalDateTime"`
}

type DateTime struct {
	Day    int `xml:"Date>Day"`
	Month  int `xml:"Date>Month"`
	Year   int `xml:"Date>Year"`
	Hour   int `xml:"Time>Hour"`
	Minute int `xml:"Time>Minute"`
	Second int `xml:"Time>Second"`
}

func (datetime DateTime) String() string {
	return fmt.Sprintf("%02d-%02d-%04d %02d:%02d:%02d", datetime.Day, datetime.Month, datetime.Year, datetime.Hour, datetime.Minute, datetime.Second)
}

func getSystemDateTime(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetSystemDateAndTime action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getSystemDateAndTimePayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getSystemDateAndTimeResponse GetSystemDateAndTimeResponse
	err = xml.Unmarshal(content, &getSystemDateAndTimeResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-15s: %s\n", "DateTimeType", getSystemDateAndTimeResponse.DateTimeType)
	fmt.Printf("%-15s: %s\n", "DaylightSavings", getSystemDateAndTimeResponse.DaylightSavings)
	fmt.Printf("%-15s: %s\n", "Timezone TZ", getSystemDateAndTimeResponse.TimeZone)
	fmt.Printf("%-15s: %s\n", "UTC DateTime", getSystemDateAndTimeResponse.UTCDateTime)
	fmt.Printf("%-15s: %s\n", "Local DateTime", getSystemDateAndTimeResponse.LocalDateTime)
	return nil
}
