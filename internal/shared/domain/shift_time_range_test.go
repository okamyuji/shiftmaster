package domain

import (
	"testing"
)

func TestNewShiftTimeRange(t *testing.T) {
	start := MustParseTimeOfDay("09:00")
	end := MustParseTimeOfDay("17:00")

	r := NewShiftTimeRange(start, end)

	if !r.Start().Equals(start) {
		t.Errorf("Start() = %v, want %v", r.Start(), start)
	}
	if !r.End().Equals(end) {
		t.Errorf("End() = %v, want %v", r.End(), end)
	}
}

func TestNewShiftTimeRangeFromStrings(t *testing.T) {
	tests := []struct {
		name      string
		startStr  string
		endStr    string
		wantStart string
		wantEnd   string
		wantErr   bool
	}{
		{"正常系_日勤", "09:00", "17:00", "09:00", "17:00", false},
		{"正常系_夜勤", "17:00", "09:00", "17:00", "09:00", false},
		{"異常系_開始不正", "25:00", "17:00", "", "", true},
		{"異常系_終了不正", "09:00", "25:00", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewShiftTimeRangeFromStrings(tt.startStr, tt.endStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewShiftTimeRangeFromStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if r.Start().String() != tt.wantStart || r.End().String() != tt.wantEnd {
					t.Errorf("NewShiftTimeRangeFromStrings() = %v-%v, want %v-%v",
						r.Start(), r.End(), tt.wantStart, tt.wantEnd)
				}
			}
		})
	}
}

func TestShiftTimeRange_CrossesDate(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string
		want  bool
	}{
		{"日勤_日跨ぎなし", "09:00", "17:00", false},
		{"夜勤_日跨ぎあり", "17:00", "09:00", true},
		{"深夜勤務_日跨ぎあり", "22:00", "06:00", true},
		{"同一時刻_日跨ぎ扱い", "09:00", "09:00", true},
		{"準夜勤_日跨ぎあり", "16:00", "01:00", true},
		{"早番_日跨ぎなし", "06:00", "14:00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewShiftTimeRangeFromStrings(tt.start, tt.end)
			if got := r.CrossesDate(); got != tt.want {
				t.Errorf("CrossesDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShiftTimeRange_DurationMinutes(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string
		want  int
	}{
		{"日勤8時間", "09:00", "17:00", 480},
		{"夜勤16時間", "17:00", "09:00", 960},
		{"短時間4時間", "10:00", "14:00", 240},
		{"深夜8時間", "22:00", "06:00", 480},
		{"準夜勤9時間", "16:00", "01:00", 540},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewShiftTimeRangeFromStrings(tt.start, tt.end)
			if got := r.DurationMinutes(); got != tt.want {
				t.Errorf("DurationMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShiftTimeRange_Overlaps_SameDay(t *testing.T) {
	tests := []struct {
		name string
		r1   [2]string
		r2   [2]string
		want bool
	}{
		{"重複あり_一部重複", [2]string{"09:00", "17:00"}, [2]string{"12:00", "20:00"}, true},
		{"重複あり_完全包含", [2]string{"09:00", "17:00"}, [2]string{"10:00", "16:00"}, true},
		{"重複あり_同一", [2]string{"09:00", "17:00"}, [2]string{"09:00", "17:00"}, true},
		{"重複なし_連続", [2]string{"09:00", "17:00"}, [2]string{"17:00", "22:00"}, false},
		{"重複なし_離れている", [2]string{"09:00", "12:00"}, [2]string{"14:00", "17:00"}, false},
		{"重複あり_境界1分重複", [2]string{"09:00", "17:01"}, [2]string{"17:00", "22:00"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, _ := NewShiftTimeRangeFromStrings(tt.r1[0], tt.r1[1])
			r2, _ := NewShiftTimeRangeFromStrings(tt.r2[0], tt.r2[1])
			if got := r1.Overlaps(r2); got != tt.want {
				t.Errorf("Overlaps() = %v, want %v", got, tt.want)
			}
			// 対称性確認
			if got := r2.Overlaps(r1); got != tt.want {
				t.Errorf("Overlaps() symmetric = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShiftTimeRange_Overlaps_WithCrossDate(t *testing.T) {
	tests := []struct {
		name string
		r1   [2]string
		r2   [2]string
		want bool
	}{
		// 夜勤同士は必ず重複（両方00:00を跨ぐ）
		{"夜勤同士_重複", [2]string{"17:00", "09:00"}, [2]string{"22:00", "06:00"}, true},
		// 日勤16:00終了と夜勤16:00開始は境界で重複
		{"日勤と夜勤_重複", [2]string{"09:00", "17:00"}, [2]string{"16:00", "01:00"}, true},
		// 日勤16:00終了と夜勤17:00開始は重複なし
		{"日勤と夜勤_重複なし", [2]string{"09:00", "16:00"}, [2]string{"17:00", "01:00"}, false},
		// 早番14:00終了と夜勤22:00開始だが、夜勤は07:00まで続き早番06:00開始と重複
		{"早番と夜勤_重複", [2]string{"06:00", "14:00"}, [2]string{"22:00", "07:00"}, true},
		// 早番14:00終了と準夜勤16:00開始は重複なし
		{"早番と準夜勤_重複なし", [2]string{"06:00", "14:00"}, [2]string{"16:00", "01:00"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, _ := NewShiftTimeRangeFromStrings(tt.r1[0], tt.r1[1])
			r2, _ := NewShiftTimeRangeFromStrings(tt.r2[0], tt.r2[1])
			if got := r1.Overlaps(r2); got != tt.want {
				t.Errorf("Overlaps() = %v, want %v (r1=%s, r2=%s)", got, tt.want, r1.String(), r2.String())
			}
		})
	}
}

func TestShiftTimeRange_OverlapsOnConsecutiveDays(t *testing.T) {
	tests := []struct {
		name string
		r1   [2]string // 今日のシフト
		r2   [2]string // 翌日のシフト
		want bool
	}{
		{"夜勤後の早番_重複あり", [2]string{"17:00", "09:00"}, [2]string{"06:00", "14:00"}, true},
		{"夜勤後の日勤_重複なし", [2]string{"17:00", "08:00"}, [2]string{"09:00", "17:00"}, false},
		{"夜勤後の遅番_重複なし", [2]string{"22:00", "06:00"}, [2]string{"14:00", "22:00"}, false},
		{"日勤後の早番_日跨ぎなし重複なし", [2]string{"09:00", "17:00"}, [2]string{"06:00", "14:00"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, _ := NewShiftTimeRangeFromStrings(tt.r1[0], tt.r1[1])
			r2, _ := NewShiftTimeRangeFromStrings(tt.r2[0], tt.r2[1])
			if got := r1.OverlapsOnConsecutiveDays(r2); got != tt.want {
				t.Errorf("OverlapsOnConsecutiveDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShiftTimeRange_Contains(t *testing.T) {
	tests := []struct {
		name string
		r    [2]string
		time string
		want bool
	}{
		{"日勤_範囲内", [2]string{"09:00", "17:00"}, "12:00", true},
		{"日勤_開始時刻", [2]string{"09:00", "17:00"}, "09:00", true},
		{"日勤_終了時刻", [2]string{"09:00", "17:00"}, "17:00", false},
		{"日勤_範囲外", [2]string{"09:00", "17:00"}, "08:00", false},
		{"夜勤_夜間", [2]string{"17:00", "09:00"}, "22:00", true},
		{"夜勤_早朝", [2]string{"17:00", "09:00"}, "06:00", true},
		{"夜勤_日中", [2]string{"17:00", "09:00"}, "12:00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewShiftTimeRangeFromStrings(tt.r[0], tt.r[1])
			time := MustParseTimeOfDay(tt.time)
			if got := r.Contains(time); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShiftTimeRange_String(t *testing.T) {
	r, _ := NewShiftTimeRangeFromStrings("09:00", "17:00")
	want := "09:00-17:00"
	if got := r.String(); got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}

func TestShiftTimeRange_Equals(t *testing.T) {
	r1, _ := NewShiftTimeRangeFromStrings("09:00", "17:00")
	r2, _ := NewShiftTimeRangeFromStrings("09:00", "17:00")
	r3, _ := NewShiftTimeRangeFromStrings("09:00", "18:00")

	if !r1.Equals(r2) {
		t.Error("同一範囲がEqualsでfalse")
	}
	if r1.Equals(r3) {
		t.Error("異なる範囲がEqualsでtrue")
	}
}
