package action

import (
	"flag"
	"fmt"
	"jch-onvif/internal/util"
	"time"
)

var SetSystemDateTimeHandler = Handler{
	ActionKey: "set-system-date-time",
	Action: func(args Arguments) error {
		dateTimeType := flag.Arg(0)
		if (dateTimeType != "Manual") && (dateTimeType != "NTP") {
			return fmt.Errorf("invalid DateTimeType: %s", dateTimeType)
		}
		daylightSaving := flag.Arg(1)
		if (daylightSaving != "true") && (daylightSaving != "false") {
			return fmt.Errorf("invalid DaylightSavings: %s", daylightSaving)
		}
		posixTimeZone := flag.Arg(2)
		if posixTimeZone == "" {
			return fmt.Errorf("timezone must be in posix format %s", posixTimeZone)
		}
		dateTime := flag.Arg(3)
		if dateTime == "" {
			return fmt.Errorf("time is not defined %s", dateTime)
		}
		return setSystemDateTime(args.NetInterface, args.ServiceUrl, dateTimeType, daylightSaving, posixTimeZone, dateTime)
	},
}

var setSystemDateAndTimePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:tds="http://www.onvif.org/ver10/device/wsdl"	
	xmlns:tt="http://www.onvif.org/ver10/schema"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/device/wsdl/SetSystemDateAndTime	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<tds:SetSystemDateAndTime>
			<tds:DateTimeType>%s</tds:DateTimeType>
			<tds:DaylightSavings>%s</tds:DaylightSavings>
			<tds:TimeZone>
				<tt:TZ>%s</tt:TZ>
			</tds:TimeZone>
			<tds:UTCDateTime>
				<tt:Time>
					<tt:Hour>%d</tt:Hour>
					<tt:Minute>%d</tt:Minute>
					<tt:Second>%d</tt:Second>
				</tt:Time>
				<tt:Date>
					<tt:Year>%d</tt:Year>
					<tt:Month>%d</tt:Month>
					<tt:Day>%d</tt:Day>
				</tt:Date>
			</tds:UTCDateTime>
		</tds:SetSystemDateAndTime>
	</soapenv:Body>
</soapenv:Envelope>
`

func setSystemDateTime(interfaceName string, serviceUrl string, dateTimeType string, daylightSaving string, posixTimezone string, isoDate string) error {
	var t time.Time
	var err error
	if isoDate == "now" {
		t = time.Now()
	} else {
		t, err = time.Parse("2006-01-02T15:04:05", isoDate)
		if err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}
	}
	utcTime := t.UTC()
	fmt.Printf("Setting date time to %s with time zone %s\n", utcTime, posixTimezone)
	fmt.Println("Performing SetSystemDateAndTime action")
	uuid := util.GenerateUUID()
	message := fmt.Sprintf(setSystemDateAndTimePayloadTemplate, uuid, dateTimeType, daylightSaving, posixTimezone, utcTime.Hour(),
		utcTime.Minute(), utcTime.Second(), utcTime.Year(), utcTime.Month(), utcTime.Day())
	_, err = util.SoapCall(interfaceName, serviceUrl, message)
	if err != nil {
		return err
	}
	return nil
}
