package fritzbox

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	. "github.com/alexruf/fritzmond/http"
)

// Fritzbox is an HTTP fritzbox for communicating with FRITZ!Box routers via TR-064 protocol.
type Fritzbox struct {
	ctx              context.Context
	digestAuthClient DigestAuthClient
	url              string
}

func New(ctx context.Context, digestAuthClient DigestAuthClient, url string) Fritzbox {
	return Fritzbox{
		ctx:              ctx,
		digestAuthClient: digestAuthClient,
		url:              url,
	}
}

func (f Fritzbox) GetCommonLinkProperties() (interface{}, error) {
	ctx, cancel := context.WithCancel(f.ctx)
	defer cancel()

	data := []byte("<?xml version=\"1.0\"?><s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><s:Body><u:GetCommonLinkProperties xmlns:u=\"urn:dslforum-org:service:WANCommonInterfaceConfig:1\"/></s:Body></s:Envelope>")
	reqUrl := f.buildUrl("/upnp/control/wancommonifconfig1")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SoapAction", "urn:dslforum-org:service:WANCommonInterfaceConfig:1#GetCommonLinkProperties")

	resp, err := f.digestAuthClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (f Fritzbox) buildUrl(path string) string {
	u, err := url.Parse(f.url)
	if err != nil {
		return ""
	}
	if strings.ToLower(u.Scheme) == "https" {
		u.Host = u.Host + ":49443"
	} else {
		u.Host = u.Host + ":49000"
	}
	u.Path = path
	return u.String()
}
