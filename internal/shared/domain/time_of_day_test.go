package domain

import (
	"testing"
)

func TestNewTimeOfDay(t *testing.T) {
	tests := []struct {
		name    string
		hour    int
		minute  int
		wantErr bool
	}{
		{"正常系_午前9時", 9, 0, false},
		{"正常系_午後5時30分", 17, 30, false},
		{"正常系_深夜0時", 0, 0, false},
		{"正常系_23時59分", 23, 59, false},
		{"境界値_時間最小", 0, 0, false},
		{"境界値_時間最大", 23, 59, false},
		{"異常系_時間負数", -1, 0, true},
		{"異常系_時間24以上", 24, 0, true},
		{"異常系_分負数", 9, -1, true},
		{"異常系_分60以上", 9, 60, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tod, err := NewTimeOfDay(tt.hour, tt.minute)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTimeOfDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tod.Hour() != tt.hour || tod.Minute() != tt.minute {
					t.Errorf("NewTimeOfDay() = %v, want %02d:%02d", tod, tt.hour, tt.minute)
				}
			}
		})
	}
}

func TestParseTimeOfDay(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"正常系_09:00", "09:00", "09:00", false},
		{"正常系_17:30", "17:30", "17:30", false},
		{"正常系_00:00", "00:00", "00:00", false},
		{"正常系_23:59", "23:59", "23:59", false},
		{"異常系_空文字", "", "", true},
		{"異常系_コロンなし", "0900", "", true},
		{"異常系_時間が文字", "aa:00", "", true},
		{"異常系_分が文字", "09:bb", "", true},
		{"異常系_時間範囲外", "25:00", "", true},
		{"異常系_分範囲外", "09:60", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tod, err := ParseTimeOfDay(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeOfDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tod.String() != tt.want {
				t.Errorf("ParseTimeOfDay() = %v, want %v", tod.String(), tt.want)
			}
		})
	}
}

func TestTimeOfDay_ToMinutes(t *testing.T) {
	tests := []struct {
		name string
		time string
		want int
	}{
		{"00:00", "00:00", 0},
		{"01:00", "01:00", 60},
		{"09:30", "09:30", 570},
		{"12:00", "12:00", 720},
		{"23:59", "23:59", 1439},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tod := MustParseTimeOfDay(tt.time)
			if got := tod.ToMinutes(); got != tt.want {
				t.Errorf("ToMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeOfDay_IsBefore(t *testing.T) {
	tests := []struct {
		name  string
		time1 string
		time2 string
		want  bool
	}{
		{"9:00 < 17:00", "09:00", "17:00", true},
		{"17:00 < 09:00", "17:00", "09:00", false},
		{"09:00 < 09:00", "09:00", "09:00", false},
		{"09:00 < 09:01", "09:00", "09:01", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := MustParseTimeOfDay(tt.time1)
			t2 := MustParseTimeOfDay(tt.time2)
			if got := t1.IsBefore(t2); got != tt.want {
				t.Errorf("IsBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeOfDay_Add(t *testing.T) {
	tests := []struct {
		name    string
		time    string
		minutes int
		want    string
	}{
		{"加算_30分", "09:00", 30, "09:30"},
		{"加算_60分", "09:00", 60, "10:00"},
		{"加算_日跨ぎ", "23:00", 120, "01:00"},
		{"減算_30分", "09:30", -30, "09:00"},
		{"減算_日跨ぎ", "00:30", -60, "23:30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tod := MustParseTimeOfDay(tt.time)
			result := tod.Add(tt.minutes)
			if result.String() != tt.want {
				t.Errorf("Add(%d) = %v, want %v", tt.minutes, result.String(), tt.want)
			}
		})
	}
}

func TestTimeOfDay_Equals(t *testing.T) {
	t1 := MustParseTimeOfDay("09:00")
	t2 := MustParseTimeOfDay("09:00")
	t3 := MustParseTimeOfDay("09:01")

	if !t1.Equals(t2) {
		t.Error("同一時刻がEqualsでfalse")
	}
	if t1.Equals(t3) {
		t.Error("異なる時刻がEqualsでtrue")
	}
}
