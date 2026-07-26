package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMMYYYYDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantStr string
	}{
		{
			name:    "valid",
			input:   `"01-2026"`,
			wantErr: false,
			wantStr: "01-2026",
		},
		{
			name:    "valid december",
			input:   `"12-2025"`,
			wantErr: false,
			wantStr: "12-2025",
		},
		{
			name:    "invalid format",
			input:   `"2026-01"`,
			wantErr: true,
		},
		{
			name:    "invalid month",
			input:   `"13-2026"`,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   `""`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d MMYYYYDate
			err := d.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantStr, d.Time.Format(DATE_FORMAT))
		})
	}
}

func TestMMYYYYDate_MarshalJSON(t *testing.T) {
	d := MMYYYYDate{
		Time: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	b, err := d.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"01-2026"`, string(b))
}
