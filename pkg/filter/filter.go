package filter

import (
	"SFC-Scheduler/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// 检查pod所需带宽是否小于节点可用带宽，满足返回true，否则返回false
// 这里可以考虑把不满足带宽的边暂时移除
func BandwidthFilter(pod *v1.Pod, nodeInfo *framework.NodeInfo) bool {
	minBandwidth := utils.GetPodInfoFromLabel(pod, "minBandwidth")
	podBandwidth := utils.StringtoFloatBandwidth(minBandwidth)
	nodeBandwidth := utils.GetNodeBandwidth(nodeInfo, "avBandwidth")
	klog.V(3).Infof("the pod bandwidth is %v, the node bandwidth is %v", podBandwidth, nodeBandwidth)
	if podBandwidth <= nodeBandwidth {
		return true
	} else {
		return false
	}
}
