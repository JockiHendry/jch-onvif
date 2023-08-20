package action

import (
	"context"
	"encoding/xml"
	"fmt"
	"golang.org/x/sys/unix"
	"jch-onvif/internal/util"
	"net"
	"net/netip"
	"strings"
	"sync"
	"syscall"
	"time"
)

var DiscoverHandler = Handler{
	ActionKey: "discover",
	Action: func(args Arguments) error {
		fmt.Println("Performing multicast discovery")
		return runDiscovery(args.NetInterface, []string{"239.255.255.250:3702"}, true)
	},
}

var DiscoverDirectHandler = Handler{
	ActionKey: "discover-direct",
	Action: func(args Arguments) error {
		fmt.Println("Performing direct discovery to IP ", args.TargetIP)
		return runDiscovery(args.NetInterface, []string{fmt.Sprintf("%s:%d", args.TargetIP, 3702)}, false)
	},
}

var DiscoverRangeHandler = Handler{
	ActionKey: "discover-range",
	Action: func(args Arguments) error {
		fmt.Println("Performing discovery on the following IP range ", args.TargetRange)
		prefix, err := netip.ParsePrefix(args.TargetRange)
		if err != nil {
			return fmt.Errorf("failed to parse target IP range %s: %w", args.TargetRange, err)
		}
		var ips []string
		for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
			lastDigit := addr.As4()[3]
			if (lastDigit > 1) && (lastDigit < 254) {
				ips = append(ips, fmt.Sprintf("%s:%d", addr.String(), 3702))
			}
		}
		return runDiscovery(args.NetInterface, ips, true)
	},
}

var probePayloadTemplate = `
<soapenv:Envelope
	xmlns:wsa="http://www.w3.org/2005/08/addressing"
	xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery"
	xmlns:dn="http://www.onvif.org/ver10/network/wsdl"
	xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
	<soapenv:Header>
		<wsa:Action>
			http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe
		</wsa:Action>
		<wsa:MessageID>
			urn:uuid:%s
		</wsa:MessageID>
	</soapenv:Header>
	<soapenv:Body>
		<d:Probe>
			<d:Types>dn:NetworkVideoTransmitter</d:Types>
		</d:Probe>
	</soapenv:Body>
</soapenv:Envelope>
`

type ProbeResponse struct {
	XMLName             xml.Name            `xml:"Envelope"`
	ProbeResponseHeader OnvifResponseHeader `xml:"Header"`
	ProbeMatches        []ProbeMatch        `xml:"Body>ProbeMatches>ProbeMatch"`
}

type ProbeMatch struct {
	EndpointReference ProbeEndpointReference `xml:"EndpointReference"`
	Types             string                 `xml:"Types"`
	Scopes            string                 `xml:"Scopes"`
	XAddrs            string                 `xml:"XAddrs"`
	MetadataVersion   string                 `xml:"MetadataVersion"`
}

type ProbeEndpointReference struct {
	Address string `xml:"Address"`
}

var listenConfig = net.ListenConfig{
	Control: func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			if err != nil {
				fmt.Println("Error creating socket: ", err)
			}
		})
	},
}

func listen(src *net.UDPAddr, uuid string, waitForever bool, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	b := make([]byte, 65536)
	conn, err := listenConfig.ListenPacket(context.Background(), "udp4", src.String())
	if err != nil {
		fmt.Printf("Error listening to packets: %s", err)
		return
	}
	expectedUrnUUID := fmt.Sprintf("urn:uuid:%s", uuid)
	for {
		_, addr, err := conn.ReadFrom(b)
		if err != nil {
			fmt.Printf("Failed to read bytes from local UDP addr %s: %s\n", src, err)
			continue
		}
		fmt.Println("Receiving response from", addr, ":")
		var probeResponse ProbeResponse
		err = xml.Unmarshal(b, &probeResponse)
		if err != nil {
			fmt.Printf("Failed to parse XML response: %s\n", err)
			continue
		}
		var urnUUID string
		if probeResponse.ProbeResponseHeader.RelatesTo == "" {
			urnUUID = strings.TrimSpace(probeResponse.ProbeResponseHeader.MessageID)
		} else {
			urnUUID = strings.TrimSpace(probeResponse.ProbeResponseHeader.RelatesTo)
		}
		if urnUUID != expectedUrnUUID {
			fmt.Printf("Invalid UUID %s doesn't match %s", urnUUID, expectedUrnUUID)
		} else {
			for _, match := range probeResponse.ProbeMatches {
				fmt.Println("Ref     :", match.EndpointReference.Address)
				fmt.Println("Type    :", match.Types)
				fmt.Println("Scopes  :", match.Scopes)
				fmt.Println("XAddrs  :", match.XAddrs)
				fmt.Println("Version :", match.MetadataVersion)
			}
		}
		if !waitForever {
			break
		}
	}
}

func sendMessage(src *net.UDPAddr, targetIPs []string, message string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	for _, targetIP := range targetIPs {
		fmt.Println("Sending probe to IP", targetIP)
		conn, err := listenConfig.ListenPacket(context.Background(), "udp4", src.String())
		if err != nil {
			fmt.Printf("Error listening to packets: %s\n", err)
			continue
		}
		targetUDPAddr, err := net.ResolveUDPAddr("udp4", targetIP)
		if err != nil {
			fmt.Printf("Can't resolve UDP address: %s\n", err)
			continue
		}
		_, err = conn.WriteTo([]byte(message), targetUDPAddr)
		if err != nil {
			fmt.Printf("Failed to send multicast payload: %s\n", err)
			continue
		}
		err = conn.Close()
		if err != nil {
			fmt.Printf("Failed to close probe connection: %s\n", err)
			continue
		}
		time.Sleep(3 * time.Second)
	}
}

func runDiscovery(interfaceName string, targetIPs []string, waitForever bool) error {
	localIP := util.GetLocalIP(interfaceName)
	localUdpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", localIP, 12345))
	if err != nil {
		return fmt.Errorf("can't resolve local address: %w", err)
	}
	uuid := util.GenerateUUID()
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go listen(localUdpAddr, uuid, waitForever, wg)
	go sendMessage(localUdpAddr, targetIPs, fmt.Sprintf(probePayloadTemplate, uuid), wg)
	wg.Wait()
	return nil
}
