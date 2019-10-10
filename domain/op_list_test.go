package domain

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func strRef(s string) *string { return &s }

func intRef(i int) *int { return &i }

func Test_milestoneEq(t *testing.T) {
	type args struct {
		l *milestone
		r *milestone
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil vs nil",
			args: args{
				l: nil,
				r: nil,
			},
			want: false,
		},
		{
			name: "A vs nil",
			args: args{
				l: &milestone{},
				r: nil,
			},
			want: false,
		},
		{
			name: "A vs A",
			args: args{
				l: &milestone{
					&github.Milestone{
						Title: strRef("poppoe"),
					},
				},
				r: &milestone{
					&github.Milestone{
						Title: strRef("poppoe"),
					},
				},
			},
			want: true,
		},
		{
			name: "A vs B",
			args: args{
				l: &milestone{&github.Milestone{
					Title: strRef("poppoe"),
				},
				},
				r: &milestone{&github.Milestone{
					Title: strRef("ubobobo"),
				},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.l.eq(tt.args.r); got != tt.want {
				t.Errorf("eqMilestone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMilestoneOpsList(t *testing.T) {
	type args struct {
		sourceMilestones []*github.Milestone
		targetMilestones []*github.Milestone
	}
	tests := []struct {
		name string
		args args
		want MilestoneOpsList
	}{
		{
			name: "source=empty target=empty",
			args: args{
				sourceMilestones: []*github.Milestone{},
				targetMilestones: []*github.Milestone{},
			},
			want: nil,
		},
		{
			name: "source=[A] target=empty",
			args: args{
				sourceMilestones: []*github.Milestone{
					&github.Milestone{
						Title: strRef("poppoe"),
					},
				},
				targetMilestones: []*github.Milestone{},
			},
			want: MilestoneOpsList([]*MilestoneOp{
				&MilestoneOp{
					Kind: OpCreate,
					Milestone: &github.Milestone{
						Title: strRef("poppoe"),
					},
				},
			}),
		},
		{
			name: "source=[A,B] target=[A]",
			args: args{
				sourceMilestones: []*github.Milestone{
					&github.Milestone{
						Title: strRef("poppoe1"),
					},
					&github.Milestone{
						Title: strRef("poppoe2"),
					},
				},
				targetMilestones: []*github.Milestone{
					&github.Milestone{
						Title: strRef("poppoe1"),
					},
				},
			},
			want: MilestoneOpsList([]*MilestoneOp{
				&MilestoneOp{
					Kind: OpCreate,
					Milestone: &github.Milestone{
						Title: strRef("poppoe2"),
					},
				},
			}),
		},
		{
			name: "source=[A] target=[A]",
			args: args{
				sourceMilestones: []*github.Milestone{
					&github.Milestone{
						Title: strRef("poppoe"),
					},
				},
				targetMilestones: []*github.Milestone{
					&github.Milestone{
						Title: strRef("poppoe"),
					},
				},
			},
			want: MilestoneOpsList([]*MilestoneOp{}),
		},
		{
			name: "source=[A] target=[A']",
			args: args{
				sourceMilestones: []*github.Milestone{
					&github.Milestone{
						Title:       strRef("poppoe"),
						Description: strRef("master"),
					},
				},
				targetMilestones: []*github.Milestone{
					&github.Milestone{
						Title:       strRef("poppoe"),
						Description: strRef("old"),
					},
				},
			},
			want: MilestoneOpsList([]*MilestoneOp{
				&MilestoneOp{
					Kind: OpUpdate,
					Milestone: &github.Milestone{
						Title:       strRef("poppoe"),
						Description: strRef("master"),
					},
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMilestoneOpsList(tt.args.sourceMilestones, tt.args.targetMilestones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMilestoneOpsList() = %s, want %s", got, tt.want)
			}
		})
	}
}
