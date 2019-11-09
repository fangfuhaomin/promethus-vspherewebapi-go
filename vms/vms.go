package vms

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)
//虚拟机属性
var (
	vmsUsageCpu = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vms_cpu_usage",
		Help: "Current CPU usage MHZ",
	}, []string{"vm","vmip","cpu"}, )

	vmsMemorycapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vms_memory_total",
		Help: "Total virtual machine memory(M)",
	}, []string{"vm","vmip","mem"}, )

	vmsMemoryUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vms_memory_usage",
		Help: "Current memory usage MHZ (M)",
	}, []string{"vm","vmip","mem"}, )


	vmsDiskCapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vms_disk_capacity",
		Help: "磁盘最大容量 单位GB",
	}, []string{"vm","vmip","disk"}, )

	vmsDiskFreeSpace = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vms_disk_freespace",
		Help: "磁盘空闲空间 单位GB",
	}, []string{"vm","vmip","disk"}, )
	vmsPoweredstate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vms_Poweredstate",
		Help: "虚拟机电源状态 0是关 1是开",
	}, []string{"vm","vmip","powered"}, )

)

func init() {
	// 注册指标
	prometheus.MustRegister(vmsUsageCpu)
	prometheus.MustRegister(vmsMemorycapacity)
	prometheus.MustRegister(vmsMemoryUsage)
	prometheus.MustRegister(vmsDiskCapacity)
	prometheus.MustRegister(vmsDiskFreeSpace)
	prometheus.MustRegister(vmsPoweredstate)
}


//这里递归查找所有虚拟机，通过虚拟机视图来发现
// 这里的 c 就是上面登录后 client 的 Client 属性
func GetVmsInfo(c *vim25.Client) {
	//c := vcconnect.Vccon().Client
	ctx := context.Background()
	m := view.NewManager(c)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		panic(err)
	}
	defer v.Destroy(ctx)
	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary","guest","triggeredAlarmState"}, &vms)
	if err != nil {
		panic(err)
	}
	for _,vm := range vms {
		//虚拟机名字
		vmName := vm.Summary.Config.Name
		vmip := vm.Guest.IpAddress

		//虚拟机disk
		for _,vmds := range vm.Guest.Disk {
			//虚拟机disk 磁盘最大容量 GB
			vmsDiskCapacity.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"disk":vmds.DiskPath}).Set(float64(vmds.Capacity/1024/1024/1024))

			//虚拟机disk 磁盘空闲容量 GB
			vmsDiskFreeSpace.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"disk":vmds.DiskPath}).Set(float64(vmds.FreeSpace/1024/1024/1024))
		}

		//虚拟机当前使用CPU MHZ
		vmsUsageCpu.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"cpu":"cpu"}).Set(float64(vm.Summary.QuickStats.OverallCpuUsage))

		//虚拟机最大内存capacity
		vmsMemorycapacity.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"mem":"mem"}).Set(float64(vm.Summary.Config.MemorySizeMB))

		//虚拟机当前使用内存
		vmsMemoryUsage.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"mem":"mem"}).Set(float64(vm.Summary.QuickStats.OverallCpuUsage))
		//fmt.Println(vmName,vm.Summary.Runtime.PowerState)
		//虚拟机当前使用内存
		if vm.Summary.Runtime.PowerState == "poweredOn" {
			vmsPoweredstate.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"powered":"powered"}).Set(float64(1))
		}else {
			vmsPoweredstate.With(prometheus.Labels{"vm":vmName,"vmip":vmip,"powered":"powered"}).Set(float64(0))
		}


	}
}
