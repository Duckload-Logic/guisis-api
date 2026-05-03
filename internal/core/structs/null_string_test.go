package structs

import (
	"database/sql"
	"encoding/json"
	"testing"
)

func TestNullableString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    NullableString
		wantErr bool
	}{
		{
			name: "valid string",
			data: `"hello"`,
			want: NullableString{String: "hello", Valid: true},
		},
		{
			name: "null value",
			data: `null`,
			want: NullableString{Valid: false},
		},
		{
			name:    "invalid json",
			data:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ns NullableString
			err := json.Unmarshal([]byte(tt.data), &ns)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"UnmarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				if ns.String != tt.want.String || ns.Valid != tt.want.Valid {
					t.Errorf("UnmarshalJSON() got = %v, want %v", ns, tt.want)
				}
			}
		})
	}
}

func TestNullableString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		ns   NullableString
		want string
	}{
		{
			name: "valid string",
			ns:   NullableString{String: "hello", Valid: true},
			want: `"hello"`,
		},
		{
			name: "null value",
			ns:   NullableString{Valid: false},
			want: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.ns)
			if err != nil {
				t.Fatalf("MarshalJSON() unexpected error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf(
					"MarshalJSON() got = %s, want %s",
					string(got),
					tt.want,
				)
			}
		})
	}
}

func TestStringToNullableString(t *testing.T) {
	t.Run("empty string returns invalid", func(t *testing.T) {
		got := StringToNullableString("")
		if got.Valid {
			t.Error("expected invalid NullableString for empty input")
		}
	})

	t.Run("non-empty string returns valid", func(t *testing.T) {
		got := StringToNullableString("test")
		if !got.Valid || got.String != "test" {
			t.Errorf("got %v, want test (valid)", got)
		}
	})
}

func TestSqlNullConversion(t *testing.T) {
	t.Run("FromSqlNull", func(t *testing.T) {
		sn := sql.NullString{String: "test", Valid: true}
		got := FromSqlNull(sn)
		if got.String != sn.String || got.Valid != sn.Valid {
			t.Errorf("got %v, want %v", got, sn)
		}
	})

	t.Run("ToSqlNull", func(t *testing.T) {
		ns := NullableString{String: "test", Valid: true}
		got := ToSqlNull(ns)
		if got.String != ns.String || got.Valid != ns.Valid {
			t.Errorf("got %v, want %v", got, ns)
		}
	})
}

func TestNullableInt64_JSON(t *testing.T) {
	t.Run("Unmarshal null", func(t *testing.T) {
		var ni NullableInt64
		err := json.Unmarshal([]byte("null"), &ni)
		if err != nil || ni.Valid {
			t.Errorf("got %v, err %v", ni, err)
		}
	})

	t.Run("Unmarshal valid", func(t *testing.T) {
		var ni NullableInt64
		err := json.Unmarshal([]byte("123"), &ni)
		if err != nil || !ni.Valid || ni.Int64 != 123 {
			t.Errorf("got %v, err %v", ni, err)
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		ni := NullableInt64{Int64: 456, Valid: true}
		got, _ := json.Marshal(ni)
		if string(got) != "456" {
			t.Errorf("got %s", string(got))
		}
	})
}
