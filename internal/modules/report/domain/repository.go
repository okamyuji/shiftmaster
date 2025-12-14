// Package domain レポートドメイン層
package domain

import (
	"context"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ReportRepository レポートリポジトリインターフェース
type ReportRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Report, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]Report, error)
	// FindByType 種別で検索
	FindByType(ctx context.Context, organizationID sharedDomain.ID, reportType ReportType) ([]Report, error)
	// Save 保存
	Save(ctx context.Context, report *Report) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// SummaryRepository 集計リポジトリインターフェース
type SummaryRepository interface {
	// GetMonthlySummary 月次集計取得
	GetMonthlySummary(ctx context.Context, organizationID sharedDomain.ID, year, month int) (*MonthlySummary, error)
	// GetStaffSummary スタッフ集計取得
	GetStaffSummary(ctx context.Context, staffID sharedDomain.ID, year, month int) (*StaffSummary, error)
	// GetDailySummaries 日別集計取得
	GetDailySummaries(ctx context.Context, organizationID sharedDomain.ID, year, month int) ([]DailySummary, error)
}

// ReportGenerator レポート生成インターフェース
type ReportGenerator interface {
	// GenerateScheduleReport 勤務表レポート生成
	GenerateScheduleReport(ctx context.Context, scheduleID sharedDomain.ID) ([]byte, error)
	// GenerateSummaryReport 集計レポート生成
	GenerateSummaryReport(ctx context.Context, summary *MonthlySummary) ([]byte, error)
	// GenerateActualReport 実績レポート生成
	GenerateActualReport(ctx context.Context, organizationID sharedDomain.ID, year, month int) ([]byte, error)
}
