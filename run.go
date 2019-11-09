package main

import (
	"flag"
	"fmt"
	"github.com/fangfuhaomin/promethus-vspherewebapi-go/exsihost"
	"github.com/fangfuhaomin/promethus-vspherewebapi-go/vcconnect"
	"github.com/fangfuhaomin/promethus-vspherewebapi-go/vms"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)



func main() {
	var (
		// 用户
		user string
		// 密码
		password string
		// 主机名
		ip string
		//metrics port
		metricsport string
		)
	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量
	flag.StringVar(&user, "user", "", "用户名,默认为空")
	flag.StringVar(&password, "password", "", "密码,默认为空")
	flag.StringVar(&ip, "vcip", "127.0.0.1", "主机名,默认 127.0.0.1")
	flag.StringVar(&metricsport, "metricsport", "8080", "metrics端口,默认 8080")
	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()
	// 打印
	fmt.Printf("请使用格式  ./XXX   -user XX  -password XX  -vcip XX  -metricsport XX \n")
	fmt.Printf("当前参数 -user %s -password %s -vcip%s  -metricsport %s \n", user, password, ip, metricsport)
	c := vcconnect.Vccon(user,password,ip).Client
	fmt.Printf("web open url hostip:%s/metrics",metricsport)

	go func() {
		for {
			vms.GetVmsInfo(c)
			exsihost.GetExsiInfo(c)

			//time.Sleep(time.Second*30)
		}
	}()


	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+metricsport, nil))
}
