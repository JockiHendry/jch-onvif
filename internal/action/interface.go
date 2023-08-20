package action

type Handler struct {
	ActionKey string
	Action    func(Arguments) error
}

type Arguments struct {
	NetInterface string
	TargetIP     string
	TargetRange  string
	ServiceUrl   string
	Args         []string
}

type Address struct {
	Type        string `xml:"Type"`
	IPv4Address string `xml:"IPv4Address"`
	IPv6Address string `xml:"IPv6Address"`
}

func (addr Address) String() string {
	if addr.Type == "IPv4" {
		return addr.IPv4Address
	} else {
		return addr.IPv6Address
	}
}
