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

## To Do

- [x] issues
  - [x] labels
  - [x] milestones
  - [ ] users (mapping)
- [x] Pull Requests (placeholder issues)
