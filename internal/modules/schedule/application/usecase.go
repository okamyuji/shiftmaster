// Package application 勤務表アプリケーション層
package application

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"shiftmaster/internal/modules/schedule/domain"
	shiftDomain "shiftmaster/internal/modules/shift/domain"
	staffDomain "shiftmaster/internal/modules/staff/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// ScheduleUseCase 勤務表ユースケース
type ScheduleUseCase struct {
	scheduleRepo  domain.ScheduleRepository
	entryRepo     domain.ScheduleEntryRepository
	shiftTypeRepo shiftDomain.ShiftTypeRepository
	staffRepo     staffDomain.StaffRepository
	optimizer     domain.ScheduleOptimizer
	logger        *slog.Logger
}

// NewScheduleUseCase 勤務表ユースケース生成
func NewScheduleUseCase(
	scheduleRepo domain.ScheduleRepository,
	entryRepo domain.ScheduleEntryRepository,
	shiftTypeRepo shiftDomain.ShiftTypeRepository,
	staffRepo staffDomain.StaffRepository,
	optimizer domain.ScheduleOptimizer,
	logger *slog.Logger,
) *ScheduleUseCase {
	return &ScheduleUseCase{
		scheduleRepo:  scheduleRepo,
		entryRepo:     entryRepo,
		shiftTypeRepo: shiftTypeRepo,
		staffRepo:     staffRepo,
		optimizer:     optimizer,
		logger:        logger,
	}
}

// Create 勤務表作成
func (u *ScheduleUseCase) Create(ctx context.Context, input *CreateScheduleInput) (*ScheduleOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	orgID, err := sharedDomain.ParseID(input.OrganizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	// 既存チェック
	existing, err := u.scheduleRepo.FindByTargetMonth(ctx, orgID, input.TargetYear, input.TargetMonth)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeConflict, "同じ対象月の勤務表が既に存在します")
	}

	now := time.Now()
	schedule := &domain.Schedule{
		ID:             sharedDomain.NewID(),
		OrganizationID: orgID,
		TargetYear:     input.TargetYear,
		TargetMonth:    input.TargetMonth,
		Status:         domain.StatusDraft,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.scheduleRepo.Save(ctx, schedule); err != nil {
		u.logger.Error("勤務表作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("勤務表作成完了", "schedule_id", schedule.ID)
	return ToScheduleOutput(schedule), nil
}

// GetByID IDで勤務表取得
func (u *ScheduleUseCase) GetByID(ctx context.Context, id string) (*ScheduleOutput, error) {
	scheduleID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	schedule, err := u.scheduleRepo.FindByIDWithEntries(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, sharedDomain.ErrNotFound
	}

	output := ToScheduleOutput(schedule)

	// エントリにスタッフ名・シフト種別名を設定
	for i := range output.Entries {
		entry := &output.Entries[i]
		// スタッフ名取得
		if u.staffRepo != nil && entry.StaffID != "" {
			staffID, parseErr := sharedDomain.ParseID(entry.StaffID)
			if parseErr == nil {
				staff, staffErr := u.staffRepo.FindByID(ctx, staffID)
				if staffErr == nil && staff != nil {
					entry.StaffName = staff.LastName + " " + staff.FirstName
				}
			}
		}
		// シフト種別名取得
		if u.shiftTypeRepo != nil && entry.ShiftTypeID != "" {
			shiftTypeID, parseErr := sharedDomain.ParseID(entry.ShiftTypeID)
			if parseErr == nil {
				shiftType, stErr := u.shiftTypeRepo.FindByID(ctx, shiftTypeID)
				if stErr == nil && shiftType != nil {
					entry.ShiftTypeName = shiftType.Name
					entry.ShiftTypeCode = shiftType.Code
				}
			}
		}
	}

	return output, nil
}

// List 勤務表一覧取得
func (u *ScheduleUseCase) List(ctx context.Context, organizationID string) (*ScheduleListOutput, error) {
	orgID, err := sharedDomain.ParseID(organizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	schedules, err := u.scheduleRepo.FindByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ScheduleOutput, len(schedules))
	for i, s := range schedules {
		outputs[i] = *ToScheduleOutput(&s)
	}

	return &ScheduleListOutput{
		Schedules: outputs,
		Total:     len(outputs),
	}, nil
}

// CreateEntry エントリ作成
func (u *ScheduleUseCase) CreateEntry(ctx context.Context, input *CreateEntryInput) (*ScheduleEntryOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	scheduleID, err := sharedDomain.ParseID(input.ScheduleID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "勤務表IDが不正です")
	}

	staffID, err := sharedDomain.ParseID(input.StaffID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "スタッフIDが不正です")
	}

	targetDate, err := time.Parse("2006-01-02", input.TargetDate)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象日の形式が不正です")
	}

	// 勤務表存在確認
	schedule, err := u.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeNotFound, "勤務表が見つかりません")
	}

	// 対象日が勤務表の対象月内かチェック
	scheduleStart := time.Date(schedule.TargetYear, time.Month(schedule.TargetMonth), 1, 0, 0, 0, 0, time.UTC)
	scheduleEnd := scheduleStart.AddDate(0, 1, -1)
	if targetDate.Before(scheduleStart) || targetDate.After(scheduleEnd) {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "対象日が勤務表の対象月内ではありません")
	}

	var shiftTypeID *sharedDomain.ID
	if input.ShiftTypeID != "" {
		id, err := sharedDomain.ParseID(input.ShiftTypeID)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフト種別IDが不正です")
		}
		shiftTypeID = &id
	}

	now := time.Now()
	entry := &domain.ScheduleEntry{
		ID:          sharedDomain.NewID(),
		ScheduleID:  scheduleID,
		StaffID:     staffID,
		TargetDate:  targetDate,
		ShiftTypeID: shiftTypeID,
		IsConfirmed: false,
		Note:        input.Note,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.entryRepo.Save(ctx, entry); err != nil {
		u.logger.Error("エントリ作成失敗", "error", err)
		return nil, err
	}

	u.logger.Info("エントリ作成完了", "entry_id", entry.ID)
	return ToScheduleEntryOutput(entry), nil
}

// UpdateEntry エントリ更新
func (u *ScheduleUseCase) UpdateEntry(ctx context.Context, input *UpdateEntryInput) (*ScheduleEntryOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	entryID, err := sharedDomain.ParseID(input.ID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	entry, err := u.entryRepo.FindByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, sharedDomain.ErrNotFound
	}

	if input.ShiftTypeID != "" {
		shiftTypeID, err := sharedDomain.ParseID(input.ShiftTypeID)
		if err != nil {
			return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "シフト種別IDが不正です")
		}
		entry.ShiftTypeID = &shiftTypeID
	} else {
		entry.ShiftTypeID = nil
	}

	entry.IsConfirmed = input.IsConfirmed
	entry.Note = input.Note
	entry.UpdatedAt = time.Now()

	if err := u.entryRepo.Save(ctx, entry); err != nil {
		u.logger.Error("エントリ更新失敗", "error", err)
		return nil, err
	}

	return ToScheduleEntryOutput(entry), nil
}

// BulkUpdateEntries 一括エントリ更新
func (u *ScheduleUseCase) BulkUpdateEntries(ctx context.Context, input *BulkUpdateEntriesInput) error {
	scheduleID, err := sharedDomain.ParseID(input.ScheduleID)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "勤務表IDが不正です")
	}

	schedule, err := u.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return sharedDomain.ErrNotFound
	}

	now := time.Now()
	entries := make([]domain.ScheduleEntry, 0, len(input.Entries))

	for _, e := range input.Entries {
		staffID, err := sharedDomain.ParseID(e.StaffID)
		if err != nil {
			continue
		}

		targetDate, err := time.Parse("2006-01-02", e.TargetDate)
		if err != nil {
			continue
		}

		var shiftTypeID *sharedDomain.ID
		if e.ShiftTypeID != "" {
			id, err := sharedDomain.ParseID(e.ShiftTypeID)
			if err == nil {
				shiftTypeID = &id
			}
		}

		entry := domain.ScheduleEntry{
			ID:          sharedDomain.NewID(),
			ScheduleID:  scheduleID,
			StaffID:     staffID,
			TargetDate:  targetDate,
			ShiftTypeID: shiftTypeID,
			IsConfirmed: false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		entries = append(entries, entry)
	}

	if err := u.entryRepo.SaveBatch(ctx, entries); err != nil {
		u.logger.Error("一括エントリ更新失敗", "error", err)
		return err
	}

	u.logger.Info("一括エントリ更新完了", "schedule_id", scheduleID, "count", len(entries))
	return nil
}

// Publish 勤務表公開
func (u *ScheduleUseCase) Publish(ctx context.Context, id string) (*ScheduleOutput, error) {
	scheduleID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	schedule, err := u.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, sharedDomain.ErrNotFound
	}

	now := time.Now()
	schedule.Status = domain.StatusPublished
	schedule.PublishedAt = &now
	schedule.UpdatedAt = now

	if err := u.scheduleRepo.Save(ctx, schedule); err != nil {
		u.logger.Error("勤務表公開失敗", "error", err)
		return nil, err
	}

	u.logger.Info("勤務表公開完了", "schedule_id", scheduleID)
	return ToScheduleOutput(schedule), nil
}

// Validate 勤務表検証
func (u *ScheduleUseCase) Validate(ctx context.Context, id string) (*ValidateResult, error) {
	scheduleID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	schedule, err := u.scheduleRepo.FindByIDWithEntries(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, sharedDomain.ErrNotFound
	}

	// シフト種別マップを取得
	shiftTypeMap, err := u.buildShiftTypeMap(ctx, schedule.OrganizationID)
	if err != nil {
		u.logger.Error("シフト種別マップ取得失敗", "error", err)
		return nil, err
	}

	// 違反リストを収集
	violations := make([]ViolationOutput, 0)

	// スタッフごとにエントリをグルーピング
	staffEntries := u.groupEntriesByStaff(schedule.Entries)

	for staffID, entries := range staffEntries {
		// 日付でソート
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].TargetDate.Before(entries[j].TargetDate)
		})

		// 連続勤務チェック
		consecutiveViolations := u.checkConsecutiveWorkDays(staffID, entries, shiftTypeMap)
		violations = append(violations, consecutiveViolations...)

		// 夜勤連続チェック
		nightViolations := u.checkConsecutiveNightShifts(staffID, entries, shiftTypeMap)
		violations = append(violations, nightViolations...)

		// シフト間隔チェック
		intervalViolations := u.checkShiftInterval(staffID, entries, shiftTypeMap)
		violations = append(violations, intervalViolations...)

		// 月間夜勤回数チェック
		monthlyNightViolations := u.checkMonthlyNightShiftLimit(staffID, entries, shiftTypeMap)
		violations = append(violations, monthlyNightViolations...)
	}

	// シフト未割り当てチェック
	unassignedViolations := u.checkUnassignedShifts(schedule.Entries)
	violations = append(violations, unassignedViolations...)

	u.logger.Info("勤務表検証完了", "schedule_id", scheduleID, "violation_count", len(violations))

	return &ValidateResult{
		IsValid:    len(violations) == 0,
		Violations: violations,
	}, nil
}

// buildShiftTypeMap 組織のシフト種別マップを構築
func (u *ScheduleUseCase) buildShiftTypeMap(ctx context.Context, organizationID sharedDomain.ID) (map[string]*shiftDomain.ShiftType, error) {
	shiftTypes, err := u.shiftTypeRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*shiftDomain.ShiftType)
	for i := range shiftTypes {
		result[shiftTypes[i].ID.String()] = &shiftTypes[i]
	}
	return result, nil
}

// groupEntriesByStaff エントリをスタッフIDでグルーピング
func (u *ScheduleUseCase) groupEntriesByStaff(entries []domain.ScheduleEntry) map[string][]domain.ScheduleEntry {
	result := make(map[string][]domain.ScheduleEntry)
	for _, entry := range entries {
		staffID := entry.StaffID.String()
		result[staffID] = append(result[staffID], entry)
	}
	return result
}

// checkConsecutiveWorkDays 連続勤務日数チェック 7日以上連続で勤務している場合に警告
func (u *ScheduleUseCase) checkConsecutiveWorkDays(
	staffID string,
	entries []domain.ScheduleEntry,
	shiftTypeMap map[string]*shiftDomain.ShiftType,
) []ViolationOutput {
	const maxConsecutiveDays = 6
	violations := make([]ViolationOutput, 0)

	consecutiveCount := 0
	var consecutiveStart time.Time

	for i, entry := range entries {
		// シフトが割り当てられているか確認
		if entry.ShiftTypeID == nil {
			consecutiveCount = 0
			continue
		}

		// 休日シフトかどうか確認
		shiftType, ok := shiftTypeMap[entry.ShiftTypeID.String()]
		if !ok || shiftType.IsHoliday {
			consecutiveCount = 0
			continue
		}

		// 連続勤務のカウント
		if consecutiveCount == 0 {
			consecutiveStart = entry.TargetDate
			consecutiveCount = 1
		} else {
			prevEntry := entries[i-1]
			daysDiff := entry.TargetDate.Sub(prevEntry.TargetDate).Hours() / 24
			if daysDiff == 1 {
				consecutiveCount++
			} else {
				consecutiveStart = entry.TargetDate
				consecutiveCount = 1
			}
		}

		// 連続勤務日数が上限を超えた場合
		if consecutiveCount > maxConsecutiveDays {
			violations = append(violations, ViolationOutput{
				Type:     "consecutive_work",
				Message:  fmt.Sprintf("%d日連続勤務です（上限: %d日）", consecutiveCount, maxConsecutiveDays),
				StaffID:  staffID,
				Date:     consecutiveStart.Format("2006-01-02"),
				Severity: "warning",
			})
		}
	}

	return violations
}

// checkConsecutiveNightShifts 連続夜勤チェック 3日以上連続で夜勤の場合に警告
func (u *ScheduleUseCase) checkConsecutiveNightShifts(
	staffID string,
	entries []domain.ScheduleEntry,
	shiftTypeMap map[string]*shiftDomain.ShiftType,
) []ViolationOutput {
	const maxConsecutiveNights = 2
	violations := make([]ViolationOutput, 0)

	consecutiveCount := 0
	var consecutiveStart time.Time

	for i, entry := range entries {
		if entry.ShiftTypeID == nil {
			consecutiveCount = 0
			continue
		}

		shiftType, ok := shiftTypeMap[entry.ShiftTypeID.String()]
		if !ok || !shiftType.IsNightShift {
			consecutiveCount = 0
			continue
		}

		// 連続夜勤のカウント
		if consecutiveCount == 0 {
			consecutiveStart = entry.TargetDate
			consecutiveCount = 1
		} else {
			prevEntry := entries[i-1]
			daysDiff := entry.TargetDate.Sub(prevEntry.TargetDate).Hours() / 24
			if daysDiff == 1 {
				consecutiveCount++
			} else {
				consecutiveStart = entry.TargetDate
				consecutiveCount = 1
			}
		}

		// 連続夜勤が上限を超えた場合
		if consecutiveCount > maxConsecutiveNights {
			violations = append(violations, ViolationOutput{
				Type:     "consecutive_night",
				Message:  fmt.Sprintf("%d日連続夜勤です（上限: %d日）", consecutiveCount, maxConsecutiveNights),
				StaffID:  staffID,
				Date:     consecutiveStart.Format("2006-01-02"),
				Severity: "error",
			})
		}
	}

	return violations
}

// checkShiftInterval シフト間隔チェック 夜勤明けの翌日が日勤の場合に警告
func (u *ScheduleUseCase) checkShiftInterval(
	staffID string,
	entries []domain.ScheduleEntry,
	shiftTypeMap map[string]*shiftDomain.ShiftType,
) []ViolationOutput {
	const minIntervalHours = 11 // 最小インターバル11時間
	violations := make([]ViolationOutput, 0)

	for i := 1; i < len(entries); i++ {
		prevEntry := entries[i-1]
		currEntry := entries[i]

		// 両方ともシフトが割り当てられている必要がある
		if prevEntry.ShiftTypeID == nil || currEntry.ShiftTypeID == nil {
			continue
		}

		prevShift, prevOk := shiftTypeMap[prevEntry.ShiftTypeID.String()]
		currShift, currOk := shiftTypeMap[currEntry.ShiftTypeID.String()]

		if !prevOk || !currOk {
			continue
		}

		// 休日シフトはスキップ
		if prevShift.IsHoliday || currShift.IsHoliday {
			continue
		}

		// 連続する日付かどうか確認
		daysDiff := currEntry.TargetDate.Sub(prevEntry.TargetDate).Hours() / 24
		if daysDiff != 1 {
			continue
		}

		// 前のシフト終了時刻から次のシフト開始時刻までの間隔を計算
		prevEndHour := prevShift.EndTime.Hour()
		prevEndMin := prevShift.EndTime.Minute()
		currStartHour := currShift.StartTime.Hour()
		currStartMin := currShift.StartTime.Minute()

		// 夜勤（日跨ぎ）の場合の調整
		var intervalHours float64
		if prevShift.IsNightShift {
			// 夜勤の場合、終了時刻は翌日
			// 24時間 - 前のシフト終了時刻 + 次のシフト開始時刻ではなく
			// 次のシフト開始時刻 - 前のシフト終了時刻（翌日扱い）
			intervalHours = float64(currStartHour) + float64(currStartMin)/60 - float64(prevEndHour) - float64(prevEndMin)/60
			if intervalHours < 0 {
				intervalHours += 24
			}
		} else {
			// 日勤の場合
			intervalHours = 24 - float64(prevEndHour) - float64(prevEndMin)/60 + float64(currStartHour) + float64(currStartMin)/60
		}

		if intervalHours < float64(minIntervalHours) {
			violations = append(violations, ViolationOutput{
				Type:     "shift_interval",
				Message:  fmt.Sprintf("シフト間隔が%.1f時間です（最小: %d時間）", intervalHours, minIntervalHours),
				StaffID:  staffID,
				Date:     currEntry.TargetDate.Format("2006-01-02"),
				Severity: "warning",
			})
		}
	}

	return violations
}

// checkMonthlyNightShiftLimit 月間夜勤回数チェック 8回以上の場合に警告
func (u *ScheduleUseCase) checkMonthlyNightShiftLimit(
	staffID string,
	entries []domain.ScheduleEntry,
	shiftTypeMap map[string]*shiftDomain.ShiftType,
) []ViolationOutput {
	const maxMonthlyNights = 8
	violations := make([]ViolationOutput, 0)

	nightCount := 0
	for _, entry := range entries {
		if entry.ShiftTypeID == nil {
			continue
		}

		shiftType, ok := shiftTypeMap[entry.ShiftTypeID.String()]
		if !ok {
			continue
		}

		if shiftType.IsNightShift {
			nightCount++
		}
	}

	if nightCount > maxMonthlyNights {
		violations = append(violations, ViolationOutput{
			Type:     "monthly_night_limit",
			Message:  fmt.Sprintf("月間夜勤回数が%d回です（上限: %d回）", nightCount, maxMonthlyNights),
			StaffID:  staffID,
			Date:     "",
			Severity: "warning",
		})
	}

	return violations
}

// checkUnassignedShifts シフト未割り当てチェック
func (u *ScheduleUseCase) checkUnassignedShifts(entries []domain.ScheduleEntry) []ViolationOutput {
	violations := make([]ViolationOutput, 0)

	for _, entry := range entries {
		if entry.ShiftTypeID == nil && !entry.IsConfirmed {
			violations = append(violations, ViolationOutput{
				Type:     "unassigned_shift",
				Message:  "シフトが未割り当てです",
				StaffID:  entry.StaffID.String(),
				Date:     entry.TargetDate.Format("2006-01-02"),
				Severity: "info",
			})
		}
	}

	return violations
}

// Delete 勤務表削除
func (u *ScheduleUseCase) Delete(ctx context.Context, id string) error {
	scheduleID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	schedule, err := u.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return sharedDomain.ErrNotFound
	}

	if schedule.Status == domain.StatusPublished {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "公開済みの勤務表は削除できません")
	}

	// エントリ削除
	if err := u.entryRepo.DeleteBySchedule(ctx, scheduleID); err != nil {
		return err
	}

	// 勤務表削除
	if err := u.scheduleRepo.Delete(ctx, scheduleID); err != nil {
		u.logger.Error("勤務表削除失敗", "error", err)
		return err
	}

	u.logger.Info("勤務表削除完了", "schedule_id", scheduleID)
	return nil
}
