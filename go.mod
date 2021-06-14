module github.com/flanksource/konfig-manager

go 1.16

require (
	github.com/flanksource/commons v1.5.6
	github.com/flanksource/kommons v0.18.0
	github.com/go-logr/logr v0.3.0
	github.com/hairyhenderson/gomplate/v3 v3.6.0
	github.com/magiconair/properties v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.8.3
)

replace k8s.io/client-go => k8s.io/client-go v0.20.4
