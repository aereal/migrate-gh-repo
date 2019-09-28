package domain

import "testing"

func TestUserAliasResolver_AssumeResolved(t *testing.T) {
	type args struct {
		fromUserName string
	}
	tests := []struct {
		name          string
		resolver      *UserAliasResolver
		args          args
		wantAliasName string
		wantAliased   bool
	}{
		{
			name:          "no alias",
			resolver:      NewUserAliasResolver(nil),
			args:          args{fromUserName: "aereal"},
			wantAliasName: "aereal",
			wantAliased:   false,
		},
		{
			name:          "alias",
			resolver:      NewUserAliasResolver(map[string]string{"aereal": "noreal"}),
			args:          args{fromUserName: "aereal"},
			wantAliasName: "noreal",
			wantAliased:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAliasName, gotAliased := tt.resolver.AssumeResolved(tt.args.fromUserName)
			if gotAliasName != tt.wantAliasName {
				t.Errorf("UserAliasResolver.AssumeResolved() gotAliasName = %v, want %v", gotAliasName, tt.wantAliasName)
			}
			if gotAliased != tt.wantAliased {
				t.Errorf("UserAliasResolver.AssumeResolved() gotAliased = %v, want %v", gotAliased, tt.wantAliased)
			}
		})
	}
}
