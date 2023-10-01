package kubeconfig

import (
	"fmt"

	k8sclientcmd "k8s.io/client-go/tools/clientcmd"
)

type Kubeconfig struct {
	Name           string
	Path           string
	Content        string
	CurrentContext string
}

func New(name, path string) (*Kubeconfig, error) {
	// load kubeconfig and check if it is valid (maybe we should check before trying to load to speed up?)
	kubeconfig, err := k8sclientcmd.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	content, err := k8sclientcmd.Write(*kubeconfig)
	if err != nil {
		return nil, err
	}

	return &Kubeconfig{
		Name:           name,
		Path:           path,
		Content:        string(content),
		CurrentContext: kubeconfig.CurrentContext,
	}, nil
}

// String prints a Kubeconfig struct
func (k *Kubeconfig) String() string {
	return fmt.Sprintf("name:%s, path:%s, file:%s\n", k.Name, k.Path, k.Content)
}
