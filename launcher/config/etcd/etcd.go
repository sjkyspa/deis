package config
import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/launcher/config/model"
	"net/url"
)

type EtcdBackend struct {
	client *etcd.Client
}

func (eb *EtcdBackend) Get(key string) (string, error) {
	sort, recursive := true, false
	resp, err := eb.client.Get(key, sort, recursive)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (eb *EtcdBackend) GetWithDefault(key, defaultValue string) (string, error) {
	sort, recursive := true, false
	resp, err := eb.client.Get(key, sort, recursive)
	if err != nil {
		etcdErr, ok := err.(*etcd.EtcdError)
		if ok && etcdErr.ErrorCode == 100 {
			return defaultValue, nil
		}
		return "", err
	}
	return resp.Node.Value, nil
}

func (eb *EtcdBackend) Set(key, value string) (string, error) {
	resp, err := eb.client.Set(key, value, 0) // don't use TTLs
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (eb *EtcdBackend) SetWithTTL(key string, value string, ttl uint64) (string, error) {
	resp, err := eb.client.Update(key, value, ttl)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (eb *EtcdBackend) Delete(key string) error {
	_, err := eb.client.Delete(key, false)
	return err
}

func (eb *EtcdBackend) GetRecursive(key string) ([]*model.ConfigNode, error) {
	resp, err := eb.client.Get(key, true, true)
	if err != nil {
		return nil, err
	}

	nodes := traverseNode(resp.Node)
	return nodes, nil
}


func singleNodeToConfigNode(node *etcd.Node) *model.ConfigNode {
	key := model.ConfigNode{
		Key:        node.Key,
		Expiration: node.Expiration,
	}

	if node.Dir != true && node.Key != "" {
		key.Value = node.Value
	}

	return &key
}

func traverseNode(node *etcd.Node) []*model.ConfigNode {
	var serviceKeys []*model.ConfigNode

	if len(node.Nodes) > 0 {
		for _, nodeChild := range node.Nodes {
			serviceKeys = append(serviceKeys, traverseNode(nodeChild)...)
		}
	} else {
		key := singleNodeToConfigNode(node)
		if key.Key != "" {
			serviceKeys = append(serviceKeys, key)
		}
	}

	return serviceKeys
}

func NewEtcdBackend(ep url.URL) (*EtcdBackend, error) {
	return &EtcdBackend{client: etcd.NewClient([]string{ep.String()})}, nil
}

