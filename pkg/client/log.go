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
	"net/http"

	"github.com/parnurzeal/gorequest"
	"moul.io/http2curl"

	"bk-iam-cli/pkg/logger"
)

// AsCurlCommand returns a string representing the runnable `curl' command
// version of the request.
func AsCurlCommand(request *gorequest.SuperAgent) (string, error) {
	req, err := request.MakeRequest()
	if err != nil {
		return "", err
	}

	// 脱敏, 去掉-H 中 Authorization
	req.Header.Del("Authorization")
	req.Header.Set("X-Bk-App-Secret", "*****")

	cmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return "", err
	}
	return cmd.String(), nil
}

func logFailHTTPRequest(request *gorequest.SuperAgent, response gorequest.Response, errs []error, data responseBody) {
	dump, err := AsCurlCommand(request)
	if err != nil {
		logger.Error("component request AsCurlCommand fail")
	}

	status := -1
	if response != nil {
		status = response.StatusCode
	}

	requestID := response.Header.Get("X-Request-Id")

	responseBodyError := data.Error()

	if len(errs) != 0 || response.StatusCode != http.StatusOK || responseBodyError != nil {
		message := "-"
		if responseBodyError != nil {
			message = responseBodyError.Error()
		}
		logger.Error("[http request fail] %s! status=`%d`, errs=`%v`, request_id=`%s`, request=`%s`",
			message, status, errs, requestID, dump)
	} else {
		logger.Debug("[http request] %s! status=`%d`, errs=`%v`, request_id=`%s`, request=`%s`",
			status, errs, requestID, dump)
	}
}
