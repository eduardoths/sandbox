# Consensus-transactions

There are two transactions this project could have

1. Client-side transactions when a client needs to perform multiple tasks and may commit or rollback
2. Consensus transactions to ensure that a data is either created/updated or not across all nodes

For now we will ignore client-side transactions to make the project simpler and let the client handle it's own transactions with update/deletes

## Consensus transactions
## History of changes

To ensure that it is viable to rollback a update, we need a history of changes made
In this way we could have a storage that stores current commited data
and a storage that saves all the history according to a transaction id

Let's understand the create/update/delete commit and rollback cycles in the next section

## Setting data transactions
First the the leader send it's followers a request to set the data
accompanied by a numerical sequential id.
Each command is attomic and can be identified by it's id to rollback or commit a operation.

The followers will store that id and the message the leader sent then. So we will have the following data on history storag:
TRANSACTION_HISTORY-STORAGE
0: SET K0 V1

At the same time the followers will also save the value on the transaction storage:

CURRENT_TRANSACTION-STORAGE
0: AWAITING
This way it's possible to know that the transaction-0 is still awaiting a commit

And the commited storage will be empty:
COMMITED-STORAGE

To ensure an easy understanding of a key history of changes, there should be a fourth storage containing transactions commited by key

COMMITS_HISTORY-STORAGE

### Rollback
When the leader tells it's followers to rollback a transaction they will only delete the current_transaction
this way we could have have something like:


TRANSACTION_HISTORY-STORAGE
0: SET K0 V1
1: ROLLBACK 0

CURRENT_TRANSACTION-STORAGE

COMMITED-STORAGE

COMMITS_HISTORY-STORAGE

### Commit
If for instance a commit is approved first the transaction should commit the value changed, then update the commits history.

if it fails to commit, due to an error such as lost internet connection, the follower should be out of sync and unavailable until it can sync with the leader

After saving the commited value with a metadata containing all 

TRANSACTION_HISTORY-STORAGE
0: SET K0 V1
1: COMMIT 0

CURRENT_TRANSACTION-STORAGE

COMMITED-STORAGE
K0: V1

COMMITS_HISTORY-STORAGE
K0: 0,x,y,z... // all the commited transactions that went through it
