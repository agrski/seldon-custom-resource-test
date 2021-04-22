package main

import (
  "context"
  "errors"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"

  seldonapi "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
  seldonclientset "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
  seldonscheme "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned/scheme"
  seldondeployment "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned/typed/machinelearning.seldon.io/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  corev1 "k8s.io/api/core/v1"
  "k8s.io/apimachinery/pkg/watch"
  "k8s.io/client-go/kubernetes"
  restclient "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
)

const k8sNamespace = "seldon"

func getConfig() (*restclient.Config, error) {
  kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

  config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
  return config, err
}

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
  config, err := getConfig()
  if err != nil {
    return nil, err
  }

  clientset, err := seldonclientset.NewForConfig(config)
  if err != nil {
    return nil, err
  }

  return clientset.MachinelearningV1().SeldonDeployments(k8sNamespace), nil
}

func awaitDeploymentAvailability(
  deploymentClient seldondeployment.SeldonDeploymentInterface,
  name string,
) error {
  watcher, err := deploymentClient.Watch(context.TODO(), metav1.ListOptions{})
  if err != nil {
    return err
  }

  for e := range watcher.ResultChan() {
    switch e.Type {
    case watch.Modified:
      updatedDeployment := e.Object.(*seldonapi.SeldonDeployment)

      if updatedDeployment.ObjectMeta.Name == name &&
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

  return errors.New(fmt.Sprintf("Deployment '%s' did not become ready", name))
}

func createDeployment(
  deploymentClient seldondeployment.SeldonDeploymentInterface,
  deployment *seldonapi.SeldonDeployment,
) error {
  _, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
  return err
}

func scaleDeployment(
  deploymentClient seldondeployment.SeldonDeploymentInterface,
  name string,
  replicas int32,
) error {
  deployment, err := deploymentClient.Get(
    context.TODO(),
    name,
    metav1.GetOptions{},
  )
  if err != nil {
    return err
  }

  deployment.Spec.Replicas = &replicas

  _, err = deploymentClient.Update(
    context.TODO(),
    deployment,
    metav1.UpdateOptions{},
  )
  return err
}

func manageDeploymentLifecycle() error {
  manifest, err := getResourceManifest()
  if err != nil {
    return err
  }

  deployment, err := getSeldonDeployment(manifest)
  if err != nil {
    return err
  }

  deploymentName := deployment.ObjectMeta.Name

  deploymentClient, err := getSeldonDeploymentsClient()
  if err != nil {
    return err
  }

  err = createDeployment(deploymentClient, deployment)
  if err != nil {
    return err
  }

  err = awaitDeploymentAvailability(deploymentClient, deploymentName)
  if err != nil {
    return err
  }

  fmt.Println("Deployment created successfully")

  desiredReplicas := int32(2)
  err = scaleDeployment(deploymentClient, deploymentName, desiredReplicas)
  if err != nil {
    return err
  }

  err = awaitDeploymentAvailability(deploymentClient, deploymentName)
  if err != nil {
    return err
  }

  fmt.Printf("Deployment scaled to %v replicas\n", desiredReplicas)

  deploymentClient.Delete(
    context.TODO(),
    deployment.ObjectMeta.Name,
    metav1.DeleteOptions{},
  )

  fmt.Println("Deployment deleted")

  return nil
}

func describeEvents() error {
  config, err := getConfig()
  if err != nil {
    return err
  }

  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    return err
  }

  eventsClient := clientset.CoreV1().Events(k8sNamespace)

  watcher, err := eventsClient.Watch(context.TODO(), metav1.ListOptions{})
  if err != nil {
    return err
  }

  for e := range watcher.ResultChan() {
    event := e.Object.(*corev1.Event)
    // FIXME - resource name should not be hard-coded.
    if event.InvolvedObject.Name == "seldon-model" {
      fmt.Println(event.Message)
    }
  }

  return nil
}

func main() {
  go func() {
    describeEvents()
  }()

  err := manageDeploymentLifecycle()
  if err != nil {
    log.Fatal(err)
  }
}

