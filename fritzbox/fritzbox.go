package fritzbox

import (
	"bytes"
	"context"
	"encoding/xml"
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

func (f Fritzbox) GetCommonLinkProperties() (*CommonLinkProperties, error) {
	ctx, cancel := context.WithCancel(f.ctx)
	defer cancel()

	var result soapEnvelope
	if err := f.executeRequest(ctx, requests[getCommonLinkProperties], &result); err != nil {
		return nil, err
	}
	return result.Body.CommonLinkProperties, nil
}

func (f Fritzbox) GetTotalBytesSent() (*TotalBytesSent, error) {
	ctx, cancel := context.WithCancel(f.ctx)
	defer cancel()

	var result soapEnvelope
	if err := f.executeRequest(ctx, requests[getTotalBytesSent], &result); err != nil {
		return nil, err
	}
	return result.Body.TotalBytesSent, nil
}

func (f Fritzbox) GetTotalBytesReceived() (*TotalBytesReceived, error) {
	ctx, cancel := context.WithCancel(f.ctx)
	defer cancel()

	var result soapEnvelope
	if err := f.executeRequest(ctx, requests[getTotalBytesReceived], &result); err != nil {
		return nil, err
	}
	return result.Body.TotalBytesReceived, nil
}

func (f Fritzbox) GetTotalPacketsSent() (*TotalPacketsSent, error) {
	ctx, cancel := context.WithCancel(f.ctx)
	defer cancel()

	var result soapEnvelope
	if err := f.executeRequest(ctx, requests[getTotalPacketsSent], &result); err != nil {
		return nil, err
	}
	return result.Body.TotalPacketsSent, nil
}

func (f Fritzbox) GetTotalPacketsReceived() (*TotalPacketsReceived, error) {
	ctx, cancel := context.WithCancel(f.ctx)
	defer cancel()

	var result soapEnvelope
	if err := f.executeRequest(ctx, requests[getTotalPacketsReceived], &result); err != nil {
		return nil, err
	}
	return result.Body.TotalPacketsReceived, nil
}

func (f Fritzbox) executeRequest(ctx context.Context, params requestParams, out interface{}) error {
	reqUrl := f.buildUrl(params.path)
	data := buildRequestBody(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SoapAction", params.Urn+"#"+params.Action.String())

	resp, err := f.digestAuthClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(body) > 0 {
		if err := xml.Unmarshal(body, out); err != nil {
			return err
		}
	}
	return nil
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
