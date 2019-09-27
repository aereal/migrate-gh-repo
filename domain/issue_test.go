package domain

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestNewIssueOpsList(t *testing.T) {
	type args struct {
		sourceIssues []*github.Issue
		targetIssues []*github.Issue
	}
	tests := []struct {
		name string
		args args
		want IssueOpsList
	}{
		{
			name: "source=empty target=empty",
			args: args{
				sourceIssues: []*github.Issue{},
				targetIssues: []*github.Issue{},
			},
			want: nil,
		},
		{
			name: "source=[A] target=empty",
			args: args{
				sourceIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe"),
					},
				},
				targetIssues: []*github.Issue{},
			},
			want: IssueOpsList([]*IssueOp{
				&IssueOp{
					Kind: OpCreate,
					Issue: &github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe"),
					},
				},
			}),
		},
		{
			name: "source=[A,B] target=[A]",
			args: args{
				sourceIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe1"),
					},
					&github.Issue{
						Number: intRef(2),
						Title:  strRef("poppoe2"),
					},
				},
				targetIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe1"),
					},
				},
			},
			want: IssueOpsList([]*IssueOp{
				&IssueOp{
					Kind: OpCreate,
					Issue: &github.Issue{
						Number: intRef(2),
						Title:  strRef("poppoe2"),
					},
				},
			}),
		},
		{
			name: "source=[A] target=[A]",
			args: args{
				sourceIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe"),
					},
				},
				targetIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe"),
					},
				},
			},
			want: IssueOpsList([]*IssueOp{}),
		},
		{
			name: "source=[A] target=[A']",
			args: args{
				sourceIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe"),
						Body:   strRef("master"),
					},
				},
				targetIssues: []*github.Issue{
					&github.Issue{
						Number: intRef(1),
						Title:  strRef("poppoe"),
						Body:   strRef("old"),
					},
				},
			},
			want: IssueOpsList([]*IssueOp{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIssueOpsList(tt.args.sourceIssues, tt.args.targetIssues); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMilestoneOpsList() = %s, want %s", got, tt.want)
			}
		})
	}
}
