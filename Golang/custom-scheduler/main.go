package main

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const PluginName = "CustomScheduler"

func main() {

	// 创建调度器命令
	command := app.NewSchedulerCommand(
		app.WithPlugin("CustomPlugin", New),
	)

	// 运行调度器
	if err := command.Execute(); err != nil {
		panic(err)
	}

}

// CustomCPUScheduler 是调度插件的结构体
type CustomScheduler struct {
	handle framework.Handle
}

// New 创建插件实例
func New(obj k8sruntime.Object, h framework.Handle) (framework.Plugin, error) {
	return &CustomScheduler{handle: h}, nil
}

// Name 返回插件的名字
func (c *CustomScheduler) Name() string {
	return PluginName
}

// PreFilter 实现预过滤逻辑 检查 Pod 是否指定了必须的标签。
func (p *CustomScheduler) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	requiredLabel := "app"
	if _, exists := pod.Labels[requiredLabel]; !exists {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Pod must have label %s", requiredLabel))
	}
	return framework.NewStatus(framework.Success)
}

// Filter 实现节点过滤逻辑
func (c *CustomScheduler) Filter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	// 获取节点的可用资源信息
	node := nodeInfo.Node()
	allocatable := node.Status.Allocatable

	// 获取节点 CPU 剩余资源（以 milliCPU 为单位）
	cpuRemaining := allocatable.Cpu().MilliValue()

	// 检查剩余 CPU 是否满足要求
	if cpuRemaining < 2000 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Node %s has insufficient CPU", node.Name))
	}

	return framework.NewStatus(framework.Success)
}

// PostFilter 实现后过滤逻辑 处理无可用节点的情况，尝试打印警告或执行其他逻辑。
func (p *CustomScheduler) PostFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, filteredNodes []*framework.NodeInfo) *framework.Status {
	if len(filteredNodes) == 0 {
		// 记录日志或触发其他动作
		fmt.Printf("No available nodes for Pod %s\n", pod.Name)
	}
	return framework.NewStatus(framework.Success)
}

// PreScore 实现预打分逻辑，检查 Pod 是否指定了必须的标签。
func (p *CustomScheduler) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := p.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Node %s not found", nodeName))
	}

	cpuAvailable := nodeInfo.Allocatable.MilliCPU - nodeInfo.Requested.MilliCPU
	return cpuAvailable, framework.NewStatus(framework.Success)
}

// Reserve 实现预订逻辑，为 Pod 预留资源。
func (p *CustomScheduler) Reserve(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	fmt.Printf("Reserving resources for Pod %s on Node %s\n", pod.Name, nodeName)
	return framework.NewStatus(framework.Success)
}

// Permit 实现许可逻辑，检查 Pod 是否需要等待管理员批准。
func (p *CustomScheduler) Permit(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (*framework.Status, time.Duration) {
	if pod.Annotations["require-approval"] == "true" {
		// 等待管理员批准
		return framework.NewStatus(framework.Wait), 10 * time.Second
	}
	return framework.NewStatus(framework.Success), 0
}

// PreBind 实现预绑定逻辑，检查 Pod 是否指定了必须的标签。
func (p *CustomScheduler) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	fmt.Printf("PreBind check for Pod %s on Node %s\n", pod.Name, nodeName)
	return framework.NewStatus(framework.Success)
}

// Bind 实现绑定逻辑，将 Pod 绑定到节点。
func (p *CustomScheduler) Bind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	binding := &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{Name: pod.Name, Namespace: pod.Namespace},
		Target: v1.ObjectReference{
			Kind: "Node",
			Name: nodeName,
		},
	}
	err := p.handle.ClientSet().CoreV1().Pods(pod.Namespace).Bind(ctx, binding, metav1.CreateOptions{})
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to bind Pod %s to Node %s", pod.Name, nodeName))
	}
	return framework.NewStatus(framework.Success)
}

// Unreserve 是插件的 Unreserve 方法
func (c *CustomScheduler) Unreserve(ctx context.Context, state *framework.CycleState, pod *framework.PodInfo, nodeName string) {
	fmt.Printf("Unreserving Pod: %s from Node: %s\n", pod.Pod.Name, nodeName)
	// 在这里添加你的清理逻辑
	// 这里可以加入你特定的业务逻辑，例如：
	// - 日志记录
	// - 发送通知
	// - 资源回滚操作
	// - 其他业务逻辑

}
