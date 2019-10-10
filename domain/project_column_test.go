package domain

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestNewProjectColumnOpsList(t *testing.T) {
	type args struct {
		sourceColumns []*github.ProjectColumn
		targetColumns []*github.ProjectColumn
		sourceProject *github.Project
		targetProject *github.Project
	}
	tests := []struct {
		name string
		args args
		want ProjectColumnOpsList
	}{
		{
			name: "empty",
			args: args{
				sourceColumns: []*github.ProjectColumn{},
				targetColumns: []*github.ProjectColumn{},
			},
			want: nil,
		},
		{
			name: "source -> empty",
			args: args{
				sourceColumns: []*github.ProjectColumn{
					&github.ProjectColumn{
						Name: strRef("To Do"),
					},
				},
				targetColumns: []*github.ProjectColumn{},
			},
			want: ProjectColumnOpsList{
				&ProjectColumnOp{
					Kind: OpCreate,
					ProjectColumn: &github.ProjectColumn{
						Name: strRef("To Do"),
					},
				},
			},
		},
		{
			name: "same name",
			args: args{
				sourceColumns: []*github.ProjectColumn{
					&github.ProjectColumn{
						Name: strRef("To Do"),
					},
				},
				targetColumns: []*github.ProjectColumn{
					&github.ProjectColumn{
						Name: strRef("To Do"),
					},
				},
			},
			want: ProjectColumnOpsList{
				&ProjectColumnOp{
					Kind: OpUpdate,
					ProjectColumn: &github.ProjectColumn{
						Name: strRef("To Do"),
					},
					TargetProjectColumn: &github.ProjectColumn{
						Name: strRef("To Do"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProjectColumnOpsList(tt.args.sourceColumns, tt.args.targetColumns, tt.args.sourceProject, tt.args.targetProject); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProjectColumnOpsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
