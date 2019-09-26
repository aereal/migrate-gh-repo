package domain

import (
	"testing"

	"github.com/google/go-github/github"
)

func strRef(s string) *string { return &s }

func Test_eqMilestone(t *testing.T) {
	type args struct {
		l *github.Milestone
		r *github.Milestone
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
				l: &github.Milestone{},
				r: nil,
			},
			want: false,
		},
		{
			name: "A vs A",
			args: args{
				l: &github.Milestone{
					Title: strRef("poppoe"),
				},
				r: &github.Milestone{
					Title: strRef("poppoe"),
				},
			},
			want: true,
		},
		{
			name: "A vs B",
			args: args{
				l: &github.Milestone{
					Title: strRef("poppoe"),
				},
				r: &github.Milestone{
					Title: strRef("ubobobo"),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eqMilestone(tt.args.l, tt.args.r); got != tt.want {
				t.Errorf("eqMilestone() = %v, want %v", got, tt.want)
			}
		})
	}
}

