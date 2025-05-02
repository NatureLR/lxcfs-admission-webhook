package v1

import (
	"encoding/json"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

// -v /var/lib/lxcfs/proc/cpuinfo:/proc/cpuinfo:ro
// -v /var/lib/lxcfs/proc/diskstats:/proc/diskstats:ro
// -v /var/lib/lxcfs/proc/meminfo:/proc/meminfo:ro
// -v /var/lib/lxcfs/proc/stat:/proc/stat:ro
// -v /var/lib/lxcfs/proc/swaps:/proc/swaps:ro
// -v /var/lib/lxcfs/proc/uptime:/proc/uptime:ro
// -v /var/lib/lxcfs/proc/loadavg:/proc/loadavg:ro
func LxcPatch(pod *corev1.Pod) {
	// 节点:容器
	mountMap := map[string]string{
		"/var/lib/lxcfs/proc/cpuinfo":   "/proc/cpuinfo",
		"/var/lib/lxcfs/proc/diskstats": "/proc/diskstats",
		"/var/lib/lxcfs/proc/meminfo":   "/proc/meminfo",
		"/var/lib/lxcfs/proc/stat":      "/proc/stat",
		"/var/lib/lxcfs/proc/swaps":     "/proc/swaps",
		"/var/lib/lxcfs/proc/uptime":    "/proc/uptime",
		"/var/lib/lxcfs/proc/loadavg":   "/proc/loadavg",
	}

	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}
	for hostPath, containerPath := range mountMap {
		name := "lfxfs-" + filepath.Base(hostPath)

		// 创建 VolumeMount
		volumes = append(volumes, corev1.Volume{
			Name:         name,
			VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: hostPath}}},
		)
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      name,
			MountPath: containerPath,
			ReadOnly:  true,
		})
	}

	pod.Spec.Volumes = append(pod.Spec.Volumes, volumes...)
	var cs []corev1.Container
	for _, container := range pod.Spec.Containers {
		container.VolumeMounts = append(container.VolumeMounts, volumeMounts...)
		cs = append(cs, container)
	}
	pod.Spec.Containers = cs
	json.NewEncoder(os.Stdout).Encode(pod)
}
