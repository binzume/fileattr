# ファイルの属性を保存・復元するやつ

[![Build Status](https://github.com/binzume/fileattr/actions/workflows/test.yaml/badge.svg)](https://github.com/binzume/fileattr/actions)

ファイルのタイムスタンプや属性をtsvに保存しておいて正確に復元するコマンドラインツール．

異なるファイルシステムを経由したり，zipで固めたり，gitでファイルを管理するときにファイルのタイムスタンプが失われたり丸められたりするのが困る場合に．

## Usage

```bash
go install github.com/binzume/fileattr@latest

fileattr -m save -l hoge.tsv ./path_to_target
fileattr -m compare -l hoge.tsv ./path_to_target
fileattr -m restore -l hoge.tsv ./path_to_target
```

tsvにはナノ秒単位のタイムスタンプが記録されますが，実際の精度はOSやファイルシステムに依存します．

## Windows

- save: 作成日時，更新日時，アクセス日時，パーミッション
- restore: 作成日時，更新日時，アクセス日時，パーミッション

## Linux

inodeのctimeは通常は書き換えられないので作成日時は復元しません．

- save: 作成日時，更新日時，アクセス日時，パーミッション
- restore: 更新日時，アクセス日時，パーミッション

## Other platforms

- save: 更新日時，パーミッション
- restore: 更新日時，アクセス日時，パーミッション
