// SPDX-License-Identifier: Apache-2.0

package salt

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type wrappedTokenResp struct {
	Return []struct {
		Status    string `json:"status"`
		WrapToken string `json:"wrap_token"`
	} `json:"return"`
}

type wrappedKeyPayload struct {
	Client string         `json:"client"`
	Fun    string         `json:"fun"`
	Kwarg  map[string]any `json:"kwarg"`
}

func (s Client) WrappedPrivateKey(minionID string) string {
	payload := wrappedKeyPayload{
		Client: "runner",
		Fun:    "wrapped_key_gen_accept.wrapped_key_gen_accept",
		Kwarg: map[string]any{
			"fqdn": minionID,
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	request, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewBuffer(jsonData))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return ""
	}

	createResp, err := s.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(createResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tokenResp wrappedTokenResp
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		log.Fatal(err)
	}

	return tokenResp.Return[0].WrapToken
}
