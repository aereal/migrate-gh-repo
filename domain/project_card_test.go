package domain

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestNewProjectCardOpsList(t *testing.T) {
	type args struct {
		sourceCards  []*github.ProjectCard
		targetCards  []*github.ProjectCard
		sourceColumn *github.ProjectColumn
		targetColumn *github.ProjectColumn
	}
	tests := []struct {
		name string
		args args
		want ProjectCardOpsList
	}{
		{
			name: "empty",
			args: args{
				sourceCards: []*github.ProjectCard{},
				targetCards: []*github.ProjectCard{},
			},
			want: nil,
		},
		{
			name: "source -> empty",
			args: args{
				sourceCards: []*github.ProjectCard{
					&github.ProjectCard{
						Note: strRef("poppoe"),
					},
				},
				targetCards: []*github.ProjectCard{},
				targetColumn: &github.ProjectColumn{
					Name: strRef("To Do"),
				},
			},
			want: ProjectCardOpsList{
				&ProjectCardOp{
					Kind: OpCreate,
					ProjectCard: &github.ProjectCard{
						Note: strRef("poppoe"),
					},
					ProjectColumn: &github.ProjectColumn{
						Name: strRef("To Do"),
					},
				},
			},
		},
		{
			name: "same name",
			args: args{
				sourceCards: []*github.ProjectCard{
					&github.ProjectCard{
						Note: strRef("poppoe"),
					},
				},
				targetCards: []*github.ProjectCard{
					&github.ProjectCard{
						Note: strRef("poppoe"),
					},
				},
				targetColumn: &github.ProjectColumn{
					Name: strRef("To Do"),
				},
			},
			want: ProjectCardOpsList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProjectCardOpsList(tt.args.sourceCards, tt.args.targetCards, tt.args.sourceColumn, tt.args.targetColumn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProjectCardOpsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
