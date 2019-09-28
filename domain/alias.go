package domain

type UserAliasResolver struct {
	mapping map[string]string
}

func NewUserAliasResolver(mapping map[string]string) *UserAliasResolver {
	return &UserAliasResolver{mapping: mapping}
}

func (r *UserAliasResolver) AssumeResolved(fromUserName string) (aliasName string, aliased bool) {
	aliasName, aliased = r.mapping[fromUserName]
	if !aliased {
		aliasName = fromUserName
	}
	return
}
