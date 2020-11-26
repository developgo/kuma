package vips

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	mesh_core "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/manager"

	"github.com/kumahq/kuma/pkg/core/resources/model"

	config_manager "github.com/kumahq/kuma/pkg/core/config/manager"
	config_model "github.com/kumahq/kuma/pkg/core/resources/apis/system"
	"github.com/kumahq/kuma/pkg/core/resources/store"
)

type Persistence struct {
	configManager   config_manager.ConfigManager
	resourceManager manager.ReadOnlyResourceManager
}

const template = "kuma-%s-dns-vips"

var re = regexp.MustCompile(fmt.Sprintf(template, `(.*)`))

func ConfigKey(mesh string) string {
	return fmt.Sprintf(template, mesh)
}

func MeshFromConfigKey(name string) (string, bool) {
	match := re.FindStringSubmatch(name)
	if len(match) < 2 {
		return "", false
	}
	mesh := match[1]
	return mesh, true
}

func NewPersistence(resourceManager manager.ReadOnlyResourceManager, configManager config_manager.ConfigManager) *Persistence {
	return &Persistence{
		resourceManager: resourceManager,
		configManager:   configManager,
	}
}

func (m *Persistence) Get() (List, error) {
	resourceList := &config_model.ConfigResourceList{}
	if err := m.configManager.List(context.Background(), resourceList); err != nil {
		return nil, err
	}

	var errs error
	vips := List{}
	for _, resource := range resourceList.Items {
		if _, ok := MeshFromConfigKey(resource.Meta.GetName()); !ok {
			continue
		}
		if resource.Spec.Config == "" {
			continue
		}
		v := List{}
		if err := json.Unmarshal([]byte(resource.Spec.Config), &v); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		vips.Append(v)
	}
	return vips, nil
}

func (m *Persistence) GetByMesh(mesh string) (List, error) {
	name := fmt.Sprintf(template, mesh)
	vips := List{}
	resource := &config_model.ConfigResource{}
	err := m.configManager.Get(context.Background(), resource, store.GetByKey(name, ""))
	if err != nil {
		if store.IsResourceNotFound(err) {
			return vips, nil
		}
		return nil, err
	}

	if resource.Spec.Config == "" {
		return vips, nil
	}

	err = json.Unmarshal([]byte(resource.Spec.Config), &vips)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal")
	}

	return vips, nil
}

func (m *Persistence) Set(mesh string, vips List) error {
	ctx := context.Background()
	name := fmt.Sprintf(template, mesh)
	resource := &config_model.ConfigResource{}
	create := false
	if err := m.configManager.Get(ctx, resource, store.GetByKey(name, model.NoMesh)); err != nil {
		if !store.IsResourceNotFound(err) {
			return err
		}
		create = true
	}

	jsonBytes, err := json.Marshal(vips)
	if err != nil {
		return errors.Wrap(err, "unable to marshall VIP list")
	}
	resource.Spec.Config = string(jsonBytes)

	if create {
		meshRes := &mesh_core.MeshResource{}
		if err := m.resourceManager.Get(ctx, meshRes, store.GetByKey(mesh, model.NoMesh)); err != nil {
			return err
		}
		return m.configManager.Create(ctx, resource, store.CreateByKey(name, model.NoMesh), store.CreateWithOwner(meshRes))
	} else {
		return m.configManager.Update(ctx, resource)
	}
}