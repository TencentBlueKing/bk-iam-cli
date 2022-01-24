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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TencentBlueKing/gopkg/conv"
	"github.com/parnurzeal/gorequest"

	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/util"
)

const (
	bkIAMVersion = "1"
)

type Method string

var (
	POST Method = "POST"
	GET  Method = "GET"
)

type IAMBackendResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (r *IAMBackendResponse) Error() error {
	if r.Code == 0 {
		return nil
	}

	return fmt.Errorf("response error[code=`%d`,  message=`%s`]", r.Code, r.Message)
}

func (r *IAMBackendResponse) String() string {
	return fmt.Sprintf("response[code=`%d`, message=`%s`, data=`%s`]",
		r.Code, r.Message, conv.BytesToString(r.Data))
}

type IAMBackendClient interface {
	Ping() error
	Healthz() (err error)
	Version() (version map[string]interface{}, err error)

	ListSystems() (data []map[string]interface{}, err error)

	QueryModel(system string) (map[string]interface{}, error)
	QueryAction(system string) (map[string]interface{}, error)
	QuerySubject(_type string, id string) (map[string]interface{}, error)
	QueryPolicy(system, subjectType, subjectID, action string, force bool, debug bool) (map[string]interface{}, error)

	QueryCachePolicy(system, subjectType, subjectID, action string) (map[string]interface{}, error)
	QueryCacheExpression(pks []int) (map[string]interface{}, error)

	PolicyGet(policyID int64) (data map[string]interface{}, err error)
	PolicyList(body interface{}) (data map[string]interface{}, err error)
	PolicySubjects(policyIDs []int64) (data []map[string]interface{}, err error)
}

type iamBackendClient struct {
	Host string

	System    string
	appCode   string
	appSecret string

	isApiDebugEnabled bool
	isApiForceEnabled bool
}

func NewIAMBackendClient(host string, system string, appCode string, appSecret string) IAMBackendClient {
	return &iamBackendClient{
		Host: host,

		System:    system,
		appCode:   appCode,
		appSecret: appSecret,

		// will add ?debug=true in url, for debug api/policy, show the details
		isApiDebugEnabled: os.Getenv("IAM_API_DEBUG") == "true" || os.Getenv("BKAPP_IAM_API_DEBUG") == "true",
		// will add ?force=true in url, for api/policy run without cache(all data from database)
		isApiForceEnabled: os.Getenv("IAM_API_FORCE") == "true" || os.Getenv("BKAPP_IAM_API_FORCE") == "true",
	}
}

func (c *iamBackendClient) call(
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

	headers := map[string]string{
		"X-BK-APP-CODE":    c.appCode,
		"X-BK-APP-SECRET":  c.appSecret,
		"X-Bk-IAM-Version": bkIAMVersion,
	}

	url := fmt.Sprintf("%s%s", c.Host, path)
	result := IAMBackendResponse{}
	start := time.Now()
	callbackFunc := NewMetricCallback("IAMBackend", start)

	// request := gorequest.New().Timeout(callTimeout).Post(url).Type("json")
	request := gorequest.New().Timeout(callTimeout).Type("json")
	switch method {
	case POST:
		request = request.Post(url).Send(data)
	case GET:
		request = request.Get(url).Query(data)
	}

	if c.isApiDebugEnabled {
		request.QueryData.Add("debug", "true")
	}
	if c.isApiForceEnabled {
		request.QueryData.Add("force", "true")
	}

	// set headers
	for key, value := range headers {
		request.Header.Set(key, value)
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

func (c *iamBackendClient) callWithReturnMapData(
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

func (c *iamBackendClient) callWithReturnSliceMapData(
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

func (c *iamBackendClient) Ping() (err error) {
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

func (c *iamBackendClient) Healthz() (err error) {
	url := fmt.Sprintf("%s%s", c.Host, "/healthz")

	resp, body, errs := gorequest.New().Timeout(10 * time.Second).Get(url).End()
	if len(errs) != 0 {
		return fmt.Errorf("healthz fail! errs=%v", errs)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthz fail! status_code=%d, body=%s", resp.StatusCode, body)
	}
	return nil
}

func (c *iamBackendClient) Version() (version map[string]interface{}, err error) {
	url := fmt.Sprintf("%s%s", c.Host, "/version")

	resp, body, errs := gorequest.New().Timeout(10 * time.Second).Get(url).EndBytes()
	if len(errs) != 0 {
		err = fmt.Errorf("version fail! errs=%v", errs)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("version fail! status_code=%d, body=%s", resp.StatusCode, body)
		return
	}

	err = json.Unmarshal(body, &version)
	if err != nil {
		err = fmt.Errorf("unmarshal version data fail! %w", err)
		return
	}
	return
}

func (c *iamBackendClient) QueryModel(system string) (map[string]interface{}, error) {
	path := "/api/v1/debug/query/model"
	body := map[string]interface{}{
		"system": system,
	}
	data, err := c.callWithReturnMapData(GET, path, body, 10)
	return data, err
}

func (c *iamBackendClient) QueryAction(system string) (map[string]interface{}, error) {
	path := "/api/v1/debug/query/action"
	body := map[string]interface{}{
		"system": system,
	}
	data, err := c.callWithReturnMapData(GET, path, body, 10)
	return data, err
}

func (c *iamBackendClient) QuerySubject(_type string, id string) (map[string]interface{}, error) {
	path := "/api/v1/debug/query/subject"
	body := map[string]interface{}{
		"type": _type,
		"id":   id,
	}
	data, err := c.callWithReturnMapData(GET, path, body, 10)
	return data, err
}

func (c *iamBackendClient) QueryPolicy(
	system, subjectType, subjectID, action string,
	force bool,
	debug bool,
) (map[string]interface{}, error) {
	path := "/api/v1/debug/query/policy"
	body := map[string]interface{}{
		"system":       system,
		"subject_type": subjectType,
		"subject_id":   subjectID,
		"action":       action,
	}
	if force {
		body["force"] = true
	}
	if debug {
		body["debug"] = true
	}

	data, err := c.callWithReturnMapData(GET, path, body, 20)
	return data, err
}

func (c *iamBackendClient) QueryCachePolicy(
	system, subjectType, subjectID, action string,
) (map[string]interface{}, error) {
	// NOTE: action can be empty
	path := "/api/v1/debug/cache/policy"
	body := map[string]interface{}{
		"system":       system,
		"subject_type": subjectType,
		"subject_id":   subjectID,
	}

	if action != "" {
		body["action"] = action
	}

	data, err := c.callWithReturnMapData(GET, path, body, 10)
	return data, err
}

func (c *iamBackendClient) QueryCacheExpression(pks []int) (map[string]interface{}, error) {
	pkList := make([]string, 0, len(pks))
	for _, pk := range pks {
		pkList = append(pkList, strconv.Itoa(pk))
	}

	path := "/api/v1/debug/cache/expression"
	body := map[string]interface{}{
		"pks": strings.Join(pkList, ","),
	}

	data, err := c.callWithReturnMapData(GET, path, body, 10)
	return data, err
}

// query system's policies, just for the system which need to use the policies to do something

func (c *iamBackendClient) PolicyGet(policyID int64) (data map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/systems/%s/policies/%d", c.System, policyID)
	data, err = c.callWithReturnMapData(GET, path, map[string]interface{}{}, 10)
	return
}

func (c *iamBackendClient) PolicyList(body interface{}) (data map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/systems/%s/policies", c.System)
	data, err = c.callWithReturnMapData(GET, path, body, 10)
	return
}

func (c *iamBackendClient) PolicySubjects(policyIDs []int64) (data []map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/systems/%s/policies/-/subjects", c.System)

	body := map[string]interface{}{
		"ids": util.Int64ArrayToString(policyIDs, ","),
	}
	data, err = c.callWithReturnSliceMapData(GET, path, body, 10)
	return
}

func (c *iamBackendClient) ListSystems() (data []map[string]interface{}, err error) {
	path := "/api/v1/web/systems"

	body := map[string]interface{}{
		"fields": "",
	}
	data, err = c.callWithReturnSliceMapData(GET, path, body, 10)
	return
}
