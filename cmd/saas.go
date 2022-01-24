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
	"fmt"

	"github.com/spf13/cobra"
)

// saasCmd represents the saasCmd command
var saasCmd = &cobra.Command{
	Use:   "saas",
	Short: "saas debug api",
	Long: `saas login {host} {app_code} {app_secret}
saas ping
saas debug list {20210501}
saas debug get {request_id/task_id}`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`saas login {host} {app_code} {app_secret}
saas debug list {20210501}
saas debug request {request_id}`)
	},
}

func init() {
	rootCmd.AddCommand(saasCmd)
}
