# Pre-requisites

## Dependencies

* [Go](https://golang.org/doc/install)
* [Docker](https://docs.docker.com/engine/install/)
  * Allow Docker to be run as [non-root](https://docs.docker.com/engine/install/linux-postinstall/)
* [Kind](https://kind.sigs.k8s.io/docs/user/quick-start)
* [Helm](https://docs.helm.sh/docs/intro/install/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
* [Seldon-Core](https://docs.seldon.io/projects/seldon-core/en/v1.1.0/workflow/install.html)

## Environment

The `run.sh` script in this directory takes care of the below steps,
which detail how to prepare the environment for running tests and clean up
thereafter.

* Set up a k8s cluster using kind.
* Set up a namespace for Seldon to deploy into, e.g. `seldon-system`.
* Install Seldon Core.

After the tests have been run, all k8s resources can be torn down by deleting
the kind cluster.

