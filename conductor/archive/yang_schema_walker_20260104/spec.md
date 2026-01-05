# Track Specification: Implement YANG Schema Walker (Track 1)

## 1. Overview
gNMIcのプロンプト補完機能を改善するための基盤として、YANGスキーマを探索（Walk）し、動的補完に必要なリストノードの情報を抽出する `SchemaWalker` モジュールを実装します。

## 2. Functional Requirements
- **解析対象**: `yang.Entry` 形式のパース済みスキーマツリーを受け取ります。
- **リスト抽出**: スキーマツリーを再帰的に探索し、`list` ステートメントを持つ全ノードを特定します。
- **情報抽出 (ListInfo)**:
    - リストの絶対パス（例: `/interfaces/interface`）
    - キー名（例: `name`）
    - 実機問い合わせ用パス（State Path）:
        - 原則としてキーの絶対パスを使用。
        - **OpenConfig対応 (Heuristic)**: リスト直下に `state` コンテナがあり、その中にキー名と同名のLeafが存在する場合、それを優先します（例: `/interfaces/interface/state/name`）。
- **データ構造**:
    - 複合キー（Multiple Keys）を持つリストの場合、キーごとに独立した `ListInfo` 構造体を生成し、スライスとして返します。
- **インターフェース**:
    - すでにパース済みの `yang.Entry` を引数として受け取る設計とし、テストや既存のローダーとの連携を容易にします。

## 3. Non-Functional Requirements
- **依存関係**: `github.com/openconfig/goyang` を使用します。
- **パフォーマンス**: 大規模なYANGモデル（OpenConfigなど）でも効率的に動作するよう、不要な探索を避けるロジックとします。

## 4. Acceptance Criteria
- [ ] `pkg/app/schema_walker.go` が作成されている。
- [ ] 指定された `yang.Entry` から全てのリストノードが正しく抽出される。
- [ ] OpenConfigモデルの `state` コンテナを優先するロジックが期待通りに動作する。
- [ ] 複合キーを持つリストが、キーごとの `ListInfo` に正しく分解される。
- [ ] ユニットテストですべての正常系・異常系がカバーされている。

## 5. Out of Scope
- YANGファイルのディスクからの直接読み込み（既存の `gnmic` ローダーの利用を想定）。
- キャッシュの実装（本Trackの範囲外）。
- 実際のgNMI Get/Subscribeによる値の取得処理。
