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
	Name       string
	SubDomains []string
}

type Pool struct {
	Name      string `json:"name"`
	Region    string `json:"region"`
	Subdomain string `json:"subdomain"`
	Port      int    `json:"port"`
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
	Port      int    `json:"port" binding:"required"`
	Outs      []Out  `json:"outs" binding:"required"`
}

type UpdatePoolRequest struct {
	Region    *string `json:"region,omitempty" `
	Subdomain *string `json:"subdomain,omitempty"`
	Port      *int    `json:"port,omitempty"`
	Outs      *[]Out  `json:"outs,omitempty"`
}
