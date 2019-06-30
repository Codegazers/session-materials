package demosite

import (
	demov1alpha1 "github.com/jdob/demo-operator/pkg/apis/demo/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const backendPort = 8000
const backendServicePort = 30685

func (r *ReconcileDemoSite) backendDeployment(v *demov1alpha1.DemoSite) *appsv1.Deployment {
	labels := labels(v, "backend")
	size := v.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-backend",
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
						Image:           "jdob/visitors-service:latest",
						ImagePullPolicy: corev1.PullAlways,
						Name:            "visitors-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: backendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_DATABASE",
								Value: "visitors",
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

func (r *ReconcileDemoSite) backendService(v *demov1alpha1.DemoSite) *corev1.Service {
	labels := labels(v, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-backend-service",
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
				NodePort:   30685,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}
