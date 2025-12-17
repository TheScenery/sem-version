package parser

import (
	"testing"
)

func TestParseCommit(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantType  CommitType
		wantScope string
		wantDesc  string
		wantBreak bool
	}{
		{
			name:      "simple feat",
			message:   "feat: add new feature",
			wantType:  TypeFeat,
			wantScope: "",
			wantDesc:  "add new feature",
			wantBreak: false,
		},
		{
			name:      "feat with scope",
			message:   "feat(api): add new endpoint",
			wantType:  TypeFeat,
			wantScope: "api",
			wantDesc:  "add new endpoint",
			wantBreak: false,
		},
		{
			name:      "fix",
			message:   "fix: resolve bug",
			wantType:  TypeFix,
			wantScope: "",
			wantDesc:  "resolve bug",
			wantBreak: false,
		},
		{
			name:      "breaking change with !",
			message:   "feat!: breaking api change",
			wantType:  TypeFeat,
			wantScope: "",
			wantDesc:  "breaking api change",
			wantBreak: true,
		},
		{
			name:      "breaking change with scope and !",
			message:   "feat(api)!: breaking api change",
			wantType:  TypeFeat,
			wantScope: "api",
			wantDesc:  "breaking api change",
			wantBreak: true,
		},
		{
			name:      "refactor",
			message:   "refactor: clean up code",
			wantType:  TypeRefactor,
			wantScope: "",
			wantDesc:  "clean up code",
			wantBreak: false,
		},
		{
			name:      "docs",
			message:   "docs: update readme",
			wantType:  TypeDocs,
			wantScope: "",
			wantDesc:  "update readme",
			wantBreak: false,
		},
		{
			name:      "unknown type",
			message:   "random: something",
			wantType:  TypeUnknown,
			wantScope: "",
			wantDesc:  "something",
			wantBreak: false,
		},
		{
			name:      "non-conventional commit",
			message:   "just a regular commit message",
			wantType:  TypeUnknown,
			wantScope: "",
			wantDesc:  "",
			wantBreak: false,
		},
		{
			name:      "breaking change in body",
			message:   "feat: add feature\n\nBREAKING CHANGE: this breaks the API",
			wantType:  TypeFeat,
			wantScope: "",
			wantDesc:  "add feature",
			wantBreak: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCommit(tt.message)
			if got.Type != tt.wantType {
				t.Errorf("ParseCommit().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Scope != tt.wantScope {
				t.Errorf("ParseCommit().Scope = %v, want %v", got.Scope, tt.wantScope)
			}
			if got.Description != tt.wantDesc {
				t.Errorf("ParseCommit().Description = %v, want %v", got.Description, tt.wantDesc)
			}
			if got.IsBreaking != tt.wantBreak {
				t.Errorf("ParseCommit().IsBreaking = %v, want %v", got.IsBreaking, tt.wantBreak)
			}
		})
	}
}

func TestIsBumpType(t *testing.T) {
	tests := []struct {
		commitType CommitType
		want       bool
	}{
		{TypeFeat, true},
		{TypeFix, true},
		{TypeRefactor, true},
		{TypePerf, true},
		{TypeDocs, false},
		{TypeStyle, false},
		{TypeTest, false},
		{TypeChore, false},
		{TypeBuild, false},
		{TypeCI, false},
		{TypeUnknown, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.commitType), func(t *testing.T) {
			parsed := ParsedCommit{Type: tt.commitType}
			if got := parsed.IsBumpType(); got != tt.want {
				t.Errorf("IsBumpType() = %v, want %v", got, tt.want)
			}
		})
	}
}
