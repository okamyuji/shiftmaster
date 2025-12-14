package domain

import (
	"testing"

	"shiftmaster/internal/shared/domain"
)

func TestNewPosition(t *testing.T) {
	orgID := domain.NewID()

	tests := []struct {
		name    string
		orgID   domain.ID
		pname   string
		code    string
		level   int
		wantErr bool
	}{
		{"正常系", orgID, "師長", "DIR", 1, false},
		{"正常系_レベル0", orgID, "一般", "STF", 0, false},
		{"異常系_名前なし", orgID, "", "DIR", 1, true},
		{"異常系_コードなし", orgID, "師長", "", 1, true},
		{"異常系_レベル負数", orgID, "師長", "DIR", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewPosition(tt.orgID, tt.pname, tt.code, tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if pos.Name != tt.pname {
					t.Errorf("Name = %v, want %v", pos.Name, tt.pname)
				}
				if pos.Code != tt.code {
					t.Errorf("Code = %v, want %v", pos.Code, tt.code)
				}
				if pos.Level != tt.level {
					t.Errorf("Level = %v, want %v", pos.Level, tt.level)
				}
				if !pos.IsActive {
					t.Error("IsActive should be true by default")
				}
			}
		})
	}
}

func TestPosition_Update(t *testing.T) {
	orgID := domain.NewID()
	pos, _ := NewPosition(orgID, "師長", "DIR", 1)

	tests := []struct {
		name        string
		pname       string
		code        string
		description string
		level       int
		sortOrder   int
		wantErr     bool
	}{
		{"正常系", "師長更新", "DIR2", "説明", 2, 1, false},
		{"異常系_名前なし", "", "DIR", "", 1, 0, true},
		{"異常系_コードなし", "師長", "", "", 1, 0, true},
		{"異常系_レベル負数", "師長", "DIR", "", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pos.Update(tt.pname, tt.code, tt.description, tt.level, tt.sortOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPosition_IsHigherThan(t *testing.T) {
	orgID := domain.NewID()
	director, _ := NewPosition(orgID, "師長", "DIR", 1)
	manager, _ := NewPosition(orgID, "主任", "MGR", 2)
	leader, _ := NewPosition(orgID, "リーダー", "LDR", 3)
	staff, _ := NewPosition(orgID, "一般", "STF", 4)

	tests := []struct {
		name string
		p1   *Position
		p2   *Position
		want bool
	}{
		{"師長>主任", director, manager, true},
		{"師長>一般", director, staff, true},
		{"主任>師長", manager, director, false},
		{"主任>リーダー", manager, leader, true},
		{"一般>リーダー", staff, leader, false},
		{"同レベル", manager, manager, false},
		{"nil比較", director, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.IsHigherThan(tt.p2); got != tt.want {
				t.Errorf("IsHigherThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsEqualOrHigherThan(t *testing.T) {
	orgID := domain.NewID()
	director, _ := NewPosition(orgID, "師長", "DIR", 1)
	manager, _ := NewPosition(orgID, "主任", "MGR", 2)

	tests := []struct {
		name string
		p1   *Position
		p2   *Position
		want bool
	}{
		{"師長>=主任", director, manager, true},
		{"主任>=主任", manager, manager, true},
		{"主任>=師長", manager, director, false},
		{"nil比較", director, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.IsEqualOrHigherThan(tt.p2); got != tt.want {
				t.Errorf("IsEqualOrHigherThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_Deactivate(t *testing.T) {
	orgID := domain.NewID()
	pos, _ := NewPosition(orgID, "師長", "DIR", 1)

	if !pos.IsActive {
		t.Error("初期状態でIsActiveがfalse")
	}

	pos.Deactivate()
	if pos.IsActive {
		t.Error("Deactivate後もIsActiveがtrue")
	}

	pos.Activate()
	if !pos.IsActive {
		t.Error("Activate後もIsActiveがfalse")
	}
}
