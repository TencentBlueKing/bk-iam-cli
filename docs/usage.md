使用文档


## 调试后台

注意, 这里 `IAM_HOST` 是权限中心后台地址

### 1. login

```bash
$ ./bk-iam-cli login http://{IAM_HOST} bk_iam {bk_iam_saas_app_secret}
INFO: success
```

### 2. check health

```bash
$ ./bk-iam-cli ping
INFO: pong

$ ./bk-iam-cli healthz
INFO: ok


$ ./bk-iam-cli version
INFO: success
{
  "buildTime": "2022-01-13_08:55:17",
  "commit": "4941e9d055b0a6bbde3e70d768b64078b818c841",
  "date": "2022-01-18T16:15:07.114089921+08:00",
  "env": "stage",
  "goVersion": "go version go1.17.3 linux/amd64",
  "timestamp": 1642493707,
  "version": "1.10.0"
}
```

### 3. query

switch to the system

```bash
$ ./bk-iam-cli use {system_id}
INFO: success
```

query system's permission model

```bash
$ ./bk-iam-cli query model
{
}
```

query system's actions

```bash
$ ./bk-iam-cli query action
{
    "actions": [],
    "pks": {},
}
```

query subject's basic info

```bash
$ ./bk-iam-cli query subject user tom
{
  "departments": [
    {
      "groups": [
        {
          "pk": 159041,
          "policy_expired_at": 1649591084
        }
      ],
      "id": "2871",
      "name": "部门1",
      "pk": 121346,
      "type": "department"
    }
  ],
  "errs": {},
  "groups": [
    {
      "pk": 168966,
      "policy_expired_at": 4102444800
    }
  ],
  "subject": {
    "id": "tom",
    "pk": 93162,
    "type": "user"
  }
}
```

query group's basic info

```bash
$ ./bk-iam-cli query subject group 2
{
  "departments": [],
  "errs": {},
  "groups": [],
  "subject": {
    "id": "2",
    "pk": 105970,
    "type": "group"
  }
}
```

query subject's policies, all actions

```bash
$ ./bk-iam-cli query policy user tom project_view
{
  "field": "project.id",
  "op": "in",
  "value": [
    "8",
    "14",
    "15",
    "16",
    "23",
    "21",
    "100133"
  ]
}
```

### 4. cache

list subject's policy in cache

```bash
$ ./bk-iam-cli cache policy user tom
{
  "actions": [
    {
      "ID": "common_flow_create",
      "PK": 18,
      "System": "bk_sops"
    },
    {
      "ID": "project_view",
      "PK": 2,
      "System": "bk_sops"
    }
  ],
  "errs": [
    null
  ],
  "keys": [
    "18",
    "2"
  ],
  "subject_pk": 86769
}
```

get the specific action's policy in cache

```bash
$ ./bk-iam-cli cache policy user tom project_view
{
  "action_pk": 2,
  "errs": [
    null,
    null
  ],
  "expressions": [],
  "notInCache": true,
  "policies": [],
  "subject_pk": 86769
}
```

get the spcific expression in cache

```bash
$ ./bk-iam-cli cache expression 11332
{
  "err": null,
  "expressions": [],
  "noCachePKs": [
    11332
  ],
  "pks": [
    11332
  ]
}
```

## 调试SaaS

### 1. login

注意, 这里 `IAM_HOST` 是权限中心SaaS的访问地址

```bash
$ ./bk-iam-cli saas login http://{IAM_SAAS_HOST} bk_iam {bk_iam_saas_app_secret}
INFO: success
```

### 2. check health

```bash
$ ./bk-iam-cli saas ping
INFO: pong
```


### 3. debug

查询20210101这一天 SaaS 的所有 debug 信息

```bash
$ ./bk-iam-cli saas debug list 20210101
[
  {
    "id": "cdc04cd3-c91f-4a98-9772-f237e313e90c",
    "type": "task",
    "name": "backend.apps.role.tasks.role_group_expire_remind",
    "exc": "",
    "stack": [],
  }
]
```

通过 request_id/task_id 查询单个请求的 Debug 信息

```bash
$ ./bk-iam-cli saas debug get 205310e3fe5548059ad386d7969b8161
{
  "id": "205310e3fe5548059ad386d7969b8161",  # request_id
  "type": "api",
  "path": "/api/v1/accounts/user/",  #  请求path
  "method": "post",  # 请求method
  "data": {},  # request data
  "exc": "",  # 异常信息
  "stack": [],  # 调用链信息
}
```


