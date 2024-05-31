# community
去中心化社区

## local 
本机节点服务，

udp打洞，udp/tcp数据传输,http服务

安装后异步启动 打洞，数据接受和同步，http web服务

## burrow
打洞服务，常部署在公网


## 数据包定义

第一个字节: 
- 0 : [来自burrow] 节点表
- 1 : [来自节点] 数据包
- 2 : [来自burrow] 返回节点自己的地址

