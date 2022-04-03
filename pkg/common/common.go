package common

import (
	"SFC-Scheduler/pkg/graph"
	"SFC-Scheduler/pkg/pod"
	cl "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type SchedulerInterface struct {
	Clientset *cl.Clientset
}

var (
	// LatencyGraph 权重为时延的无向图
	LatencyGraph = graph.NewGraph()
	// Id id表示当前service的唯一标志，每个service不同
	Id = 0
	// ServiceMap 表示每个service对应放置的Node，key是每个service的唯一key，value是该service所放置的node
	ServiceMap = make(map[string]string)

	// Scheduler client-go的接口
	Scheduler SchedulerInterface = newScheduler()

	// AllocatedPods 已经分配结束的PodList
	AllocatedPods = pod.CreatePodList("scheduledPods")
)

func init() {
	LatencyGraph.AddEdge("node1", "node2", 10)
	klog.V(3).Infof("graph init complete: ")
	klog.V(3).Infoln(LatencyGraph)
}

func newScheduler() SchedulerInterface {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.V(3).Infof("fail to get config")
	}
	clientSet, err := cl.NewForConfig(config)
	if err != nil {
		klog.V(3).Infof("fail to get clientSet")
	}

	scheduler := SchedulerInterface{
		Clientset: clientSet,
	}
	return scheduler
}
