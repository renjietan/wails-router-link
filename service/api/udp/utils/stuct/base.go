package udp_utils_struct

type Base struct {
	Data  any   `json:"data"`
	Range []int `json:"rg"`
	Size  int   `json:"length"`
}
