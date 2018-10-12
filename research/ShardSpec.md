# General Sharding Specification

## Shard Definition

In the proposed network, shards are defined as groups of transaction-verifying nodes. The size of a shard is dynamic, and calculated at runtime. In said shards, each sibling node processes the exact same amount of information, and stores the same portion of the state. Should any node in said shard disobey any rule set forth by the protocol, that node's stake will be burned and will be banned from the shard. In reference to aforementioned stake, one would find [Proof-of-stake Outline](ProofOfStake.md) quite helpful.

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