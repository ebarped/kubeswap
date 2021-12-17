#!/bin/bash
kubeconfig_count=5
for (( c=0; c<${kubeconfig_count}; c++ )); do
    echo "Inserting kubeconfig ${c}"
    ./dist/kubeswap_linux_amd64/kubeswap add --name test-${c} --kubeconfig test/kubeconfig.yml --db /tmp/test
done