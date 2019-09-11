package gitlab_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.ozon.ru/arubtsov/ozt-token/internal/repository_client"
)

var _ repository_client.Client = new(Client)

const (
	urlTpl                = `https://%s/api/v4/projects/%d/variables`
	variableCreateBodyTpl = `key=%s&value=%s&protected=false&masked=%s&environment_scope=*`

	urlUpdateTpl          = `https://%s/api/v4/projects/%d/variables/%s`
	variableUpdateBodyTpl = `value=%s&protected=false&masked=%s&environment_scope=*`

	authHeader = `PRIVATE-TOKEN`
)

type (
	Client struct {
		host  string
		token string
	}

	ErrResponse struct {
		Message map[string][]string `json:"message"`
		Err     string              `json:"error"`
	}
)

func New(host, token string) *Client {
	return &Client{host: host, token: token}
}

func (c *Client) UpdateVariable(ctx context.Context, projectId int, v repository_client.Variable) error {
	strMasked := "false"
	if v.Masked {
		strMasked = "true"
	}
	buff := new(bytes.Buffer)
	_, err := fmt.Fprintf(buff, variableUpdateBodyTpl, v.Val, strMasked)
	if err != nil {
		return errors.Wrap(err, "cant write data to variableUpdateBodyTpl")
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(urlUpdateTpl, c.host, projectId, v.Key), buff)
	if err != nil {
		return errors.Wrap(err, "cant create UpdateVariable request")
	}
	req.Header.Set(authHeader, c.token)
	req = req.WithContext(ctx)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cant execute UpdateVariable request")
	}

	badResp := new(ErrResponse)
	if r.Body != nil {
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(badResp)
	}

	if r.StatusCode != http.StatusOK {
		errMsg := badResp.Err
		if errMsg == "" && badResp.Message != nil {
			errMsg = fmt.Sprintf("%+v", badResp.Message)
		}
		return errors.Errorf(`receive incorrect response: [%d], [%q]`, r.StatusCode, errMsg)
	}

	return nil
}

func (c *Client) Variables(ctx context.Context, projectId int) ([]repository_client.Variable, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlTpl, c.host, projectId), nil)
	if err != nil {
		return nil, errors.Wrap(err, "err on making request for get variables")
	}

	req.Header.Set(authHeader, c.token)
	req = req.WithContext(ctx)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "err on getting variables")
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	if r.StatusCode != http.StatusOK {
		badResp := new(ErrResponse)
		if r.Body != nil {
			json.NewDecoder(r.Body).Decode(badResp)
		}
		errMsg := badResp.Err
		if errMsg == "" && badResp.Message != nil {
			errMsg = fmt.Sprintf("%+v", badResp.Message)
		}
		return nil, errors.Errorf(`getting variables: receive incorrect response: [%d], [%q]`, r.StatusCode, errMsg)
	}

	result := make([]repository_client.Variable, 0)
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) CreateVariable(ctx context.Context, projectId int, v repository_client.Variable) error {
	strMasked := "false"
	if v.Masked {
		strMasked = "true"
	}
	buff := new(bytes.Buffer)
	_, err := fmt.Fprintf(buff, variableCreateBodyTpl, v.Key, v.Val, strMasked)
	if err != nil {
		return errors.Wrap(err, "cant write data to variableCreateBodyTpl")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(urlTpl, c.host, projectId), buff)
	if err != nil {
		return err
	}
	req.Header.Set(authHeader, c.token)
	req = req.WithContext(ctx)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cant execute request CreateVariable")
	}

	badResp := new(ErrResponse)
	if r.Body != nil {
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(badResp)
	}

	if r.StatusCode != http.StatusCreated {
		errMsg := badResp.Err
		if errMsg == "" && badResp.Message != nil {
			errMsg = fmt.Sprintf("%+v", badResp.Message)
		}
		return errors.Errorf(`receive incorrect response: [%d], [%q]`, r.StatusCode, errMsg)
	}

	return nil
}
