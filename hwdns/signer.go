package hwdns

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const dateFormat = "20060102T150405Z"

type Signer struct {
	ak string
	sk string
}

func NewSigner(ak, sk string) *Signer {
	return &Signer{ak: ak, sk: sk}
}

func (s *Signer) Sign(req *http.Request) error {
	now := time.Now().UTC().Format(dateFormat)
	req.Header.Set("X-Sdk-Date", now)

	var bodyBytes []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		bodyBytes = b
		req.Body = io.NopCloser(strings.NewReader(string(b)))
	}
	bodyHash := hexHash(bodyBytes)
	req.Header.Set("X-Sdk-Content-Sha256", bodyHash)

	signedHeaders, canonicalHeaders := buildCanonicalHeaders(req)

	canonicalURI := req.URL.EscapedPath()
	if canonicalURI == "" {
		canonicalURI = "/"
	}
	if !strings.HasSuffix(canonicalURI, "/") {
		canonicalURI += "/"
	}
	canonicalQuery := req.URL.Query().Encode()
	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedHeaders,
		bodyHash,
	}, "\n")

	stringToSign := "SDK-HMAC-SHA256\n" + now + "\n" + hexHash([]byte(canonicalRequest))
	sig := hex.EncodeToString(hmacSHA256([]byte(s.sk), []byte(stringToSign)))

	req.Header.Set("Authorization", fmt.Sprintf(
		"SDK-HMAC-SHA256 Access=%s, SignedHeaders=%s, Signature=%s",
		s.ak, signedHeaders, sig,
	))
	return nil
}

func buildCanonicalHeaders(req *http.Request) (signedHeaders, canonical string) {
	// Huawei Cloud requires x-sdk-content-sha256 to be included in signed headers
	candidates := []string{"content-type", "host", "x-sdk-content-sha256", "x-sdk-date"}
	var present []string
	for _, k := range candidates {
		if k == "host" {
			present = append(present, k)
			continue
		}
		if req.Header.Get(k) != "" {
			present = append(present, k)
		}
	}
	sort.Strings(present)

	var sb strings.Builder
	for _, k := range present {
		var v string
		if k == "host" {
			v = req.Host
			if v == "" {
				v = req.URL.Host
			}
		} else {
			v = strings.TrimSpace(req.Header.Get(k))
		}
		sb.WriteString(k + ":" + v + "\n")
	}
	return strings.Join(present, ";"), sb.String()
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func hexHash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
