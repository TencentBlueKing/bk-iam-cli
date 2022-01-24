/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-权限中心Cli
 * (BlueKing-IAM-Cli) available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"

	"bk-iam-cli/pkg/logger"
)

type IAMSaaSClient interface {
	Ping() error

	ListDebug(ymd string) (data []map[string]interface{}, err error)
	GetDebug(request_id string) (data map[string]interface{}, err error)
}

type iamSaaSClient struct {
	Host string

	appCode   string
	appSecret string
}

func NewIAMSaaSClient(host string, appCode string, appSecret string) IAMSaaSClient {
	return &iamSaaSClient{
		Host: host,

		appCode:   appCode,
		appSecret: appSecret,
	}
}

func (c *iamSaaSClient) call(
	method Method,
	path string,
	data interface{},
	timeout int64,
	responseData interface{},
) error {
	callTimeout := time.Duration(timeout) * time.Second
	if timeout == 0 {
		callTimeout = defaultTimeout
	}

	url := fmt.Sprintf("%s%s", c.Host, path)
	result := IAMBackendResponse{}
	start := time.Now()
	callbackFunc := NewMetricCallback("IAMBackend", start)

	// request := gorequest.New().Timeout(callTimeout).Post(url).Type("json")
	request := gorequest.New().Timeout(callTimeout).Type("json").SetBasicAuth(c.appCode, c.appSecret)
	switch method {
	case POST:
		request = request.Post(url).Send(data)
	case GET:
		request = request.Get(url).Query(data)
	}

	// do request
	resp, _, errs := request.
		EndStruct(&result, callbackFunc)

	duration := time.Since(start)

	logFailHTTPRequest(request, resp, errs, &result)

	logger.Debug("do http request: method=`%s`, url=`%s`, data=`%s`", method, url, data)
	logger.Debug("http request result: %v", result.String())
	logger.Debug("http request took %v ms", float64(duration/time.Millisecond))

	if len(errs) != 0 {
		return fmt.Errorf("gorequest errors=`%s`", errs)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gorequest statusCode is %d not 200", resp.StatusCode)
	}
	if result.Code != 0 {
		return errors.New(result.Message)
	}

	err := json.Unmarshal(result.Data, responseData)
	if err != nil {
		return fmt.Errorf("http request response body data not valid: %w, data=`%v`", err, result.Data)
	}

	return nil
}

func (c *iamSaaSClient) callWithReturnSliceMapData(
	method Method,
	path string,
	data interface{},
	timeout int64,
) ([]map[string]interface{}, error) {
	var responseData []map[string]interface{}
	err := c.call(method, path, data, timeout, &responseData)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	return responseData, nil
}

func (c *iamSaaSClient) callWithReturnMapData(
	method Method,
	path string,
	data interface{},
	timeout int64,
) (map[string]interface{}, error) {
	var responseData map[string]interface{}
	err := c.call(method, path, data, timeout, &responseData)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return responseData, nil
}

func (c *iamSaaSClient) Ping() (err error) {
	url := fmt.Sprintf("%s%s", c.Host, "/ping")

	resp, _, errs := gorequest.New().Timeout(5 * time.Second).Get(url).EndBytes()
	if len(errs) != 0 {
		return fmt.Errorf("ping fail! errs=%v", errs)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping fail! status_code=%d", resp.StatusCode)
	}
	return nil
}

func (c *iamSaaSClient) ListDebug(ymd string) (data []map[string]interface{}, err error) {
	path := "/api/v1/debug/"

	body := map[string]interface{}{
		"day": ymd,
	}
	data, err = c.callWithReturnSliceMapData(GET, path, body, 10)
	return
}

func (c *iamSaaSClient) GetDebug(request_id string) (data map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/debug/%s/", request_id)
	data, err = c.callWithReturnMapData(GET, path, map[string]interface{}{}, 10)
	return
}
