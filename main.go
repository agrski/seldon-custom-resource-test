package main

import (
  "fmt"
  "os"
  "path/filepath"

  seldonclientset "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
  seldondeployment "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned/typed/machinelearning.seldon.io/v1"
  "k8s.io/client-go/tools/clientcmd"
)

const k8sNamespace = "seldon"

func getSeldonDeploymentsClient() seldondeployment.SeldonDeploymentInterface {
  kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

  config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
  if err != nil {
    panic(err)
  }

  clientset, err := seldonclientset.NewForConfig(config)
  if err != nil {
    panic(err)
  }

  return clientset.MachinelearningV1().SeldonDeployments(k8sNamespace)
}

func main() {
  getSeldonDeploymentsClient()
  fmt.Println("Created Seldon client")
}

