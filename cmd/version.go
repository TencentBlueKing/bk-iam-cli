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

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "call /version to check the version of iam backend",
	Long:  `call /version to check the version of iam backend`,
	Run: func(cmd *cobra.Command, args []string) {
		credential := storage.NewCredential(backendCredentialFile)
		host, appCode, appSecret, err := credential.Read()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		client := client.NewIAMBackendClient(host, "", appCode, appSecret)

		version, err := client.Version()
		if err != nil {
			logger.Error("version check fail: %s", host, err.Error())
			return
		}

		logger.Info("success")
		logger.PrettyJson(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
