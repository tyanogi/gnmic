# Track Specification: Refactor Suggestion Engine to Config-based Fetching

## 1. Overview
gNMIc プロンプトモードにおける YANG スキーマの全解析による自動リスト抽出を廃止し、ユーザーが指定した設定ファイルに基づいた明示的な事前取得（Prefetch）にリファクタリングします。これにより、複雑な解析ロジックによる不安定さを排除し、確実に必要な補完候補を取得できる仕組みを提供します。

## 2. Functional Requirements
- **構造化設定ファイルベースの取得**:
    - 専用の設定ファイルから取得対象を読み込みます。形式は以下の通りとします：
      ```yaml
      suggestions:
        - path: "/interfaces/interface"
          key: "name"
          state-path: "/interfaces/interface[name=*]/state/name"
      ```
- **CLI フラグの追加**:
    - 新しいフラグ `--suggestions-file` を追加。このフラグが指定された場合のみ設定ファイルを読み込みます。
- **Suggestion Engine のリファクタリング**:
    - `SchemaWalker` による YANG 解析ロジックを完全に削除します。
    - 起動時に指定されたファイルから `SuggestionConfig` を読み込み、`SuggestionEngine` に渡して非同期取得を開始します。
- **補完ロジックとの統合**:
    - 取得された値は、指定された `path` と `key` を組みわせたキャッシュキー（例: `/interfaces/interface::name`）で保存され、既存の `findMatchedXPATH` 関数で利用されます。

## 3. Non-Functional Requirements
- **保守性**: 自動探索コードの削除により `pkg/app` 以下のコードを軽量化します。
- **信頼性**: ユーザーがデバイスの仕様に合わせて最適なクエリパスを記述できるため、NotFound エラーなどを回避できます。

## 4. Acceptance Criteria
- [ ] `SchemaWalker` 関連のコードが完全に削除されている。
- [ ] `--suggestions-file` で指定した設定ファイルに基づいて gNMI Get が発行される。
- [ ] 取得された値がキャッシュに保存され、プロンプト入力時にサジェスト（補完）される。
- [ ] 設定ファイルが指定されない、またはファイルが存在しない場合にエラーにならず通常通りプロンプトが起動すること。

## 5. Out of Scope
- YANG スキーマからの自動リスト探索。
- 複数のパスを1つの Get リクエストにまとめる最適化。
