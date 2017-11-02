package middlewares

import "net/http"
//单个中单件
/*
func main() {
	Stack := &MiddlewareStack{}

	// Add middleware `auth` to stack
	Stack.Use(&middlewares.Middleware{
		Name: "auth",
		// Insert middleware `auth` after middleware `session` if it exists
		InsertAfter: []string{"session"},
		// Insert middleware `auth` before middleare `authorization` if it exists
		InsertBefore: []string{"authorization"},
	})

	// Remove middleware `cookie` from stack
	Stack.Remove("cookie")

	mux := http.NewServeMux()
	http.ListenAndServe(":9000", Stack.Apply(mux))
}
*/
type Middleware struct {
	Name         string
	Handler      func(http.Handler) http.Handler
	InsertAfter  []string  //此中间件插入到哪些中间件之后，在Middlewares sort中有用
	InsertBefore []string  //此中间件插入到哟些中间件之前
	Requires     []string  //此中间件需要哪些其它的中间件
}
