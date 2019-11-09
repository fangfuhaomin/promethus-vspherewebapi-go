# promethus-vspherewebapi-go

当前需求有100多台EXSI 和 多台VCenter，写了一个通用的prometheus vmware go客户端
之前使用python写过，因为python有依赖包的问题，所以又写了一个go，使用起来方便。

#主要监控指标，
EXSI：
    cpu使用
    内存最大
    内存剩余
    磁盘剩余
    磁盘总值



虚拟机：
    cpu使用
    内存总值
    内存剩余
    磁盘剩余
    磁盘总值
    

#使用方法，

./govsphere -user XX  -password XX  -vcip XX  -metricsport XX

user:EXSI或者Vcenter的用户名
password：EXSI或者Vcenter的密码
vcip：EXSI或者Vcenter的IP
metricsport：在你电脑上开启的端口

#打开方式
http://你电脑的ip:metricsport