/*
 * Copyright (c) 2023. NetApp, Inc. All Rights Reserved.
 */

package connector

import (
	"context"
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/api/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/NetApp-Polaris/astra-connector-operator/app/conf"
	"github.com/NetApp-Polaris/astra-connector-operator/app/deployer/model"
	"github.com/NetApp-Polaris/astra-connector-operator/common"
	v1 "github.com/NetApp-Polaris/astra-connector-operator/details/operator-sdk/api/v1"
)

type AstraConnectNatlessDeployer struct{}

func NewAstraConnectorNatlessDeployer() model.Deployer {
	return &AstraConnectNatlessDeployer{}
}

// GetDeploymentObjects returns a Astra Connect Deployment object
func (d *AstraConnectNatlessDeployer) GetDeploymentObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	log := ctrllog.FromContext(ctx)
	ls := LabelsForAstraConnectClient(common.AstraConnectName, m.Spec.Labels)

	var imageRegistry string
	var containerImage string
	var connectorImage string
	if m.Spec.ImageRegistry.Name != "" {
		imageRegistry = m.Spec.ImageRegistry.Name
	} else {
		imageRegistry = common.DefaultImageRegistry
	}

	if m.Spec.AstraConnect.Image != "" {
		containerImage = m.Spec.AstraConnect.Image
	} else {
		containerImage = common.ConnectorImageTag
	}

	connectorImage = fmt.Sprintf("%s/astra-connector:%s", imageRegistry, containerImage)
	log.Info("Using AstraConnector image", "image", connectorImage)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.AstraConnectName,
			Namespace: m.Namespace,
			Annotations: map[string]string{
				"container.seccomp.security.alpha.kubernetes.io/pod": "runtime/default",
			},
		},
		Spec: appsv1.DeploymentSpec{
			// TODO remove option to set replica count in CRD. This should always only-ever be 1
			Replicas: &m.Spec.AstraConnect.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: connectorImage,
						Name:  common.AstraConnectName,
						Env: []corev1.EnvVar{
							{
								Name:  "LOG_LEVEL", // todo should this match what operator is
								Value: "trace",
							},
							{
								Name:  "NATS_DISABLED",
								Value: "true",
							},
							{
								Name:  "API_TOKEN_SECRET_REF",
								Value: m.Spec.Astra.TokenRef,
							},
							{
								Name:  "ASTRA_CONTROL_URL",
								Value: m.Spec.NatsSyncClient.CloudBridgeURL,
							},
							{
								Name:  "ACCOUNT_ID",
								Value: m.Spec.Astra.AccountId,
							},
							{
								Name:  "CLOUD_ID",
								Value: m.Spec.Astra.CloudId,
							},
							{
								Name:  "CLUSTER_ID",
								Value: m.Spec.Astra.ClusterId,
							},
							{
								Name:  "HOST_ALIAS_IP",
								Value: m.Spec.NatsSyncClient.HostAliasIP,
							},
							{
								Name:  "SKIP_TLS_VALIDATION",
								Value: strconv.FormatBool(m.Spec.Astra.SkipTLSValidation),
							},
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("0.1"),
							},
						},
						SecurityContext: conf.GetSecurityContext(),
					}},
					ServiceAccountName: common.AstraConnectName,
				},
			},
		},
	}

	if m.Spec.ImageRegistry.Secret != "" {
		dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: m.Spec.ImageRegistry.Secret,
			},
		}
	}
	return []client.Object{dep}, model.NonMutateFn, nil
}

// GetServiceObjects returns an Astra-Connect Service object
func (d *AstraConnectNatlessDeployer) GetServiceObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	return nil, model.NonMutateFn, nil
}

// GetConfigMapObjects returns a ConfigMap object for Astra Connect
func (d *AstraConnectNatlessDeployer) GetConfigMapObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Namespace,
			Name:      common.AstraConnectName,
		},
		Data: map[string]string{
			//"nats_url":            GetNatsURL(m),
			"skip_tls_validation": strconv.FormatBool(m.Spec.Astra.SkipTLSValidation),
		},
	}
	return []client.Object{configMap}, model.NonMutateFn, nil
}

// GetServiceAccountObjects returns a ServiceAccount object for Astra Connect
func (d *AstraConnectNatlessDeployer) GetServiceAccountObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.AstraConnectName,
			Namespace: m.Namespace,
		},
	}
	return []client.Object{sa}, model.NonMutateFn, nil
}

func (d *AstraConnectNatlessDeployer) GetClusterRoleObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: common.AstraConnectName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "persistentvolumes", "nodes", "pods", "services"},
				Verbs:     []string{"watch", "list", "get"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"update", "watch", "list", "get"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"csidrivers"},
				Verbs:     []string{"watch", "list", "get"},
			},
			{
				APIGroups: []string{"snapshot.storage.k8s.io"},
				Resources: []string{"volumesnapshotclasses"},
				Verbs:     []string{"watch", "list", "get"},
			},
			{
				APIGroups: []string{"trident.netapp.io"},
				Resources: []string{"tridentversions", "tridentorchestrators"},
				Verbs:     []string{"watch", "list", "get"},
			},
			{
				APIGroups: []string{"astra.netapp.io"},
				Resources: []string{
					"applications",
					"appmirrorrelationships",
					"appmirrorupdates",
					"appvaults",
					"autosupportbundles",
					"autosupportbundleschedules",
					"backups",
					"backupinplacerestores",
					"backuprestores",
					"exechooks",
					"exechooksruns",
					"pvccopies",
					"pvcerases",
					"resourcebackups",
					"resourcedeletes",
					"resourcerestores",
					"resourcesummaryuploads",
					"resticvolumebackups",
					"resticvolumerestores",
					"schedules",
					"shutdownsnapshots",
					"snapshots",
					"snapshotinplacerestores",
					"snapshotrestores",
					"astraconnectors",
				},
				Verbs: []string{"watch", "list", "get"},
			},
			{
				APIGroups: []string{"security.openshift.io"},
				Resources: []string{"securitycontextconstraints"},
				Verbs:     []string{"use"},
			},
		},
	}
	return []client.Object{clusterRole}, model.NonMutateFn, nil
}

func (d *AstraConnectNatlessDeployer) GetClusterRoleBindingObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: common.AstraConnectName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      common.AstraConnectName,
				Namespace: m.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     common.AstraConnectName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	return []client.Object{clusterRoleBinding}, model.NonMutateFn, nil
}

// GetRoleObjects returns a ConfigMapRole object for Astra Connect
func (d *AstraConnectNatlessDeployer) GetRoleObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.AstraConnectName,
			Namespace: m.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "update", "list", "get"},
			},
			{
				APIGroups: []string{"astra.netapp.io"},
				Resources: []string{
					"applications",
					"appmirrorrelationships",
					"appmirrorupdates",
					"appvaults",
					"autosupportbundles",
					"autosupportbundleschedules",
					"backups",
					"backupinplacerestores",
					"backuprestores",
					"exechooks",
					"exechooksruns",
					"pvccopies",
					"pvcerases",
					"resourcebackups",
					"resourcedeletes",
					"resourcerestores",
					"resourcesummaryuploads",
					"resticvolumebackups",
					"resticvolumerestores",
					"schedules",
					"shutdownsnapshots",
					"snapshots",
					"snapshotinplacerestores",
					"snapshotrestores",
					"astraconnectors",
				},
				Verbs: []string{"create", "update", "delete"},
			},
		},
	}
	return []client.Object{role}, model.NonMutateFn, nil
}

// GetRoleBindingObjects returns a ConfigMapRoleBinding object
func (d *AstraConnectNatlessDeployer) GetRoleBindingObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.AstraConnectName,
			Namespace: m.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      common.AstraConnectName,
				Namespace: m.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     common.AstraConnectName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	return []client.Object{roleBinding}, model.NonMutateFn, nil
}

// NIL RESOURCES BELOW
func (d *AstraConnectNatlessDeployer) GetStatefulSetObjects(m *v1.AstraConnector, ctx context.Context) ([]client.Object, controllerutil.MutateFn, error) {
	return nil, model.NonMutateFn, nil
}
