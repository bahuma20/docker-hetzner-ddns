package ip

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	"matthias-kutz.com/hetzner-ddns/pkg/request"
)

type FritzBox struct {
	IpVersion       string
	FritzBoxAddress string
}

func (i FritzBox) getHost() string {
	return i.FritzBoxAddress
}

func (i FritzBox) IsOnline() bool {
	ipifyUrl := url.URL{
		Scheme: "http",
		Host:   i.getHost(),
	}

	_, err := http.Get(ipifyUrl.String())
	return err == nil
}

func (i FritzBox) Request() (IP, error) {

	requestUrl := url.URL{
		Scheme: "http",
		Host:   i.getHost() + ":49000",
		Path:   "igdupnp/control/WANIPConn1",
	}

	body := "<?xml version=\"1.0\" encoding=\"utf-8\"?><s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"> <s:Body> <u:GetExternalIPAddress xmlns:u=\"urn:schemas-upnp-org:service:WANIPConnection:1\" /></s:Body></s:Envelope>"

	respBody, err := request.Request(http.MethodPost, requestUrl,
		map[string]string{
			"Content-Type": "text/xml; charset=\"utf-8\"",
			"SOAPAction":   "urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress",
		},
		[]byte(body))

	if err != nil {
		return IP{}, &ProviderNotAvailableError{ProviderName: i.getHost()}
	}

	r, _ := regexp.Compile("<NewExternalIPAddress>(.*)<\\/NewExternalIPAddress>")

	matches := r.FindStringSubmatch(string(respBody))

	detectedIp := matches[1]

	log.Printf("Got IP address %s from FritzBox %s", detectedIp, i.getHost())

	ip := IP{
		Value:  detectedIp,
		Source: i.getHost(),
	}

	return ip, nil
}
