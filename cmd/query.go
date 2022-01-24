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

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [model/action/subject/policy]",
	Short: "Query data of model/action/subject/policy",
	Long: `Query data of model/action/subject/policy
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("query type required, e.g. model/action/subject/policy")
		}
		validTypes := map[string]struct{}{
			"model":   {},
			"action":  {},
			"subject": {},
			"policy":  {},
		}
		_, hit := validTypes[args[0]]
		if !hit {
			return errors.New("query type not valid, should be one of model/action/subject/policy")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 问题: 参数怎么传?
		//

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
		// query model   查询系统权限模型
		case "model":
			data, err := client.QueryModel(system)
			if err != nil {
				logger.Error("query model fail!", err.Error())
				return
			}
			logger.PrettyJson(data)
		// query action  查询系统action列表
		case "action":
			data, err := client.QueryAction(system)
			if err != nil {
				logger.Error("query action fail!", err.Error())
				return
			}
			logger.PrettyJson(data)
		// query subject 查询subject机器上级关系(部门/部门-组/组); 参数: type=user&id=x
		case "subject":
			if len(args) != 3 {
				logger.Error("query subject {subject_type} {subject_id}")
				return
			}
			_type := args[1]
			id := args[2]

			// fmt.Println("type", _type, "id", id)
			if _type != "user" && _type != "group" {
				logger.Error("query subject {subject_type} {subject_id}, subject_type should be user or group")
				return
			}
			if _, err := strconv.Atoi(id); err != nil && _type == "group" {
				logger.Error("query subject {subject_type} {subject_id}, subject_id should be an integer")
				return
			}

			data, err := client.QuerySubject(_type, id)
			if err != nil {
				logger.Error("query action fail!", err.Error())
				return
			}
			logger.PrettyJson(data)
		// query policy  查询策略; 参数: subject_type=&subject_id=&action=; 以及&force=1&debug=1
		case "policy":
			if len(args) != 4 {
				logger.Error("query policy {subject_type} {subject_id} {action}")
				return
			}
			subjectType := args[1]
			subjectID := args[2]
			action := args[3]

			// fmt.Println("type", subjectType, "id", subjectID)
			if subjectType != "user" && subjectType != "group" {
				logger.Error("query policy {subject_type} {subject_id} {action}, subject_type should be user or group")
				return
			}
			if _, err := strconv.Atoi(subjectID); err != nil && subjectType == "group" {
				logger.Error("query policy {subject_type} {subject_id} {action}, subject_id should be an integer")
				return
			}

			// TODO: debug & force 怎么搞?
			data, err := client.QueryPolicy(system, subjectType, subjectID, action, false, true)
			if err != nil {
				logger.Error("query action fail!", err.Error())
				return
			}
			logger.PrettyJson(data)
		default:
			logger.Error("not support yet")
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
