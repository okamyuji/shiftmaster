// Package domain 勤務希望ドメイン層
package domain

import (
	"context"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// RequestPeriodRepository 勤務希望受付期間リポジトリインターフェース
type RequestPeriodRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*RequestPeriod, error)
	// FindAll 全件取得
	FindAll(ctx context.Context) ([]RequestPeriod, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]RequestPeriod, error)
	// FindByTargetMonth 対象年月で検索
	FindByTargetMonth(ctx context.Context, organizationID sharedDomain.ID, year, month int) (*RequestPeriod, error)
	// FindActive 受付中の期間を取得
	FindActive(ctx context.Context, organizationID sharedDomain.ID) ([]RequestPeriod, error)
	// Save 保存
	Save(ctx context.Context, period *RequestPeriod) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// ShiftRequestRepository 勤務希望リポジトリインターフェース
type ShiftRequestRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*ShiftRequest, error)
	// FindByPeriodID 期間IDで検索
	FindByPeriodID(ctx context.Context, periodID sharedDomain.ID) ([]ShiftRequest, error)
	// FindByStaffID スタッフIDで検索
	FindByStaffID(ctx context.Context, staffID sharedDomain.ID) ([]ShiftRequest, error)
	// FindByPeriodAndStaff 期間とスタッフで検索
	FindByPeriodAndStaff(ctx context.Context, periodID, staffID sharedDomain.ID) ([]ShiftRequest, error)
	// FindByPeriodAndDate 期間と日付で検索
	FindByPeriodAndDate(ctx context.Context, periodID sharedDomain.ID, date time.Time) ([]ShiftRequest, error)
	// CountByPeriodAndStaff スタッフの希望数カウント
	CountByPeriodAndStaff(ctx context.Context, periodID, staffID sharedDomain.ID) (int, error)
	// CountByPeriodAndDate 日付の希望数カウント
	CountByPeriodAndDate(ctx context.Context, periodID sharedDomain.ID, date time.Time) (int, error)
	// Save 保存
	Save(ctx context.Context, request *ShiftRequest) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
	// DeleteByPeriodAndStaff 期間とスタッフで削除
	DeleteByPeriodAndStaff(ctx context.Context, periodID, staffID sharedDomain.ID) error
}
