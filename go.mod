module github.com/flanksource/konfig-manager

go 1.16

require (
	github.com/flanksource/commons v1.5.6
	github.com/flanksource/kommons v0.18.0
	github.com/hairyhenderson/gomplate/v3 v3.6.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.3
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/apimachinery v0.20.4
)

replace k8s.io/client-go => k8s.io/client-go v0.20.4
