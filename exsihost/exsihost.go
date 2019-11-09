package exsihost

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"log"
)

var (

	exsiCpuCapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exsi_cpu_capacity",
		Help: "cpu Capacity MHZ",
	}, []string{"exsi"},
	)
	exsiCpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exsi_cpu_usage",
		Help: "exsi cpu usage MHZ",
	}, []string{"exsi"},
	)
	exsiMemCapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exsi_mem_capacity",
		Help: "exsi 内存最大 G",
	}, []string{"exsi"},
	)
	exsiMemUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exsi_mem_usage",
		Help: "exsi 内存使用 G",
	}, []string{"exsi"},
	)
	exsidsCapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exsi_disk_capacity",
		Help: "exsi 存储最大值 G",
	}, []string{"exsi","exsidisk"},
	)
	exsidsFreeSpace = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "exsi_disk_FreeSpace",
		Help: "exsi 存储未使用 G",
	}, []string{"exsi","exsidisk"},
	)
)


func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(exsiCpuCapacity)
	prometheus.MustRegister(exsiCpuUsage)
	prometheus.MustRegister(exsiMemCapacity)
	prometheus.MustRegister(exsiMemUsage)
	prometheus.MustRegister(exsidsCapacity)
	prometheus.MustRegister(exsidsFreeSpace)
}

func GetExsiInfo(c *vim25.Client) {
	//c := vcconnect.Vccon().Client
	ctx := context.Background()
	m := view.NewManager(c)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Destroy(ctx)
	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary", "vm" , "datastore"}, &hss)
	if err != nil {
		log.Fatal(err)
	}

	for _,hs := range hss{
		var exsiHostName = hs.Summary.Config.Name
		//var cpuCapacity = float64(int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores))
		//var cpuUsage = float64(hs.Summary.QuickStats.OverallCpuUsage)
		//var memCapacity = float64(hs.Summary.Hardware.MemorySize / (1024 * 1024))
		//var memUsage = float64(hs.Summary.QuickStats.OverallMemoryUsage)

		//exsi最大cpu
		exsiCpuCapacity.With(prometheus.Labels{"exsi":exsiHostName}).Set(float64(int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)))


		//exsi当前使用cpu
		exsiCpuUsage.With(prometheus.Labels{"exsi":exsiHostName}).Set(float64(hs.Summary.QuickStats.OverallCpuUsage))


		//exsi最大内存
		exsiMemCapacity.With(prometheus.Labels{"exsi":exsiHostName}).Set(float64(hs.Summary.Hardware.MemorySize / (1024 * 1024)))


		//exsi当前使用内存
		exsiMemUsage.With(prometheus.Labels{"exsi":exsiHostName}).Set(float64(hs.Summary.QuickStats.OverallMemoryUsage))



		pc := property.DefaultCollector(c)
		//获取vm虚拟机的磁盘信息，先转换类型，vm.Datastore转船成mo.Datastore{}
		var hdss []mo.Datastore

		err := pc.Retrieve(ctx, hs.Datastore, []string{"info","summary"}, &hdss)
		if err != nil {
			fmt.Println(exsiHostName + "  get exsi datastore info error ", err)
		}else {
			//虚拟机有多个的存储
			for _, ds := range hdss {
				//虚拟机磁盘总大小
				exsidsCapacity.With(prometheus.Labels{"exsi":exsiHostName,"exsidisk":ds.Summary.Name}).Set(float64(ds.Summary.Capacity / 1024 / 1024 / 1024))
				//未使用空间
				exsidsFreeSpace.With(prometheus.Labels{"exsi":exsiHostName,"exsidisk":ds.Summary.Name}).Set(float64(ds.Summary.FreeSpace / 1024/ 1024 / 1024))
			}
		}
	}





}