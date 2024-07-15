# scim-patch

[![Go Reference](https://pkg.go.dev/badge/github.com/ivixvi/scim-patch.svg)](https://pkg.go.dev/github.com/ivixvi/scim-patch)

SCIM2.0 Patch 操作のGo言語実装です


> [!CAUTION]
> 安定していないため、本番利用できる状態ではありません。


# 概要

SCIM2.0のPatch操作の仕様の幅が広く、また、IdP毎の差異を吸収するのも大変です。
そこで、本ライブラリが「Patchによるスキーマの操作」を一通り吸収します。

直接アプリケーションのデータを操作しないため、総合的な処理やデータストアは冗長になると考えられます。
しかし、その代わりにSCIMのスキーマとアプリケーションで利用しているスキーマ間のマッピングを考えるだけでよくなり、密結合を下げる助けになります。

また、本ライブラリはSchemaやfilterの扱いのため、以下のSCIM関連実装に依存しています。

- https://github.com/elimity-com/scim
- https://github.com/scim2/filter-parser

### 想定される使用例

構想としては以下のIssueが近く、この例のような形で実装できることを期待しています。
https://github.com/elimity-com/scim/issues/171

現在の実装における利用例は [example](./_example/README-ja.md) をご確認ください。
