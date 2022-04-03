package plugins

//
//import (
//	//"fmt"
//	"context"
//	v1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/runtime"
//	"k8s.io/klog/v2"
//	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
//)
//
//// 插件名称
//const Name = "sfc-plugin"
//
//type Args struct {
//	FavoriteColor  string `json:"favorite_color,omitempty"`
//	FavoriteNumber int    `json:"favorite_number,omitempty"`
//	ThanksTo       string `json:"thanks_to,omitempty"`
//}
//
//type Sample struct {
//	args   *Args
//	handle framework.FrameworkHandle
//}
//
//func (s *Sample) Name() string {
//	return Name
//}
//
//func (s *Sample) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) *framework.Status {
//	klog.V(3).Infof("I'm LC, prefilter pod: %v", pod.Name)
//	return framework.NewStatus(framework.Success, "")
//}
//
////func (s *Sample) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, node *framework.NodeInfo) *framework.Status {
////	klog.V(3).Infof("I'm LC, filter pod: %v, node: %v", pod.Name, node.Node().Name)
////	klog.V(3).Infof("Node information: alloCPU: %v, alloMem: %v", node.Allocatable.MilliCPU, node.Allocatable.Memory)
////	return framework.NewStatus(framework.Success, "")
////}
////
////func (s *Sample) PostFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, filteredNodeStatusMap framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status) {
////	klog.V(3).Infof("I'm LC, postfilter pod: %v", pod.Name)
////	return &framework.PostFilterResult{}, framework.NewStatus(framework.Success)
////}
//
//func (s *Sample) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
//	klog.V(3).Infof("I'm LC, prebind node info: %+v", nodeName)
//	return framework.NewStatus(framework.Success, "")
//}
//
////type PluginFactory = func(configuration *runtime.Unknown, f FrameworkHandle) (Plugin, error)
//
//func New(configuration runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
//	args := &Args{}
//	klog.V(3).Infof("get plugin config args: %+v", args)
//	return &Sample{
//		args:   args,
//		handle: f,
//	}, nil
//}
