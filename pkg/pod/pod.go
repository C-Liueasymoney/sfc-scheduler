package pod

import (
	"k8s.io/klog/v2"
)

type Pod struct {
	name          string
	namespace     string
	key           string
	minBandwidth  float64
	NodeAllocated string
	next          *Pod
}

type PodList struct {
	Name    string
	Head    *Pod
	Current *Pod
}

func CreatePodList(name string) *PodList {
	return &PodList{
		Name: name,
	}
}

func (p *PodList) AddPod(name string, namespace string, key string, minBandwidth float64, nodeAllocated string) error {
	s := &Pod{
		name:          name,
		namespace:     namespace,
		minBandwidth:  minBandwidth,
		key:           key,
		NodeAllocated: nodeAllocated,
	}
	if p.Head == nil {
		p.Head = s
	} else {
		currentPod := p.Head
		for currentPod.next != nil {
			currentPod = currentPod.next
		}
		currentPod.next = s
	}
	p.Current = s
	return nil
}

// removePod() removes a particular element from the list
func (p *PodList) RemovePod(name string) *PodList {

	// List is empty, cannot remove Pod
	if p.IsEmpty() {
		klog.V(3).Infof("list is empty")
	}

	// auxiliary variable
	temp := p.Head

	// If head holds the Pod to be deleted
	if temp != nil && temp.name == name {
		p.Head = temp.next // Changed head
		return p
	}

	// Search for the Pod to be deleted, keep track of the
	// previous Pod as we need to change temp.next
	for temp != nil {
		if temp.next.name == name {
			temp.next = temp.next.next
			return p
		}
		temp = temp.next
	}
	return p
}

// showAllPods() prints all elements on the list
func (p *PodList) ShowAllPods() error {
	currentPod := p.Head
	if currentPod == nil {
		klog.V(3).Infof("PodList is empty.")
		return nil
	}
	klog.V(3).Infof("%v \n", currentPod.name)
	for currentPod.next != nil {
		currentPod = currentPod.next
		klog.V(3).Infof("%v \n", currentPod.name)
	}

	return nil
}

//start() returns the first/head element
func (p *PodList) Start() *Pod {
	p.Current = p.Head
	return p.Current
}

//next() returns the next element on the list
func (p *PodList) Next() *Pod {
	p.Current = p.Current.next
	return p.Current
}

// IsEmpty() returns true if the list is empty
func (p *PodList) IsEmpty() bool {
	if p.Head == nil {
		return true
	}
	return false
}

// getSize() returns the linked list size
func (p *PodList) GetSize() int {
	size := 1
	last := p.Head
	for {
		if last == nil || last.next == nil {
			break
		}
		last = last.next
		size++
	}
	return size
}

//func main() {
//	podlist := createPodList("test")
//	podlist.addPod("pod1", "default", "1", 10, "node1")
//	podlist.addPod("pod2", "default", "2", 10, "node3")
//	podlist.addPod("pod3", "default", "3", 10, "node2")
//	fmt.Println(podlist.getSize())
//	fmt.Println(podlist.showAllPods())
//}
