---
paths:
  - ".github/workflows/**"
---

# リリースワークフロー

## 全体の流れ

リリースは [tagpr](https://github.com/Songmu/tagpr) によって自動化されている。

```
[main] PR がマージされる
  → tagpr が tagpr-from-v* ブランチを force-push で作成・更新
  → [tagpr-from-v*] build-release.yml が Docker image をビルドして action.yml の digest を更新・コミット
  → [PR: tagpr-from-v* → main] verify-release-image.yml が action.yml の digest を検証
  → PR をマージ
  → [main] tagpr がマージコミットに vX.Y.Z タグを打つ
  → [tag: vX.Y.Z] release.yml が GitHub Release を作成
```

## ワークフローファイル一覧

| ファイル | トリガー | 役割 |
|---------|---------|------|
| `build-release.yml` | `tagpr-from-*` へのプッシュ | Docker image をビルドし、digest を action.yml に書き込む |
| `verify-release-image.yml` | main への pull_request | action.yml の digest と実際のイメージが一致するか検証 |
| `release.yml` | タグ `v*.*.*` のプッシュ | GitHub Release の作成、メジャーバージョンタグのプッシュ |

## build-release の無限ループ防止

build-release は action.yml を更新してコミットするため、それ自体が push イベントをトリガーする。job-level の `if` 条件で無限ループを防いでいる：

```yaml
if: "github.event.forced || !startsWith(github.event.head_commit.message, 'chore: pin Docker image digest')"
```

- `github.event.forced == true` → tagpr の force-push → **実行**（新しいコードが入ったので再ビルドが必要）
- `github.event.forced == false` かつメッセージが `chore: pin Docker image digest` → build-release 自身のコミット → **スキップ**

注意：Docker image の digest はビルドのたびにタイムスタンプで変わるため、digest の比較による冪等チェックは機能しない。

## 鶏と卵問題：build-release が使うイメージ

`build-release.yml` の "Create GitHub App Token" ステップは `uses: ./` でブランチ上の `action.yml` を参照する。つまり：

- トークン取得に使われるイメージは**常に一世代前**（前回のビルドで書き込まれた digest のイメージ）
- 新しいイメージはその後の "Build and Push Docker Image" ステップで初めてビルドされる

これは仕様。新しい機能が main にマージされた後の最初の build-release run では、トークン取得は旧イメージで行われる。

## イメージの中身を調べる方法

digest からソースコミットを特定するには：

1. そのイメージをビルドした `build-release` の成功 run を探す：`gh run list --json workflowName,conclusion,headSha,url`
2. `headSha` がビルド時の `tagpr-from-*` ブランチの HEAD コミット
3. `git log --oneline <sha> --first-parent` で含まれる main コミットを確認できる

## main ブランチの Required status checks

Ruleset "Protect main" で以下のジョブが Required checks として設定されている：

| コンテキスト名 | ワークフロー | 役割 |
|--------------|------------|------|
| `test` | `test.yml` | Go ユニットテスト・ビルド確認 |
| `actionlint` | `ghalint.yml` | GitHub Actions ワークフローの静的解析 |
| `ghalint` | `ghalint.yml` | GitHub Actions セキュリティ lint |
| `verify` | `verify-release-image.yml` | action.yml の Docker image digest 検証 |

### paths フィルタの注意事項

trigger レベルの `paths` / `paths-ignore` は使わないこと。パスが一致しないときワークフロー自体が起動せず、required check が PR に現れなくなりマージ不可になる。

paths フィルタが必要な場合は `dorny/paths-filter` を `changes` ジョブで使い、後続ジョブを `if` 条件で skip する方式を使う。GitHub は `if` で skip されたジョブを "Skipped"（＝ passing）として扱い、required check の要件を満たす。
