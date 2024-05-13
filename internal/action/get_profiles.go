package action

import (
	"encoding/xml"
	"fmt"
	"jch-onvif/internal/util"
)

var GetProfilesHandler = Handler{
	ActionKey: "get-profiles",
	Action: func(args Arguments) error {
		return getProfiles(args.NetInterface, args.ServiceUrl)
	},
}

var getProfilesPayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:trt="http://www.onvif.org/ver10/media/wsdl"	
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://www.onvif.org/ver10/media/wsdl/GetProfiles	
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<trt:GetProfiles />	
	</soapenv:Body>
</soapenv:Envelope>
`

type GetProfilesResponse struct {
	XMLName  xml.Name            `xml:"Envelope"`
	Header   OnvifResponseHeader `xml:"Header"`
	Profiles []Profile           `xml:"Body>GetProfilesResponse>Profiles"`
}

type Profile struct {
	Token                     string                    `xml:"token,attr"`
	Name                      string                    `xml:"Name"`
	VideoSourceConfiguration  VideoSourceConfiguration  `xml:"VideoSourceConfiguration"`
	AudioSourceConfiguration  AudioSourceConfiguration  `xml:"AudioSourceConfiguration"`
	VideoEncoderConfiguration VideoEncoderConfiguration `xml:"VideoEncoderConfiguration"`
	AudioEncoderConfiguration AudioEncoderConfiguration `xml:"AudioEncoderConfiguration"`
}

type VideoSourceConfiguration struct {
	Token  string `xml:"token,attr"`
	Name   string `xml:"Name"`
	Bounds Bounds `xml:"Bounds"`
}

type AudioSourceConfiguration struct {
	Token string `xml:"token,attr"`
	Name  string `xml:"Name"`
}

type Bounds struct {
	Height int `xml:"height,attr"`
	Width  int `xml:"width,attr"`
	X      int `xml:"x,attr"`
	Y      int `xml:"y,attr"`
}

type VideoEncoderConfiguration struct {
	Token            string  `xml:"token,attr"`
	Name             string  `xml:"Name"`
	Encoding         string  `xml:"Encoding"`
	ResolutionWidth  string  `xml:"Resolution>Width"`
	ResolutionHeight string  `xml:"Resolution>Height"`
	Quality          string  `xml:"Quality"`
	FrameRateLimit   string  `xml:"RateControl>FrameRateLimit"`
	EncodingInterval string  `xml:"RateControl>EncodingInterval"`
	BitrateLimit     string  `xml:"RateControl>BitrateLimit"`
	MulticastAddress Address `xml:"Multicast>Address"`
	MulticastPort    string  `xml:"Multicast>Port"`
	SessionTimeout   string  `xml:"SessionTimeout"`
}

type AudioEncoderConfiguration struct {
	Token            string  `xml:"token,attr"`
	Name             string  `xml:"Name"`
	Encoding         string  `xml:"Encoding"`
	Bitrate          string  `xml:"Bitrate"`
	SampleRate       string  `xml:"SampleRate"`
	MulticastAddress Address `xml:"Multicast>Address"`
	MulticastPort    string  `xml:"Multicast>Port"`
	SessionTimeout   string  `xml:"SessionTimeout"`
}

func (bounds Bounds) String() string {
	return fmt.Sprintf("%d x %d (%d,%d)", bounds.Width, bounds.Height, bounds.X, bounds.Y)
}

func getProfiles(interfaceName string, serviceUrl string) error {
	fmt.Println("Performing GetProfiles action")
	uuid := util.GenerateUUID()
	content, err := util.SoapCall(interfaceName, serviceUrl, fmt.Sprintf(getProfilesPayloadTemplate, uuid))
	if err != nil {
		return err
	}
	var getProfilesResponse GetProfilesResponse
	err = xml.Unmarshal(content, &getProfilesResponse)
	if err != nil {
		return fmt.Errorf("error while parsing XML: %w", err)
	}
	fmt.Println()
	for _, profile := range getProfilesResponse.Profiles {
		fmt.Printf("Token [%s]\n", profile.Token)
		fmt.Printf("================\n")
		fmt.Printf("%-20s: %s\n", "Name", profile.Name)
		fmt.Printf("%-20s: %s\n", "Video Source Token", profile.VideoSourceConfiguration.Token)
		fmt.Printf("%-20s: %s\n", "  Name", profile.VideoSourceConfiguration.Name)
		fmt.Printf("%-20s: %s\n", "  Bounds", profile.VideoSourceConfiguration.Bounds.String())
		fmt.Printf("%-20s: %s\n", "Audio Source Token", profile.AudioSourceConfiguration.Token)
		fmt.Printf("%-20s: %s\n", "  Name", profile.AudioSourceConfiguration.Name)
		fmt.Printf("%-20s: %s\n", "Video Encoder Token", profile.VideoEncoderConfiguration.Token)
		fmt.Printf("%-20s: %s\n", "  Name", profile.VideoEncoderConfiguration.Name)
		fmt.Printf("%-20s: %s\n", "  Encoding", profile.VideoEncoderConfiguration.Encoding)
		fmt.Printf("%-20s: %s x %s\n", "  Resolution", profile.VideoEncoderConfiguration.ResolutionWidth, profile.VideoEncoderConfiguration.ResolutionHeight)
		fmt.Printf("%-20s: %s\n", "  Quality", profile.VideoEncoderConfiguration.Quality)
		fmt.Printf("%-20s: %s\n", "  FrameRate Limit", profile.VideoEncoderConfiguration.FrameRateLimit)
		fmt.Printf("%-20s: %s\n", "  Encoding Interval", profile.VideoEncoderConfiguration.EncodingInterval)
		fmt.Printf("%-20s: %s\n", "  BitRate Limit", profile.VideoEncoderConfiguration.BitrateLimit)
		fmt.Printf("%-20s: %s\n", "  Multicast Address", profile.VideoEncoderConfiguration.MulticastAddress)
		fmt.Printf("%-20s: %s\n", "  Multicast Port", profile.VideoEncoderConfiguration.MulticastPort)
		fmt.Printf("%-20s: %s\n", "  Session Timeout", profile.VideoEncoderConfiguration.SessionTimeout)
		fmt.Printf("%-20s: %s\n", "Audio Encoder Token", profile.AudioEncoderConfiguration.Token)
		fmt.Printf("%-20s: %s\n", "  Name", profile.AudioEncoderConfiguration.Name)
		fmt.Printf("%-20s: %s\n", "  Encoding", profile.AudioEncoderConfiguration.Encoding)
		fmt.Printf("%-20s: %s\n", "  BitRate", profile.AudioEncoderConfiguration.Bitrate)
		fmt.Printf("%-20s: %s\n", "  SampleRate", profile.AudioEncoderConfiguration.SampleRate)
		fmt.Printf("%-20s: %s\n", "  Multicast Address", profile.AudioEncoderConfiguration.MulticastAddress)
		fmt.Printf("%-20s: %s\n", "  Multicast Port", profile.AudioEncoderConfiguration.MulticastPort)
		fmt.Printf("%-20s: %s\n", "  Session Timeout", profile.AudioEncoderConfiguration.SessionTimeout)
		fmt.Println()
		fmt.Println()
	}
	return nil
}
