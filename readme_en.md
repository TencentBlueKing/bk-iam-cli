![](docs/resource/img/bk_iam_en.png)
---

[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/TencentBlueKing/bk-iam-cli/blob/master/LICENSE.txt) [![Release Version](https://img.shields.io/badge/release-1.0.0-brightgreen.svg)](https://github.com/TencentBlueKing/bk-iam-cli/releases) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TencentBlueKing/bk-iam-cli/pulls)

## Overview

BK-IAM-CLI is an command line tool for debug, it is used for debugging and analysis, and can get the model, policy, expression and cache data from the backend.

## Getting started

### Installation

```bash
$ git clone git@github.com:TencentBlueKing/bk-iam-cli.git
$ cd bk-iam-cli

# go 1.17 required
$ go version                                                                                                                                                  146 ↵ wukunliang@wklken-MacBook-Pro
go version go1.17.5 darwin/arm64

$ make dep
go mod tidy
go mod vendor

$ make build
go build -mod=vendor .

$ ls bk-iam-cli
bk-iam-cli*
```

### Usage

- [usage doc](docs/usage.md)

## Roadmap

- [release log](release.md)

## IAM Repos

- [TencentBlueKing/bk-iam](https://github.com/TencentBlueKing/bk-iam)
- [TencentBlueKing/bk-iam-saas](https://github.com/TencentBlueKing/bk-iam-saas)
- [TencentBlueKing/bk-iam-search-engine](https://github.com/TencentBlueKing/bk-iam-search-engine)
- [TencentBlueKing/bk-iam-cli](https://github.com/TencentBlueKing/bk-iam-cli)
- [TencentBlueKing/iam-python-sdk](https://github.com/TencentBlueKing/iam-python-sdk)
- [TencentBlueKing/iam-go-sdk](https://github.com/TencentBlueKing/iam-go-sdk)
- [TencentBlueKing/iam-php-sdk](https://github.com/TencentBlueKing/iam-php-sdk)

## Support

- [bk forum](https://bk.tencent.com/s-mart/community)
- [bk DevOps online video tutorial(In Chinese)](https://cloud.tencent.com/developer/edu/major-100008)
- Contact us, technical exchange QQ group:

<img src="https://github.com/Tencent/bk-PaaS/raw/master/docs/resource/img/bk_qq_group.png" width="250" hegiht="250" align=center />


## BlueKing Community

- [BK-CI](https://github.com/Tencent/bk-ci)：a continuous integration and continuous delivery system that can easily present your R & D process to you.
- [BK-BCS](https://github.com/Tencent/bk-bcs)：a basic container service platform which provides orchestration and management for micro-service business.
- [BK-BCS-SaaS](https://github.com/Tencent/bk-bcs-saas)：a SaaS provides users with highly scalable, flexible and easy-to-use container products and services.
- [BK-PaaS](https://github.com/Tencent/bk-PaaS)：an development platform that allows developers to create, develop, deploy and manage SaaS applications easily and quickly.
- [BK-SOPS](https://github.com/Tencent/bk-sops)：an lightweight scheduling SaaS  for task flow scheduling and execution through a visual graphical interface. 
- [BK-CMDB](https://github.com/Tencent/bk-cmdb)：an enterprise-level configuration management platform for assets and applications.

## Contributing

If you have good ideas or suggestions, please let us know by Issues or Pull Requests and contribute to the Blue Whale Open Source Community.

## License

Based on the MIT protocol. Please refer to [LICENSE](LICENSE.txt)
