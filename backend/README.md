# nycu-meeting-room
meeting room project

# Setup
## set up environment
1. install go with 1.22.1
2. install makefile
3. go install github.com/go-delve/delve/cmd/dlv@latest (for debug mode in docker, but not yet ready)

## makefile cli
> see makefile to get more detail
- make: run main api server
- make gqlgen: regenerate graphQL API  


# branch policy
- Do not merge into `main` without a PR and a review.
- When developing any feature, create a new branch named `feat/$(your_branch_name)`.

# coding convention
## golang
- file name: hello_world.go
- package name: short and lower case
- function name: upper case camel in public, lower case camel in private
- put your test in the same folder
## mongoDb
- lower case camel

## graphQL
- https://medium.com/@andriiandriiets/graphql-standards-and-practices-da3246dfb619

- 