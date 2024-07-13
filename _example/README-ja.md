# サンプルサーバー

インメモリのデータストアを利用したサンプルサーバーです。
PATCHを試す上で必要な以下4エンドポイントが実装されています。

- POST /Users 
- GET /Users/{id}
- PATCH /Users/{id}
- DELETE /Users/{id}

## 使い方

```shell
$ cd path/to/scim-patch/example
$ go run .\...
```
