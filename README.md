# Kubeswap
Tool to manage multiple kubeconfig files and swap between clusters easily

## Why kubeswap
If you interact with a lot of kubernetes clusters/contexts, and you dont want to manage it in a big single kubeconfig file (merging kubeconfigs is tedious...), this is the right tool for you!

Just throw your kubeconfig files inside `$HOME/.kube/`, and kubeswap will manage it for you.

## Basic usage
Basically, you will use 2 commands:
- **kubeswap**: scans your `$HOME/.kube` dir and shows you a pretty interactive list to choose the desired kubeconfig
- **kubeswap \<name\>**: directly select the kubeconfig with that name from your `$HOME/.kube/` dir

- Extra: `kubeswap none` will remove the current **kubeconfig**.

## Advanced usage
Besides the basic usage, kubeswap has a key-value store, so you can:
- Add/delete kubeconfigs to/from the db
- List the kubeconfigs stored
- Select one to use
- Much more... (not really)

I have implemented the store with 2 objectives:
- **Portability**: you can use this db to store all your kubeconfigs and carry them with you
- **Backup/Restore**: you can use the db to backup/restore the kubeconfigs

To use the store, check the help :)

## Help
```bash
kubeswap help
```

## Tips
Use some shell/program that shows you your current k8s cluster/context.

## Quickstart
Basic (without store):
- Interactive list:
```bash
kubeswap
```
- Select one kubeconfig from your `$HOME/.kube/` directory:
```bash
kubeswap <filename>
```

Advanced (with store):
- add:
```bash
kubeswap add --name test --kubeconfig test/kubeconfig.yml --db /tmp/kubeswap.db
```
- list:
```bash
kubeswap list --db /tmp/kubeswap.db
```
- print:
```bash
kubeswap print -n test --db /tmp/kubeswap.db
```
- printall:
```bash
kubeswap printall --db /tmp/kubeswap.db
```
- use:
```bash
kubeswap use -n test --db /tmp/kubeswap.db
```
- delete:
```bash
kubeswap delete -n test --db /tmp/kubeswap.db
```

## TODO
- If the file to check is a folder just skip it
- Test in windows & mac
- status command: add concurrency
- Compress the db into a single file, to enable
  - simplicity: the user has a single file with all the database, not a directory
  - backup/restore: easier to backup, restore or move between machines
