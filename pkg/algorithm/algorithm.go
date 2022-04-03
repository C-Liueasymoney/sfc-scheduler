package algorithm

import (
	"SFC-Scheduler/pkg/common"
	pl "SFC-Scheduler/pkg/pod"
	"SFC-Scheduler/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"math"
)

func NodeSelector(pod *v1.Pod, nodeInfo *framework.NodeInfo) (int64, error) {
	klog.V(3).Infof("*****************算法开始运行**************")
	// 提取label信息
	sourceNode := utils.GetPodInfoFromLabel(pod, "sourceNode")
	targetNode := utils.GetPodInfoFromLabel(pod, "targetNode")
	sfcName := utils.GetPodInfoFromLabel(pod, "sfcName")
	totalSerString := utils.GetPodInfoFromLabel(pod, "totalService")
	totalSer := utils.StringToInt(totalSerString)
	lifeCycleString := utils.GetPodInfoFromLabel(pod, "lifeCycle")
	lifeCycle := utils.StringToInt(lifeCycleString)
	minBandwidth := utils.GetPodInfoFromLabel(pod, "minBandwidth")
	podBandwidth := utils.StringtoFloatBandwidth(minBandwidth)
	delayString := utils.GetPodInfoFromLabel(pod, "delay")
	delay := utils.StringToInt(delayString)

	appName := utils.GetPodInfoFromLabel(pod, "app")
	positionString := utils.GetPodInfoFromLabel(pod, "chainPosition")
	position := utils.StringToInt(positionString)

	nextApp := ""
	prevApp := ""
	var appList []string

	if position == 1 {
		nextApp = utils.GetPodInfoFromLabel(pod, "nextService")
		appList = []string{nextApp}
	} else if position == totalSer {
		prevApp = utils.GetPodInfoFromLabel(pod, "prevService")
		appList = []string{prevApp}
	} else {
		nextApp = utils.GetPodInfoFromLabel(pod, "nextService")
		prevApp = utils.GetPodInfoFromLabel(pod, "prevService")
		appList = []string{prevApp, nextApp}
	}
	reqCpuString := utils.GetPodInfoFromLabel(pod, "cpu")
	reqCpu := utils.StringToInt(reqCpuString)
	reqMemString := utils.GetPodInfoFromLabel(pod, "mem")
	reqMem := utils.StringToInt(reqMemString)

	policy := ""

	if delay > 50 {
		policy = "sensitive"
	} else {
		policy = "tolerate"
	}

	klog.V(3).Infof("-----------Get方法获取node信息测试-----------")
	nc := utils.GetNodeResource(nodeInfo.Node().GetName(), "cpu")
	if nc == "error" {
		klog.V(3).Infof("获取cpu错误")
	} else {
		ncpu := utils.StringToInt(nc)
		klog.V(3).Infof("get方法获取的cpu为：%v", ncpu)
	}
	nm := utils.GetNodeResource(nodeInfo.Node().GetName(), "memory")
	if nm == "error" {
		klog.V(3).Infof("获取mem错误")
	} else {
		nmem := utils.StringToInt(nm)
		klog.V(3).Infof("get方法获取的cpu为：%v", nmem)
	}
	klog.V(3).Infof("-----------原生方法获取node信息测试-----------")
	klog.V(3).Infof("原生获取的cpu为%v", nodeInfo.Allocatable.MilliCPU)
	klog.V(3).Infof("原生获取的mem为%v", nodeInfo.Allocatable.Memory)

	klog.V(3).Infof("appName is %v", appName)
	klog.V(3).Infof("targetNode is %v", targetNode)
	klog.V(3).Infof("minBandwidth is %v", minBandwidth)
	klog.V(3).Infof("policy is %v", policy)
	klog.V(3).Infof("podBandwidth is %v", podBandwidth)

	if policy == "tolerate" {
		klog.V(3).Infof("start the policy: %v", policy)
		// 这里考虑一下如果有相同delay的情况怎么解决
		delay, path := common.LatencyGraph.GetPath(nodeInfo.Node().Name, targetNode)

		klog.V(3).Infof("the node: %v 's delay is %v", nodeInfo.Node().Name, delay)
		klog.V(3).Infoln(path)

		if path != nil {
			return 10000 - int64(delay), nil
		}
		return 0, nil
	} else if policy == "sensitive" {
		klog.V(3).Infof("start the policy: %v", policy)
		podList := pl.CreatePodList(sfcName)

		klog.V(3).Infof("the current id is: %v", common.Id)
		// 在之前创建过的所有Pod里遍历
		for i := 1; i <= common.Id; i++ {
			// 在当前sfc总Pod数遍历
			for j := 1; j <= totalSer; j++ {
				// 在之前加入的当前pod前后pod遍历
				for _, app := range appList {
					// 跳过当前pod
					if j != position {
						key := utils.GetKey(common.Id, app, sfcName, j, totalSer)
						allocatedNode, ok := common.ServiceMap[key]
						if ok {
							klog.V(3).Infof("key: %v, has allocated on node: %v", key, allocatedNode)
							err := podList.AddPod("pod", "dijkstra", key, 0.0, allocatedNode)
							if err != nil {
								klog.V(3).Infof("err on addPod: %v", err)
							}
						}
					}
				}
			}
		}

		// 如果podlist不为空，代表此条SFC有之前已经部署好的pod作为参照，可根据此信息进行当前pod的部署
		if !podList.IsEmpty() {

			score, _ := calculateShortPath(nodeInfo, podList)
			return score, nil
		} else { // podList为空就说明此条SFC之前无部署好的pod，就使用location策略
			if targetNode != "Any" {
				delay, path := common.LatencyGraph.GetPath(nodeInfo.Node().Name, targetNode)
				if path != nil {
					return 10000 - int64(delay), nil
				}
			}
		}
	}

	klog.V(3).Infof("no meet policy")
	return 0, nil
}

func calculateShortPath(nodeInfo *framework.NodeInfo, podList *pl.PodList) (int64, map[string]float64) {
	klog.V(3).Infof("start calculate the short path in podList: %v", podList.Name)
	delayCost := make(map[string]float64)
	minCost := math.MaxFloat64

	podList.Start()
	// 遍历podList
	for podList.Current != nil {
		// 求当前node到PodList中已分配的pod所在节点的delay
		delay, _ := common.LatencyGraph.GetPath(nodeInfo.Node().Name, podList.Current.NodeAllocated)
		// 累加上当前node到podlist中所有节点的delay，也就是如果把pod分配到当前节点所增加的delay
		prevValue := utils.GetValue(delayCost, nodeInfo.Node().Name)
		delayCost[nodeInfo.Node().Name] = prevValue + float64(delay)
		minCost = delayCost[nodeInfo.Node().Name]
		podList.Next()
	}
	klog.V(3).Infof("the node: %v minDelay is %v", nodeInfo.Node().Name, minCost)

	return 10000 - int64(minCost), delayCost

}
