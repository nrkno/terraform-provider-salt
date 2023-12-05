// SPDX-License-Identifier: Apache-2.0

package salt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
)

type wrappedKeyData struct {
	Client string         `json:"client"`
	Fun    string         `json:"fun"`
	Kwarg  map[string]any `json:"kwarg"`
}

type wrappedTokenResp struct {
	Return []struct {
		Status    string `json:"status"`
		WrapToken string `json:"wrap_token"`
	} `json:"return"`
}

type keyDeleteData struct {
	Client string `json:"client"`
	Fun    string `json:"fun"`
	Match  string `json:"match"`
}

type keyDeleteResp struct {
	Return []struct {
		Data struct {
			Fun     string `json:"fun"`
			Success bool   `json:"success"`
		} `json:"data"`
	} `json:"return"`
}

func (s Client) WrappedPrivateKey(ctx context.Context, minionID string) (string, error) {
	payload := wrappedKeyData{
		Client: "runner",
		Fun:    "wrapped_key_gen_accept.wrapped_key_gen_accept",
		Kwarg: map[string]any{
			"fqdn": minionID,
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: unable to marshal data: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewBuffer(jsonData))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: unable to create requests: %w", err)
	}

	createResp, err := s.Do(request)
	if createResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got error from Salt API: %s", createResp.Status)
	}

	if err != nil {
		return "", fmt.Errorf("unable to do request: %w", err)
	}

	body, err := io.ReadAll(createResp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %w", err)
	}

	var tokenResp wrappedTokenResp
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Response body: %s", body))
		return "", fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return tokenResp.Return[0].WrapToken, nil
}

func (s Client) DeleteKey(ctx context.Context, minionID string) error {
	var payload []keyDeleteData
	p := keyDeleteData{
		Client: "wheel",
		Fun:    "key.delete",
		Match:  minionID,
	}
	payload = append(payload, p)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewBuffer(jsonData))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	deleteResp, err := s.Do(request)
	if deleteResp.StatusCode != http.StatusOK {
		return fmt.Errorf("got error from Salt API: %s", deleteResp.Status)
	}

	if err != nil {
		return err
	}

	body, err := io.ReadAll(deleteResp.Body)
	if err != nil {
		return err
	}

	var d keyDeleteResp
	err = json.Unmarshal(body, &d)
	if err != nil {
		return err
	}

	return nil
}
