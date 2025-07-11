package kubeconfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
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

// Reachable checks if the cluster is reachable by network and if has
// permissions to get "server version"
func (k *Kubeconfig) Reachable() bool {
	// kubeconfigs from the db have no path, store them to load and check
	if k.Path == "" {
		f, err := os.CreateTemp("", k.Name)
		if err != nil {
			log.Fatalf("error creating tmp file to store kubeconfig: %s\n", err)
		}
		defer os.Remove(f.Name())

		_, err = f.WriteString(k.Content)
		if err != nil {
			log.Fatalf("error writing kubeconfig to tmp file: %s\n", err)
		}

		// Get the absolute path
		absPath, err := filepath.Abs(f.Name())
		if err != nil {
			log.Fatalf("error getting absolute path of tmp file: %s\n", err)
		}

		k.Path = absPath
	}

	config, err := clientcmd.BuildConfigFromFlags("", k.Path)
	if err != nil {
		log.Fatalf("error loading kubeconfig: %s\n", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	_, err = discoveryClient.ServerVersion()

	return err == nil
}
