package main

import (
  "errors"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"

  seldonapi "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
  seldonclientset "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
  seldonscheme "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned/scheme"
  seldondeployment "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned/typed/machinelearning.seldon.io/v1"
  "k8s.io/client-go/tools/clientcmd"
)

const k8sNamespace = "seldon"

func getSeldonDeployment() (*seldonapi.SeldonDeployment, error) {
  var fileName string
  flag.StringVar(&fileName, "filename", "", "Name of file containing Seldon Core custom resource")
  flag.Parse()

  if fileName == "" {
    return nil, errors.New("filename must be provided")
  }

  fileContents, err := ioutil.ReadFile(fileName)
  if err != nil {
    return nil, errors.New("File must exist and be readable")
  }

  decode := seldonscheme.Codecs.UniversalDeserializer().Decode
  obj, _, err := decode(fileContents, nil, nil)
  if err != nil {
    return nil, errors.New("Unable to decode file contents")
  }

  deployment := obj.(*seldonapi.SeldonDeployment)

  return deployment, nil
}

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

