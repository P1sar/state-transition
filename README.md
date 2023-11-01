# Flashbots worktest

&nbsp;

Tasks consist of two parts:

1. Transition with Geth State Snapshot & Mempool-Dumpster. Located in transtition module [Uncompleted]
2. Block tracer [Completed]. Uses `trace_block` RPC method to get all transactions from the block and then writes all touched addresses to ./seen.txt if --snapshot provided checks touched addresses against snapshot and writes matched addresses to ./seen.tx and not matched to ./unseen.tx

### Tracer usage:
build
```
go build -o trace
```

to get cli help
```
./trace --help
```
```
Flags:
      --block string      block number
  -h, --help              help for this command
      --rpc string        URL of rpc to fetch the traces
      --snapshot string   Path to the snapshot
```
If --snapshot flag provided script will check snapshot block with provided block will return error if they missmatch.

