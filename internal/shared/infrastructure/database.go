// Package infrastructure 共有インフラストラクチャ層
package infrastructure

import (
	"context"

	"github.com/uptrace/bun"
)

// RunInTransaction トランザクション実行
func RunInTransaction(ctx context.Context, db *bun.DB, fn func(ctx context.Context, tx bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Pagination ページネーション設定
type Pagination struct {
	// Page ページ番号 1始まり
	Page int
	// PerPage 1ページあたりの件数
	PerPage int
}

// Offset オフセット計算
func (p Pagination) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

// Limit リミット値
func (p Pagination) Limit() int {
	if p.PerPage < 1 {
		return 20
	}
	if p.PerPage > 100 {
		return 100
	}
	return p.PerPage
}

// NewPagination ページネーション生成
func NewPagination(page, perPage int) Pagination {
	return Pagination{
		Page:    page,
		PerPage: perPage,
	}
}
