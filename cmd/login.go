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

package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"bk-iam-cli/pkg/client"
	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/storage"
)

const backendCredentialFile = ".credential"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via app_code/app_secret of IAM",
	Long: `Login via app_code/app_secret of IAM. 
The login credentials will be encrypted and store at current dir.
And you should login every 1 hour.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("login http://{iam_host} {app_code} {app_secret}")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// NOTE:
		// 执行login, 传入host地址/app_code/app_secret, 调用后台ping(可达), 并调用一个接口确认app_code/app_secret正确性;
		// 如果正确, 将信息加密保存到本地, 2h过期
		// 过期后需要重新登录
		// 如果不正确, 提示用户;

		// TODO: 可以执行clear, 清理掉login状态
		// ? 问题: 怎么同时login到后台和saas

		host := args[0]
		appCode := args[1]
		appSecret := args[2]

		if !strings.HasPrefix(host, "http://") {
			host = fmt.Sprintf("http://%s", host)
		}

		// 1. host is connectable : /ping
		client := client.NewIAMBackendClient(host, "", appCode, appSecret)

		err := client.Ping()
		if err != nil {
			logger.Error("connect to host %s fail! %s\n", host, err.Error())
			return
		}

		// 2. the app_code/app_secret is valid: /api/v1/web/systems
		_, err = client.ListSystems()
		if err != nil {
			logger.Error("app_code or app_secret invalid")
			return
		}

		// 3. create the credential
		credential := storage.NewCredential(backendCredentialFile)
		err = credential.Write(host, appCode, appSecret)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		// 4. success
		logger.Info("success")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
