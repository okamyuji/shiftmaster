// Package domain 勤務表ドメイン層
package domain

import (
	"time"

	"shiftmaster/internal/shared/domain"
)

// ShiftTimeInfo シフト時間情報
// 検証に必要なシフト種別の時間情報
type ShiftTimeInfo struct {
	// ID シフト種別ID
	ID domain.ID
	// StartTime 開始時刻 "HH:MM"
	StartTime string
	// EndTime 終了時刻 "HH:MM"
	EndTime string
	// IsWorkingShift 勤務シフトかどうか 公休・有給等はfalse
	IsWorkingShift bool
}

// CrossesDate 日跨ぎ判定
func (s ShiftTimeInfo) CrossesDate() bool {
	startTOD, err1 := domain.ParseTimeOfDay(s.StartTime)
	endTOD, err2 := domain.ParseTimeOfDay(s.EndTime)
	if err1 != nil || err2 != nil {
		return false
	}
	return endTOD.ToMinutes() <= startTOD.ToMinutes()
}

// ToShiftTimeRange ShiftTimeRangeへ変換
func (s ShiftTimeInfo) ToShiftTimeRange() (domain.ShiftTimeRange, error) {
	return domain.NewShiftTimeRangeFromStrings(s.StartTime, s.EndTime)
}

// ScheduleValidator 勤務表検証ドメインサービス
type ScheduleValidator struct{}

// NewScheduleValidator 勤務表検証サービス生成
func NewScheduleValidator() *ScheduleValidator {
	return &ScheduleValidator{}
}

// ValidateEntries 勤務表エントリ検証
func (v *ScheduleValidator) ValidateEntries(
	entries []ScheduleEntry,
	shiftInfoMap map[domain.ID]ShiftTimeInfo,
) []ConstraintViolation {
	var violations []ConstraintViolation

	// 時間重複チェック
	timeOverlapViolations := v.validateNoTimeOverlap(entries, shiftInfoMap)
	violations = append(violations, timeOverlapViolations...)

	return violations
}

// ValidateNoTimeOverlap 同一スタッフの時間重複検証
func (v *ScheduleValidator) validateNoTimeOverlap(
	entries []ScheduleEntry,
	shiftInfoMap map[domain.ID]ShiftTimeInfo,
) []ConstraintViolation {
	var violations []ConstraintViolation

	// スタッフごとにグループ化
	byStaff := make(map[domain.ID][]entryWithShift)
	for _, e := range entries {
		if e.ShiftTypeID == nil {
			continue
		}
		info, ok := shiftInfoMap[*e.ShiftTypeID]
		if !ok || !info.IsWorkingShift {
			continue // 公休・有給等はスキップ
		}

		timeRange, err := info.ToShiftTimeRange()
		if err != nil {
			continue
		}

		byStaff[e.StaffID] = append(byStaff[e.StaffID], entryWithShift{
			entry:     e,
			shiftInfo: info,
			timeRange: timeRange,
		})
	}

	// スタッフごとに重複チェック
	for staffID, staffEntries := range byStaff {
		staffViolations := v.checkStaffTimeOverlaps(staffID, staffEntries)
		violations = append(violations, staffViolations...)
	}

	return violations
}

// entryWithShift エントリとシフト情報のペア
type entryWithShift struct {
	entry     ScheduleEntry
	shiftInfo ShiftTimeInfo
	timeRange domain.ShiftTimeRange
}

// checkStaffTimeOverlaps スタッフの時間重複チェック
func (v *ScheduleValidator) checkStaffTimeOverlaps(
	staffID domain.ID,
	entries []entryWithShift,
) []ConstraintViolation {
	var violations []ConstraintViolation

	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			e1, e2 := entries[i], entries[j]

			if v.hasTimeOverlap(e1, e2) {
				date := e1.entry.TargetDate
				violations = append(violations, ConstraintViolation{
					ConstraintType: "time_overlap",
					Message:        "同一スタッフに時間が重複するシフトが割り当てられています",
					StaffID:        &staffID,
					Date:           &date,
					Severity:       "error",
				})
			}
		}
	}

	return violations
}

// hasTimeOverlap 2つのエントリ間の時間重複判定
func (v *ScheduleValidator) hasTimeOverlap(e1, e2 entryWithShift) bool {
	// 同一日の場合
	if isSameDate(e1.entry.TargetDate, e2.entry.TargetDate) {
		return e1.timeRange.Overlaps(e2.timeRange)
	}

	// 連続日の場合
	daysDiff := daysBetween(e1.entry.TargetDate, e2.entry.TargetDate)

	// e1が前日、e2が翌日で、e1が日跨ぎの場合
	if daysDiff == 1 && e1.shiftInfo.CrossesDate() {
		return e1.timeRange.OverlapsOnConsecutiveDays(e2.timeRange)
	}

	// e2が前日、e1が翌日で、e2が日跨ぎの場合
	if daysDiff == -1 && e2.shiftInfo.CrossesDate() {
		return e2.timeRange.OverlapsOnConsecutiveDays(e1.timeRange)
	}

	return false
}

// isSameDate 同一日判定
func isSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// daysBetween 日数差 t2 - t1
func daysBetween(t1, t2 time.Time) int {
	// 時刻を正規化（日付のみ比較）
	d1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.UTC)
	return int(d2.Sub(d1).Hours() / 24)
}

// ValidateConsecutiveWorkDays 連続勤務日数チェック
func (v *ScheduleValidator) ValidateConsecutiveWorkDays(
	entries []ScheduleEntry,
	shiftInfoMap map[domain.ID]ShiftTimeInfo,
	maxConsecutiveDays int,
) []ConstraintViolation {
	var violations []ConstraintViolation

	// スタッフごとにグループ化
	byStaff := make(map[domain.ID][]ScheduleEntry)
	for _, e := range entries {
		if e.ShiftTypeID == nil {
			continue
		}
		info, ok := shiftInfoMap[*e.ShiftTypeID]
		if !ok || !info.IsWorkingShift {
			continue
		}
		byStaff[e.StaffID] = append(byStaff[e.StaffID], e)
	}

	for staffID, staffEntries := range byStaff {
		// 日付でソート済みと仮定
		consecutive := 1
		for i := 1; i < len(staffEntries); i++ {
			prev := staffEntries[i-1]
			curr := staffEntries[i]

			if daysBetween(prev.TargetDate, curr.TargetDate) == 1 {
				consecutive++
				if consecutive > maxConsecutiveDays {
					date := curr.TargetDate
					violations = append(violations, ConstraintViolation{
						ConstraintType: "consecutive_work_days",
						Message:        "連続勤務日数が上限を超えています",
						StaffID:        &staffID,
						Date:           &date,
						Severity:       "warning",
					})
				}
			} else {
				consecutive = 1
			}
		}
	}

	return violations
}

// ValidateNightShiftInterval 夜勤間隔チェック
func (v *ScheduleValidator) ValidateNightShiftInterval(
	entries []ScheduleEntry,
	shiftInfoMap map[domain.ID]ShiftTimeInfo,
	minIntervalDays int,
) []ConstraintViolation {
	var violations []ConstraintViolation

	// スタッフごとにグループ化
	byStaff := make(map[domain.ID][]ScheduleEntry)
	for _, e := range entries {
		if e.ShiftTypeID == nil {
			continue
		}
		info, ok := shiftInfoMap[*e.ShiftTypeID]
		if !ok {
			continue
		}
		// 日跨ぎシフトを夜勤とみなす
		if info.CrossesDate() {
			byStaff[e.StaffID] = append(byStaff[e.StaffID], e)
		}
	}

	for staffID, nightShifts := range byStaff {
		for i := 1; i < len(nightShifts); i++ {
			prev := nightShifts[i-1]
			curr := nightShifts[i]

			interval := daysBetween(prev.TargetDate, curr.TargetDate)
			if interval < minIntervalDays && interval > 0 {
				date := curr.TargetDate
				violations = append(violations, ConstraintViolation{
					ConstraintType: "night_shift_interval",
					Message:        "夜勤間隔が短すぎます",
					StaffID:        &staffID,
					Date:           &date,
					Severity:       "warning",
				})
			}
		}
	}

	return violations
}
