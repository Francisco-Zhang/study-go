package handler

// HelloServiceName 加前缀是为了解决名称冲突问题
const HelloServiceName = "handler/HelloService"

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello," + request
	return nil
}
