package http

import "net/http"

func (c *Controller) SayHello(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(c.HelloSayerUserCase.SayHello(r.Context(), name)))
}
