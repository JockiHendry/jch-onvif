package action

import (
	"encoding/xml"
	"flag"
	"fmt"
	"jch-onvif/internal/util"
	"os/exec"
)

var GetStreamUriHandler = Handler{
	ActionKey: "get-stream-uri",
	Action: func(args Arguments) error {
		profileToken := flag.Arg(0)
		if profileToken == "" {
			return fmt.Errorf("invalid profile token %s", profileToken)
		}
		transport := flag.Arg(1)
		var stream string
		if transport != "" {
			if (transport != "UDP") && (transport != "TCP") && (transport != "RTSP") && (transport != "HTTP") {
				return fmt.Errorf("invalid transport %s", transport)
			}
			stream = flag.Arg(2)
			if (stream != "") && (stream != "RTP-Unicast") && (stream != "RTP-Multicast") {
				return fmt.Errorf("invalid stream %s", stream)
			}
		}
		return getStreamUri(args.NetInterface, args.ServiceUrl, profileToken, transport, stream)
	},
}

var getStreamUriPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:trt="http://www.onvif.org/ver10/media/wsdl"
	xmlns:tt="http://www.onvif.org/ver10/schema"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/media/wsdl/GetStreamUri	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<trt:GetStreamUri>
			<trt:StreamSetup>
				%s
				<tt:Transport>
					<tt:Protocol>%s</tt:Protocol>
				</tt:Transport>
			</trt:StreamSetup>
			<trt:ProfileToken>%s</trt:ProfileToken>
		</trt:GetStreamUri>
	</soapenv:Body>
</soapenv:Envelope>
`

type GetStreamUriResponse struct {
	XMLName             xml.Name            `xml:"Envelope"`
	Header              OnvifResponseHeader `xml:"Header"`
	Uri                 string              `xml:"Body>GetStreamUriResponse>MediaUri>Uri"`
	InvalidAfterConnect string              `xml:"Body>GetStreamUriResponse>MediaUri>InvalidAfterConnect"`
	InvalidAfterReboot  string              `xml:"Body>GetStreamUriResponse>MediaUri>InvalidAfterReboot"`
	Timeout             string              `xml:"Body>GetStreamUriResponse>MediaUri>Timeout"`
}

func getStreamUri(interfaceName string, serviceUrl string, profileToken string, transport string, stream string) error {
	fmt.Println("Performing GetStreamUri action")
	uuid := util.GenerateUUID()
	streamTag := ""
	if stream != "" {
		streamTag = "<tt:Stream>" + stream + "</tt:Stream>"
	}
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getStreamUriPayloadTemplate, uuid, streamTag, transport, profileToken))
	if err != nil {
		return err
	}
	var getStreamUriResponse GetStreamUriResponse
	err = xml.Unmarshal(content, &getStreamUriResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	fmt.Printf("%-20s: %s\n", "Uri", getStreamUriResponse.Uri)
	fmt.Printf("%-20s: %s\n", "InvalidAfterConnect", getStreamUriResponse.InvalidAfterConnect)
	fmt.Printf("%-20s: %s\n", "InvalidAfterReboot", getStreamUriResponse.InvalidAfterReboot)
	fmt.Printf("%-20s: %s\n\n", "Timeout", getStreamUriResponse.Timeout)
	path, err := exec.LookPath("vlc")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, getStreamUriResponse.Uri)
	return cmd.Run()
}
