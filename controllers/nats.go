package controllers

import (
	"context"
	"fmt"

	cachev1 "github.com/NetApp/astraagent-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// StatefulsetForNats returns a NATS Statefulset object
func (r *AstraAgentReconciler) StatefulsetForNats(m *cachev1.AstraAgent, ctx context.Context) (*appsv1.StatefulSet, error) {
	log := ctrllog.FromContext(ctx)
	ls := labelsForNats(NatsName)

	var replicas int32
	if m.Spec.Nats.Size > 2 {
		replicas = m.Spec.Nats.Size
	} else {
		log.Info("Defaulting the Nats replica size", "size", NatsDefaultSize)
		replicas = NatsDefaultSize
	}

	var natsImage string
	if m.Spec.Nats.Image != "" {
		natsImage = m.Spec.Nats.Image
	} else {
		log.Info("Defaulting the Nats image", "image", NatsDefaultImage)
		natsImage = NatsDefaultImage
	}

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NatsName,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: NatsClusterServiceName,
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: NatsServiceAccountName,
					Containers: []corev1.Container{{
						Image: natsImage,
						Name:  NatsName,
						Ports: []corev1.ContainerPort{
							{
								Name:          "client",
								ContainerPort: NatsClientPort,
							},
							{
								Name:          "cluster",
								ContainerPort: NatsClusterPort,
							},
							{
								Name:          "monitor",
								ContainerPort: NatsMonitorPort,
							},
							{
								Name:          "metrics",
								ContainerPort: NatsMetricsPort,
							},
						},
						Command: []string{"nats-server", "--config", "/etc/nats-config/nats.conf"},
						Env: []corev1.EnvVar{
							{
								Name:  "CLUSTER_ADVERTISE",
								Value: fmt.Sprintf("%s.nats.%s.svc", NatsName, m.Namespace),
							},
							{
								Name:  "POD_NAME",
								Value: NatsName,
							}, {
								Name:  "POD_NAMESPACE",
								Value: m.Namespace,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      NatsVolumeName,
								MountPath: "/etc/nats-config",
							},
							{
								Name:      "pid",
								MountPath: "/var/run/nats",
							},
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: NatsVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: NatsConfigMapName,
									},
								},
							},
						},
						{
							Name: "pid",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium: "",
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

// ClusterServiceForNats returns a cluster Service object for Nats
func (r *AstraAgentReconciler) ClusterServiceForNats(m *cachev1.AstraAgent) (*corev1.Service, error) {
	ls := labelsForNats(NatsName)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NatsClusterServiceName,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "",
			Ports: []corev1.ServicePort{
				{
					Name: "client",
					Port: NatsClientPort,
				},
				{
					Name: "cluster",
					Port: NatsClusterPort,
				},
				{
					Name: "monitor",
					Port: NatsMonitorPort,
				},
				{
					Name: "metrics",
					Port: NatsMetricsPort,
				},
				{
					Name: "gateways",
					Port: NatsGatewaysPort,
				},
			},
			Selector: ls,
		},
	}
	// Set astraAgent instance as the owner and controller
	err := ctrl.SetControllerReference(m, service, r.Scheme)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// ConfigMapForNats returns a ConfigMap object for nats
func (r *AstraAgentReconciler) ConfigMapForNats(m *cachev1.AstraAgent) (*corev1.ConfigMap, error) {
	natsConf := "pid_file: \"/var/run/nats/nats.pid\"\nhttp: %d\n\ncluster {\n  port: %d\n  routes [\n    nats://nats-0.nats-cluster:%d\n    nats://nats-1.nats-cluster:%d\n    nats://nats-2.nats-cluster:%d\n  ]\n\n  cluster_advertise: $CLUSTER_ADVERTISE\n  connect_retries: 30\n}\n"
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Namespace,
			Name:      NatsConfigMapName,
		},
		Data: map[string]string{
			"nats.conf": fmt.Sprintf(natsConf, NatsMonitorPort, NatsClusterPort, NatsClusterPort, NatsClusterPort, NatsClusterPort),
		},
	}
	err := ctrl.SetControllerReference(m, configMap, r.Scheme)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// ServiceAccountForNats returns a ServiceAccount object for nats
func (r *AstraAgentReconciler) ServiceAccountForNats(m *cachev1.AstraAgent) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NatsServiceAccountName,
			Namespace: m.Namespace,
			Labels:    labelsForNats(NatsName),
		},
	}
	err := ctrl.SetControllerReference(m, sa, r.Scheme)
	if err != nil {
		return nil, err
	}
	return sa, nil
}

// ServiceForNats returns a Service object for nats
func (r *AstraAgentReconciler) ServiceForNats(m *cachev1.AstraAgent) (*corev1.Service, error) {
	ls := labelsForNats(NatsName)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NatsName,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: NatsName,
					Port: NatsClientPort,
				},
			},
			Selector: ls,
		},
	}
	// Set astraAgent instance as the owner and controller
	err := ctrl.SetControllerReference(m, service, r.Scheme)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// labelsForNats returns the labels for selecting the nats resources
// belonging to the given astraAgent CR name.
func labelsForNats(name string) map[string]string {
	return map[string]string{"app": name}
}

// GetNatsURL returns the nats URL
func (r *AstraAgentReconciler) GetNatsURL(m *cachev1.AstraAgent) string {
	natsURL := fmt.Sprintf("nats://%s.%s:%d", NatsName, m.Namespace, NatsClientPort)
	return natsURL
}
