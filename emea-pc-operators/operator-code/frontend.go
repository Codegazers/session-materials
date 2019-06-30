package demosite

import (
	"strconv"

	demov1alpha1 "github.com/jdob/demo-operator/pkg/apis/demo/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const frontendPort = 3000
const frontendServicePort = 30686

func (r *ReconcileDemoSite) frontendDeployment(v *demov1alpha1.DemoSite) *appsv1.Deployment {
	labels := labels(v, "frontend")
	size := int32(1)
	host := v.Spec.MinikubeIP

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-frontend",
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "jdob/visitors-webui:latest",
						Name:  "visitors-webui",
						Ports: []corev1.ContainerPort{{
							ContainerPort: frontendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "REACT_APP_SERVICE_HOST",
								Value: host,
							},
							{
								Name:  "REACT_APP_SERVICE_PORT",
								Value: strconv.Itoa(backendServicePort),
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcileDemoSite) frontendService(v *demov1alpha1.DemoSite) *corev1.Service {
	labels := labels(v, "frontend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-frontend-service",
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       frontendPort,
				TargetPort: intstr.FromInt(frontendPort),
				NodePort:   frontendServicePort,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	log.Info("Service Spec", "Service.Name", s.ObjectMeta.Name)

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}
