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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TencentBlueKing/gopkg/conv"
	"github.com/TencentBlueKing/gopkg/cryptography"
)

// TODO: encrypted credential

const (
	cryptoKey   = "C4QSNKR4GNPIZAH3B0RPWAIV29E7QZ66"
	aesGcmNonce = "KC9DvYrNGnPW"
)

func newCrypto() (cryptography.Crypto, error) {
	c, err := cryptography.NewAESGcm([]byte(cryptoKey), []byte(aesGcmNonce))
	if err != nil {
		return nil, fmt.Errorf("cryptos key error: %w", err)
	}
	return c, nil
}

func encryptToBase64(c cryptography.Crypto, plaintext string) string {
	encryptedText := c.Encrypt(conv.StringToBytes(plaintext))
	return base64.StdEncoding.EncodeToString(encryptedText)
}

func decryptFromBase64(c cryptography.Crypto, encryptedTextB64 string) (plainText string, err error) {
	var encryptedText []byte
	encryptedText, err = base64.StdEncoding.DecodeString(encryptedTextB64)
	if err != nil {
		return
	}

	var plaintextBytes []byte
	plaintextBytes, err = c.Decrypt(encryptedText)
	if err != nil {
		return
	}

	return conv.BytesToString(plaintextBytes), err
}

func encryptCredential(host, appCode, appSecret string, expire int64) (string, error) {
	plain := fmt.Sprintf("%s,%s,%s,%d", host, appCode, appSecret, expire)
	c, err := newCrypto()
	if err != nil {
		return "", err
	}
	return encryptToBase64(c, plain), nil
}

func decryptCredential(cs string) (host, appCode, appSecret string, expiration int64, err error) {
	var c cryptography.Crypto
	c, err = newCrypto()
	if err != nil {
		return
	}

	var credential string
	credential, err = decryptFromBase64(c, cs)
	if err != nil {
		return
	}

	parts := strings.Split(credential, ",")
	if len(parts) != 4 {
		err = fmt.Errorf("invalid credential")
		return
	}

	host = parts[0]
	appCode = parts[1]
	appSecret = parts[2]

	expiration, err = strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid expiration")
		return
	}
	return
}

type Credential struct {
	file string
}

func NewCredential(file string) *Credential {
	return &Credential{
		file: file,
	}
}

func (c *Credential) Write(host, appCode, appSecret string) error {
	f, err := os.Create(c.file)
	if err != nil {
		return fmt.Errorf("create credential file fail! %w", err)
	}
	defer f.Close()

	ts := time.Now().Unix()
	expire := ts + 60*60

	cs, err := encryptCredential(host, appCode, appSecret, expire)
	if err != nil {
		return fmt.Errorf("encrypt credential fail! %w", err)
	}
	_, err = f.WriteString(cs)
	if err != nil {
		return fmt.Errorf("write credential fail! %w", err)
	}

	return nil
}

func (c *Credential) Read() (host, appCode, appSecret string, err error) {
	if _, err = os.Stat(c.file); os.IsNotExist(err) {
		err = fmt.Errorf("please login first")
		return
	}

	dat, err := ioutil.ReadFile(c.file)
	if err != nil {
		err = fmt.Errorf("read credential fail! %w", err)
		return
	}

	var expiration int64
	host, appCode, appSecret, expiration, err = decryptCredential(string(dat))
	if time.Now().Unix() > expiration {
		err = fmt.Errorf("credential expires, please login again")
		return
	}

	return
}

func (c *Credential) RemoveCredential() error {
	return os.Remove(c.file)
}
