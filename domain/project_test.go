package domain

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestNewProjectOpsList(t *testing.T) {
	type args struct {
		sourceProjects []*github.Project
		targetProjects []*github.Project
	}
	tests := []struct {
		name string
		args args
		want ProjectOpsList
	}{
		{
			name: "empty",
			args: args{
				sourceProjects: []*github.Project{},
				targetProjects: []*github.Project{},
			},
			want: nil,
		},
		{
			name: "source <=> empty",
			args: args{
				sourceProjects: []*github.Project{
					&github.Project{
						Name: strRef("kanban"),
					},
				},
				targetProjects: []*github.Project{},
			},
			want: ProjectOpsList{
				&ProjectOp{
					Kind: OpCreate,
					Project: &github.Project{
						Name: strRef("kanban"),
					},
				},
			},
		},
		{
			name: "same name",
			args: args{
				sourceProjects: []*github.Project{
					&github.Project{
						Name: strRef("kanban"),
					},
				},
				targetProjects: []*github.Project{
					&github.Project{
						Name: strRef("kanban"),
					},
				},
			},
			want: ProjectOpsList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProjectOpsList(tt.args.sourceProjects, tt.args.targetProjects); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProjectOpsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
