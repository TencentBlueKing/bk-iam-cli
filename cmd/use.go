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

	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/storage"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [system_id]",
	Short: "The current system you want to query",
	Long: `The current system you want to query, like "use bk_paas"
After that, the all query commands will query the data of system bk_paas
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 切换到哪个系统, 例如use bk_paas, 当前session中system切换到bk_paas
		system := args[0]

		err := storage.WriteUseSystem(system)
		if err != nil {
			logger.Error("Use system fail: %s", err.Error())
			return
		}

		logger.Info("success")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
