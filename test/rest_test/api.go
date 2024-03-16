package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type Result struct {
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Response map[string]interface{} `json:"response"`
}

const (
	GET    = 1
	POST   = 2
	DELETE = 3
	PUT    = 4
)

type ApiClient struct {
	Client      *resty.Client
	Address     string
	Token       string
	TokenExpire time.Time
}

func NewApiClient(address, userName, passWord string) (*ApiClient, error) { // 构造API客户端
	client := resty.New()
	//client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	param := map[string]interface{}{
		"username": userName,
		"password": passWord,
	}
	requestUrl := fmt.Sprintf("http://%s:8000/login", address)
	result := &Result{}
	resp, err := client.R().SetBody(param).
		SetResult(result).Post(requestUrl)
	if err != nil {
		logrus.Errorf("Login Error: %+v\n", err)
		return nil, err
	}
	status := resp.Result().(*Result).Code
	if status == 200 {
		result := resp.Result().(*Result).Response
		token, ok1 := result["token"].(string)
		expire, ok2 := result["expire"].(string)
		if ok1 && ok2 {
			token = "Tangula " + token
			tokenExpire, _ := time.Parse("2006-01-02 15:04:05", expire)
			return &ApiClient{client, address, token, tokenExpire}, nil
		}
	}
	err = errors.New("ERROR: login failed")
	logrus.Errorf("Error: %+v\n", err)
	return nil, err
}

func (Api *ApiClient) SendGetRequest(c *resty.Request, param map[string]string, requestUrl string) (
	body map[string]interface{}, err error) { // 发起get请求
	result := &Result{}
	resp, err := c.SetQueryParams(param).SetHeader("Authorization", Api.Token).SetResult(result).Get(requestUrl)
	if err != nil {
		logrus.Errorf("Get error %+v\n", err)
		return body, err
	}
	//logrus.Infof("Body:%+v\n", string(resp.Body()))
	status := resp.Result().(*Result).Code
	if status == 200 {
		return resp.Result().(*Result).Response, nil
	} else {
		var temp Result
		_ = json.Unmarshal(resp.Body(), &temp)
		errInfo := fmt.Sprintf("请求接口%s失败(%s)!", requestUrl, temp.Message)
		err = errors.New(errInfo)
		logrus.Errorf("Get error %+v\n", err)
		return body, err
	}
}

func (Api *ApiClient) SendPostRequest(c *resty.Request, param map[string]interface{}, requestUrl string) (
	body map[string]interface{}, err error) { // 发起post请求
	result := &Result{}
	resp, err := c.SetBody(param).
		SetHeader("Authorization", Api.Token).SetResult(result).Post(requestUrl)
	if err != nil {
		logrus.Errorf("Post error %+v\n", err)
		return body, err
	}
	//logrus.Infof("Body:%+v\n", string(resp.Body()))
	status := resp.Result().(*Result).Code
	if status == 200 {
		return resp.Result().(*Result).Response, nil
	} else {
		var temp Result
		_ = json.Unmarshal(resp.Body(), &temp)
		errInfo := fmt.Sprintf("请求接口%s失败(%s)!", requestUrl, temp.Message)
		err = errors.New(errInfo)
		logrus.Errorf("Get error %+v\n", err)
		return body, err
	}
}

func (Api *ApiClient) SendPutRequest(c *resty.Request, param map[string]interface{}, requestUrl string) (
	body map[string]interface{}, err error) { // 发起put请求
	result := &Result{}
	resp, err := c.SetBody(param).
		SetHeader("Authorization", Api.Token).SetResult(result).Put(requestUrl)
	if err != nil {
		logrus.Errorf("Post error %+v\n", err)
		return body, err
	}
	//logrus.Infof("Body:%+v\n", string(resp.Body()))
	status := resp.Result().(*Result).Code
	if status == 200 {
		return resp.Result().(*Result).Response, nil
	} else {
		var temp Result
		_ = json.Unmarshal(resp.Body(), &temp)
		errInfo := fmt.Sprintf("请求接口%s失败(%s)!", requestUrl, temp.Message)
		err = errors.New(errInfo)
		logrus.Errorf("Get error %+v\n", err)
		return body, err
	}
}

func (Api *ApiClient) SendDeleteRequest(c *resty.Request, param map[string]interface{}, requestUrl string) (
	body map[string]interface{}, err error) { // 发起delete请求
	result := &Result{}
	resp, err := c.SetBody(param).
		SetHeader("Authorization", Api.Token).SetResult(result).Delete(requestUrl)
	if err != nil {
		logrus.Errorf("Error %+v\n", err)
		return body, err
	}
	//logrus.Infof("Body:%+v\n", string(resp.Body()))
	status := resp.Result().(*Result).Code
	if status == 200 {
		return resp.Result().(*Result).Response, nil
	} else {
		var temp Result
		_ = json.Unmarshal(resp.Body(), &temp)
		errInfo := fmt.Sprintf("请求接口%s失败(%s)!", requestUrl, temp.Message)
		err = errors.New(errInfo)
		logrus.Errorf("Get error %+v\n", err)
		return body, err
	}
}

func (Api *ApiClient) SendRequest(method int, param map[string]interface{}, url string) (body map[string]interface{}, err error) {
	requestUrl := fmt.Sprintf("http://%s:8000%s", Api.Address, url)
	switch method {
	case GET:
		// 类型转换
		getParam, err := util.MapInterfaceToMapString(param)
		if err != nil {
			return body, err
		}
		return Api.SendGetRequest(Api.Client.R(), getParam, requestUrl)
	case POST:
		return Api.SendPostRequest(Api.Client.R(), param, requestUrl)
	case PUT:
		return Api.SendPutRequest(Api.Client.R(), param, requestUrl)
	case DELETE:
		return Api.SendDeleteRequest(Api.Client.R(), param, requestUrl)
	}
	return body, errors.New("support method")
}
