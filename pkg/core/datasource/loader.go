package datasource

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"

	system_proto "github.com/kumahq/kuma/api/system/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/apis/system"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
)

type SecretLoader interface {
	Get(ctx context.Context, mesh string, secretName string) ([]byte, error)
}

type storeSecretLoader struct {
	resourceManager manager.ReadOnlyResourceManager
}

func (l *storeSecretLoader) Get(ctx context.Context, mesh string, secret string) ([]byte, error) {
	if l.resourceManager == nil {
		return nil, errors.New("no resource manager")
	}
	resource := system.NewSecretResource()
	if err := l.resourceManager.Get(ctx, resource, core_store.GetByKey(secret, mesh)); err != nil {
		return nil, err
	}
	return resource.Spec.GetData().GetValue(), nil
}

func NewStoreSecretLoader(resourceManager manager.ReadOnlyResourceManager) SecretLoader {
	return &storeSecretLoader{
		resourceManager: resourceManager,
	}
}

var _ SecretLoader = &storeSecretLoader{}

type meshSecretLoader struct {
	secrets map[string][]byte
}

func (l *meshSecretLoader) Get(ctx context.Context, mesh string, secret string) ([]byte, error) {
	if l.secrets == nil {
		return nil, errors.New("secrets not initialized")
	}
	data, ok := l.secrets[secret]
	if !ok {
		return nil, fmt.Errorf("secret not found: name=\"%s\" mesh=\"%s\"", secret, mesh)
	}

	return data, nil
}

func NewMeshSecretLoader(secrets map[string][]byte) SecretLoader {
	return &meshSecretLoader{
		secrets: secrets,
	}
}

var _ SecretLoader = &meshSecretLoader{}

type Loader interface {
	Load(ctx context.Context, mesh string, source *system_proto.DataSource) ([]byte, error)
}

type loader struct {
	secretLoader SecretLoader
}

var _ Loader = &loader{}

func NewDataSourceLoader(secretLoader SecretLoader) Loader {
	return &loader{
		secretLoader: secretLoader,
	}
}

func (l *loader) Load(ctx context.Context, mesh string, source *system_proto.DataSource) ([]byte, error) {
	var data []byte
	var err error
	switch source.GetType().(type) {
	case *system_proto.DataSource_Secret:
		data, err = l.secretLoader.Get(ctx, mesh, source.GetSecret())
	case *system_proto.DataSource_Inline:
		data, err = source.GetInline().GetValue(), nil
	case *system_proto.DataSource_InlineString:
		data, err = []byte(source.GetInlineString()), nil
	case *system_proto.DataSource_File:
		data, err = os.ReadFile(source.GetFile())
	default:
		return nil, errors.New("unsupported type of the DataSource")
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not load data")
	}
	return data, nil
}
