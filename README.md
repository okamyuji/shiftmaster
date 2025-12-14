# ShiftMaster

勤務表作成システム - 汎用シフト管理アプリケーション

## 概要

ShiftMasterはスタッフの勤務希望を効率的に管理し、最適な勤務表を作成するシステムです。将来的なAI自動作成機能の追加を見据えた拡張可能な設計を採用しています。

## スクリーンショット

### ダッシュボード

- クイックアクション（スタッフ追加、シフト追加、勤務表作成、レポート確認）
- ナビゲーションメニュー

### 勤務表管理

- 月別勤務表一覧
- 勤務表作成・編集
- 条件違反チェック

### スタッフ管理

- スタッフ一覧（フィルタリング機能付き）
- スタッフ追加・編集
- チーム所属管理

## 機能一覧

### 1. 認証・認可

- JWTトークンベース認証
- ロールベースアクセス制御（admin/manager/user）
- HTTP-Only Cookieによるセキュアなトークン管理
- bcryptによるパスワードハッシュ化

### 2. ユーザー管理（管理者専用）

- ユーザー一覧・追加・編集・削除
- ロール（権限）設定
- ログイン履歴確認

### 3. 勤務希望の自動受付

- マルチデバイス対応（PC、スマートフォン）
- 柔軟な申告形式
    - 必須 / できれば の優先度設定
    - 割り当てる / 避ける の選択
    - 複数シフトからの選択
- 希望数の上限設定
- 受付期間の自動開始・終了

### 4. 勤務表作成

- 手動勤務表作成
- 月別勤務表管理
- 条件設定管理
    - チーム、スキル配置
    - 配置人数ルール
    - シフトの並び、連続、間隔
- 前月実績考慮
- 条件違反チェック
- AI自動作成インターフェース（将来拡張用）

### 5. スタッフ管理

- スタッフ情報管理
- 組織・部署・チーム階層管理
- 雇用形態管理（正社員/パート/契約/派遣）
- スキル・資格管理

### 6. シフト種別管理

- シフト種別定義（日勤/夜勤/早番/遅番/公休/有給等）
- 勤務時間・休憩時間設定
- 色分け表示

### 7. 勤務実績管理（予定）

- 実働時間・拘束時間記録
- 有給休暇消化管理
- 各種集計機能

### 8. 要件監視（予定）

- 人員配置シミュレーション
- 充足度診断
- 補充必要度診断

### 9. 帳票出力（予定）

- 勤務予定実績表
- 各種集計表
- PDF/Excel出力

## アーキテクチャ

### Modular Monolith + DDD

```text
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│                 net/http + HTMX + Alpine.js                  │
├─────────────────────────────────────────────────────────────┤
│                    Application Layer                         │
│                       Use Cases                              │
├──────────┬──────────┬──────────┬──────────┬────────────────┤
│   Auth   │  Staff   │  Shift   │ Schedule │    Request     │
│  Module  │  Module  │  Module  │  Module  │    Module      │
├──────────┴──────────┴──────────┴──────────┴────────────────┤
│                   Infrastructure Layer                       │
│                    Bun ORM + PostgreSQL                      │
└─────────────────────────────────────────────────────────────┘
```

### モジュール構成

| モジュール | 責務 |
|-----------|------|
| Auth | JWT認証、ログイン/ログアウト |
| User | ユーザー管理、ロール管理 |
| Staff | スタッフ、チーム、部署、組織管理 |
| Shift | シフト種別、勤務ルール定義 |
| Request | 勤務希望申告、受付期間管理 |
| Schedule | 勤務表作成、エントリ管理、条件検証 |
| Report | 実績管理、集計、帳票出力 |

### データ階層構造

```text
Organization（組織）
└── Department（部署）
    └── Team（チーム）
        └── Staff（スタッフ）
```

## 技術スタック

| 項目 | 技術 |
|------|------|
| 言語 | Go 1.25+ |
| データベース | PostgreSQL 16+ |
| ORM | Bun |
| HTTP | net/http (標準) |
| テンプレート | html/template (標準) |
| ログ | slog (標準) |
| 認証 | JWT (github.com/golang-jwt/jwt/v5) |
| フロントエンド | HTMX, Alpine.js, Tailwind CSS |
| コンテナ | Docker, Docker Compose |

## セットアップ

### 必要条件

- Go 1.25+
- Docker / Docker Compose
- Node.js 20+ (Tailwind CSS ビルド用)

### クイックスタート

```bash
# リポジトリクローン
git clone https://github.com/okamyuji/shiftmaster
cd shiftmaster

# Docker Compose で全サービス起動（ホットリロード付き）
docker compose up -d

# ログ確認
docker compose logs -f app

# ブラウザでアクセス
open http://localhost:8080
```

初期ログイン情報:

- メール: `admin@example.com`
- パスワード: `Password123$!`

### ローカル開発

```bash
# 依存関係インストール
go mod download
npm install

# データベースのみ起動
docker compose up -d db

# マイグレーション実行
go run cmd/migrate/main.go up

# 初期データ投入
go run cmd/seed/main.go

# Tailwind CSS ビルド（ウォッチモード）
npm run build:css

# サーバー起動
go run cmd/server/main.go
```

### 開発コマンド

```bash
# 全サービス起動
docker compose up -d

# ログ確認
docker compose logs -f

# アプリケーションログのみ
docker compose logs -f app

# 停止
docker compose down

# 停止（データベースも削除）
docker compose down -v
```

## 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| DATABASE_URL | PostgreSQL接続URL | postgres://shiftmaster:shiftmaster@localhost:5432/shiftmaster?sslmode=disable |
| SERVER_PORT | サーバーポート | 8080 |
| SERVER_HOST | サーバーホスト | 0.0.0.0 |
| LOG_LEVEL | ログレベル | info |
| JWT_SECRET | JWT署名秘密鍵 | (要設定) |

## API エンドポイント

### 認証

| Method | Path | 説明 |
|--------|------|------|
| GET | /login | ログインページ |
| POST | /login | ログイン処理 |
| POST | /logout | ログアウト |

### ユーザー管理（管理者専用）

| Method | Path | 説明 |
|--------|------|------|
| GET | /admin/users | ユーザー一覧 |
| GET | /admin/users/new | ユーザー追加フォーム |
| POST | /admin/users | ユーザー作成 |
| GET | /admin/users/{id}/edit | ユーザー編集 |
| PUT | /admin/users/{id} | ユーザー更新 |
| DELETE | /admin/users/{id} | ユーザー削除 |

### スタッフ認証

| Method | Path | 説明 |
|--------|------|------|
| GET | /staffs | スタッフ一覧 |
| GET | /staffs/{id} | スタッフ詳細 |
| POST | /staffs | スタッフ作成 |
| PUT | /staffs/{id} | スタッフ更新 |
| DELETE | /staffs/{id} | スタッフ削除 |

### チーム

| Method | Path | 説明 |
|--------|------|------|
| GET | /teams | チーム一覧 |
| POST | /teams | チーム作成 |
| PUT | /teams/{id} | チーム更新 |
| DELETE | /teams/{id} | チーム削除 |

### 認証でのシフト種別

| Method | Path | 説明 |
|--------|------|------|
| GET | /shifts | シフト種別一覧 |
| POST | /shifts | シフト種別作成 |
| PUT | /shifts/{id} | シフト種別更新 |
| DELETE | /shifts/{id} | シフト種別削除 |

### 勤務希望

| Method | Path | 説明 |
|--------|------|------|
| GET | /requests | 受付期間一覧 |
| POST | /requests | 受付期間作成 |
| GET | /requests/{id} | 受付期間詳細 |
| POST | /requests/{id}/open | 受付開始 |
| POST | /requests/{id}/close | 受付終了 |
| GET | /requests/{period_id}/entries | 勤務希望一覧 |
| POST | /requests/{period_id}/entries | 勤務希望作成 |
| DELETE | /requests/entries/{id} | 勤務希望削除 |

### 勤務表

| Method | Path | 説明 |
|--------|------|------|
| GET | /schedules | 勤務表一覧 |
| GET | /schedules/{id} | 勤務表詳細 |
| POST | /schedules | 勤務表作成 |
| PUT | /schedules/{id}/entries/{entry_id} | エントリ更新 |
| POST | /schedules/{id}/publish | 勤務表公開 |
| POST | /schedules/{id}/validate | 条件検証 |
| DELETE | /schedules/{id} | 勤務表削除 |

### レポート（予定）

| Method | Path | 説明 |
|--------|------|------|
| GET | /reports/summary | 集計レポート |
| GET | /reports/export | 帳票出力 |

## テスト

```bash
# 全テスト実行
go test ./...

# カバレッジ付き
go test -cover ./...

# 特定パッケージ
go test ./internal/modules/staff/...

# 詳細ログ付き
go test -v ./...
```

## 初期データ

`go run cmd/seed/main.go` で以下のデータが作成されます:

### ユーザー

- システム管理者（<admin@example.com> / Password123$!）

### 組織構造

- 組織: サンプル病院
- 部署: 看護部
- チーム: 看護1チーム

### シフト種別

| 名前 | コード | 開始時間 | 終了時間 |
|------|--------|----------|----------|
| 日勤 | D | 08:30 | 17:30 |
| 夜勤 | N | 17:00 | 翌09:00 |
| 早番 | E | 06:30 | 15:30 |
| 遅番 | L | 12:00 | 21:00 |
| 公休 | OFF | - | - |
| 有給 | PTO | - | - |

### 組織構造でのスタッフ

- 山田花子（EMP001）

## 将来の拡張

### AI自動作成機能

```go
// ScheduleOptimizer インターフェース
type ScheduleOptimizer interface {
    Optimize(ctx context.Context, schedule *Schedule, constraints []Constraint) (*OptimizationResult, error)
}
```

制約条件:

- 連続勤務日数上限
- 夜勤連続回数上限
- シフト間隔（インターバル）
- 月間夜勤回数上限
- 必要人員配置

### 要件監視機能

- 看護配置基準チェック
- 様式9自動作成

## ライセンス

MIT License
