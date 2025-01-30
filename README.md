# Slack Word Cloud Generator

Slackのメッセージを解析してワードクラウドを生成するフルスタックアプリケーションです。

## 機能

### バックエンド
- Slackチャンネルからのメッセージ取得
- 日本語テキストの形態素解析
- ワードクラウドデータの生成
- CSVおよびJSONファイルの出力

### フロントエンド
- CSVファイルのドラッグ&ドロップまたは選択によるアップロード
- インタラクティブなワードクラウド表示
- 単語の出現頻度に応じたサイズ・色の動的調整
- レスポンシブデザイン対応

## システム構成

```
slack-word-analyzer/
├── backend/
│   ├── cmd/
│   │   ├── getmessage/     # Slackメッセージ取得コマンド
│   │   └── wordcloud/      # ワードクラウド生成コマンド
│   ├── pkg/
│   │   ├── slack/          # Slack API共通コード
│   │   └── wordcloud/      # ワードクラウド生成共通コード
│   └── data/               # 生成データの保存先
└── frontend/
    ├── src/
    │   ├── app/            # Next.js アプリケーションルート
    │   ├── components/
    │   │   ├── ui/         # 共通UIコンポーネント
    │   │   ├── layout/     # レイアウトコンポーネント
    │   │   └── wordcloud/  # ワードクラウド関連コンポーネント
    │   └── lib/
    │       ├── types/      # 型定義
    │       └── utils/      # ユーティリティ関数
    └── public/             # 静的ファイル
```

## 必要要件

### バックエンド
- Go 1.21以上
- Slack Bot Token（以下のスコープが必要）
  - channels:history
  - channels:read
  - users:read
  - groups:history（プライベートチャンネル用）

### フロントエンド
- Node.js 18.0以上
- npm 9.0以上

## セットアップ

### 1. バックエンドのセットアップ

```bash
cd backend
go mod init your/project/path
go get -u github.com/slack-go/slack
go get -u github.com/ikawaha/kagome/v2
go get -u github.com/ikawaha/kagome-dict/ipa
```

### 2. フロントエンドのセットアップ

```bash
cd frontend
npm install
```

### 3. Slack Appの設定
- [Slack API](https://api.slack.com/apps)で新しいアプリを作成
- 必要なBot Token Scopesを追加
- アプリをワークスペースにインストール
- Bot User OAuth Tokenを取得

## 使用方法

### 1. メッセージの取得（バックエンド）

```bash
cd backend
go run cmd/getmessage/main.go \
  -token "xoxb-your-bot-token" \
  -channel "C1234567890" \
  -output "../data"
```

### 2. ワードクラウドの生成（バックエンド）

```bash
go run cmd/wordcloud/main.go \
  -input "./data/messages_general_20240130_120000.csv" \
  -output "./data/wordcloud.json"
```

### 3. フロントエンドの起動

```bash
cd frontend
npm run dev
```

## データフォーマット

### CSVファイル（Slackメッセージ）
```csv
Timestamp,UserID,Username,Message,ThreadTS
2024-01-30 12:00:00,U123,user1,こんにちは,
```

### JSONファイル（ワードクラウドデータ）
```json
[
  {
    "text": "単語",
    "count": 10,
    "fontSize": 32,
    "color": "#0000FF"
  }
]
```

## エラーハンドリング

### バックエンド
- 認証エラー：トークンとスコープの確認
- チャンネルアクセスエラー：Botの参加確認
- レートリミットエラー：API呼び出し間隔の調整

### フロントエンド
- ファイル形式エラー：CSVファイルのみ許可
- ファイルサイズエラー：5MB以下に制限
- データ処理エラー：エラーメッセージの表示

## 開発ガイドライン

### バックエンド
1. パッケージの役割を明確に分離
2. エラーは適切にラップして日本語で具体的に
3. 設定値は構造体のフィールドとして公開

### フロントエンド
1. コンポーネントの責務を明確に分離
2. TypeScriptの型定義を徹底
3. レスポンシブデザインに対応
4. アクセシビリティに配慮

## パフォーマンス考慮事項
- バックエンド：大量メッセージ処理時のメモリ管理
- フロントエンド：大規模データセット表示時の最適化

## セキュリティ
- APIトークンは環境変数で管理
- フロントエンドでのファイルバリデーション
- 適切なエラーログ設定

## ライセンス
MIT License

## 貢献方法
1. このリポジトリをフォーク
2. 新しいブランチを作成
3. 変更をコミット
4. プルリクエストを作成
