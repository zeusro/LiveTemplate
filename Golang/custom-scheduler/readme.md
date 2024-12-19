Kubernetes 调度器插件通过配置 Scheduling Framework 接口来扩展功能。调度插件可以分为以下几类：
    1.    Filter: 用于过滤不符合条件的节点。
    2.    Score: 为符合条件的节点打分。
    3.    PreFilter: 在节点筛选之前执行预处理逻辑。
    4.    PostFilter: 在节点过滤之后执行。
    5.    Reserve: 临时“预留”节点资源。
    6.    Permit: 控制 Pod 是否可以被调度。
    7.    PreBind: 在 Pod 被绑定到节点之前执行。
    8.    Bind: 自定义 Pod 绑定行为。


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