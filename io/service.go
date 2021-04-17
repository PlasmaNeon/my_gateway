package io

import (
	"github.com/gin-gonic/gin"
	"my_gateway/public"
)

type ServiceAddHTTPInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"Service Description" example:"" validate:"required,min=1,max=255"`

	RuleType      int    `json:"rule_type" form:"rule_type" comment:"Connect Rule Type" example:"" validate:"min=0,max=1"`
	Rule          string `json:"rule" form:"rule" comment:"Connect Rule(Prefix/Domain)" example:"" validate:"required,is_valid_rule"`
	NeedHttps     int    `json:"need_https" form:"need_https" comment:"Need Https?" example:"" validate:"min=0,max=1"`
	NeedStripURI  int    `json:"need_strip_uri" form:"need_strip_uri" comment:"NeedStripURI" example:"" validate:"min=0,max=1"`
	NeedWebsocket int    `json:"need_websocket" form:"need_websocket" comment:"need_websocket" example:"" validate:"min=0,max=1"`
	URLRewrite    string `json:"url_rewrite" form:"url_rewrite" comment:"url_rewrite" example:"" validate:"is_valid_url_rewrite"`
	HeadTransform string `json:"head_transform" form:"head_transform" comment:"head_transform" example:"" validate:""`

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"open_auth" example:"" validate:"min=0,max=1"`
	BlackList         string `json:"black_list" form:"black_list" comment:"black_list" example:"" validate:""`
	WhiteList         string `json:"white_list" form:"white_list" comment:"white_list" example:"" validate:""`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"client_ip_flow_limit" example:"" validate:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"service_flow_limit" example:"" validate:"min=0"`
	RoundType         int    `json:"round_type" form:"round_type" comment:"round_type" example:"" validate:"min=0,max=3"`

	IPList                 string `json:"ip_list" form:"ip_list" comment:"ip_list" example:"" validate:"required,is_valid_ip_list"`
	WeightList             string `json:"weight_list" form:"weight_list" comment:"weight_list" example:"" validate:"required,is_valid_weight_list"`
	UpstreamConnectTimeOut int    `json:"upstream_connect_time_out" form:"upstream_connect_time_out" comment:"upstream_connect_time_out" example:"" validate:"min=0"`
	UpstreamHeaderTimeOut  int    `json:"upstream_header_time_out" form:"upstream_header_time_out" comment:"upstream_header_time_out" example:"" validate:"min=0"`
	UpstreamIdleTimeOut    int    `json:"upstream_idle_time_out" form:"upstream_idle_time_out" comment:"upstream_idle_time_out" example:"" validate:"min=0"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"upstream_max_idle" example:"" validate:"min=0"`
}

//Validation process
func (param *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceDeleteInput struct {
	ID int `json:"id" form:"id" comment:"Service id" example:"56" validate:"required"`
}

//Validation process
func (param *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"Keyword" example:"" validate:""`
	PageNo   int    `json:"page_no" form:"page_no" comment:"page number" example:"1" validate:"required"`
	PageSize int    `json:"page_size" form:"page_size" comment:"number of entries per page" example:"20" validate:"required"`
}

//Validation process
func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateHTTPInput struct {
	ID          int    `json:"id" form:"id" comment:"Service id" example:"56" validate:"required,min=1"`
	ServiceName string `json:"service_name" form:"service_name" comment:"Service Name" example:"test_add_http" validate:"required,is_valid_service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"Service Description" example:"test_add_http" validate:"required,min=1,max=255"`

	RuleType      int    `json:"rule_type" form:"rule_type" comment:"Connect Rule Type" example:"" validate:"min=0,max=1"`
	Rule          string `json:"rule" form:"rule" comment:"Connect Rule(Prefix/Domain)" example:"" validate:"required,is_valid_rule"`
	NeedHttps     int    `json:"need_https" form:"need_https" comment:"Need Https?" example:"" validate:"min=0,max=1"`
	NeedStripURI  int    `json:"need_strip_uri" form:"need_strip_uri" comment:"NeedStripURI" example:"" validate:"min=0,max=1"`
	NeedWebsocket int    `json:"need_websocket" form:"need_websocket" comment:"need_websocket" example:"" validate:"min=0,max=1"`
	URLRewrite    string `json:"url_rewrite" form:"url_rewrite" comment:"url_rewrite" example:"" validate:"is_valid_url_rewrite"`
	HeadTransform string `json:"head_transform" form:"head_transform" comment:"head_transform" example:"" validate:""`

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"open_auth" example:"" validate:"min=0,max=1"`
	BlackList         string `json:"black_list" form:"black_list" comment:"black_list" example:"" validate:""`
	WhiteList         string `json:"white_list" form:"white_list" comment:"white_list" example:"" validate:""`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"client_ip_flow_limit" example:"" validate:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"service_flow_limit" example:"" validate:"min=0"`
	RoundType         int    `json:"round_type" form:"round_type" comment:"round_type" example:"" validate:"min=0,max=3"`

	IPList                 string `json:"ip_list" form:"ip_list" comment:"ip_list" example:"127.0.0.1:80" validate:"required,is_valid_ip_list"`
	WeightList             string `json:"weight_list" form:"weight_list" comment:"weight_list" example:"50" validate:"required,is_valid_weight_list"`
	UpstreamConnectTimeOut int    `json:"upstream_connect_time_out" form:"upstream_connect_time_out" comment:"upstream_connect_time_out" example:"" validate:"min=0"`
	UpstreamHeaderTimeOut  int    `json:"upstream_header_time_out" form:"upstream_header_time_out" comment:"upstream_header_time_out" example:"" validate:"min=0"`
	UpstreamIdleTimeOut    int    `json:"upstream_idle_time_out" form:"upstream_idle_time_out" comment:"upstream_idle_time_out" example:"" validate:"min=0"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"upstream_max_idle" example:"" validate:"min=0"`
}

func (param *ServiceUpdateHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddTCPInput struct {
	//ID            int64  `json:"id" form:"id" comment:"Service id" example:"56" validate:"required,min=1"`
	ServiceName string `json:"service_name" form:"service_name" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"Service Description" example:"" validate:"required,min=1,max=255"`

	HeadTransform string `json:"head_transform" form:"head_transform" comment:"head_transform" example:"" validate:""`

	Port              int    `json:"port" form:"port" comment:"port" example:"" validate:"required,min=8001,max=9999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"open_auth" example:"" validate:"min=0,max=1"`
	BlackList         string `json:"black_list" form:"black_list" comment:"black_list" example:"" validate:""`
	WhiteList         string `json:"white_list" form:"white_list" comment:"white_list" example:"" validate:""`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"client_ip_flow_limit" example:"" validate:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"service_flow_limit" example:"" validate:"min=0"`
	RoundType         int    `json:"round_type" form:"round_type" comment:"round_type" example:"" validate:"min=0,max=3"`

	WhiteHostName string `json:"white_host_name" form:"white_host_name" comment:"white_host_nam" validate:"is_valid_ip_list"`
	IPList        string `json:"ip_list" form:"ip_list" comment:"ip_list" example:"" validate:"required,is_valid_ip_list"`
	WeightList    string `json:"weight_list" form:"weight_list" comment:"weight_list" example:"" validate:"required,is_valid_weight_list"`
	ForbidList    string `json:"forbid_list" form:"forbid_list" comment:"forbid_list" example:"" validate:"required,is_valid_ip_list"`
}

// BindValidParam
// Validation process
func (param *ServiceAddTCPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateTCPInput struct {
	ID          int    `json:"id" form:"id" comment:"Service id" example:"56" validate:"required,min=1"`
	ServiceName string `json:"service_name" form:"service_name" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"Service Description" example:"" validate:"required,min=1,max=255"`

	HeadTransform string `json:"head_transform" form:"head_transform" comment:"head_transform" example:"" validate:""`

	Port              int    `json:"port" form:"port" comment:"port" example:"" validate:"required,min=8001,max=9999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"open_auth" example:"" validate:"min=0,max=1"`
	BlackList         string `json:"black_list" form:"black_list" comment:"black_list" example:"" validate:""`
	WhiteList         string `json:"white_list" form:"white_list" comment:"white_list" example:"" validate:""`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"client_ip_flow_limit" example:"" validate:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"service_flow_limit" example:"" validate:"min=0"`
	RoundType         int    `json:"round_type" form:"round_type" comment:"round_type" example:"" validate:"min=0,max=3"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"white_host_nam" validate:"is_valid_ip_list"`
	IPList            string `json:"ip_list" form:"ip_list" comment:"ip_list" example:"" validate:"required,is_valid_ip_list"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"weight_list" example:"" validate:"required,is_valid_weight_list"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"forbid_list" example:"" validate:"required,is_valid_ip_list"`
}

func (param *ServiceUpdateTCPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddGRPCInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"Service Description" example:"" validate:"required,min=1,max=255"`

	HeadTransform string `json:"head_transform" form:"head_transform" comment:"head_transform" example:"" validate:""`

	Port              int    `json:"port" form:"port" comment:"port" example:"" validate:"required,min=8001,max=9999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"open_auth" example:"" validate:"min=0,max=1"`
	BlackList         string `json:"black_list" form:"black_list" comment:"black_list" example:"" validate:""`
	WhiteList         string `json:"white_list" form:"white_list" comment:"white_list" example:"" validate:""`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"client_ip_flow_limit" example:"" validate:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"service_flow_limit" example:"" validate:"min=0"`
	RoundType         int    `json:"round_type" form:"round_type" comment:"round_type" example:"" validate:"min=0,max=3"`

	WhiteHostName string `json:"white_host_name" form:"white_host_name" comment:"white_host_nam" validate:"is_valid_ip_list"`
	IPList        string `json:"ip_list" form:"ip_list" comment:"ip_list" example:"" validate:"required,is_valid_ip_list"`
	WeightList    string `json:"weight_list" form:"weight_list" comment:"weight_list" example:"" validate:"required,is_valid_weight_list"`
	ForbidList    string `json:"forbid_list" form:"forbid_list" comment:"forbid_list" example:"" validate:"required,is_valid_ip_list"`
}

//Validation process
func (param *ServiceAddGRPCInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateGRPCInput struct {
	ID          int    `json:"id" form:"id" comment:"Service id" example:"56" validate:"required,min=1"`
	ServiceName string `json:"service_name" form:"service_name" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"Service Description" example:"" validate:"required,min=1,max=255"`

	HeadTransform string `json:"head_transform" form:"head_transform" comment:"head_transform" example:"" validate:""`

	Port              int    `json:"port" form:"port" comment:"port" example:"" validate:"required,min=8001,max=9999"`
	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"open_auth" example:"" validate:"min=0,max=1"`
	BlackList         string `json:"black_list" form:"black_list" comment:"black_list" example:"" validate:""`
	WhiteList         string `json:"white_list" form:"white_list" comment:"white_list" example:"" validate:""`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"client_ip_flow_limit" example:"" validate:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"service_flow_limit" example:"" validate:"min=0"`
	RoundType         int    `json:"round_type" form:"round_type" comment:"round_type" example:"" validate:"min=0,max=3"`

	WhiteHostName string `json:"white_host_name" form:"white_host_name" comment:"white_host_nam" validate:"is_valid_ip_list"`
	IPList        string `json:"ip_list" form:"ip_list" comment:"ip_list" example:"" validate:"required,is_valid_ip_list"`
	WeightList    string `json:"weight_list" form:"weight_list" comment:"weight_list" example:"" validate:"required,is_valid_weight_list"`
	ForbidList    string `json:"forbid_list" form:"forbid_list" comment:"forbid_list" example:"" validate:"required,is_valid_ip_list"`
}

func (param *ServiceUpdateGRPCInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListItemOutput struct {
	ID          int    `json:"id" form:"id"`
	ServiceName string `json:"service_name" form:"service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc"`
	LoadType    int    `json:"load_type" form:"load_type"`
	ServiceAddr string `json:"service_addr" form:"service_addr"`
	QPS         int    `json:"qps" form:"qps"`
	QPD         int    `json:"qpd" form:"qpd"`
	TotalNode   int    `json:"total_node" form:"total_node"`
}

type ServiceListOutput struct {
	Total int                     `json:"total" form:"total" comment:"Total entries"`
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"list number"`
}
type ServiceStatOutput struct {
	Today     []int `json:"today" form:"today" comment:"Today"`
	Yesterday []int `json:"yesterday" form:"yesterday" comment:"Yesterday"`
}
