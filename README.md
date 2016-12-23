## tidb-tools

tidb-tools are some useful tool collections for [TiDB](https://github.com/pingcap/tidb).


## How to build

```
make build # build all tools

make importer # build importer

make syncer # build syncer

make checker # build checker

make loader # build loader
```

When build successfully, you can find the binary in bin directory.

## Tool list

[importer](./importer)

A tool for generating and inserting datas to database which is compatible with MySQL protocol, like MySQL, TiDB.

[syncer](./syncer)

A tool for syncing source database data to target database which is compatible with MySQL protocol, like MySQL, TiDB.

[checker](./checker)

A tool for checking the compatibility of an existed MySQL database with TiDB.

[loader](./loader)

A tool for multi-threaded restoration of mydumper.

[scripts](./scripts)

Collection of auxiliary scripts.

## License
Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.
