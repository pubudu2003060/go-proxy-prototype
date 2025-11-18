package models

type Country struct {
	Code string
}

type Region struct {
	RName     string
	Countries []Country
	Workers   []Worker
}

type Worker struct {
	WName string
	Pools []*Pool
}

type Pool struct {
	Name      string `json:"name"`
	Continent string `json:"continent"`
	Tag       string `json:"tag"`
	Subdomain string `json:"subdomain"`
	CC3       string `json:"cc3"`
	PortStart int    `json:"port_start"`
	PortEnd   int    `json:"port_end"`
	Flag      int    `json:"flag"`
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
	Continent string `json:"continent" binding:"required"`
	Tag       string `json:"tag" binding:"required"`
	Subdomain string `json:"subdomain" binding:"required"`
	CC3       string `json:"cc3"`
	PortStart int    `json:"port_start" binding:"required"`
	PortEnd   int    `json:"port_end" binding:"required"`
	Outs      []Out  `json:"outs" binding:"required"`
}

type UpdatePoolRequest struct {
	Continent *string `json:"continent,omitempty" `
	Tag       *string `json:"tag,omitempty"`
	Subdomain *string `json:"subdomain,omitempty"`
	CC3       *string `json:"cc3,omitempty"`
	PortStart *int    `json:"port_start,omitempty"`
	PortEnd   *int    `json:"port_end,omitempty"`
	Outs      *[]Out  `json:"outs,omitempty"`
}
