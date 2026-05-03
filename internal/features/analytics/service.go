package analytics

import (
	"context"
	"math"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo  *Repository
	redis *datastore.RedisClient
}

func NewService(
	repo *Repository,
	redis *datastore.RedisClient,
) *Service {
	return &Service{repo: repo, redis: redis}
}

func (s *Service) GetDashboard(
	ctx context.Context,
	year int,
	courseID int,
) (*DashboardResponse, error) {
	total, err := s.repo.GetTotalStudents(ctx, year, courseID)
	if err != nil {
		return nil, err
	}

	dashboard := &DashboardResponse{
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
		NatureOfSchooling:    []DemographicStatDTO{},
		QuietStudyPlace:      []DemographicStatDTO{},
		GenderDistribution:   []DemographicStatDTO{},
	}

	if total > 0 {
		// Demographic data
		rawGender, _ := s.repo.GetGenderStats(ctx, year, courseID)
		dashboard.GenderDistribution = s.mapToDTO(rawGender, total)

		rawAges, _ := s.repo.GetAgeStats(ctx, year, courseID)
		dashboard.AgeDistribution = s.mapToDTO(rawAges, total)

		rawCivilStatus, _ := s.repo.GetCivilStatusStats(ctx, year, courseID)
		dashboard.CivilStatus = s.mapToDTO(rawCivilStatus, total)

		rawReligions, _ := s.repo.GetReligionStats(ctx, year, courseID)
		dashboard.Religions = s.mapToDTO(rawReligions, total)

		rawCityAddress, _ := s.repo.GetCityAddressStats(ctx, year, courseID)
		dashboard.CityAddress = s.mapToDTO(rawCityAddress, total)

		// Economic/Social data
		rawMonthlyIncome, _ := s.repo.GetMonthlyIncomeStats(ctx, year, courseID)
		dashboard.MonthlyIncome = s.mapToDTO(rawMonthlyIncome, total)

		rawOrdinalPosition, _ := s.repo.GetOrdinalPositionStats(
			ctx,
			year,
			courseID,
		)
		dashboard.OrdinalPosition = s.mapToDTO(rawOrdinalPosition, total)

		rawQuietPlace, _ := s.repo.GetQuietStudyPlaceStats(ctx, year, courseID)
		dashboard.QuietStudyPlace = s.mapToDTO(rawQuietPlace, total)

		// Family data
		rawFatherEd, _ := s.repo.GetFatherEducationStats(ctx, year, courseID)
		dashboard.FatherEducation = s.mapToDTO(rawFatherEd, total)

		rawMotherEd, _ := s.repo.GetMotherEducationStats(ctx, year, courseID)
		dashboard.MotherEducation = s.mapToDTO(rawMotherEd, total)

		rawParentsMarital, _ := s.repo.GetParentsMaritalStatusStats(
			ctx,
			year,
			courseID,
		)
		dashboard.ParentsMaritalStatus = s.mapToDTO(rawParentsMarital, total)

		// Academic data
		rawHSGWA, _ := s.repo.GetHSGWAStats(ctx, year, courseID)
		dashboard.HighSchoolGWA = s.mapToDTO(rawHSGWA, total)

		rawElem, _ := s.repo.GetElementaryStats(ctx, year, courseID)
		dashboard.Elementary = s.mapToDTO(rawElem, total)

		rawJHS, _ := s.repo.GetJuniorHighStats(ctx, year, courseID)
		dashboard.JuniorHigh = s.mapToDTO(rawJHS, total)

		rawSHS, _ := s.repo.GetSeniorHighStats(ctx, year, courseID)
		dashboard.SeniorHigh = s.mapToDTO(rawSHS, total)

		rawNature, _ := s.repo.GetNatureOfSchoolingStats(ctx, year, courseID)
		dashboard.NatureOfSchooling = s.mapToDTO(rawNature, total)
	}

	return dashboard, nil
}

func (s *Service) GetAdminDashboard(
	ctx context.Context,
	timeRange string,
	source string,
) (*AdminDashboardResponse, error) {
	totalStudents, err := s.repo.GetTotalStudents(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	totalReports, err := s.repo.GetTotalReports(ctx)
	if err != nil {
		return nil, err
	}

	totalAppointments, err := s.repo.GetTotalAppointments(ctx)
	if err != nil {
		return nil, err
	}

	totalSlips, err := s.repo.GetTotalSlips(ctx)
	if err != nil {
		return nil, err
	}

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
		TotalReports:      totalReports,
		TotalAppointments: totalAppointments,
		TotalSlips:        totalSlips,
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

		dtos = append(dtos, dto)
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
