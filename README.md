# kubeswap
Tool to manage multiple kubeconfig files and swap between clusters easily

# commands
- kubeswap: displays the list of kubeconfigs in the db, same as subcommand use but without specifying a name

## Subcommands
- list: lists all the kubeconfigs
- add -n \<name\> -f \<kubeconfig\>: adds a new kubeconfig identified as \<name\> from \<kubeconfig\> file to the db
- delete -n \<name\>: removes a kubeconfig from the db
- use -n \<name\>: modify current kubeconfig to the kubeconfig identified by name
- print -n \<name\>: print the kubeconfig identified by name
- printall: prints all the kubeconfig from the db
- status: checks if the clusters referenced by each kubeconfig are reachable
  - \<name>\: checks if the cluster of the \<name\> kubeconfig is reachable

## Tip
I like to create an alias to kubeswap:
`alias ks=kubeswap`

## Common flags
- --loglevel: sets loglevel (info/debug)
- --kubeconfig: specify location of kubeconfig file
- --db: specify location of the database file

## Technologies
- cobra: cli library
- fzf: fuzzy finding to enable more easisly
- pogreb: key-value database to store the files
- zerolog: structured (and fast) logger

## Differences with other projects
- [kubecm](https://github.com/sunny0826/kubecm): kubecm uses a single kubeconfig file. This projects uses a key-value DB to store multiple separated kubeconfigs.

## Usage
### Add
```
kubeswap add --name test --kubeconfig test-kubeconfig.yml
```

## Test
- Init:
```bash
make clean build test
```
- add:
```bash
./dist/kubeswap_linux_amd64/kubeswap add --name test --kubeconfig test/kubeconfig.yml --db /tmp/test
```
- list:
```bash
./dist/kubeswap_linux_amd64/kubeswap list --db /tmp/test
```
- print:
```bash
./dist/kubeswap_linux_amd64/kubeswap print -n test-0 --db /tmp/test
```
- printall:
```bash
./dist/kubeswap_linux_amd64/kubeswap printall --db /tmp/test
```
- delete:
```bash
./dist/kubeswap_linux_amd64/kubeswap delete -n test-0 --db /tmp/test
```

## TODO
- Get homeDir in windows, linux & mac