package main

type InteractionRequest struct {
	Id            string      `json:"id"`
	ApplicationId string      `json:"application_id"`
	Type          uint        `json:"type"`
	Data          interface{} `json:"data"`
}

type InteractionResponse struct {
	Type uint        `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type Commands struct {
	Name        string `json:"name"`
	Type        uint   `json:"type"`
	Description string `json:"description"`
}
