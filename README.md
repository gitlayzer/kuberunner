![kuberunner](https://github.com/gitlayzer/kuberunner-dashboard/assets/77761224/d94d0f6d-aa64-4179-956b-2dcd5544b9c1)

|    名称    |   版本   |    发布时间    |  发布人   | 是否更新UI |
| :--------: |:------:|:----------:| :-------: |:------:|
| kuberunner | v0.0.1 | 2023/03/10 | Gitlayzer |   ×    |
| kuberunner | v0.0.2 | 2023/03/13 | Gitlayzer |   ×    |
| kuberunner | v0.0.3 | 2023/03/13 | Gitlayzer |   √    |
| kuberunner | v0.0.4 | 2023/06/03 | Gitlayzer |   √    |

```shell
本次更新主要针对的路由注册方法，重构了路由注册的思路，通过模块化的形式进行路由注册，可具体到某个资源的所有操作
```
```yaml
# API监听的端口
ListenAddress: ":80"
# 终端服务监听的端口
WsListenAddress: ":81"
# 前端登录验证的账号
Username: "admin"
# 前端登录验证的密码
Password: "admin"
# 可配置多个kubeconfig实现多集群管理
KubeConfigs:
  cluster1: "C:\\Users\\Administrator\\.kube\\config"
```

