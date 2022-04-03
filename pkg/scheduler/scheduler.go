package scheduler

import (
	"SFC-Scheduler/pkg/algorithm"
	"SFC-Scheduler/pkg/common"
	"SFC-Scheduler/pkg/filter"
	"SFC-Scheduler/pkg/sort"
	"SFC-Scheduler/pkg/utils"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"strconv"
)

const (
	Name = "sfc-scheduler"
)

type Sfc struct {
	handle framework.FrameworkHandle
}

//var _ framework.PreFilterPlugin = &Sfc{}
var _ framework.FilterPlugin = &Sfc{}

var _ framework.QueueSortPlugin = &Sfc{}
var _ framework.ScorePlugin = &Sfc{}
var _ framework.ScoreExtensions = &Sfc{}
var _ framework.PreBindPlugin = &Sfc{}

func (s *Sfc) Name() string {
	return Name
}

func (s *Sfc) Less(podInfo1, podInfo2 *framework.QueuedPodInfo) bool {
	klog.V(3).Infof("================== Sort Start ================")
	return sort.Less(podInfo1, podInfo2)
}

func (s *Sfc) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	klog.V(3).Infof("================== Filter Start Node: %v ================", nodeInfo.Node().Name)

	if !filter.BandwidthFilter(pod, nodeInfo) {
		return framework.NewStatus(framework.Unschedulable, "Node: %v doesn't have enough bandwidth meet pod: %v", nodeInfo.Node().Name, pod.Name)
	}
	return framework.NewStatus(framework.Success, "")
}

func (s *Sfc) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(3).Infof("================== Score Start Node: %v ================", nodeName)
	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err == nil {
		score, _ := algorithm.NodeSelector(p, nodeInfo)
		return score, framework.NewStatus(framework.Success, fmt.Sprintf("node: %v, score: %v", nodeInfo.Node().Name, score))
	}
	return 0, framework.NewStatus(framework.Error, fmt.Sprintf("get nodeInfo fail"))
}

func (s *Sfc) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	klog.V(3).Infof("================== NormalizeScore Start ================")
	var heighest int64 = 0
	lowest := scores[0].Score

	for _, nodeScore := range scores {
		if nodeScore.Score > heighest {
			heighest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	if heighest == lowest {
		lowest--
	}

	for i, nodeScore := range scores {
		scores[i].Score = (nodeScore.Score - lowest) * framework.MaxNodeScore / (heighest - lowest)
		klog.V(3).Infof("shcedule pod: %v, the normalized score in node: %v is %v", p.Name, scores[i].Name, scores[i].Score)
	}
	return framework.NewStatus(framework.Success, "")
}

func (s *Sfc) ScoreExtensions() framework.ScoreExtensions {
	return s
}

// preBind只会在一个node上运行

func (s *Sfc) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("================== PreBind Start ================")
	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err == nil {
		klog.V(3).Infof("bind node: %v", nodeInfo.Node().Name)

		appName := utils.GetPodInfoFromLabel(pod, "app")
		sfcName := utils.GetPodInfoFromLabel(pod, "sfcName")
		positionString := utils.GetPodInfoFromLabel(pod, "chainPosition")
		totalSerString := utils.GetPodInfoFromLabel(pod, "totalService")
		position := utils.StringToInt(positionString)
		totalSer := utils.StringToInt(totalSerString)
		// 将pod产生的key与对应的node加入servicemap中
		common.Id++
		key := utils.GetKey(common.Id, appName, sfcName, position, totalSer)
		utils.AddService(key, nodeName)

		klog.V(3).Infof("now the serviceMap include: %v", common.ServiceMap)

		minBandwidth := utils.GetPodInfoFromLabel(pod, "minBandwidth")
		podBandwidth := utils.StringtoFloatBandwidth(minBandwidth)
		// 更新节点带宽
		nodeBandwidth := utils.GetNodeBandwidth(nodeInfo, "avBandwidth")
		value := nodeBandwidth - podBandwidth
		label := strconv.FormatFloat(value, 'f', 2, 64)
		if err := utils.UpdateBandwidthLabel(ctx, label, nodeInfo); err != nil {
			klog.V(3).Infof("find error when updating bandwidth label: %V", err)
		}

		// 在podList中添加此Pod
		if err = common.AllocatedPods.AddPod(pod.Name, pod.Namespace, key, podBandwidth, nodeName); err != nil {
			klog.V(3).Infof("find error when add pod to list: %v", err)
		} else {
			klog.V(3).Infof("the current podList is:")
			err = common.AllocatedPods.ShowAllPods()
			if err != nil {
				klog.V(3).Infof("fail to show podList")
			}
		}

		// 如果一条SFC都部署完毕，发送部署完毕的信息

		return framework.NewStatus(framework.Success, fmt.Sprintf("success allocate, and nodeBandwidth updated!!!"))
	}
	return framework.NewStatus(framework.Error, "get nodeInfo fail")
}

func New(configuration runtime.Object, h framework.FrameworkHandle) (framework.Plugin, error) {
	return &Sfc{
		handle: h,
	}, nil
}
