package models

type Country struct {
	Name string
	Code string
}

type Region struct {
	Name      string
	Countries []Country
	Pools     []Pool
}

type Worker struct {
	Name      string
	SubDomain string
}

type Pool struct {
	Name      string `json:"name"`
	Region    string `json:"region"`
	Subdomain string `json:"subdomain"`
	PortStart int    `json:"port_start"`
	PortEnd   int    `json:"port_end"`
	Outs      []Out  `json:"outs"`
}

type Out struct {
	Format       string `json:"format"`
	UpstreamPort int    `json:"upstream_port"`
	Domain       string `json:"domain"`
	Weight       int    `json:"weight"`
}

type CreatePoolRequest struct {
	Name      string `json:"name" binding:"required"`
	Region    string `json:"region" binding:"required"`
	Subdomain string `json:"subdomain" binding:"required"`
	PortStart int    `json:"port_start" binding:"required"`
	PortEnd   int    `json:"port_end" binding:"required"`
	Outs      []Out  `json:"outs" binding:"required"`
}

type UpdatePoolRequest struct {
	Region    *string `json:"region,omitempty" `
	Subdomain *string `json:"subdomain,omitempty"`
	PortStart *int    `json:"port_start,omitempty"`
	PortEnd   *int    `json:"port_end,omitempty"`
	Outs      *[]Out  `json:"outs,omitempty"`
}
