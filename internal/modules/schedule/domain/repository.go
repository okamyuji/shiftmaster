// Package domain 勤務表ドメイン層
package domain

import (
	"context"
	"time"

	sharedDomain "shiftmaster/internal/shared/domain"
)

// ScheduleRepository 勤務表リポジトリインターフェース
type ScheduleRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*Schedule, error)
	// FindByIDWithEntries エントリ付きで検索
	FindByIDWithEntries(ctx context.Context, id sharedDomain.ID) (*Schedule, error)
	// FindByOrganizationID 組織IDで検索
	FindByOrganizationID(ctx context.Context, organizationID sharedDomain.ID) ([]Schedule, error)
	// FindByTargetMonth 対象年月で検索
	FindByTargetMonth(ctx context.Context, organizationID sharedDomain.ID, year, month int) (*Schedule, error)
	// Save 保存
	Save(ctx context.Context, schedule *Schedule) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}

// ScheduleEntryRepository 勤務表エントリリポジトリインターフェース
type ScheduleEntryRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*ScheduleEntry, error)
	// FindByScheduleID 勤務表IDで検索
	FindByScheduleID(ctx context.Context, scheduleID sharedDomain.ID) ([]ScheduleEntry, error)
	// FindByScheduleAndStaff 勤務表とスタッフで検索
	FindByScheduleAndStaff(ctx context.Context, scheduleID, staffID sharedDomain.ID) ([]ScheduleEntry, error)
	// FindByScheduleAndDate 勤務表と日付で検索
	FindByScheduleAndDate(ctx context.Context, scheduleID sharedDomain.ID, date time.Time) ([]ScheduleEntry, error)
	// Save 保存
	Save(ctx context.Context, entry *ScheduleEntry) error
	// SaveBatch 一括保存
	SaveBatch(ctx context.Context, entries []ScheduleEntry) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
	// DeleteBySchedule 勤務表のエントリ全削除
	DeleteBySchedule(ctx context.Context, scheduleID sharedDomain.ID) error
}

// ActualRecordRepository 勤務実績リポジトリインターフェース
type ActualRecordRepository interface {
	// FindByID IDで検索
	FindByID(ctx context.Context, id sharedDomain.ID) (*ActualRecord, error)
	// FindByEntryID エントリIDで検索
	FindByEntryID(ctx context.Context, entryID sharedDomain.ID) (*ActualRecord, error)
	// FindBySchedule 勤務表の実績検索
	FindBySchedule(ctx context.Context, scheduleID sharedDomain.ID) ([]ActualRecord, error)
	// Save 保存
	Save(ctx context.Context, record *ActualRecord) error
	// Delete 削除
	Delete(ctx context.Context, id sharedDomain.ID) error
}
