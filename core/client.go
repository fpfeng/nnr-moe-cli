package core

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
)

type URL string

const (
	URLServers    URL = "https://nnr.moe/api/servers"
	URLRules      URL = "https://nnr.moe/api/rules"
	URLAddRule    URL = "https://nnr.moe/api/rules/add"
	URLDeleteRule URL = "https://nnr.moe/api/rules/del"
	URLEditRule   URL = "https://nnr.moe/api/rules/edit"
)

type NNRMoe struct {
	Token       string
	HttpRequest *resty.Request
}

func New(token string) *NNRMoe {
	return &NNRMoe{
		Token: token,
		HttpRequest: resty.New().R().
			SetHeader("content-type", "application/json").
			SetHeader("token", token),
	}
}

func sendRequest[TB TypeRequestBody, TR TypeResponse](request *resty.Request, url URL, body TB, result TR) error {
	resp, err := request.SetBody(body).Post(string(url))
	if err != nil {
		log.Fatalf("Call %v failed with error: %v", url, err)
		return err
	}
	json.Unmarshal(resp.Body(), &result)
	if !result.IsStatusOK() {
		rbody, _ := json.Marshal(body)
		log.Fatalf("Call %v failed, try execute again with non-china proxy if cloudflare protect was in response.\nRequest: %v\nResponse: %v", url, string(rbody), resp.String())
		return errors.New("Call api failed.")
	}
	return nil
}

func (nnrMoe *NNRMoe) Servers() (result *ResponseServers) {
	result = &ResponseServers{
		ResponseBase: ResponseBase{Status: 0},
		Data:         make(ServerList, 0),
	}
	sendRequest(nnrMoe.HttpRequest, URLServers, "", result)
	return
}

func (nnrMoe *NNRMoe) Rules() (result *ResponseRules) {
	result = &ResponseRules{
		ResponseBase: ResponseBase{Status: 0},
		Data:         make(RuleList, 0),
	}
	sendRequest(nnrMoe.HttpRequest, URLRules, "", result)
	return
}

func (nnrMoe *NNRMoe) DeleteRule(deleteRule *RequestDeleteRule) (result *ResponseDeleteRule) {
	result = &ResponseDeleteRule{
		ResponseBase: ResponseBase{Status: 0},
	}
	sendRequest(nnrMoe.HttpRequest, URLDeleteRule, deleteRule, result)
	return
}

func (nnrMoe *NNRMoe) EditRule(editedRule *RequestEditedRule) (result *ResponseAddOrEditRule) {
	result = &ResponseAddOrEditRule{
		ResponseBase: ResponseBase{Status: 0},
	}
	sendRequest(nnrMoe.HttpRequest, URLEditRule, editedRule, result)
	return
}

func (nnrMoe *NNRMoe) AddRule(rule *RequestAddRule) (result *ResponseAddOrEditRule) {
	result = &ResponseAddOrEditRule{
		ResponseBase: ResponseBase{Status: 0},
	}
	sendRequest(nnrMoe.HttpRequest, URLAddRule, rule, result)
	return
}
