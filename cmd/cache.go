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
	"strconv"

	"github.com/spf13/cobra"

	"bk-iam-cli/pkg/client"
	"bk-iam-cli/pkg/logger"
	"bk-iam-cli/pkg/storage"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Query policy or expression from cache",
	Long: `Query policy or expression from cache
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("query type required, e.g. policy/expression")
		}
		validTypes := map[string]struct{}{
			"policy":     {},
			"expression": {},
		}
		_, hit := validTypes[args[0]]
		if !hit {
			return errors.New("cache type not valid, should be one of policy/expression")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		system, err := storage.ReadUseSystem()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		credential := storage.NewCredential(backendCredentialFile)
		host, appCode, appSecret, err := credential.Read()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		client := client.NewIAMBackendClient(host, "", appCode, appSecret)

		switch args[0] {
		case "policy":

			// cache-query policy 查询缓存中的策略; 参数: subject_type=&subject_id=; 不带action则展示列表, 带action展示详情
			if len(args) < 3 {
				logger.Error("cache policy {subject_type} {subject_id} [{action_id}]")
				return
			}

			subjectType := args[1]
			subjectID := args[2]
			action := ""
			if len(args) == 4 {
				action = args[3]
			}
			data, err := client.QueryCachePolicy(system, subjectType, subjectID, action)
			if err != nil {
				logger.Error("cache policy fail!", err.Error())
				return
			}
			// NOTE: notInCache=false, 可能是in cache but expired
			logger.PrettyJson(data)
		case "expression":
			// cache-query expression 查询缓存中的表达式; 参数: pks=1,2,3,4
			if len(args) < 2 {
				logger.Error("cache expression {pk} required")
				return
			}

			pks := make([]int, 0, len(args[1:]))
			for _, spk := range args[1:] {
				pk, err := strconv.Atoi(spk)
				if err != nil {
					logger.Error("pk should be an integer")
					return
				}
				pks = append(pks, pk)
			}

			data, err := client.QueryCacheExpression(pks)
			if err != nil {
				logger.Error("cache expression pks fail!", err.Error())
				return
			}
			logger.PrettyJson(data)
		default:
			logger.Error("not support yet")
		}
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
