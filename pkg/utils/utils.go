package utils

import (
	"SFC-Scheduler/pkg/common"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func GetPodInfoFromLabel(pod *v1.Pod, label string) string {
	value, ok := pod.Labels[label]
	if ok {
		value = pod.Labels[label]
		return value
	}
	return "Any"
}

func GetNodeBandwidth(nodeInfo *framework.NodeInfo, avBandwidth string) float64 {
	if nodeBandwidth, ok := nodeInfo.Node().Labels[avBandwidth]; ok {
		nodeBandwidth, err := strconv.ParseFloat(nodeBandwidth, 64)
		if err == nil {
			return nodeBandwidth
		}
	}
	return math.MaxFloat64
}

func StringtoFloatBandwidth(minBandwidth string) float64 {
	bandwidth, err := strconv.ParseFloat(minBandwidth, 64)
	if err == nil {
		return bandwidth
	}
	return 0.250 // Default Value: 250 Kbit/s
}

func StringToInt(value string) int {
	if newValue, err := strconv.Atoi(value); err == nil {
		return newValue
	}
	return 1
}

func GetKey(id int, appName string, sfcName string, position int, totalSer int) string {
	return strconv.Itoa(id) + "-" + appName + "-" + sfcName + "-" + strconv.Itoa(position) + "-" + strconv.Itoa(totalSer)
}

func GetValue(shortPathCost map[string]float64, key string) float64 {
	return shortPathCost[key]
}

func AddService(key string, nodeName string) {
	common.ServiceMap[key] = nodeName
	klog.V(3).Infof("ServiceMap add item: key: %v; value: %v", key, common.ServiceMap[key])
}

func UpdateBandwidthLabel(ctx context.Context, label string, nodeInfo *framework.NodeInfo) error {
	nodeLabels := nodeInfo.Node().GetLabels()
	prevLabel := nodeLabels["avBandwidth"]
	nodeLabels["avBandwidth"] = label

	nodeInfo.Node().SetLabels(nodeLabels)

	klog.V(3).Infof("Node: %v updating bandwidth label: previous is %v, now is %v", nodeInfo.Node().Name, prevLabel, label)

	if _, err := common.Scheduler.Clientset.CoreV1().Nodes().Update(ctx, nodeInfo.Node(), metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("failed to update label!!!,reason is : %v", err)
	}
	return nil
}

func GetNodeResource(nodeName string, resourceType string) string {
	resp, err := http.Get("http://192.168.27.2:8201//service/deployment/getNodeResource/" + nodeName + "/" + resourceType)
	if err != nil {
		return "error"
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func PostInfo(status string) {
	b := status
	resp, _ :=
		http.Post("http://192.168.27.2:8201//service/deployment/receive",
			"application/x-www-form-urlencoded", strings.NewReader("heel="+string(b)))

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}
