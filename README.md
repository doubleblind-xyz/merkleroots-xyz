## About merkleroots.xyz

merkleroots.xyz is a centralized repository of [Merkle Trees](https://en.wikipedia.org/wiki/Merkle_tree) that supports free API read and write access. If you are a developer, please also check out our client libaries (todo js, python).

## API

Merkle trees are identified by a root value (`root: string`). Each root is associated with an ordered list of leaf values (`leaves: string[]`).

|Method|Endpoint|Description|Request Body|Response Body|
|--|--|--|--|--|
|GET|`api.merkleroots.xyz/tree/:root`|Returns the tree rooted at `root` in array form (see Tree Layout)||`{ nodes: string[] }`|
|POST|`api.merkleroots.xyz/tree`|Creates a tree whose leaves are `leaves` and returns its `root`. If a tree with `root` exists already, the tree with more nodes is kept.|`{ leaves: string[] }`|`string`|
|GET|`api.merkleroots.xyz/alias/:name`|Returns the value that `name` points to||`string`|
|POST|`api.merkleroots.xyz/alias`|Point `name` to the value `root`. Returns `Ok` or `Forbidden`|`{ name: string; root: string }`|`string`|

## Hash Function
All merkle trees created by merkleroots.xyz use the [`Poseidon`](https://www.poseidon-hash.info/) hash function, using [`circomlibjs`](https://github.com/iden3/circomlibjs/blob/5164544558570f934d72d40c70779fc745350a0e/src/poseidon_reference.js)'s implementation. No other hash functions are supported at this time.

## Tree Layout

We define a merkle tree over an ordered list of leaf nodes to be a [complete](https://en.wikipedia.org/wiki/Binary_tree#Types_of_binary_trees) binary tree, where the value of a parent node is the hash of its two child nodes. A merkle tree with n leaves will be represented by an array of size 2n, whose root is at index 1 and whose leaves are at indices stored n...2n-1.

See [Implicit data strucutre](https://en.wikipedia.org/wiki/Binary_tree#Arrays) and the merkle.xyz source code for more details (coming soon).

## Aliases

The aliasing feature allows you to shorten a root value into a human-readable string.
Namespaced aliases, which follow the format `@namespace/name`, are writable only by authorized users (TODO).
You can use namespaced aliases to simultaneously name and sign a merkle root.

All aliases are publicly readable. Only namespaced aliases are mutable.

## Local Development

run postgres server:
```
brew install postgres
brew services start postgres
```

set up postgres:
```
brew info postgres
psql; create user mroots; create database mroots; # linux: sudo su postgres before psql
```

run server:

```
go run src/server.go
# or with gin
cd src; gin run
```

test server:
```
curl -X POST 127.0.0.1:3000/tree -d @request.json
curl 127.0.0.1:3000/tree/18413384112310044160158834857773087012258255807213533751159036388968578675718
```