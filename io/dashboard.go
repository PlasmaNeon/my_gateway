package io

type DashboardSumDataOutput struct {
	ServiceNum      int `json:"service_num"`
	AppNum          int `json:"app_num"`
	CurrentQPS      int `json:"current_qps"`
	TodayRequestNum int `json:"today_request_num"`
}

type DashboardGlobalFlowStatOutput struct {
	Today     []int `json:"today"`
	Yesterday []int `json:"yesterday"`
}

type ServiceCount struct {
	LoadType int    `json:"load_type"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type DashboardGlobalServiceCountOutput struct {
	Services []string       `json:"services"`
	Count    []ServiceCount `json:"count"`
}
