# Kubeswap
Tool to manage multiple kubeconfig files and swap between clusters easily

```yaml

  ██   ██ ██    ██ ██████  ███████ ███████ ██     ██  █████  ██████  
  ██  ██  ██    ██ ██   ██ ██      ██      ██     ██ ██   ██ ██   ██ 
  █████   ██    ██ ██████  █████   ███████ ██  █  ██ ███████ ██████  
  ██  ██  ██    ██ ██   ██ ██           ██ ██ ███ ██ ██   ██ ██      
  ██   ██  ██████  ██████  ███████ ███████  ███ ███  ██   ██ ██

Manage your kubeconfig files easily

Usage:
  kubeswap [command]

Available Commands:
  add         Adds a new kubeconfig to the database
  completion  Generate completion script
  delete      Deletes a kubeconfig from the database
  help        Help about any command
  list        Lists all the kubeconfigs in the db
  print       Prints the content of the kubeconfig referenced by <name>
  printall    Prints the content of all the kubeconfigs from the db
  version     Print the version number

Flags:
      --db string         db file path (default "$HOME/.kube/kubeswap.db")
  -h, --help              help for kubeswap
      --loglevel string   loglevel (info/debug) (default "info")

Use "kubeswap [command] --help" for more information about a command.
```

# Commands
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

## Technologies
- cobra: cli library
- fzf: fuzzy finding to enable more easisly
- pogreb: key-value database to store the files
- zerolog: structured (and fast) logger


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
- Add shell completion
- Test in windows & linux
- use command:
  - use: prints a list of the keys and allows the user to select the desired kubeconfig
  - use -n \<name\>: modify current kubeconfig to the kubeconfig identified by name
- status command:
  - status: checks if the clusters referenced by each kubeconfig are reachable
    - \<name>\: checks if the cluster of \<name\> kubeconfig is reachable
- Compress the db into a single file, to enable
  - simplicity: the user has a single file with all the database, not a directory
  - backup/restore: easier to backup, restore or move between machines

## Differences with other projects
- [kubecm](https://github.com/sunny0826/kubecm): kubecm uses a single kubeconfig file, meanwhile this projects uses a key-value DB to store multiple kubeconfigs, but separated one from another. I use a lot of different kubeconfig files, and some of them are ephemeral (lifetime under 1 day), so it makes no point to merge them in a "master" kubeconfig.