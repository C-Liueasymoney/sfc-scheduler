package sort

import (
	"SFC-Scheduler/pkg/utils"
	"k8s.io/klog/v2"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

func Less(podInfo1, podInfo2 *framework.QueuedPodInfo) bool {
	return getPodPriority(podInfo1) < getPodPriority(podInfo2)
}

// 首先如果选择了节点资源优先算法可以按照Pod所需资源降序排序
func getPodPriority(podInfo *framework.QueuedPodInfo) int {
	if pos := utils.GetPodInfoFromLabel(podInfo.Pod, "chainPosition"); pos != "Any" {
		klog.V(3).Infof("Pod name: %v, Position: %v", podInfo.Pod.Name, pos)
		p := utils.StringToInt(pos)
		return p
	}
	return 0
}
