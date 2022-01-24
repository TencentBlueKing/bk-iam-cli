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

package logger

import (
	"fmt"
	"os"

	"github.com/TylerBrock/colorjson"
	"github.com/gookit/color"
)

func Debug(format string, args ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		color.Debug.Tips(format, args...)
	}
}

func Info(format string, args ...interface{}) {
	color.Info.Tips(format, args...)
}

func Warn(format string, args ...interface{}) {
	color.Warn.Tips(format, args...)
}

func Error(format string, args ...interface{}) {
	color.Error.Tips(format, args...)
}

func PrettyJson(obj interface{}) {
	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(obj)
	// return string(s)

	fmt.Println(string(s))
}
