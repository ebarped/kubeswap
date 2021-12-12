# kubeswap
Tool to manage multiple kubeconfig files and swap between clusters easily

## Commands
- setup: creates the file that holds the database. By default, is created at ~/.kube/ks.db 
    --db \<path\>: create the database file on different location
- list: lists all the kubeconfigs
- add \<name\> -f \<kubeconfig\>: adds a new kubeconfig identified as \<name\> from \<kubeconfig\> file to the db
- remove \<name\>: removes a kubeconfig from the db
- use \<name\>: modify current kubeconfig to the kubeconfig identified by name
- print \<name\>: print the kubeconfig identified by name
- status: checks if the clusters referenced by each kubeconfig are reachable
  - \<name>\: checks if the cluster of the \<name\> kubeconfig is reachable

## Tip
I like to create an alias to kubeswap:
`alias ks=kubeswap`

## Common flags
- --log-level: sets loglevel (info/debug)
- --db: specify location of the database file

## Technologies
- cobra: cli library
- fzf: fuzzy finding to enable more easisly
- badger: to store the kubeconfig files
- zerolog: structured (and fast) logger

## Differences with other projects
- [kubecm](https://github.com/sunny0826/kubecm): kubecm uses a single kubeconfig file. This projects uses a key-value DB to store multiple separated kubeconfigs.

## Usage
### Add
```
kubeswap add --name test --kubeconfig test-kubeconfig.yml
```

## TODO