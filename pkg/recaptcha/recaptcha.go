package recaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const siteVerifyEndpoint = "https://www.google.com/recaptcha/api/siteverify"

// Verifier for recaptcha requests
type Verifier struct {
	Secret string
}

// Response is the response of the recaptcha verification
type Response struct {
	Success      bool      `json:"success"`
	ChallengedAt time.Time `json:"challenge_ts"`
	Hostname     string    `json:"hostname"`
	ErrorCodes   []string  `json:"error-codes"`
}

// Verify a given recaptcha request
func (v *Verifier) Verify(response, remoteIP string) (*Response, error) {
	resp, err := http.PostForm(
		siteVerifyEndpoint,
		url.Values{
			"secret":   {v.Secret},
			"response": {response},
			"remoteip": {remoteIP},
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	vResponse := Response{
		ErrorCodes: make([]string, 0),
	}
	if err := json.NewDecoder(resp.Body).Decode(&vResponse); err != nil {
		return nil, err
	}

	if len(vResponse.ErrorCodes) > 0 {
		return nil, errors.WithMessagef(err, "error codes: %s", strings.Join(vResponse.ErrorCodes, ", "))
	}

	return &vResponse, nil
}
