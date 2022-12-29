package models

type InputData struct {
	ClientID int `json:"clientId"`
	Puntos   int `json:"puntos"`
}

type InputPostData struct {
	ClientID   int    `json:"clientId"`
	NroFactura string `json:"nroFactura"`
	Vlrfactura int    `json:"vlrfactura"`
	Puntos     int    `json:"puntos"`
}

type Item struct {
	ClientID int
	Points   int
}
