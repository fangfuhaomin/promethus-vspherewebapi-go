package vcconnect

import (
	"context"
	"github.com/vmware/govmomi"
	"log"
	"net/url"
)


func Vccon(user string,password string,ip string) *govmomi.Client {
	u := &url.URL{
		Scheme: "https",
		Host:   ip,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(user, password)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		log.Fatal(err)
	}
	//打印vsphere连接
	return client
}