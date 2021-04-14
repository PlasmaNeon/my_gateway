package db

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"Basic info."`
	HTTPRule      *HTTPRule      `json:"http_rule" description:"HTTP rules."`
	TCPRule       *TCPRule       `json:"tcp_rule"  description:"TCP rules."`
	GRPCRule      *GRPCRule      `json:"grpc_rule"  description:"GRPC rules."`
	LoadBalance   *LoadBalance   `json:"load_balance"  description:"Load balance."`
	AccessControl *AccessControl `json:"access_control"  description:"Access control."`
}
