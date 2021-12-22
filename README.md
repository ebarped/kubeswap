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
  use         Select kubeconfig to use
  version     Print the version number

Flags:
      --db string         db file path (default "$HOME/.kube/kubeswap.db")
  -h, --help              help for kubeswap
      --loglevel string   loglevel (info/debug) (default "info")

Use "kubeswap [command] --help" for more information about a command.
```

# Commands
- kubeswap: displays the list of kubeconfigs in the db, same as subcommand use but without specifying a name

## Tip
I like to create an alias to kubeswap:
`alias ks=kubeswap`

## Technologies
- cobra: cli library
- pogreb: key-value database to store the files
- zerolog: structured (and fast) logger
- bubbletea: 

## Quickstart (test basic usage)
- Init:
```bash
make clean build test
```
- add:
```bash
./dist/kubeswap_linux_amd64/kubeswap add --name test --kubeconfig test/kubeconfig.yml
```
- list:
```bash
./dist/kubeswap_linux_amd64/kubeswap list
```
- print:
```bash
./dist/kubeswap_linux_amd64/kubeswap print -n test-0
```
- printall:
```bash
./dist/kubeswap_linux_amd64/kubeswap printall
```
- delete:
```bash
./dist/kubeswap_linux_amd64/kubeswap delete -n test-0
```
- use:
```bash
./dist/kubeswap_linux_amd64/kubeswap use -n test-1
```
- use interactive:
```bash
./dist/kubeswap_linux_amd64/kubeswap use
```

## TODO
- Test in windows & mac
- status command:
  - status: checks if the clusters referenced by each kubeconfig are reachable
    - \<name>\: checks if the cluster of \<name\> kubeconfig is reachable
- Compress the db into a single file, to enable
  - simplicity: the user has a single file with all the database, not a directory
  - backup/restore: easier to backup, restore or move between machines
- Move all the complexity of the interactive list into the tui internal package