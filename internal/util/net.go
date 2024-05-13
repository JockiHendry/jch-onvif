package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Code    string   `xml:"Body>Fault>Code>Value"`
	Subcode string   `xml:"Body>Fault>Subcode>Value"`
	Reason  string   `xml:"Body>Fault>Reason>Text"`
}

func (e ErrorResponse) Error() string {
	var str strings.Builder
	str.WriteString("SOAP Error")
	if e.Code != "" {
		str.WriteString(fmt.Sprintf(" Code [%s]", e.Code))
	}
	if e.Subcode != "" {
		str.WriteString(fmt.Sprintf(" Subcode [%s]", e.Subcode))
	}
	if e.Reason != "" {
		str.WriteString(fmt.Sprintf(" Reason [%s]", e.Reason))
	}
	return str.String()
}

func SoapCall(interfaceName string, serviceUrl string, body string) ([]byte, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", GetLocalIP(interfaceName).String(), 0))
	if err != nil {
		fmt.Println("Failed to resolve local address: ", err)
		return nil, err
	}
	dialer := &net.Dialer{LocalAddr: addr}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}
	req, _ := http.NewRequest("POST", serviceUrl, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "text/xml")
	if httpDigestAuthentication != nil {
		req.Header.Set("Authorization", httpDigestAuthentication.AuthorizationHeaderValue("POST", req.URL.Path))
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("soap request failed: %w", err)
	}
	fmt.Println("Response Status:", response.Status)
	content, _ := io.ReadAll(response.Body)
	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing response: %w", err)
	}
	if response.StatusCode != 200 {
		var errorResponse ErrorResponse
		err = xml.Unmarshal(content, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("problem while parsing error response: %w", err)
		}
		return nil, errorResponse
	}
	return content, nil
}

func GetDefaultNetworkInterface() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Failed to retrieve default network interface:", err)
		return "lo"
	}
	if len(interfaces) == 0 {
		fmt.Println("Can't find default network interface")
		return "lo"
	}
	return interfaces[0].Name
}

func GetLocalIP(interfaceName string) net.IP {
	networkInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println("Failed to retrieve network interface ", interfaceName, ":", err)
		return nil
	}
	networkAddresses, err := networkInterface.Addrs()
	if err != nil {
		fmt.Println("Failed to retrieve network interface addresses ", interfaceName, ":", err)
		return nil
	}
	for _, addr := range networkAddresses {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			return addr.(*net.IPNet).IP
		}
	}
	return nil
}
