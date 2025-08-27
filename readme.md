# TKView

TestKube Terminal User Interface (TUI).

Using [`bubbletea`](https://github.com/charmbracelet/bubbletea).

## Non-goals

Rather than try to define what TKView is, instead here is a list of all the things TKView should never be:
- A way to install, initialise, or update: Testkube or any related components. Instead, use `helm`, `kustomize`, `kubectl apply`,
  or any other standard Kubernetes installation tooling.
- A way to create, edit, or delete: `TestWorkflows`, `TestWorkflowTemplates`, `Webhooks`, `Triggers`, or any other Testkube related resources. 
  You should use Testkube Enterprise or the above standard Kubernetes tooling.
- A tool for migrating Testkube resources from older versions to newer ones.
  The `testkube` CLI tool can perform that action for you.

## Contributing

- See the TODO list of outstanding items below, pick one up and get to it!
- Abide by the lint rules!
  - Mostly, you should be fine if you are following the Go code review comments:
    - [Code Review Comments](https://go.dev/wiki/CodeReviewComments)
    - [Test Comments](https://go.dev/wiki/TestComments)
- Abide by the [Go project structure](https://go.dev/doc/modules/layout)!
  - Nothing in here should be exported, so all packages live in `internal`, or the `main` package.
- Tools should use the `go tool` directive where possible.
  - Some tools don't like this, so it can be easier to use a separate `go.mod` file in its own package.
    See [`internal/tools/golangci-lint`](internal/tools/golangci-lint/go.mod) as an example.

### Working with complex development commands

The [`plz`](plz) script can be used to document and run some of the more complex commands.
There is no need to add a case to it for every command, for example:
```shell
"build")
  go build ./...
```
would be unecessary, whereas:
```shell
"build")
  CGO_ENABLED=1 go build -ldflags '-s -w -linkmode external -extldflags "-fno-PIC -static"' -o foobar
```
would warrant inclusion in the script.

There is no need for `make` here at the moment.

## TODO

- [ ] API Authentication
  - UI and/or CLI flags
  - Initial view before starting the main application.
- [ ] Select Organisation and Environment
  - On Selected choose that Organisation and Environment combination.
- [ ] Show Agent statuses
  - No selection required, display as in web UI.
- [ ] Show Executions
  - How to show status?
  - Allow cancel.
- [ ] Start Execution
- [ ] Dive into granular Execution Step status

### Example UI

This is an indicator of where we want to end up:
```
╔Environments════════════════╗╔Agents══════════════════════════════════════════╗
║Testkube                    ║║Name              │Version      │Last Seen      ║
║│──Organisation A           ║║──────────────────│─────────────│───────────────║
║│  │──Environment 1         ║║foo-bar-baz       │2.1.168      │3 minutes ago  ║
║│  │──Environment 2         ║║tkagent-47        │1.56.734     │5 days ago     ║
║│──Organisation B           ║║                  │             │               ║
║   │──Environment 3         ║║                  │             │               ║
╚════════════════════════════╝╚════════════════════════════════════════════════╝
╔Executions════════════════════════════════════════════════════════════════════╗
║Name                             │Number│Status       │Started       │Duration║
║─────────────────────────────────│──────│─────────────│──────────────│────────║
║postman-workflow-smoke-cron-test │   293│Successful   │4 minutes ago │  23.41s║
║postman-workflow-smoke-cron-test │   292│Successful   │14 minutes ago│  22.36s║
║postman-workflow-smoke-cron-test │   291│Aborted      │24 minutes ago│   2.11s║
║postman-workflow-smoke-cron-test │   290│Successful   │34 minutes ago│  22.54s║
║postman-workflow-smoke-cron-test │   289│Failed       │44 minutes ago│  57.36s║
║postman-workflow-smoke-cron-test │   288│Successful   │54 minutes ago│  25.12s║
║postman-workflow-smoke-cron-test │   287│Successful   │1 hour ago    │  24.32s║
║postman-workflow-smoke-cron-test │   286│Successful   │1 hour ago    │  21.43s║
║                                 │      │             │              │        ║
║                                 │      │             │              │        ║
║                                 │      │             │              │        ║
║                                 │      │             │              │        ║
╚══════════════════════════════════════════════════════════════════════════════╝
```