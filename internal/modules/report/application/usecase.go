// Package application レポートアプリケーション層
package application

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"shiftmaster/internal/modules/report/domain"
	sharedDomain "shiftmaster/internal/shared/domain"
)

// ReportUseCase レポートユースケース
type ReportUseCase struct {
	reportRepo  domain.ReportRepository
	summaryRepo domain.SummaryRepository
	generator   domain.ReportGenerator
	logger      *slog.Logger
}

// NewReportUseCase レポートユースケース生成
func NewReportUseCase(
	reportRepo domain.ReportRepository,
	summaryRepo domain.SummaryRepository,
	generator domain.ReportGenerator,
	logger *slog.Logger,
) *ReportUseCase {
	return &ReportUseCase{
		reportRepo:  reportRepo,
		summaryRepo: summaryRepo,
		generator:   generator,
		logger:      logger,
	}
}

// Generate レポート生成
func (u *ReportUseCase) Generate(ctx context.Context, input *GenerateReportInput) (*ReportOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	orgID, err := sharedDomain.ParseID(input.OrganizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	reportType := domain.ReportType(input.Type)
	title := fmt.Sprintf("%s %d年%d月", reportType.Label(), input.TargetYear, input.TargetMonth)

	now := time.Now()
	report := &domain.Report{
		ID:             sharedDomain.NewID(),
		OrganizationID: orgID,
		Type:           reportType,
		Title:          title,
		TargetYear:     input.TargetYear,
		TargetMonth:    input.TargetMonth,
		Status:         domain.ReportStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.reportRepo.Save(ctx, report); err != nil {
		u.logger.Error("レポート作成失敗", "error", err)
		return nil, err
	}

	// 非同期で生成処理 実際の実装ではgoroutineやジョブキュー使用
	go u.generateAsync(context.Background(), report)

	u.logger.Info("レポート生成開始", "report_id", report.ID)
	return ToReportOutput(report), nil
}

// generateAsync 非同期レポート生成
func (u *ReportUseCase) generateAsync(ctx context.Context, report *domain.Report) {
	// 状態更新 生成中
	report.Status = domain.ReportStatusGenerating
	report.UpdatedAt = time.Now()
	if err := u.reportRepo.Save(ctx, report); err != nil {
		u.logger.Error("レポート状態更新失敗", "error", err)
		return
	}

	var err error
	switch report.Type {
	case domain.ReportTypeSummary:
		_, err = u.summaryRepo.GetMonthlySummary(ctx, report.OrganizationID, report.TargetYear, report.TargetMonth)
	}

	now := time.Now()
	if err != nil {
		report.Status = domain.ReportStatusFailed
		u.logger.Error("レポート生成失敗", "error", err, "report_id", report.ID)
	} else {
		report.Status = domain.ReportStatusCompleted
		report.GeneratedAt = &now
		u.logger.Info("レポート生成完了", "report_id", report.ID)
	}

	report.UpdatedAt = now
	if err := u.reportRepo.Save(ctx, report); err != nil {
		u.logger.Error("レポート状態更新失敗", "error", err)
	}
}

// GetByID IDでレポート取得
func (u *ReportUseCase) GetByID(ctx context.Context, id string) (*ReportOutput, error) {
	reportID, err := sharedDomain.ParseID(id)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	report, err := u.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, sharedDomain.ErrNotFound
	}

	return ToReportOutput(report), nil
}

// List レポート一覧取得
func (u *ReportUseCase) List(ctx context.Context, organizationID string) ([]ReportOutput, error) {
	orgID, err := sharedDomain.ParseID(organizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	reports, err := u.reportRepo.FindByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ReportOutput, len(reports))
	for i, r := range reports {
		outputs[i] = *ToReportOutput(&r)
	}

	return outputs, nil
}

// Delete レポート削除
func (u *ReportUseCase) Delete(ctx context.Context, id string) error {
	reportID, err := sharedDomain.ParseID(id)
	if err != nil {
		return sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "IDが不正です")
	}

	report, err := u.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return err
	}
	if report == nil {
		return sharedDomain.ErrNotFound
	}

	if err := u.reportRepo.Delete(ctx, reportID); err != nil {
		u.logger.Error("レポート削除失敗", "error", err)
		return err
	}

	u.logger.Info("レポート削除完了", "report_id", reportID)
	return nil
}

// GetMonthlySummary 月次集計取得
func (u *ReportUseCase) GetMonthlySummary(ctx context.Context, organizationID string, year, month int) (*MonthlySummaryOutput, error) {
	orgID, err := sharedDomain.ParseID(organizationID)
	if err != nil {
		return nil, sharedDomain.NewDomainError(sharedDomain.ErrCodeValidation, "組織IDが不正です")
	}

	summary, err := u.summaryRepo.GetMonthlySummary(ctx, orgID, year, month)
	if err != nil {
		return nil, err
	}
	if summary == nil {
		// 空の集計を返す
		summary = &domain.MonthlySummary{
			TargetYear:         year,
			TargetMonth:        month,
			StaffSummaries:     []domain.StaffSummary{},
			DailySummaries:     []domain.DailySummary{},
			ShiftTypeSummaries: []domain.ShiftTypeSummary{},
		}
	}

	return ToMonthlySummaryOutput(summary), nil
}
