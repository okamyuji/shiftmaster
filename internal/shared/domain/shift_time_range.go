// Package domain 共有ドメイン型定義
package domain

// ShiftTimeRange シフト時間範囲 Value Object
// 夜勤等の日跨ぎシフトに対応
type ShiftTimeRange struct {
	start TimeOfDay
	end   TimeOfDay
}

// NewShiftTimeRange シフト時間範囲生成
func NewShiftTimeRange(start, end TimeOfDay) ShiftTimeRange {
	return ShiftTimeRange{start: start, end: end}
}

// NewShiftTimeRangeFromStrings 文字列からシフト時間範囲生成
func NewShiftTimeRangeFromStrings(startStr, endStr string) (ShiftTimeRange, error) {
	start, err := ParseTimeOfDay(startStr)
	if err != nil {
		return ShiftTimeRange{}, err
	}
	end, err := ParseTimeOfDay(endStr)
	if err != nil {
		return ShiftTimeRange{}, err
	}
	return NewShiftTimeRange(start, end), nil
}

// Start 開始時刻取得
func (r ShiftTimeRange) Start() TimeOfDay {
	return r.start
}

// End 終了時刻取得
func (r ShiftTimeRange) End() TimeOfDay {
	return r.end
}

// CrossesDate 日跨ぎ判定
// 終了時刻が開始時刻より前または同じ場合は日跨ぎ
func (r ShiftTimeRange) CrossesDate() bool {
	return r.end.ToMinutes() <= r.start.ToMinutes()
}

// DurationMinutes 時間長（分）
func (r ShiftTimeRange) DurationMinutes() int {
	if r.CrossesDate() {
		// 日跨ぎ：(24時間 - 開始時刻) + 終了時刻
		return (1440 - r.start.ToMinutes()) + r.end.ToMinutes()
	}
	return r.end.ToMinutes() - r.start.ToMinutes()
}

// Overlaps 時間重複判定（同一日での重複）
// 注意: このメソッドは同一日に2つのシフトが入る場合の重複を判定
// 連続日での重複はOverlapsOnConsecutiveDaysを使用
func (r ShiftTimeRange) Overlaps(other ShiftTimeRange) bool {
	// 両方とも日跨ぎでない場合（通常のケース）
	if !r.CrossesDate() && !other.CrossesDate() {
		return r.start.ToMinutes() < other.end.ToMinutes() &&
			other.start.ToMinutes() < r.end.ToMinutes()
	}

	// 片方だけが日跨ぎの場合
	if r.CrossesDate() && !other.CrossesDate() {
		// rが日跨ぎ（例: 17:00-09:00）、otherは通常（例: 09:00-16:00）
		// otherがrの開始時刻以降（当日部分）または終了時刻以前（翌日部分）なら重複
		// 当日部分: r.start <= other.end (17:00以降にotherの終了があるか)
		// 翌日部分: other.start < r.end (otherの開始が09:00より前か)
		return other.end.ToMinutes() > r.start.ToMinutes() ||
			other.start.ToMinutes() < r.end.ToMinutes()
	}

	if !r.CrossesDate() && other.CrossesDate() {
		// rが通常、otherが日跨ぎ
		// 同様のロジック
		return r.end.ToMinutes() > other.start.ToMinutes() ||
			r.start.ToMinutes() < other.end.ToMinutes()
	}

	// 両方とも日跨ぎの場合（例: 2つの夜勤）
	// 必ず重複する（00:00を跨ぐ時間帯が両方にある）
	return true
}

// OverlapsOnConsecutiveDays 連続日での重複判定
// thisDateのrと、nextDateのotherの重複を判定
func (r ShiftTimeRange) OverlapsOnConsecutiveDays(other ShiftTimeRange) bool {
	if !r.CrossesDate() {
		return false // 日跨ぎでなければ翌日と重複しない
	}

	// rの翌日部分（0:00〜r.end）と otherの開始部分の重複
	// rの翌日終了時刻 < otherの開始時刻 なら重複なし
	return other.start.ToMinutes() < r.end.ToMinutes()
}

// Contains 指定時刻が範囲内か判定
func (r ShiftTimeRange) Contains(t TimeOfDay) bool {
	if r.CrossesDate() {
		// 日跨ぎ：開始以降 または 終了以前
		return t.ToMinutes() >= r.start.ToMinutes() || t.ToMinutes() < r.end.ToMinutes()
	}
	return t.ToMinutes() >= r.start.ToMinutes() && t.ToMinutes() < r.end.ToMinutes()
}

// Equals 等価判定
func (r ShiftTimeRange) Equals(other ShiftTimeRange) bool {
	return r.start.Equals(other.start) && r.end.Equals(other.end)
}

// IsZero ゼロ値判定
func (r ShiftTimeRange) IsZero() bool {
	return r.start.IsZero() && r.end.IsZero()
}

// String 文字列変換
func (r ShiftTimeRange) String() string {
	return r.start.String() + "-" + r.end.String()
}
