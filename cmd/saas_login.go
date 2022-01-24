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
	"time"

	"github.com/spf13/cobra"

	"bk-iam-cli/pkg/client"
	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/storage"
)

const saasCredentialFile = ".saas-credential"

var saasLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via app_code/app_secret of IAM SaaS",
	Long: `Login via app_code/app_secret of IAM SaaS. 
The login credentials will be encrypted and store at current dir.
And you should login every 1 hour.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("saas login http://{iam_saas_host} {app_code} {app_secret}")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		appCode := args[1]
		appSecret := args[2]

		// 1. host is connectable : /ping
		client := client.NewIAMSaaSClient(host, appCode, appSecret)

		err := client.Ping()
		if err != nil {
			logger.Error("connect to host %s fail! %s\n", host, err.Error())
			return
		}

		// 2. the app_code/app_secret is valid: /api/v1/web/systems
		day := time.Now().Format("20210101")
		_, err = client.ListDebug(day)
		if err != nil {
			logger.Error("app_code or app_secret invalid")
			return
		}

		// 3. create the credential
		credential := storage.NewCredential(saasCredentialFile)
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
	saasCmd.AddCommand(saasLoginCmd)
}
