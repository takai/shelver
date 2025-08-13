package parser

import (
	"testing"
)

func TestParseGroup_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		prefix   string
		want     string
		wantOK   bool
	}{
		{
			name:     "empty filename",
			filename: "",
			prefix:   "",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "only extension",
			filename: ".txt",
			prefix:   "",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "only digits",
			filename: "123.txt",
			prefix:   "",
			want:     "1",
			wantOK:   true,
		},
		{
			name:     "single character before digits",
			filename: "a99.txt",
			prefix:   "",
			want:     "a",
			wantOK:   true,
		},
		{
			name:     "multiple separators",
			filename: "test---___...   99.txt",
			prefix:   "",
			want:     "test",
			wantOK:   true,
		},
		{
			name:     "prefix at start of filename",
			filename: "p123.txt",
			prefix:   "p",
			want:     "p",
			wantOK:   true,
		},
		{
			name:     "prefix same as entire group",
			filename: "p-p99.txt",
			prefix:   "p",
			want:     "p",
			wantOK:   true,
		},
		{
			name:     "very long numbers",
			filename: "test-123456789012345.txt",
			prefix:   "",
			want:     "test",
			wantOK:   true,
		},
		{
			name:     "numbers in middle",
			filename: "test123file.txt",
			prefix:   "",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "special characters in prefix",
			filename: "test-[special]99.txt",
			prefix:   "[special]",
			want:     "test",
			wantOK:   true,
		},
		{
			name:     "regex metacharacters in prefix",
			filename: "test-.*+99.txt",
			prefix:   ".*+",
			want:     "test",
			wantOK:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseGroup(tt.filename, tt.prefix)
			if ok != tt.wantOK {
				t.Errorf("ParseGroup() ok = %v, want %v", ok, tt.wantOK)
				return
			}
			if got != tt.want {
				t.Errorf("ParseGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
