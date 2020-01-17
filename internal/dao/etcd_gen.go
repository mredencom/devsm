package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"kratos-demo/api"

	"os"
)

var s = wire.NewSet(NewClient)
// AppID your appid, ensure unique.
const AppID = "demo.service" // NOTE: example

func init(){
	// NOTE: 注意这段代码，表示要使用etcd进行服务发现 ,其他事项参考discovery的说明
	// NOTE: 在启动应用时，可以通过flag(-etcd.endpoints) 或者 环境配置(ETCD_ENDPOINTS)指定etcd节点
	// NOTE: 如果需要自己指定配置时 需要同时设置DialTimeout 与 DialOptions: []grpc.DialOption{grpc.WithBlock()}
	//resolver.Register(etcd.Builder(nil))
}

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (api.DemoClient, error) {
	client := warden.NewClient(cfg, opts...)
	// 这里使用etcd scheme
	conn, err := client.Dial(context.Background(), "etcd://default/"+AppID)
	if err != nil {
		return nil, err
	}
	// 注意替换这里：
	// NewDemoClient方法是在"api"目录下代码生成的
	// 对应proto文件内自定义的service名字，请使用正确方法名替换
	return api.NewDemoClient(conn), nil
}

func EtcdNew() {
	ip := "" // NOTE: 必须拿到您实例节点的真实IP，
	port := "" // NOTE: 必须拿到您实例grpc监听的真实端口，warden默认监听9000
	hn, _ := os.Hostname()
	dis := discovery.New(nil)
	ins := &naming.Instance{
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		AppID:    "your app id",
		Hostname: hn,
		Addrs: []string{
			"grpc://" + ip + ":" + port,
		},
	}
	cancel, err := dis.Register(context.Background(), ins)
	if err != nil {
		panic(err)
	}
	cancel()
}
