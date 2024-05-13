package util

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
)
import "encoding/hex"

var httpDigestAuthentication *HttpDigestAuthentication

type HttpDigestAuthentication struct {
	Username string
	Password string
	Realm    string
	Nonce    string
}

func (a *HttpDigestAuthentication) CalculateResponse(method string, uri string) string {
	ha1 := calculateMD5(fmt.Sprintf("%s:%s:%s", a.Username, a.Realm, a.Password))
	ha2 := calculateMD5(fmt.Sprintf("%s:%s", method, uri))
	return calculateMD5(fmt.Sprintf("%s:%s:%s", ha1, a.Nonce, ha2))
}

func (a *HttpDigestAuthentication) AuthorizationHeaderValue(method string, uri string) string {
	return fmt.Sprintf("Digest username=\"%s\", realm=\"%s\", nonce=\"%s\", uri=\"%s\", response=\"%s\"",
		a.Username, a.Realm, a.Nonce, uri, a.CalculateResponse(method, uri))
}

func calculateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func InitHttpDigestAuthentication(username string, password string, serviceUrl string) error {
	body := `
		<soapenv:Envelope
			xmlns:wsa="http://www.w3.org/2005/08/addressing"
			xmlns:trt="http://www.onvif.org/ver10/media/wsdl"	
			xmlns:soapenv="http://www.w3.org/2003/05/soap-envelope">
			<soapenv:Header>
				<wsa:Action>
					http://www.onvif.org/ver10/media/wsdl/GetProfiles	
				</wsa:Action>	
			</soapenv:Header>
			<soapenv:Body>
				<trt:GetProfiles />	
			</soapenv:Body>
		</soapenv:Envelope>
	`
	res, err := http.Post(serviceUrl, "text/xml", bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("soap request failed: %w", err)
	}
	if res.StatusCode != 401 {
		return fmt.Errorf("expected status code is 401 but found %d", res.StatusCode)
	}
	wwwAuthenticateHeader := res.Header.Get("WWW-Authenticate")
	if wwwAuthenticateHeader == "" {
		return fmt.Errorf("can't find WWW-Authenticate header")
	}
	var (
		nonce string
		realm string
	)
	for _, attr := range strings.Split(wwwAuthenticateHeader[6:], ",") {
		f := strings.Split(attr, "=")
		attrName := strings.TrimSpace(f[0])
		attrValue := strings.Trim(strings.TrimSpace(f[1]), "\"")
		if attrName == "nonce" {
			nonce = attrValue
		} else if attrName == "realm" {
			realm = attrValue
		}
	}
	if nonce == "" {
		return fmt.Errorf("can't parse nonce from %s", wwwAuthenticateHeader)
	}
	httpDigestAuthentication = &HttpDigestAuthentication{
		Username: username,
		Password: password,
		Realm:    realm,
		Nonce:    nonce,
	}
	return nil
}
