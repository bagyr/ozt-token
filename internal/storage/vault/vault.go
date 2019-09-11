package vault

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"

	"gitlab.ozon.ru/arubtsov/ozt-token/internal/storage"
	"gitlab.ozon.ru/platform/errors"
)

const tokenCreateURL = "/v1/auth/token/create/"

type (
	Client struct {
		api *api.Client
	}

	requestMeta struct {
		ProjectId int `json:"ProjectId"`
	}
	requestCreateToken struct {
		Meta      requestMeta `json:"meta"`
		TTLSec    int         `json:"ttl"`
		Renewable bool        `json:"renewable"`
	}
)

func New(api *api.Client) *Client {
	return &Client{api: api}
}

func (c *Client) NewAuthToken(ctx context.Context, projectId int, role string, ttl time.Duration) (t storage.Token, err error) {
	r := requestCreateToken{
		Meta:   requestMeta{ProjectId: projectId},
		TTLSec: int(ttl.Seconds()),
	}

	req := c.api.NewRequest(http.MethodPost, tokenCreateURL+role)

	err = req.SetJSONBody(r)
	if err != nil {
		return
	}

	resp, err := c.api.RawRequestWithContext(ctx, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	var respBody *api.Secret
	if resp.Body != nil {
		respBody, err = api.ParseSecret(resp.Body)
		if err != nil {
			return
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("receive incorrect response status: [%d], body: [%+v]", resp.StatusCode, respBody)
		return
	}

	t.Secret = respBody.Auth.ClientToken
	t.TTL = time.Second * time.Duration(respBody.Auth.LeaseDuration)
	return
}
