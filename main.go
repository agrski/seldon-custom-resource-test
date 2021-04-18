package main

import (
  "fmt"
  "os"
  "path/filepath"

  seldon "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
  "k8s.io/client-go/tools/clientcmd"
)

func main() {
  kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

  config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
  if err != nil {
    panic(err)
  }

  clientset, err := seldon.NewForConfig(config)
  if err != nil {
    panic(err)
  }

  clientset.MachinelearningV1().SeldonDeployments("seldon")

  fmt.Println("Created Seldon client")
}

