
指定使用某种调度器来处理资源有2种方式：

## 使用现有调度器插件

你可以通过编写调度器扩展（例如 scheduler extender）来为你的自定义资源实现调度策略。调度器扩展允许你将外部系统集成到 Kubernetes 的调度过程中。

调度器扩展是一种 HTTP 服务，Kubernetes 调度器会向该服务发送请求，以查询是否可以调度特定的 Pod 或资源。
	1.	编写调度器扩展服务：
	•	它可以是一个 HTTP 服务，接收调度请求并根据你的自定义逻辑返回调度结果。
	•	例如，可以根据自定义资源的字段来决定如何调度 Pods。
	2.	配置 Kubernetes 调度器使用扩展服务：
	•	通过配置 --scheduler-extender-config 参数，指定自定义的调度器扩展服务。


```json
{
  "apiVersion": "v1",
  "kind": "Extender",
  "urlPrefix": "http://<extender-service-url>/scheduler/extend",
  "filterVerb": "filter",
  "prioritizeVerb": "prioritize",
  "weight": 1,
  "nodeCacheCapable": false
}
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-scheduler
  namespace: kube-system
spec:
  template:
    spec:
      containers:
      - name: kube-scheduler
        image: k8s.gcr.io/kube-scheduler:v1.27.0
        command:
        - /bin/kube-scheduler
        - --config=/etc/kubernetes/scheduler/scheduler-config.yaml
        - --scheduler-extender-config=/etc/kubernetes/scheduler/extender-config.json
        volumeMounts:
        - name: scheduler-config
          mountPath: /etc/kubernetes/scheduler
      volumes:
      - name: scheduler-config
        configMap:
          name: scheduler-config
```          

## 创建自定义调度器
Kubernetes 调度器插件通过配置 Scheduling Framework 接口来扩展功能。调度插件可以分为以下几类：
    1.    Filter: 用于过滤不符合条件的节点。
    2.    Score: 为符合条件的节点打分。
    3.    PreFilter: 在节点筛选之前执行预处理逻辑。
    4.    PostFilter: 在节点过滤之后执行。
    5.    Reserve: 临时“预留”节点资源。
    6.    Permit: 控制 Pod 是否可以被调度。
    7.    PreBind: 在 Pod 被绑定到节点之前执行。
    8.    Bind: 自定义 Pod 绑定行为。

    The scheduling cycle selects a node for the Pod, and the binding cycle applies that decision to the cluster. Together, a scheduling cycle and binding cycle are referred to as a "scheduling context".



实现对应逻辑之后，需要
1. 将插件注册到调度框架中，通过 frameworkruntime.Registry 实现：
1. 在 Kubernetes 中通过配置 scheduler-config.yaml 加载自定义插件。
    ```yaml
    apiVersion: kubescheduler.config.k8s.io/v1beta3
    kind: KubeSchedulerConfiguration
    profiles:
      - schedulerName: default-scheduler
        plugins:
          filter:
            enabled:
              - name: CustomScheduler
    ```
1. 将自定义调度器部署为一个 Pod，使用配置文件启动调度器服务。

有些版本注册的方式可能略有不同：
```go
// 自定义插件注册表
	registry := runtime.Registry{
		customplugin.Name: customplugin.New, // 注册插件
	}

	// 创建调度器命令
	command := app.NewSchedulerCommand(
		app.WithPlugin(registry), // 使用注册表加载插件
	)

	// 运行调度器
	if err := command.Execute(); err != nil {
		panic(err)
	}
```

或者直接自定义调度策略：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  schedulerName: my-custom-scheduler
  containers:
  - name: mycontainer
    image: myimage
```


https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/