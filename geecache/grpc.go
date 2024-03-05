package geecache

import (
	"context"
	"fmt"
	"geecache/consistenthash"
	pb "geecache/geecachepb"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcdefaultBasePath = "/_geecache/"
	grpcdefaultReplicas = 50
)

// grpc客户端
type grpcGetter struct {
	baseurl string // 地址
}

func (g *grpcGetter) Get(in *pb.Request, out *pb.Response) error {
	coon, err := grpc.Dial(g.baseurl, grpc.WithInsecure()) // 第二个参数即不使用tls
	if err != nil {
		return err
	}
	defer coon.Close()
	client := pb.NewGroupCacheServiceClient(coon)
	response, err := client.Get(context.Background(), in)
	out.Value = response.Value
	return err
}

// 检查grpcGetter是否实现了PeerGetter接口
var _ PeerGetter = (*grpcGetter)(nil)

// GrpcPool用于管理节集群节点选取器，保存所有节点信息
// 获取对等节点缓存时，通过计算key的“一致性哈希值”与节点的哈希值比较来选取集群中的某个节点
type GrpcPool struct {
	pb.UnimplementedGroupCacheServiceServer

	self        string
	mu          sync.Mutex
	peers       *consistenthash.ConsistentHash
	grpcGetters map[string]*grpcGetter
}

// 实例化GrpcPool
func NewGrpcPool(self string) *GrpcPool {
	return &GrpcPool{
		self:        self,
		peers:       consistenthash.New(grpcdefaultReplicas, nil),
		grpcGetters: map[string]*grpcGetter{},
	}
}

func (p *GrpcPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers.AddTrueNode(peers...)
	for _, peer := range peers {
		p.grpcGetters[peer] = &grpcGetter{
			baseurl: peer,
		}
	}
}

// 根据key选择节点，然后返回节点对应的GrpcGetter
func (p *GrpcPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.GetTrueNode(key); peer != "" && peer != p.self {
		return p.grpcGetters[peer], true
	}
	return nil, false
}

// 检查Grpc是否实现了PeerPicker接口
var _ PeerPicker = (*GrpcPool)(nil)

// 日志
func (p *GrpcPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 作为服务器端对get方法的处理
func (p *GrpcPool) Get(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	p.Log("%s %s", in.Group, in.Key)
	response := &pb.Response{}

	group := GetGroup(in.Group)
	if group == nil {
		p.Log("no such group %v", in.Group)
		return response, fmt.Errorf("no such group %v", in.Group)
	}

	value, err := group.Get(in.Key)
	if err != nil {
		p.Log("get key %v error %v", in.Key, err)
		return response, fmt.Errorf("get key %v error %v", in.Key, err)
	}

	response.Value = value.ByteSlice()
	return response, nil
}

func (p *GrpcPool) Run() {
	lis, err := net.Listen("tcp", p.self)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterGroupCacheServiceServer(server, p)

	reflection.Register(server)
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
