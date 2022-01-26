package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	cachev1 "github.com/NetApp/astraagent-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// DeploymentForNatssyncClient returns a Natssync-client Deployment object
func (r *AstraAgentReconciler) DeploymentForNatssyncClient(m *cachev1.AstraAgent) (*appsv1.Deployment, error) {
	ls := labelsForNatssyncClient(m.Spec.NatssyncClient.Name)
	replicas := m.Spec.NatssyncClient.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.NatssyncClient.Name,
			Namespace: m.Spec.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: m.Spec.NatssyncClient.Image,
						Name:  m.Spec.NatssyncClient.Name,
						Env: []corev1.EnvVar{
							{
								Name:  "NATS_SERVER_URL",
								Value: r.GetNatsURL(m),
							},
							{
								Name:  "CLOUD_BRIDGE_URL",
								Value: m.Spec.NatssyncClient.CloudBridgeURL,
							},
							{
								Name:  "CONFIGMAP_NAME",
								Value: m.Spec.ConfigMap.Name,
							},
							{
								Name:  "POD_NAMESPACE",
								Value: m.Spec.Namespace,
							},
						},
					}},
					Volumes: []v1.Volume{
						{
							Name: fmt.Sprintf("%s-configmap-volume", m.Spec.NatssyncClient.Name),
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: m.Spec.ConfigMap.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Set astraAgent instance as the owner and controller
	err := ctrl.SetControllerReference(m, dep, r.Scheme)
	if err != nil {
		return nil, err
	}
	return dep, nil
}

// ServiceForNatssyncClient returns a Natssync-client Service object
func (r *AstraAgentReconciler) ServiceForNatssyncClient(m *cachev1.AstraAgent) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.NatssyncClient.Name,
			Namespace: m.Spec.Namespace,
			Labels: map[string]string{
				"app": m.Spec.NatssyncClient.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Port:     m.Spec.NatssyncClient.Port,
					NodePort: m.Spec.NatssyncClient.NodePort,
					Protocol: m.Spec.NatssyncClient.Protocol,
				},
			},
			Selector: map[string]string{
				"app": m.Spec.NatssyncClient.Name,
			},
		},
	}
	// Set astraAgent instance as the owner and controller
	ctrl.SetControllerReference(m, service, r.Scheme)
	return service
}

// labelsForNatssyncClient returns the labels for selecting the resources
// belonging to the given astraAgent CR name.
func labelsForNatssyncClient(name string) map[string]string {
	return map[string]string{"app": name}
}

// ConfigMap returns a astraAgent ConfigMap object
func (r *AstraAgentReconciler) ConfigMap(m *cachev1.AstraAgent) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Spec.Namespace,
			Name:      m.Spec.ConfigMap.Name,
		},
	}
	ctrl.SetControllerReference(m, configMap, r.Scheme)
	return configMap
}

// ConfigMapRole returns a astraAgent ConfigMap object
func (r *AstraAgentReconciler) ConfigMapRole(m *cachev1.AstraAgent) *rbacv1.Role {
	configMapRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Spec.Namespace,
			Name:      m.Spec.ConfigMap.RoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get", "list", "patch"},
			},
		},
	}
	ctrl.SetControllerReference(m, configMapRole, r.Scheme)
	return configMapRole
}

// ConfigMapRoleBinding returns a Natssync-Client ConfigMapRoleBinding object
func (r *AstraAgentReconciler) ConfigMapRoleBinding(m *cachev1.AstraAgent) *rbacv1.RoleBinding {
	configMapRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Spec.Namespace,
			Name:      m.Spec.ConfigMap.RoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: m.Spec.ConfigMap.ServiceAccountName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     m.Spec.ConfigMap.RoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	ctrl.SetControllerReference(m, configMapRoleBinding, r.Scheme)
	return configMapRoleBinding
}

// ServiceAccountForConfigMap returns a ServiceAccount object
func (r *AstraAgentReconciler) ServiceAccountForConfigMap(m *cachev1.AstraAgent) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.ConfigMap.ServiceAccountName,
			Namespace: m.Spec.Namespace,
		},
	}
	ctrl.SetControllerReference(m, sa, r.Scheme)
	return sa
}

func (r *AstraAgentReconciler) getNatssyncClientStatus(m *cachev1.AstraAgent, ctx context.Context) (cachev1.NatssyncClientStatus, error) {
	pods := &corev1.PodList{}
	lb := labelsForNatssyncClient(m.Spec.NatssyncClient.Name)
	listOpts := []client.ListOption{
		client.MatchingLabels(lb),
	}
	log := ctrllog.FromContext(ctx)

	if err := r.List(ctx, pods, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Namespace", m.Spec.Namespace)
		return cachev1.NatssyncClientStatus{}, err
	}

	natssyncClientStatus := cachev1.NatssyncClientStatus{
		Registered: "Unknown",
	}

	if len(pods.Items) < 1 {
		return cachev1.NatssyncClientStatus{}, errors.New("natssync-client pods not found")
	}
	nsClientPod := pods.Items[0]
	// If a pod is terminating, then we can't access the corresponding vault node's status.
	// so we break from here and return an error.
	if nsClientPod.Status.Phase != v1.PodRunning || nsClientPod.DeletionTimestamp != nil {
		return cachev1.NatssyncClientStatus{}, errors.New("natssync-client pod is terminating")
	}

	natssyncClientStatus.State = string(nsClientPod.Status.Phase)
	natsSyncClientURL := fmt.Sprintf("http://%s.%s:%d/bridge-client/1", m.Spec.NatssyncClient.Name, m.Spec.Namespace, m.Spec.NatssyncClient.Port)
	natsSyncClientRegisterURL := fmt.Sprintf("%s/register", natsSyncClientURL)
	natsSyncClientAboutURL := fmt.Sprintf("%s/about", natsSyncClientURL)
	natssyncClientRegistrationStatus, err := r.getNatssyncClientRegistrationStatus(natsSyncClientRegisterURL)
	if err != nil {
		log.Error(err, "Failed to get the registration status")
		return cachev1.NatssyncClientStatus{}, err
	}
	natssyncClientVersion, err := r.getNatssyncClientVersion(natsSyncClientAboutURL)
	if err != nil {
		log.Error(err, "Failed to get the natssync-client version")
		return cachev1.NatssyncClientStatus{}, err
	}
	natssyncClientStatus.Registered = natssyncClientRegistrationStatus
	natssyncClientStatus.Version = natssyncClientVersion
	return natssyncClientStatus, nil
}

func (r *AstraAgentReconciler) getNatssyncClientRegistrationStatus(natsSyncClientRegisterURL string) (string, error) {
	resp, err := http.Get(natsSyncClientRegisterURL)
	if err != nil {
		return "Unknown", err
	}

	type registrationResponse struct {
		LocationID string `json:"locationID,omitempty"`
	}
	var registrationResp registrationResponse
	all, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(all, &registrationResp)
	if err != nil {
		return "Unknown", err
	}

	if registrationResp.LocationID == "" {
		return "False", nil
	}
	return "True", nil
}

func (r *AstraAgentReconciler) getNatssyncClientVersion(natsSyncClientAboutURL string) (string, error) {
	resp, err := http.Get(natsSyncClientAboutURL)
	if err != nil {
		return "", err
	}

	type aboutResponse struct {
		AppVersion string `json:"appVersion,omitempty"`
	}
	var aboutResp aboutResponse
	all, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(all, &aboutResp)
	if err != nil {
		return "", err
	}
	return aboutResp.AppVersion, nil
}
