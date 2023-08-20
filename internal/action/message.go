package action

type OnvifResponseHeader struct {
	MessageID string `xml:"MessageID"`
	Action    string `xml:"Action"`
	RelatesTo string `xml:"RelatesTo"`
}
