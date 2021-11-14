package http

import (
	"crypto/md5" //nolint:gosec
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DigestAuthClient struct {
	httpClient *http.Client
	username   string
	password   string
}

func NewDigestAuthClient(httpClient *http.Client, username string, password string) DigestAuthClient {
	return DigestAuthClient{
		httpClient: httpClient,
		username:   username,
		password:   password,
	}
}

func (d DigestAuthClient) Do(req *http.Request) (*http.Response, error) {
	authorizedReq := req.Clone(req.Context())
	authorizedReq.Body, _ = req.GetBody()
	defer authorizedReq.Body.Close()

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}
	defer resp.Body.Close()

	digestParts := getDigestParts(resp)
	digestParts["uri"] = req.RequestURI
	digestParts["method"] = req.Method
	digestParts["username"] = d.username
	digestParts["password"] = d.password
	authorizedReq.Header.Set("Authorization", getDigestAuthorization(digestParts))
	return d.httpClient.Do(authorizedReq)
}

func getDigestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hash := md5.New() //nolint:gosec
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	_, _ = io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthorization(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}
