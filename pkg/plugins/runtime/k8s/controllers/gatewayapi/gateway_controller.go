package gatewayapi

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	kube_core "k8s.io/api/core/v1"
	kube_apierrs "k8s.io/apimachinery/pkg/api/errors"
	kube_meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube_runtime "k8s.io/apimachinery/pkg/runtime"
	kube_types "k8s.io/apimachinery/pkg/types"
	kube_ctrl "sigs.k8s.io/controller-runtime"
	kube_client "sigs.k8s.io/controller-runtime/pkg/client"
	kube_controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	kube_handler "sigs.k8s.io/controller-runtime/pkg/handler"
	kube_reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
	kube_source "sigs.k8s.io/controller-runtime/pkg/source"
	gatewayapi "sigs.k8s.io/gateway-api/apis/v1alpha2"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/pkg/core/resources/manager"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	mesh_k8s "github.com/kumahq/kuma/pkg/plugins/resources/k8s/native/api/v1alpha1"
	k8s_registry "github.com/kumahq/kuma/pkg/plugins/resources/k8s/native/pkg/registry"
	"github.com/kumahq/kuma/pkg/plugins/runtime/k8s/containers"
	"github.com/kumahq/kuma/pkg/plugins/runtime/k8s/controllers/gatewayapi/common"
	k8s_util "github.com/kumahq/kuma/pkg/plugins/runtime/k8s/util"
)

// GatewayReconciler reconciles a GatewayAPI Gateway object.
type GatewayReconciler struct {
	kube_client.Client
	Log logr.Logger

	Scheme          *kube_runtime.Scheme
	TypeRegistry    k8s_registry.TypeRegistry
	SystemNamespace string
	ProxyFactory    containers.DataplaneProxyFactory
	ResourceManager manager.ResourceManager
}

// Reconcile handles transforming a gateway-api Gateway into a Kuma Gateway and
// managing the status of the gateway-api objects.
func (r *GatewayReconciler) Reconcile(ctx context.Context, req kube_ctrl.Request) (kube_ctrl.Result, error) {
	gateway := &gatewayapi.Gateway{}
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		if kube_apierrs.IsNotFound(err) {
			// We don't know the mesh, but we don't need it to delete our
			// object.
			err := common.ReconcileLabelledObject(ctx, r.TypeRegistry, r.Client, req.NamespacedName, core_model.NoMesh, &mesh_proto.Gateway{}, nil)
			return kube_ctrl.Result{}, errors.Wrap(err, "could not delete owned Gateway.kuma.io")
		}

		return kube_ctrl.Result{}, err
	}

	class, err := common.GetGatewayClass(ctx, r.Client, gateway.Spec.GatewayClassName)
	if err != nil {
		return kube_ctrl.Result{}, errors.Wrap(err, "unable to retrieve GatewayClass referenced by Gateway")
	}

	if class.Spec.ControllerName != common.ControllerName {
		return kube_ctrl.Result{}, nil
	}

	gatewaySpec, listenerConditions, err := r.gapiToKumaGateway(ctx, gateway)
	if err != nil {
		return kube_ctrl.Result{}, errors.Wrap(err, "error generating Gateway.kuma.io")
	}

	ns := kube_core.Namespace{}
	if err := r.Client.Get(ctx, kube_types.NamespacedName{Name: gateway.Namespace}, &ns); err != nil {
		return kube_ctrl.Result{}, errors.Wrap(err, "unable to get Namespace of Gateway")
	}

	mesh := k8s_util.MeshOf(gateway, &ns)

	if err := common.ReconcileLabelledObject(ctx, r.TypeRegistry, r.Client, req.NamespacedName, mesh, &mesh_proto.Gateway{}, gatewaySpec); err != nil {
		return kube_ctrl.Result{}, errors.Wrap(err, "could not reconcile owned Gateway.kuma.io")
	}

	gatewayInstance, err := r.createOrUpdateInstance(ctx, gateway)
	if err != nil {
		return kube_ctrl.Result{}, errors.Wrap(err, "unable to reconcile GatewayInstance")
	}

	if err := r.updateStatus(ctx, gateway, gatewayInstance, listenerConditions); err != nil {
		return kube_ctrl.Result{}, errors.Wrap(err, "unable to update Gateway status")
	}

	return kube_ctrl.Result{}, nil
}

func (r *GatewayReconciler) createOrUpdateInstance(ctx context.Context, gateway *gatewayapi.Gateway) (*mesh_k8s.GatewayInstance, error) {
	instance := &mesh_k8s.GatewayInstance{
		ObjectMeta: kube_meta.ObjectMeta{
			Namespace: gateway.Namespace,
			Name:      gateway.Name,
		},
	}

	if _, err := kube_controllerutil.CreateOrUpdate(ctx, r.Client, instance, func() error {
		instance.Spec = mesh_k8s.GatewayInstanceSpec{
			ServiceType: kube_core.ServiceTypeLoadBalancer,
			Tags:        common.ServiceTagForGateway(kube_client.ObjectKeyFromObject(gateway)),
		}

		err := kube_controllerutil.SetControllerReference(gateway, instance, r.Scheme)
		return errors.Wrap(err, "unable to set GatewayInstance's controller reference to Gateway")
	}); err != nil {
		return nil, errors.Wrap(err, "couldn't create GatewayInstance")
	}

	return instance, nil
}

const gatewayIndexField = ".metadata.gateway"

// gatewaysForRoute returns a function that calculates which Gateways might
// be affected by changes in an HTTPRoute so they can be reconciled.
func gatewaysForRoute(l logr.Logger, client kube_client.Client) kube_handler.MapFunc {
	l = l.WithName("gatewaysForRoute")

	return func(obj kube_client.Object) []kube_reconcile.Request {
		route, ok := obj.(*gatewayapi.HTTPRoute)
		if !ok {
			l.Error(nil, "unexpected error converting to be mapped %T object to HTTPRoute", obj)
			return nil
		}

		var requests []kube_reconcile.Request
		for _, parentRef := range route.Spec.ParentRefs {
			namespace := route.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}

			requests = append(
				requests,
				kube_reconcile.Request{
					NamespacedName: kube_types.NamespacedName{Namespace: namespace, Name: string(parentRef.Name)},
				},
			)
		}

		return requests
	}
}

func (r *GatewayReconciler) SetupWithManager(mgr kube_ctrl.Manager) error {
	// This index helps us list routes that point to a Gateway in
	// attachedListenersForGateway.
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &gatewayapi.HTTPRoute{}, gatewayIndexField, func(obj kube_client.Object) []string {
		route := obj.(*gatewayapi.HTTPRoute)

		var names []string

		for _, parentRef := range route.Spec.ParentRefs {
			namespace := route.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}

			names = append(
				names,
				kube_types.NamespacedName{Namespace: namespace, Name: string(parentRef.Name)}.String(),
			)
		}

		return names
	}); err != nil {
		return err
	}

	return kube_ctrl.NewControllerManagedBy(mgr).
		For(&gatewayapi.Gateway{}).
		Owns(&mesh_k8s.Gateway{}).
		Owns(&mesh_k8s.GatewayInstance{}).
		Watches(
			&kube_source.Kind{Type: &gatewayapi.HTTPRoute{}},
			kube_handler.EnqueueRequestsFromMapFunc(gatewaysForRoute(r.Log, r.Client)),
		).
		Complete(r)
}
