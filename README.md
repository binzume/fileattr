# ファイルの属性を保存・復元するやつ

ファイルのタイムスタンプや属性をtsvに保存しておいて正確に復元するためのやつ．

異なるファイルシステムを経由したり，zipで固めたり，gitでファイルを管理するときにファイルのタイムスタンプが失われたり丸められたりするのが困る場合に．

## Usage

T.B.D.

```bash
go install github.com/binzume/fileattr

fileattr -m save -l hoge.tsv ./path_to_target
fileattr -m compare -l hoge.tsv ./path_to_target
fileattr -m restore -l hoge.tsv ./path_to_target
```

タイムスタンプはtsvにはナノ秒単位で記録されますが，実際の精度はOSやファイルシステムに依存します．

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
