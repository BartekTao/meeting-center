# nycu-meeting-room
meeting room project

# Setup
## set up environment
1. install go with 1.22.1
2. install makefile
3. go install github.com/go-delve/delve/cmd/dlv@latest (for debug mode in docker, but not yet ready)

## makefile cli
- make: run main api server
- make gqlgen: regenerate graphQL API
*more detail: see makefile*

## branch policy
- Do not merge into main without PR and Reviewing
- When develop any feature, create a new branch name feat/$(your_branch_name)
