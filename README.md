# Monorepo Test Spike

## Wanted Features

- [x] Testing should run based on where files are edited
- [x] Pull request reviews assigned appropriately
- [ ] Heroku builds off workspace
    - No official support yet, may need to create a custom Heroku build
    - https://github.com/heroku/heroku-buildpack-go/pull/484

## Detailed use case

- [x] Multiple server commands
- [x] Re-usable GO modules (PB, ACL, etc.)
- [ ] Heroku builds

## GO Workspaces

Test spike is using the official Go Workspace tutorial along with a lot of stub code.

https://go.dev/doc/tutorial/workspaces

### Testing

https://docs.github.com/en/actions/using-workflows/triggering-a-workflow#using-filters-to-target-specific-paths-for-pull-request-or-push-events

It is possible to filter by path, this should allow Assessment/Identity code changes to limit the number of tests they need to run.

Top level files may trigger all tests, depends on how well we filter files in the workflows

### Assigned Reviewer

Using Github actions filtered by folder and using (go-github)[https://github.com/google/go-github], a randomly assigned user from the correct team can be assigned to a review.  Users are grouped by currently open pull requests with random selection from the lowest assigned group.  This prevents one user from being assigned all reviews.

The example I have created still needs some work but there are no more known unknowns.

### Pros

1. Single repo will greatly reduce code duplication
2. Developer setup simplifies
3. Developer has easy visibility to all code
4. Standardization improves (code conventions, etc)
5. Developer collaboration across teams becomes simple

### Cons

1. Requires backward compatible Go module changes

#### Local GO module dependency Versioning
Versioning of included modules will need to be upgraded together
    - Backwards compatibility will be needed in all module library changes

Protobuf + ACL Examples

PB GO code updated will be consumable by both Identity + Assessment projects
    - Changes will need to be backwards compatible
    - This makes sense as the PB define the API's

ACL code is another example
    - Access to different data restricted by ACL
    - Changes will need to be backwards compatible
    - In some cases, 2 different versions of the same module would be easier to work with
