package controllers

import (
	"net/http"

	"github.com/midedickson/simple-banking-app/utils"
)

func (c *Controller) RequestNewIdempotencyKey(w http.ResponseWriter, r *http.Request) {
	key, err := c.idempotencyStore.CreateNewIdempotencyKey()
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	// set indempotent key in the response headers
	w.Header().Set("X-Idempotency-Key", key)
	// send the response with the generated idempotency key in the body
	utils.Dispatch200(w, "New Idempotency Key generated successfully", map[string]string{"idempotency_key": key})
}
