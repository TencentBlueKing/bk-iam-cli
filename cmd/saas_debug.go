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

	"github.com/spf13/cobra"

	"bk-iam-cli/pkg/client"
	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/storage"
)

var saasDebugCmd = &cobra.Command{
	Use:   "debug [list/get]",
	Short: "query debug info",
	Long: `saas debug list {day}
saas debug get {request_id/task_id}
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("debug sub-command required, e.g. list/get")
		}

		validTypes := map[string]struct{}{
			"list": {},
			"get":  {},
		}
		_, hit := validTypes[args[0]]
		if !hit {
			return errors.New("debug type not valid, should be one of list/get")
		}
		if len(args) != 2 {
			return errors.New("saas debug list {day}/saas debug get {request_id/task_id}")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		credential := storage.NewCredential(saasCredentialFile)
		host, appCode, appSecret, err := credential.Read()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		client := client.NewIAMSaaSClient(host, appCode, appSecret)

		switch args[0] {
		case "list":
			data, err := client.ListDebug(args[1])
			if err != nil {
				logger.Error("debug list fail!", err.Error())
				return
			}

			if len(data) == 0 {
				logger.Info("no debug list found!")
			} else {
				logger.PrettyJson(data)
			}
		case "get":
			data, err := client.GetDebug(args[1])
			if err != nil {
				logger.Error("debug get fail!", err.Error())
				return
			}
			logger.PrettyJson(data)
		default:
			logger.Error("not support yet")
		}
	},
}

func init() {
	saasCmd.AddCommand(saasDebugCmd)
}
