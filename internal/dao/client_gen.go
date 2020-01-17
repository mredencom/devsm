package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"kratos-demo/api"
)

type Daos struct {
	demoClient api.DemoClient
}

func NewC() (d *Daos) {
	cfg := &warden.ClientConfig{}
	paladin.Get("grpc.toml").UnmarshalTOML(cfg)
	d = &Daos{}
	var err error
	if d.demoClient, err = api.NewClient(cfg); err != nil {
		panic(err)
	}
	return
}

// 实现接口
func (d *Daos) SayHello(c context.Context, req *api.HelloReq) (resp *empty.Empty, err error) {
	if resp, err = d.demoClient.SayHello(c, req); err != nil {
		err = errors.Wrapf(err, "%v", nil)
	}
	return
}
