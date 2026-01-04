# Implementation Plan - Implement YANG Schema Walker (Track 1)

## Phase 1: 環境構築と型定義
- [x] 型定義の作成: `pkg/app/schema_walker.go` に `ListInfo` 構造体とインターフェースを定義する 2e027cc
- [x] テスト用ボイラープレートの作成: `pkg/app/schema_walker_test.go` を作成し、goyangを使ったモック用YANGエントリ生成の補助関数を用意する 3efad23
- [x] Task: Conductor - User Manual Verification 'Phase 1: 環境構築と型定義' (Protocol in workflow.md) [checkpoint: e634079]

## Phase 2: Schema Walker の核となるロジックの実装 [checkpoint: 26af52f]
- [x] TDD: 単純なリスト（シングルキー）から情報を抽出するテストを書く 820f433
- [x] 実装: `yang.Entry` を再帰的に探索し、リストノードを特定する基本ロジックを実装する 820f433
- [x] TDD: 複合キーを持つリストから複数の `ListInfo` に展開するテストを書く 820f433
- [x] 実装: 複合キーの各要素を個別の `ListInfo` に展開するロジックを実装する 820f433
- [x] Task: Conductor - User Manual Verification 'Phase 2: Schema Walker の核となるロジックの実装' (Protocol in workflow.md) 26af52f

## Phase 3: OpenConfig State Path 優先ロジックの追加
- [ ] TDD: OpenConfigスタイルの `state` コンテナを持つリストのテストケースを追加する
- [ ] 実装: ヒューリスティックに基づいた実機問い合わせ用パス（State Path）の生成ロジックを実装する
- [ ] 最終検証: 大規模な（あるいは複雑なネストを持つ）スキーマツリーでの動作確認テストを実施する
- [ ] Task: Conductor - User Manual Verification 'Phase 3: OpenConfig State Path 優先ロジックの追加' (Protocol in workflow.md)

## Phase 4: 仕上げとコードクオリティ確認
- [ ] Lint/Type Check: `golangci-lint` を実行し、コードの整合性を確認する
- [ ] ドキュメント: 実装した関数にGoDocコメントを追記する
- [ ] Task: Conductor - User Manual Verification 'Phase 4: 仕上げとコードクオリティ確認' (Protocol in workflow.md)
