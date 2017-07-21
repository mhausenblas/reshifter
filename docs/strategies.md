# Strategies

ReShifter supports different backup strategies, including:

- _raw_ dump: storing the value of every key under the well-known top-level keys such as `/registry` or `/kubernetes.io`.
- _render_ dump: writing the values to `stdout`.
- _filter_: storing the value of selected (white-listed) keys under the well-known top-level, for example only `deployment` or `service`.

You define the backup strategy using the `RS_BACKUP_STRATEGY` environment variable, with a default value of `raw`.

For example, using the CLI tool `rcli`, here's how to use the `filter` strategy:

```
# only back up objects below '/namespaces/mycoolproject'
$ RS_BACKUP_STRATEGY=filter:/namespaces/mycoolproject rcli backup create

# only back up objects which path contains 'deployment' or 'service'
$ RS_BACKUP_STRATEGY=filter:deployment,service rcli backup create
```

## Implementation

If you plan to add a new backup strategy, the steps are:

- define new reap function type in `types.go`
- add reap function implementation to `strategy.go`
- add case to `visitors.go`
- add case to `backup.go`
