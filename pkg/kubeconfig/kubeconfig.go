package kubeconfig

import (
	"fmt"

	k8sclientcmd "k8s.io/client-go/tools/clientcmd"
	k8sclientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Kubeconfig struct {
	name   string
	path   string
	config *k8sclientcmdapi.Config
}

func New(name, path string) (*Kubeconfig, error) {
	// load kubeconfig and check if it is valid
	kubeconfig, err := k8sclientcmd.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig from file")
	}

	return &Kubeconfig{
		name:   name,
		path:   path,
		config: kubeconfig,
	}, nil
}

// String prints a Kubeconfig struct
func (k *Kubeconfig) String() string {
	config, err := k.Config()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("name:%s, path:%s, file:%s\n", k.name, k.path, config)
}

// Name returns the name key of the Kubeconfig struct
func (k *Kubeconfig) Name() string {
	return k.name
}

// Config returns the kubeconfig content of the Kubeconfig struct
func (k *Kubeconfig) Config() ([]byte, error) {
	out, err := k8sclientcmd.Write(*k.config)
	if err != nil {
		return nil, err
	}
	return out, err
}
