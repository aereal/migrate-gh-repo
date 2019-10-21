# migrate-gh-repo

The script to migrate repositories with above constraints:

- labels and milestones are fully migrated to target repository
- issues that refers source created on target repository

## Run

```
# create config/default.cue; see Configuration section
go run ./
```

## Configuration

- Write your configuration to `config/default.cue`
- The spec is `config/spec.cue`
- refs. https://cuelang.org/

## Caveats

- all of assignees on source repository must have permission to triage issues on target repository
  - migrate-gh-repo currently does not support migration of collaborators and has no intention to implement that
    - Because management of collaborators and teams requires more strong and maybe dangerous permission but it is risky for us
    - You can use [Terraform][terraform] and [GitHub provider][terraform-github-provider]
  - refs. [Repository permission levels for an organization - GitHub Help][github-repository-permission]
- migration of ton of issues, labels, or milestones may cause excess of API rate limit
  - Currently only way to avoid it is update sleep duration by you
  - We have intention to resolve that issue on smart way but have no good idea; **patches/suggestions are welcome**

[github-repository-permission]: https://help.github.com/en/github/setting-up-and-managing-organizations-and-teams/repository-permission-levels-for-an-organization
[terraform]: https://www.terraform.io
[terraform-github-provider]: https://www.terraform.io/docs/providers/github/
