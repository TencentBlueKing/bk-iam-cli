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
	"github.com/spf13/cobra"

	"bk-iam-cli/pkg/client"
	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/storage"
)

var saasPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "call /ping to check if the iam SaaS service is alive",
	Long:  `call /ping to check if the iam SaaS service is alive.`,
	Run: func(cmd *cobra.Command, args []string) {

		credential := storage.NewCredential(saasCredentialFile)
		host, appCode, appSecret, err := credential.Read()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		client := client.NewIAMSaaSClient(host, appCode, appSecret)

		err = client.Ping()
		if err != nil {
			logger.Error("connect to host %s fail! %s", host, err.Error())
			return
		}

		logger.Info("pong")
	},
}

func init() {
	saasCmd.AddCommand(saasPingCmd)
}
