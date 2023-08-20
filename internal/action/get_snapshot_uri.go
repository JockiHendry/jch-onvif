package action

import (
	"encoding/xml"
	"flag"
	"fmt"
	"jch-onvif/internal/util"
	"os/exec"
)

var GetSnapshotUriHandler = Handler{
	ActionKey: "get-snapshot-uri",
	Action: func(args Arguments) error {
		profileToken := flag.Arg(0)
		if profileToken == "" {
			return fmt.Errorf("invalid profile token %s", profileToken)
		}
		return getSnapshotUri(args.NetInterface, args.ServiceUrl, profileToken)
	},
}

var getSnapshotUriPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:trt="http://www.onvif.org/ver10/media/wsdl"
	xmlns:tt="http://www.onvif.org/ver10/schema"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/media/wsdl/GetSnapshotUri	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<trt:GetSnapshotUri>	
			<trt:ProfileToken>%s</trt:ProfileToken>
		</trt:GetStreamUri>
	</soapenv:Body>
</soapenv:Envelope>
`

type GetSnapshotUriResponse struct {
	XMLName             xml.Name            `xml:"Envelope"`
	Header              OnvifResponseHeader `xml:"Header"`
	Uri                 string              `xml:"Body>GetSnapshotUriResponse>MediaUri>Uri"`
	InvalidAfterConnect string              `xml:"Body>GetSnapshotUriResponse>MediaUri>InvalidAfterConnect"`
	InvalidAfterReboot  string              `xml:"Body>GetSnapshotUriResponse>MediaUri>InvalidAfterReboot"`
	Timeout             string              `xml:"Body>GetSnapshotUriResponse>MediaUri>Timeout"`
}

func getSnapshotUri(interfaceName string, serviceUrl string, profileToken string) error {
	fmt.Println("Performing GetStreamUri action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getSnapshotUriPayloadTemplate, uuid, profileToken))
	if err != nil {
		return err
	}
	var getSnapshotUriResponse GetSnapshotUriResponse
	err = xml.Unmarshal(content, &getSnapshotUriResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-20s: %s\n", "Uri", getSnapshotUriResponse.Uri)
	fmt.Printf("%-20s: %s\n", "InvalidAfterConnect", getSnapshotUriResponse.InvalidAfterConnect)
	fmt.Printf("%-20s: %s\n", "InvalidAfterReboot", getSnapshotUriResponse.InvalidAfterReboot)
	fmt.Printf("%-20s: %s\n\n", "Timeout", getSnapshotUriResponse.Timeout)
	path, err := exec.LookPath("firefox")
	if err != nil {
		path, err = exec.LookPath("google-chrome")
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(path, getSnapshotUriResponse.Uri)
	return cmd.Run()
}
