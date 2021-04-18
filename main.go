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

func getResourceManifest() ([]byte, error) {
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

  return fileContents, nil
}

func getSeldonDeployment(manifest []byte) (*seldonapi.SeldonDeployment, error) {
  decode := seldonscheme.Codecs.UniversalDeserializer().Decode
  obj, _, err := decode(manifest, nil, nil)
  if err != nil {
    return nil, errors.New("Unable to decode file contents")
  }

  switch t := obj.(type) {
  case *seldonapi.SeldonDeployment:
    return obj.(*seldonapi.SeldonDeployment), nil
  default:
    return nil, nil
  }
}

func getSeldonDeploymentsClient() (seldondeployment.SeldonDeploymentInterface, error) {
  kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

  config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
  if err != nil {
    return nil, err
  }

  clientset, err := seldonclientset.NewForConfig(config)
  if err != nil {
    return nil, err
  }

  return clientset.MachinelearningV1().SeldonDeployments(k8sNamespace), nil
}

func main() {
  getSeldonDeploymentsClient()
  fmt.Println("Created Seldon client")
}

