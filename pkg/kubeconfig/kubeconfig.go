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

	// usage of this tool only makes sense if the kubeconfig has current-context set, so filter out other files
	if kubeconfig.CurrentContext == "" {
		return nil, fmt.Errorf("%q file should have 'current-context' field set", path)
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
