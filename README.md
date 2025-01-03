# Blockchain Parser

## About

Implement Ethereum blockchain parser that will allow to query transactions for subscribed
addresses.

## Quickstart

0. Look at config.yaml, change it to your needs `rpc`, `start_block` and `blocks_interval`

1. Start parser by Docker

```
    make start
```

or, by local run

```
    make local-start
```

2. Stop parser

```
    make down
```

3. Start tests

```
    make test
```

## Usage

1. Get current block number

```
    curl http://localhost:8080/block
```

2. Get subscribe for transactions by address

```
    curl http://localhost:8080/subscribe?address=0xe93685f3bBA03016F02bD1828BaDD6195988D950
```

3. Get transactions by address

```
    curl http://localhost:8080/transactions?address=0xe93685f3bBA03016F02bD1828BaDD6195988D950
```

## Project Structure

**cmd** - contains entry point to start the parser

**config** - app's config

**internal** - internal app logic

* **api** - contains HTTP handlers
* **parser** - core blockchain parser logic, contains
* **repository** - repository for transactions, currently in-memory
* **ethclient** - client for Ethereum RPC
* **jsonrpc** - client for JSON-RPC