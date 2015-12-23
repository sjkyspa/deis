package model

type Container struct {
	Name string
	Desc ContainerDesc
}

type ContainerDesc struct {
	Image string `yaml:"image"`
	Links []string `yaml:"links"`
	Ports []string `yaml:"ports"`
}
