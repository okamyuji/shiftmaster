package domain

import (
	"testing"

	"shiftmaster/internal/shared/domain"
)

func TestNewJobType(t *testing.T) {
	orgID := domain.NewID()

	tests := []struct {
		name    string
		orgID   domain.ID
		jname   string
		code    string
		wantErr bool
	}{
		{"正常系", orgID, "看護師", "NS", false},
		{"正常系_英語名", orgID, "Nurse", "NURSE", false},
		{"異常系_名前なし", orgID, "", "NS", true},
		{"異常系_コードなし", orgID, "看護師", "", true},
		{"異常系_両方なし", orgID, "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jt, err := NewJobType(tt.orgID, tt.jname, tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJobType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if jt.Name != tt.jname {
					t.Errorf("Name = %v, want %v", jt.Name, tt.jname)
				}
				if jt.Code != tt.code {
					t.Errorf("Code = %v, want %v", jt.Code, tt.code)
				}
				if jt.OrganizationID != tt.orgID {
					t.Errorf("OrganizationID = %v, want %v", jt.OrganizationID, tt.orgID)
				}
				if !jt.IsActive {
					t.Error("IsActive should be true by default")
				}
			}
		})
	}
}

func TestJobType_Update(t *testing.T) {
	orgID := domain.NewID()
	jt, _ := NewJobType(orgID, "看護師", "NS")

	tests := []struct {
		name        string
		jname       string
		code        string
		description string
		color       string
		sortOrder   int
		wantErr     bool
	}{
		{"正常系", "看護師更新", "NS2", "説明", "#FF0000", 1, false},
		{"正常系_色空", "看護師", "NS", "", "", 0, false},
		{"異常系_名前なし", "", "NS", "", "", 0, true},
		{"異常系_コードなし", "看護師", "", "", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jt.Update(tt.jname, tt.code, tt.description, tt.color, tt.sortOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if jt.Name != tt.jname {
					t.Errorf("Name = %v, want %v", jt.Name, tt.jname)
				}
				if jt.SortOrder != tt.sortOrder {
					t.Errorf("SortOrder = %v, want %v", jt.SortOrder, tt.sortOrder)
				}
			}
		})
	}
}

func TestJobType_Deactivate(t *testing.T) {
	orgID := domain.NewID()
	jt, _ := NewJobType(orgID, "看護師", "NS")

	if !jt.IsActive {
		t.Error("初期状態でIsActiveがfalse")
	}

	jt.Deactivate()
	if jt.IsActive {
		t.Error("Deactivate後もIsActiveがtrue")
	}

	jt.Activate()
	if !jt.IsActive {
		t.Error("Activate後もIsActiveがfalse")
	}
}
