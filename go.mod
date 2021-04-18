module seldon-custom-resource-test

go 1.14

require (
	github.com/seldonio/seldon-core/operator v0.0.0-20210412120902-357ba4d479ff // indirect
	k8s.io/client-go v12.0.0+incompatible // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.18.8
