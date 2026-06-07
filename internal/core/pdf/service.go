package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/datetime"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/gotenberg"
)

type Service struct {
	gotenberg *gotenberg.Client
}

func NewService(gotenberg *gotenberg.Client) *Service {
	return &Service{
		gotenberg: gotenberg,
	}
}

// GenerateFromContent renders an HTML template string with data and converts it
// to PDF.
func (s *Service) GenerateFromContent(
	ctx context.Context,
	tmplName string,
	tmplContent string,
	data interface{},
) ([]byte, error) {
	tmpl := template.New(tmplName).Funcs(getTemplateFuncs())

	tmpl, err := tmpl.Parse(tmplContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	pdfBytes, err := s.gotenberg.ConvertHTML(ctx, buf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to convert html to pdf: %w", err)
	}

	return pdfBytes, nil
}

func getTemplateFuncs() template.FuncMap {
	funcs := template.FuncMap{}

	// Merges multiple FuncMaps
	maps := []template.FuncMap{
		getBasicHelpers(),
		getDateHelpers(),
		getSupportHelpers(),
		getActivityHelpers(),
		getRelationHelpers(),
		getEduHelpers(),
		getConsultHelpers(),
	}

	for _, m := range maps {
		for k, v := range m {
			funcs[k] = v
		}
	}

	return funcs
}

func getBasicHelpers() template.FuncMap {
	return template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"ptrInt": func(i *int) int {
			if i == nil {
				return 0
			}
			return *i
		},
		"makeSlice": func(args ...interface{}) []interface{} {
			return args
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"iterate": func(n int) []int {
			var i []int
			for j := 0; j < n; j++ {
				i = append(i, j)
			}
			return i
		},
		"sumInt": func(slice interface{}, field string) int {
			v := reflect.ValueOf(slice)
			if v.Kind() != reflect.Slice {
				return 0
			}
			sum := 0
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				f := item.FieldByName(field)
				if f.IsValid() && f.Kind() == reflect.Int {
					sum += int(f.Int())
				}
			}
			return sum
		},
		"sumFloat": func(slice interface{}, field string) float64 {
			v := reflect.ValueOf(slice)
			if v.Kind() != reflect.Slice {
				return 0
			}
			sum := 0.0
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				f := item.FieldByName(field)
				if f.IsValid() && (f.Kind() == reflect.Float64 ||
					f.Kind() == reflect.Float32) {
					sum += f.Float()
				}
			}
			return sum
		},
	}
}

func getDateHelpers() template.FuncMap {
	return template.FuncMap{
		"formatDate": func(dateStr string) string {
			return datetime.FormatDate(dateStr)
		},
		"calculateAge": func(birthDate string) int {
			if birthDate == "" {
				return 0
			}
			if len(birthDate) > 10 {
				birthDate = birthDate[:10]
			}
			t, err := time.Parse("2006-01-02", birthDate)
			if err != nil {
				return 0
			}
			today := time.Now()
			age := today.Year() - t.Year()
			if today.YearDay() < t.YearDay() {
				age--
			}
			return age
		},
	}
}

func getSupportHelpers() template.FuncMap {
	return template.FuncMap{
		"hasSiblingSupport": func(supports interface{}, target string) bool {
			v := reflect.ValueOf(supports)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				if item.FieldByName("Name").String() == target {
					return true
				}
			}
			return false
		},
		"hasFinancialSupport": func(supports interface{}, target string) bool {
			v := reflect.ValueOf(supports)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				if item.FieldByName("Name").String() == target {
					return true
				}
			}
			return false
		},
	}
}

func getActivityHelpers() template.FuncMap {
	return template.FuncMap{
		"hasActivity": func(activities interface{}, target string) bool {
			v := reflect.ValueOf(activities)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				opt := item.FieldByName("ActivityOption")
				if opt.IsValid() {
					name := opt.FieldByName("Name").String()
					if strings.EqualFold(name, target) {
						return true
					}
				}
			}
			return false
		},
		"hasActivityCategory": func(
			activities interface{},
			target string,
			category string,
		) bool {
			v := reflect.ValueOf(activities)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				opt := item.FieldByName("ActivityOption")
				if opt.IsValid() {
					name := opt.FieldByName("Name").String()
					cat := opt.FieldByName("Category").String()
					if strings.EqualFold(name, target) &&
						strings.EqualFold(cat, category) {
						return true
					}
				}
			}
			return false
		},
		"indexHobby": func(hobbies interface{}, index int) string {
			v := reflect.ValueOf(hobbies)
			if v.Kind() != reflect.Slice {
				return ""
			}
			if index >= 0 && index < v.Len() {
				return v.Index(index).FieldByName("HobbyName").String()
			}
			return ""
		},
	}
}

func getRelationHelpers() template.FuncMap {
	return template.FuncMap{
		"getRelatedPerson": func(persons interface{}, role string) interface{} {
			v := reflect.ValueOf(persons)
			if v.Kind() != reflect.Slice {
				return nil
			}
			for i := 0; i < v.Len(); i++ {
				p := v.Index(i)
				isParent := p.FieldByName("IsParent").Bool()
				isGuardian := p.FieldByName("IsGuardian").Bool()
				relName := p.FieldByName("Relationship").
					FieldByName("Name").String()

				switch role {
				case "Father":
					if (isParent && relName == "Father") ||
						relName == "Father" {
						return p.Interface()
					}
				case "Mother":
					if (isParent && relName == "Mother") ||
						relName == "Mother" {
						return p.Interface()
					}
				case "Guardian":
					if isGuardian || relName == "Guardian" {
						return p.Interface()
					}
				}
			}
			return nil
		},
	}
}

func getEduHelpers() template.FuncMap {
	return template.FuncMap{
		"getEduBackground": func(
			schools interface{}, level string,
		) interface{} {
			v := reflect.ValueOf(schools)
			if v.Kind() != reflect.Slice {
				return nil
			}
			for i := 0; i < v.Len(); i++ {
				s := v.Index(i)
				lvl := s.FieldByName("EducationalLevel").
					FieldByName("Name").String()
				if lvl == level {
					return s.Interface()
				}
			}
			return nil
		},
	}
}

func getConsultHelpers() template.FuncMap {
	return template.FuncMap{
		"getConsultation": func(
			consults interface{}, profType string,
		) interface{} {
			v := reflect.ValueOf(consults)
			if v.Kind() != reflect.Slice {
				return nil
			}
			for i := 0; i < v.Len(); i++ {
				c := v.Index(i)
				ctype := c.FieldByName("ProfessionalType").String()
				if ctype == profType {
					return c.Interface()
				}
			}
			return nil
		},
	}
}
