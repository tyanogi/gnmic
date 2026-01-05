# Implementation Plan - Refactor Suggestion Engine to Config-based Fetching

## Phase 1: 設定ファイルの定義と読み込み
- [x] 型定義: `pkg/app/suggestion_config.go` を作成し、設定ファイルの構造体 (`SuggestionConfig`, `SuggestionItem`) を定義する 92576b5
- [x] 実装: 設定ファイルを読み込み、Goの構造体にパースするローダー関数を実装する 92576b5
- [x] テスト: サンプル YAML を読み込み、正しくパースできることを検証するユニットテストを作成する 92576b5

## Phase 2: 自動解析の廃止と SuggestionEngine の改修
- [x] リファクタリング: `pkg/app/schema_walker.go` およびそのテスト、関連する呼び出し箇所を削除する 92576b5
- [x] インターフェース変更: `SuggestionEngine.Start` メソッドが `ListInfo` (自動抽出) ではなく、新しい `SuggestionItem` (設定ファイル由来) を受け取るように変更する 92576b5
- [x] フラグ追加: `pkg/cmd/prompt.go` に `--suggestions-file` フラグを追加し、`App` 構造体で受け取れるようにする 92576b5
- [x] 統合: `App.StartSuggestionEngine` でフラグの有無を確認し、ファイルがある場合のみ読み込んでエンジンを開始するロジックに変更する 92576b5
- [x] Task: Conductor - User Manual Verification 'Phase 2: 自動解析の廃止と SuggestionEngine の改修' (Protocol in workflow.md) 92576b5

## Phase 3: 動作検証とクリーンアップ
- [x] 検証準備: 手動テスト用の `suggestions.yaml` サンプルファイルを作成する 92576b5
- [x] 手動検証: 作成した設定ファイルを使用して `gnmic prompt` を起動し、指定したパスの補完候補が実際に表示されることを確認する 95ac0df
- [x] クリーンアップ: 不要になった `ListInfo` 定義や古いロジックが残っていないか最終確認する 95ac0df
- [x] Task: Conductor - User Manual Verification 'Phase 3: 動作検証とクリーンアップ' (Protocol in workflow.md) 95ac0df
