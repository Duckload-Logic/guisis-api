package analytics

type DemographicStatDTO struct {
	Category    string  `json:"category"`
	MaleCount   int     `json:"maleCount"`
	MalePct     float64 `json:"malePct"`
	FemaleCount int     `json:"femaleCount"`
	FemalePct   float64 `json:"femalePct"`
	Total       int     `json:"total"`
	TotalPct    float64 `json:"totalPct"`
	Rank        int     `json:"rank,omitempty"`
}

type DashboardResponse struct {
	TotalStudents      int                  `json:"totalStudents"`
	GenderDistribution []DemographicStatDTO `json:"genderDistribution"`

	// Personal Information
	AgeDistribution []DemographicStatDTO `json:"ageDistribution"`
	CivilStatus     []DemographicStatDTO `json:"civilStatus"`
	Religions       []DemographicStatDTO `json:"religions"`
	CityAddress     []DemographicStatDTO `json:"cityAddress"`

	// Family & Financial Background
	MonthlyIncome        []DemographicStatDTO `json:"monthlyIncome"`
	OrdinalPosition      []DemographicStatDTO `json:"ordinalPosition"`
	FatherEducation      []DemographicStatDTO `json:"fatherEducation"`
	MotherEducation      []DemographicStatDTO `json:"motherEducation"`
	ParentsMaritalStatus []DemographicStatDTO `json:"parentsMaritalStatus"`

	// Academic Background
	HighSchoolGWA     []DemographicStatDTO `json:"highSchoolGWA"`
	Elementary        []DemographicStatDTO `json:"elementary"`
	JuniorHigh        []DemographicStatDTO `json:"juniorHigh"`
	SeniorHigh        []DemographicStatDTO `json:"seniorHigh"`
	NatureOfSchooling []DemographicStatDTO `json:"natureOfSchooling"`

	// Study Environment
	QuietStudyPlace []DemographicStatDTO `json:"quietStudyPlace"`
}

type MonthlyVisitorStatDTO struct {
	Period   string `json:"period"`   // Labels (Daily, Weekly, Monthly, Yearly)
	Month    string `json:"month"`    // BC for react-web (Monthly only)
	Logins   int    `json:"logins"`   // System Traffic specific
	Activity int    `json:"activity"` // System Traffic specific
	Count    int    `json:"count"`    // Generic count (Appointments or Logins)
}

type AdminDashboardResponse struct {
	TotalStudents     int                     `json:"totalStudents"`
	TotalReports      int                     `json:"totalReports"`
	TotalAppointments int                     `json:"totalAppointments"`
	TotalSlips        int                     `json:"totalSlips"`
	LiveSessions      int                     `json:"liveSessions"`
	MonthlyVisitors   []MonthlyVisitorStatDTO `json:"monthlyVisitors"`
}
