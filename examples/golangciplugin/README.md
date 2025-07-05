# Example Of Golangci Plugin

Here you can find an example on how to use `structfieldinitorder` as a [golangci-lint plugin][plugin].

## How to run it

Follow the steps described in the [module-plugins][plugin]:

+ Run the command `golangci-lint custom -v`.
+ Run the resulting custom binary of golangci-lint (`./custom-gcl` by default).

If everything went good, you should see:

```bash
main.go:15:7: fields for struct "Person" are not instantiated in order (structfieldinitorder)
```

[plugin]: https://golangci-lint.run/plugins/module-plugins
