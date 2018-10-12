# General Sharding Specification

## Layout of 'Shards'

```MARKDOWN
                 Network
                    +
                    |
                    |
                    |
         +---------------------+
         |          |          |
         |          |          |
         |          |          |
         |          |          |
         +          +          +
      Shard 0    Shard 1    Shard 2
         +
  +------------+
  |      |     |
  +      +     +
Node 0 Node 1 Node 2
```

Each 'shard' contains x amount of nodes, each processing all of the transactions addressed to the shard. The shard containing such nodes is a child of a 'Parent Shard.'

## Shard-level Consensus

On the issuing of a transaction by a node, the least 'busy' non-parent shard will be selected for verification. The transaction is sent via TCP to the Shard, and is then propagated throughout the shard. Each node in the destination shard then verifies the transaction, and stores it in persistent memory. Should the transaction reference another transaction not contained in the destination shard, a request will be made to the shard holding such a transaction. For further explanation on the layout of transactions, and how shards are recorded in transactions, see [transaction specification](TransactionSpecification.md).