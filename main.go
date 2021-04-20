package main

import (
  "context"
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
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/watch"
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

  switch obj.(type) {
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

func createAndWaitForDeployment(
  deploymentClient seldondeployment.SeldonDeploymentInterface,
  deployment *seldonapi.SeldonDeployment,
) error {
  deploymentName := deployment.ObjectMeta.Name

  _, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
  if err != nil {
    return err
  }

  watcher, err := deploymentClient.Watch(context.TODO(), metav1.ListOptions{})
  if err != nil {
    return err
  }

  for e := range watcher.ResultChan() {
    switch e.Type {
    case watch.Modified:
      updatedDeployment := e.Object.(*seldonapi.SeldonDeployment)

      if updatedDeployment.ObjectMeta.Name == deploymentName &&
      updatedDeployment.Status.State == seldonapi.StatusStateAvailable {
        watcher.Stop()
        return nil
      }
    case watch.Error:
      watcher.Stop()
      return errors.New("SeldonDeployment could not be created")
    case watch.Deleted:
      watcher.Stop()
      return errors.New("SeldonDeployment was deleted unexpectedly")
    }
  }

  return errors.New(fmt.Sprintf("Deployment '%s' did not become ready", deploymentName))
}

func main() {
  manifest, err := getResourceManifest()
  if err != nil {
    panic(err)
  }

  deployment, err := getSeldonDeployment(manifest)
  if err != nil {
    panic(err)
  }

  deploymentClient, err := getSeldonDeploymentsClient()
  if err != nil {
    panic(err)
  }

  err = createAndWaitForDeployment(deploymentClient, deployment)
  if err != nil {
    panic(err)
  }
  fmt.Println("Deployment created successfully")

  deploymentClient.Delete(
    context.TODO(),
    deployment.ObjectMeta.Name,
    metav1.DeleteOptions{},
  )
  fmt.Println("Deployment deleted")
}

