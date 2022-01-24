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

package storage

import (
	"fmt"
	"io/ioutil"
	"os"
)

func WriteUseSystem(system string) error {
	f, err := os.Create(".use")
	if err != nil {
		return fmt.Errorf("create storage file fail! %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(system)
	if err != nil {
		return fmt.Errorf("write storage fail! %w", err)
	}

	return nil
}

func ReadUseSystem() (system string, err error) {
	if _, err = os.Stat(".use"); os.IsNotExist(err) {
		err = fmt.Errorf("please use system first")
		return
	}

	dat, err := ioutil.ReadFile(".use")
	if err != nil {
		err = fmt.Errorf("read use system fail! %w", err)
		return
	}
	return string(dat), nil
}

func RemoveUseSystem() error {
	return os.Remove(".use")
}
