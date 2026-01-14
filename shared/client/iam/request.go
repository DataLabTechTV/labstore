package iam

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/IllumiKnowLabs/labstore/backend/pkg/iam"
)

func (c *Client) DoRequest(op iam.IAMOp, formKeyVals ...string) (*http.Response, error) {
	if len(formKeyVals)%2 != 0 {
		return nil, errors.New("formKeyVals must contain an even number of elements")
	}

	form := url.Values{}
	form.Set("Action", string(op))

	for i := 0; i < len(formKeyVals)-1; i++ {
		form.Set(formKeyVals[i], formKeyVals[i+1])
	}

	r, err := http.NewRequestWithContext(c.Ctx, http.MethodPost, c.baseURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
