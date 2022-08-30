// Copyright 2022 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/imagehelpers"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/rand"
)

func CreateDebugDeploymentForImage(imagename string, namespace string) (*appsv1.Deployment, error) {
	result, err := imagehelpers.GetImageEntryPoint(imagename)
	if err != nil {
		return nil, err
	}
	id := rand.RandStringRunes(10)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "debug-container-",
			Namespace:    namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"instance": id,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"instance": id,
					},
				},
				Spec: v1.PodSpec{
					InitContainers: []v1.Container{
						{
							Name:            "copy-entrypoint",
							Image:           "ko.local/entrypoint",
							ImagePullPolicy: v1.PullIfNotPresent,
							Command: []string{
								"cp",
								"/ko-app/entrypoint",
								"/.containerdbg",
							},
							VolumeMounts: []v1.VolumeMount{
								{
									MountPath: "/.containerdbg/",
									Name:      "shareddir",
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            "modified-pod",
							Image:           imagename,
							ImagePullPolicy: v1.PullIfNotPresent,
							Command:         append([]string{"/.containerdbg/entrypoint"}, result...),
							Env: []v1.EnvVar{
								{
									Name:  "SHARED_DIRECTORY",
									Value: "/var/run/containerdbg/daemon",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									MountPath: "/.containerdbg/",
									Name:      "shareddir",
								},
								{
									MountPath: "/var/run/containerdbg/daemon/",
									Name:      "socket-folder",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "shareddir",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "socket-folder",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/var/run/containerdbg/daemon",
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyAlways,
				},
			},
		},
	}, nil
}

func modifyContainer(container *v1.Container) error {
	container.Env = append(container.Env,
		v1.EnvVar{
			Name:  "SHARED_DIRECTORY",
			Value: "/var/run/containerdbg/daemon",
		})
	container.VolumeMounts = append(container.VolumeMounts,
		v1.VolumeMount{
			MountPath: "/.containerdbg/",
			Name:      "shareddir",
		},
		v1.VolumeMount{
			MountPath: "/var/run/containerdbg/daemon/",
			Name:      "socket-folder",
		})

	if container.Command != nil {
		container.Command = append([]string{"/.containerdbg/entrypoint"}, container.Command...)
		return nil
	}
	imagename := container.Image
	result, err := imagehelpers.GetImageEntryPoint(imagename)
	if err != nil {
		return err
	}
	container.Command = append([]string{"/.containerdbg/entrypoint"}, result...)
	return nil
}

func modifyPodSpec(podspec *v1.PodSpec) error {
	podspec.InitContainers = append(podspec.InitContainers, v1.Container{
		Name:            "copy-entrypoint",
		Image:           "ko.local/entrypoint",
		ImagePullPolicy: v1.PullIfNotPresent,
		Command: []string{
			"cp",
			"/ko-app/entrypoint",
			"/.containerdbg",
		},
		VolumeMounts: []v1.VolumeMount{
			{
				MountPath: "/.containerdbg/",
				Name:      "shareddir",
			},
		},
	})
	podspec.Volumes = append(podspec.Volumes,
		v1.Volume{
			Name: "shareddir",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
		v1.Volume{
			Name: "socket-folder",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/var/run/containerdbg/daemon",
				},
			},
		},
	)

	for i := range podspec.Containers {
		if err := modifyContainer(&podspec.Containers[i]); err != nil {
			return err
		}
	}

	return nil
}

func MutateDeployment(dep *appsv1.Deployment) error {
	return modifyPodSpec(&dep.Spec.Template.Spec)
}
