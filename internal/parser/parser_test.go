package parser

import (
	"testing"
)

func TestParseGroup_NumericOnly(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		prefix   string
		want     string
		wantOK   bool
	}{
		{
			name:     "basic numeric suffix",
			filename: "album-001.wav",
			prefix:   "",
			want:     "album",
			wantOK:   true,
		},
		{
			name:     "two digit minimum",
			filename: "album-01.wav",
			prefix:   "",
			want:     "album",
			wantOK:   true,
		},
		{
			name:     "single digit should not match",
			filename: "album-1.wav",
			prefix:   "",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "underscore separator",
			filename: "my_album_123.mp3",
			prefix:   "",
			want:     "my_album",
			wantOK:   true,
		},
		{
			name:     "dot separator",
			filename: "album.002.wav",
			prefix:   "",
			want:     "album",
			wantOK:   true,
		},
		{
			name:     "space separator",
			filename: "album 003.wav",
			prefix:   "",
			want:     "album",
			wantOK:   true,
		},
		{
			name:     "mixed separators",
			filename: "my-album_mix.004.wav",
			prefix:   "",
			want:     "my-album_mix",
			wantOK:   true,
		},
		{
			name:     "no numeric suffix",
			filename: "album.wav",
			prefix:   "",
			want:     "",
			wantOK:   false,
		},
		{
			name:     "unicode filename",
			filename: "写真-001.jpg",
			prefix:   "",
			want:     "写真",
			wantOK:   true,
		},
		{
			name:     "spaces in filename",
			filename: "my photo album 99.jpg",
			prefix:   "",
			want:     "my photo album",
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

func TestParseGroup_PrefixBased(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		prefix   string
		want     string
		wantOK   bool
	}{
		{
			name:     "basic prefix form",
			filename: "kyoto-trip p04.jpg",
			prefix:   "p",
			want:     "kyoto-trip",
			wantOK:   true,
		},
		{
			name:     "prefix with space separator",
			filename: "album p99.wav",
			prefix:   "p",
			want:     "album",
			wantOK:   true,
		},
		{
			name:     "prefix with underscore",
			filename: "data_p123.csv",
			prefix:   "p",
			want:     "data",
			wantOK:   true,
		},
		{
			name:     "prefix with dash",
			filename: "files-p45.txt",
			prefix:   "p",
			want:     "files",
			wantOK:   true,
		},
		{
			name:     "prefix with dot",
			filename: "image.p88.jpg",
			prefix:   "p",
			want:     "image",
			wantOK:   true,
		},
		{
			name:     "multi-char prefix",
			filename: "report page12.pdf",
			prefix:   "page",
			want:     "report",
			wantOK:   true,
		},
		{
			name:     "prefix not at boundary - should fall back to numeric",
			filename: "albump99.wav",
			prefix:   "p",
			want:     "albump",
			wantOK:   true,
		},
		{
			name:     "prefix inside word should fall back to numeric",
			filename: "sp02.jpg",
			prefix:   "p",
			want:     "sp",
			wantOK:   true,
		},
		{
			name:     "prefix not found - falls back to numeric",
			filename: "album-001.wav",
			prefix:   "p",
			want:     "album",
			wantOK:   true,
		},
		{
			name:     "prefix form wins over numeric",
			filename: "test-p01-02.txt",
			prefix:   "p",
			want:     "test",
			wantOK:   true,
		},
		{
			name:     "single digit after prefix should not match",
			filename: "album p1.wav",
			prefix:   "p",
			want:     "",
			wantOK:   false,
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
