package version

import (
	"testing"

	"sem-version/internal/parser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Version
		wantErr bool
	}{
		{
			name:    "simple version",
			input:   "1.2.3",
			want:    Version{Major: 1, Minor: 2, Patch: 3},
			wantErr: false,
		},
		{
			name:    "version with v prefix",
			input:   "v1.2.3",
			want:    Version{Major: 1, Minor: 2, Patch: 3},
			wantErr: false,
		},
		{
			name:    "version with prerelease",
			input:   "v1.2.3-alpha.1",
			want:    Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1"},
			wantErr: false,
		},
		{
			name:    "version with metadata",
			input:   "v1.2.3+build.123",
			want:    Version{Major: 1, Minor: 2, Patch: 3, Metadata: "build.123"},
			wantErr: false,
		},
		{
			name:    "version with prerelease and metadata",
			input:   "v1.2.3-beta.2+build.456",
			want:    Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta.2", Metadata: "build.456"},
			wantErr: false,
		},
		{
			name:    "invalid version",
			input:   "not-a-version",
			want:    Version{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		want    string
	}{
		{
			name:    "simple version",
			version: Version{Major: 1, Minor: 2, Patch: 3},
			want:    "v1.2.3",
		},
		{
			name:    "version with prerelease",
			version: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha"},
			want:    "v1.2.3-alpha",
		},
		{
			name:    "version with metadata",
			version: Version{Major: 1, Minor: 2, Patch: 3, Metadata: "build.1"},
			want:    "v1.2.3+build.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Bump(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3}

	t.Run("BumpMajor", func(t *testing.T) {
		got := v.BumpMajor()
		want := Version{Major: 2, Minor: 0, Patch: 0}
		if got != want {
			t.Errorf("BumpMajor() = %v, want %v", got, want)
		}
	})

	t.Run("BumpMinor", func(t *testing.T) {
		got := v.BumpMinor()
		want := Version{Major: 1, Minor: 3, Patch: 0}
		if got != want {
			t.Errorf("BumpMinor() = %v, want %v", got, want)
		}
	})

	t.Run("BumpPatch", func(t *testing.T) {
		got := v.BumpPatch()
		want := Version{Major: 1, Minor: 2, Patch: 4}
		if got != want {
			t.Errorf("BumpPatch() = %v, want %v", got, want)
		}
	})
}

func TestCalculateNextVersion(t *testing.T) {
	current := Version{Major: 1, Minor: 2, Patch: 3}

	tests := []struct {
		name    string
		commits []parser.ParsedCommit
		want    Version
	}{
		{
			name:    "no commits",
			commits: []parser.ParsedCommit{},
			want:    Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name: "patch bump - fix",
			commits: []parser.ParsedCommit{
				{Type: parser.TypeFix},
			},
			want: Version{Major: 1, Minor: 2, Patch: 4},
		},
		{
			name: "minor bump - feat",
			commits: []parser.ParsedCommit{
				{Type: parser.TypeFeat},
			},
			want: Version{Major: 1, Minor: 3, Patch: 0},
		},
		{
			name: "major bump - breaking",
			commits: []parser.ParsedCommit{
				{Type: parser.TypeFeat, IsBreaking: true},
			},
			want: Version{Major: 2, Minor: 0, Patch: 0},
		},
		{
			name: "mixed commits - highest wins",
			commits: []parser.ParsedCommit{
				{Type: parser.TypeFix},
				{Type: parser.TypeFeat},
				{Type: parser.TypeDocs},
			},
			want: Version{Major: 1, Minor: 3, Patch: 0},
		},
		{
			name: "breaking change wins over all",
			commits: []parser.ParsedCommit{
				{Type: parser.TypeFix},
				{Type: parser.TypeFeat},
				{Type: parser.TypeFeat, IsBreaking: true},
			},
			want: Version{Major: 2, Minor: 0, Patch: 0},
		},
		{
			name: "docs only - no bump",
			commits: []parser.ParsedCommit{
				{Type: parser.TypeDocs},
				{Type: parser.TypeChore},
			},
			want: Version{Major: 1, Minor: 2, Patch: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateNextVersion(current, tt.commits)
			if got != tt.want {
				t.Errorf("CalculateNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
