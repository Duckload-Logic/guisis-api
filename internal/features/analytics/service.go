package analytics

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/pdf"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo  *Repository
	redis *datastore.RedisClient
	pdf   *pdf.Service
}

func NewService(
	repo *Repository,
	redis *datastore.RedisClient,
	pdf *pdf.Service,
) *Service {
	return &Service{repo: repo, redis: redis, pdf: pdf}
}

func (s *Service) GetIIRAnalyticsReport(
	ctx context.Context,
	year int,
	courseID int,
) (*IIRAnalyticsReportResponse, error) {
	total, err := s.repo.GetTotalStudents(ctx, year, courseID)
	if err != nil {
		return nil, err
	}

	report := &IIRAnalyticsReportResponse{
		TotalStudents:        total,
		AgeDistribution:      []DemographicStatDTO{},
		CivilStatus:          []DemographicStatDTO{},
		Religions:            []DemographicStatDTO{},
		CityAddress:          []DemographicStatDTO{},
		MonthlyIncome:        []DemographicStatDTO{},
		OrdinalPosition:      []DemographicStatDTO{},
		FatherEducation:      []DemographicStatDTO{},
		MotherEducation:      []DemographicStatDTO{},
		ParentsMaritalStatus: []DemographicStatDTO{},
		HighSchoolGWA:        []DemographicStatDTO{},
		Elementary:           []DemographicStatDTO{},
		JuniorHigh:           []DemographicStatDTO{},
		SeniorHigh:           []DemographicStatDTO{},
		Vocational:           []DemographicStatDTO{},
		College:              []DemographicStatDTO{},
		NatureOfSchooling:    []DemographicStatDTO{},
		QuietStudyPlace:      []DemographicStatDTO{},
		GenderDistribution:   []DemographicStatDTO{},
	}

	if total > 0 {
		// Demographic data
		rawGender, _ := s.repo.GetGenderStats(ctx, year, courseID)
		report.GenderDistribution = s.mapToDTO(rawGender, total)

		rawAges, _ := s.repo.GetAgeStats(ctx, year, courseID)
		report.AgeDistribution = s.mapToDTO(rawAges, total)

		rawCivilStatus, _ := s.repo.GetCivilStatusStats(ctx, year, courseID)
		report.CivilStatus = s.mapToDTO(rawCivilStatus, total)

		rawReligions, _ := s.repo.GetReligionStats(ctx, year, courseID)
		report.Religions = s.mapToDTO(rawReligions, total)

		rawCityAddress, _ := s.repo.GetCityAddressStats(ctx, year, courseID)
		report.CityAddress = s.mapToDTO(rawCityAddress, total)

		// Economic/Social data
		rawMonthlyIncome, _ := s.repo.GetMonthlyIncomeStats(ctx, year, courseID)
		report.MonthlyIncome = s.mapToDTO(rawMonthlyIncome, total)

		rawOrdinalPosition, _ := s.repo.GetOrdinalPositionStats(
			ctx,
			year,
			courseID,
		)
		report.OrdinalPosition = s.mapToDTO(rawOrdinalPosition, total)

		rawQuietPlace, _ := s.repo.GetQuietStudyPlaceStats(ctx, year, courseID)
		report.QuietStudyPlace = s.mapToDTO(rawQuietPlace, total)

		// Family data
		rawFatherEd, _ := s.repo.GetFatherEducationStats(ctx, year, courseID)
		report.FatherEducation = s.mapToDTO(rawFatherEd, total)

		rawMotherEd, _ := s.repo.GetMotherEducationStats(ctx, year, courseID)
		report.MotherEducation = s.mapToDTO(rawMotherEd, total)

		rawParentsMarital, _ := s.repo.GetParentsMaritalStatusStats(
			ctx,
			year,
			courseID,
		)
		report.ParentsMaritalStatus = s.mapToDTO(rawParentsMarital, total)

		// Academic data
		rawHSGWA, _ := s.repo.GetHSGWAStats(ctx, year, courseID)
		report.HighSchoolGWA = s.mapToDTO(rawHSGWA, total)

		rawElem, _ := s.repo.GetElementaryStats(ctx, year, courseID)
		report.Elementary = s.mapToDTO(rawElem, total)

		rawHS, _ := s.repo.GetHighSchoolStats(ctx, year, courseID)
		report.HighSchool = s.mapToDTO(rawHS, total)

		rawJHS, _ := s.repo.GetJuniorHighStats(ctx, year, courseID)
		report.JuniorHigh = s.mapToDTO(rawJHS, total)

		rawSHS, _ := s.repo.GetSeniorHighStats(ctx, year, courseID)
		report.SeniorHigh = s.mapToDTO(rawSHS, total)

		rawVocational, _ := s.repo.GetVocationalStats(ctx, year, courseID)
		report.Vocational = s.mapToDTO(rawVocational, total)

		rawCollege, _ := s.repo.GetCollegeStats(ctx, year, courseID)
		report.College = s.mapToDTO(rawCollege, total)

		rawNature, _ := s.repo.GetNatureOfSchoolingStats(ctx, year, courseID)
		report.NatureOfSchooling = s.mapToDTO(rawNature, total)
	}

	return report, nil
}

func (s *Service) GetAdminDashboard(
	ctx context.Context,
	timeRange string,
	source string,
) (*AdminDashboardResponse, error) {
	totalStudents, _ := s.repo.GetTotalStudents(ctx, 0, 0)
	studentsTrend, _ := s.repo.GetStudentsTrend(ctx)

	totalReports, _ := s.repo.GetTotalReports(ctx)
	reportsTrend, _ := s.repo.GetReportsTrend(ctx)

	totalAppointments, _ := s.repo.GetTotalAppointments(ctx)
	appointmentsTrend, _ := s.repo.GetAppointmentsTrend(ctx)

	totalSlips, _ := s.repo.GetTotalSlips(ctx)
	slipsTrend, _ := s.repo.GetSlipsTrend(ctx)

	var err error
	var monthlyVisitors []MonthlyVisitorStatDTO
	if source == "system" {
		monthlyVisitors, err = s.repo.GetMonthlyVisitorStats(ctx, timeRange)
	} else {
		monthlyVisitors, err = s.repo.GetMonthlyAppointmentStats(ctx, timeRange)
	}

	if err != nil {
		return nil, err
	}

	// Count live sessions (session: prefix)
	liveSessions := 0
	if s.redis != nil {
		keys, err := s.redis.Keys(ctx, "session:*")
		if err == nil {
			liveSessions = len(keys)
		}
	}

	return &AdminDashboardResponse{
		TotalStudents:     totalStudents,
		StudentsTrend:     studentsTrend,
		TotalReports:      totalReports,
		ReportsTrend:      reportsTrend,
		TotalAppointments: totalAppointments,
		AppointmentsTrend: appointmentsTrend,
		TotalSlips:        totalSlips,
		SlipsTrend:        slipsTrend,
		LiveSessions:      liveSessions,
		MonthlyVisitors:   monthlyVisitors,
	}, nil
}

func (s *Service) mapToDTO(
	rawStats []DemographicStat,
	totalStudents int,
) []DemographicStatDTO {
	dtos := make([]DemographicStatDTO, 0)

	if totalStudents == 0 {
		return dtos
	}

	minRank := math.MaxFloat64
	for _, stat := range rawStats {
		dto := DemographicStatDTO{
			Category:    stat.Category,
			MaleCount:   stat.MaleCount,
			FemaleCount: stat.FemaleCount,
			Total:       stat.Total,
			Rank:        stat.RankPos,
		}

		dto.TotalPct = s.calculatePercentage(stat.Total, totalStudents)
		dto.MalePct = s.calculatePercentage(stat.MaleCount, totalStudents)
		dto.FemalePct = s.calculatePercentage(stat.FemaleCount, totalStudents)

		if stat.RankPos < minRank && stat.RankPos > 0 {
			minRank = stat.RankPos
		}

		dtos = append(dtos, dto)
	}

	for i := range dtos {
		if dtos[i].Rank == minRank {
			dtos[i].IsTop = true
		}
	}

	return dtos
}

func (s *Service) calculatePercentage(count, total int) float64 {
	if total == 0 {
		return 0
	}
	percent := (float64(count) / float64(total)) * 100
	return math.Round(percent*100) / 100
}

func (s *Service) ExportIIRAnalyticsReport(
	ctx context.Context,
	year int,
	courseID int,
) ([]byte, error) {
	data, err := s.GetIIRAnalyticsReport(ctx, year, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analytics data: %w", err)
	}

	// For simplicity in this capstone, we load the template from disk.
	// In production, this would be embedded.
	tmplPath := "internal/features/analytics/assets/report.html"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	logoPath := "internal/features/analytics/assets/pup-logo.png"
	logoBytes, err := os.ReadFile(logoPath)
	logoBase64 := ""
	if err == nil {
		logoBase64 = base64.StdEncoding.EncodeToString(logoBytes)
	}

	ageCats, agePct := getTopCategoriesText(data.AgeDistribution)
	civilCats, civilPct := getTopCategoriesText(data.CivilStatus)
	cityCats, cityPct := getTopCategoriesText(data.CityAddress)
	religionCats, religionPct := getTopCategoriesText(data.Religions)
	elemCats, elemPct := getTopCategoriesText(data.Elementary)
	hsgwaCats, hsgwaPct := getTopCategoriesText(data.HighSchoolGWA)
	fatherCats, fatherPct := getTopCategoriesText(data.FatherEducation)
	motherCats, motherPct := getTopCategoriesText(data.MotherEducation)
	parentsMaritalCats, parentsMaritalPct := getTopCategoriesText(
		data.ParentsMaritalStatus,
	)
	incomeCats, incomePct := getTopCategoriesText(data.MonthlyIncome)
	ordinalCats, ordinalPct := getTopCategoriesText(data.OrdinalPosition)

	hscaCats, hscaPct := getTopCategoriesText(data.HighSchool)
	vocCats, vocPct := getTopCategoriesText(data.Vocational)
	collegeCats, collegePct := getTopCategoriesText(data.College)
	natureCats, naturePct := getTopCategoriesText(data.NatureOfSchooling)
	quietCats, quietPct := getTopCategoriesText(data.QuietStudyPlace)

	reportData := struct {
		Data       *IIRAnalyticsReportResponse
		DateToday  string
		Year       int
		LogoBase64 string

		AgeTopCategories            string
		AgeTopPct                   float64
		CivilTopCategories          string
		CivilTopPct                 float64
		CityTopCategories           string
		CityTopPct                  float64
		ReligionTopCategories       string
		ReligionTopPct              float64
		ElemTopCategories           string
		ElemTopPct                  float64
		HSGWATopCategories          string
		HSGWATopPct                 float64
		FatherTopCategories         string
		FatherTopPct                float64
		MotherTopCategories         string
		MotherTopPct                float64
		ParentsMaritalTopCategories string
		ParentsMaritalTopPct        float64
		IncomeTopCategories         string
		IncomeTopPct                float64
		OrdinalTopCategories        string
		OrdinalTopPct               float64

		HSCATopCategories              string
		HSCATopPct                     float64
		VocTopCategories               string
		VocTopPct                      float64
		CollegeTopCategories           string
		CollegeTopPct                  float64
		NatureOfSchoolingTopCategories string
		NatureOfSchoolingTopPct        float64
		QuietStudyPlaceTopCategories   string
		QuietStudyPlaceTopPct          float64
	}{
		Data:       data,
		DateToday:  time.Now().Format("January 02, 2006"),
		Year:       year,
		LogoBase64: logoBase64,

		AgeTopCategories:            ageCats,
		AgeTopPct:                   agePct,
		CivilTopCategories:          civilCats,
		CivilTopPct:                 civilPct,
		CityTopCategories:           cityCats,
		CityTopPct:                  cityPct,
		ReligionTopCategories:       religionCats,
		ReligionTopPct:              religionPct,
		ElemTopCategories:           elemCats,
		ElemTopPct:                  elemPct,
		HSGWATopCategories:          hsgwaCats,
		HSGWATopPct:                 hsgwaPct,
		FatherTopCategories:         fatherCats,
		FatherTopPct:                fatherPct,
		MotherTopCategories:         motherCats,
		MotherTopPct:                motherPct,
		ParentsMaritalTopCategories: parentsMaritalCats,
		ParentsMaritalTopPct:        parentsMaritalPct,
		IncomeTopCategories:         incomeCats,
		IncomeTopPct:                incomePct,
		OrdinalTopCategories:        ordinalCats,
		OrdinalTopPct:               ordinalPct,

		HSCATopCategories:              hscaCats,
		HSCATopPct:                     hscaPct,
		VocTopCategories:               vocCats,
		VocTopPct:                      vocPct,
		CollegeTopCategories:           collegeCats,
		CollegeTopPct:                  collegePct,
		NatureOfSchoolingTopCategories: natureCats,
		NatureOfSchoolingTopPct:        naturePct,
		QuietStudyPlaceTopCategories:   quietCats,
		QuietStudyPlaceTopPct:          quietPct,
	}

	pdfBytes, err := s.pdf.GenerateFromContent(
		ctx,
		"analytics_report",
		string(tmplBytes),
		reportData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pdf: %w", err)
	}

	return pdfBytes, nil
}

func getTopCategoriesText(
	stats []DemographicStatDTO,
) (string, float64) {
	if len(stats) == 0 {
		return "N/A", 0
	}

	var labels []string
	var topPct float64

	for _, s := range stats {
		if s.IsTop {
			labels = append(labels, s.Category)
			topPct = s.TotalPct
		}
	}

	if len(labels) == 0 {
		return "N/A", 0
	}

	return strings.Join(labels, " and "), topPct
}
