# General Transaction Structure Specification

## General Layout

```MARKDOWN
    Transaction: {
        TransactionData: {
            Root [Transaction]: Parent transaction (synonymous to previousBlock in a blockchain)
            Nonce [uint64]: Incrementor used to calculate hashes
            Recipient [address]: Destination for transaction
            Amount [float64]: Amount to transfer
            UnspentReward [uint64]: Reward not yet spent
            Payload [[]byte]: Smart-contract data associated with transaction (e.g. method call)
            Time [time.Time]: Time transaction issued (UTC)
            Extra [[]byte]: More miscellaneous data associated with transaction
            InitialHash [[]byte]: Transaction hash
            ParentHash [[]byte]: Hash of root transaction
            StorageShardID [[]byte]: Hash of verifying shard
        },
        Contract (optional) [pointer to Contract]: Destination contract
        Verifications [uint64]: Amount of verifications
        Weight [float64]: Proportional weight according to individual weight of verifying nodes
        InitialWitness [pointer to witness]: First node to witness transaction
        SendingAccount [Account]: Account issuing transaction
        Reward [uint64]: Amount of currency to reward verifying shard nodes with
        ChainVersion [uint64]: Index in chain
    }
```

As proposed in the above diagram, each transaction contains the fields TransactionData (transaction metadata), Contract (optional destination contract), Verifications (amount of verifications), Weight (transaction weight), InitialWitness (first node witnessing transaction), SendingAccount (account issuing transaction), Reward (amount of currency distributed throughout verifying shard), and ChainVersion (an incrementing index of a transaction in the chain).