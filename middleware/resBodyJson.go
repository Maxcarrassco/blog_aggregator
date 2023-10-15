package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"

)



func Json (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		out := map[string]any{}
		body := fmt.Sprintf("%s", data)
		body = strings.Trim(body, "{}")
		resBody := strings.Split(body, ",")
		for _, v := range resBody {
			val := strings.Split(v, ":")
			fmt.Println(val[0], val[1])
			out[val[0]] = val[1]
		}
		fmt.Printf("%v\n, %v", out, out["age"])
		next.ServeHTTP(w, r)
	})
}
