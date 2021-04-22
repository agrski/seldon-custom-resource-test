module seldon-custom-resource-test

go 1.14

require (
	github.com/seldonio/seldon-core/operator v0.0.0-20210412120902-357ba4d479ff
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.18.8
