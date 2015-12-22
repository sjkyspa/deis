package config
import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/launcher/config/model"
)

type EtcdBackend struct {
	client *etcd.Client
}

func (*EtcdBackend) Get(string) (string, error) {
	return nil, nil
}

func (*EtcdBackend) GetWithDefault(string, string) (string, error) {
	return nil, nil
}

func (*EtcdBackend) Set(string, string) (string, error) {
	return nil, nil
}

func (*EtcdBackend) SetWithTTL(string, string, uint64) (string, error) {
	return nil, nil
}

func (*EtcdBackend) Delete(string) error {
	return nil, nil
}

func (*EtcdBackend) GetRecursive(string) ([]*model.ConfigNode, error) {
	return nil, nil
}


