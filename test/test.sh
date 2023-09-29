#!/bin/bash
kubeconfig_count=5
for (( c=0; c<${kubeconfig_count}; c++ )); do
    echo "Inserting kubeconfig test-${c}"
    ./dist/kubeswap_linux_amd64_v1/kubeswap add --name test-${c} --kubeconfig test/kubeconfig.yml --db /tmp/test.db
done