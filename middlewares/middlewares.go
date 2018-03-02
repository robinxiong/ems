package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

// MiddlewareStack middlewares stack
type MiddlewareStack struct {
	middlewares []*Middleware
}

//将某个middleware添加到数组中
func (stack *MiddlewareStack) Use(middleware Middleware) {
	stack.middlewares = append(stack.middlewares, &middleware)
}

// Remove remove middleware by name
func (stack *MiddlewareStack) Remove(name string) {
	registeredMiddlewares := stack.middlewares
	for idx, middleware := range registeredMiddlewares {
		if middleware.Name == name {
			if idx > 0 {
				stack.middlewares = stack.middlewares[0 : idx-1]
			} else {
				stack.middlewares = []*Middleware{}
			}

			if idx < len(registeredMiddlewares)-1 {
				stack.middlewares = append(stack.middlewares, registeredMiddlewares[idx+1:]...)
			}
		}
	}
}

// sortMiddlewares sort middlewares
// 首先检查每一个middleware的Requires项，如果缺下依赖的中间件，则报错
// 完善每一个middleware的insertAfter And insertBefore数组，使每一个中间件都能知道它之前与之后要执行的中间件
// 上一步获取了每一个中间件的相关依赖以及前后关系，然后对这个中间件相关的上下游middleware进行排序
// 对整个stack进行排序
func (stack *MiddlewareStack) sortMiddlewares() (sortedMiddlewares []*Middleware, err error) {
	var (
		errs                         []error
		middlewareNames, sortedNames []string
		middlewaresMap               = map[string]*Middleware{}
		sortMiddleware               func(m *Middleware)
	)

	for _, middleware := range stack.middlewares {
		middlewaresMap[middleware.Name] = middleware
		middlewareNames = append(middlewareNames, middleware.Name)
	}

	for _, middleware := range stack.middlewares {
		for _, require := range middleware.Requires {
			if _, ok := middlewaresMap[require]; !ok {
				errs = append(errs, fmt.Errorf("middleware %v requires %v, but it doesn't exist", middleware.Name, require))
			}
		}

		for _, insertBefore := range middleware.InsertBefore {
			if m, ok := middlewaresMap[insertBefore]; ok {
				m.InsertAfter = uniqueAppend(m.InsertAfter, middleware.Name)
			}
		}

		for _, insertAfter := range middleware.InsertAfter {
			if m, ok := middlewaresMap[insertAfter]; ok {
				m.InsertBefore = uniqueAppend(m.InsertBefore, middleware.Name)
			}
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("%v", errs)
	}

	//   C, B <- D -> E, F, A <- B -> C
	sortMiddleware = func(m *Middleware) {
		if _, found := getRIndex(sortedNames, m.Name); !found { // if not sorted
			//此中间件没有排序
			var minIndex = -1

			// sort by InsertAfter
			//先把此中间件之前的中间件进行排序
			for _, insertAfter := range m.InsertAfter {
				idx, found := getRIndex(sortedNames, insertAfter)
				if !found {
					if middleware, ok := middlewaresMap[insertAfter]; ok {
						sortMiddleware(middleware)
						idx, found = getRIndex(sortedNames, insertAfter)
					}
				}
				//如果之前的中间件已经排序，即在sortedNames中找到。
				if found && idx > minIndex {
					minIndex = idx
				}
			}

			// sort by InsertBefore
			// 遍历之后的中间件，如果之后的中间件出件在排序中, 而之后的中间件的位置要比它之前的位置小（idx < minIndex)
			// 则将这个之后的中间件从sortedNames中删除，然后对这个之后的中间件进行重排
			for _, insertBefore := range m.InsertBefore {
				if idx, found := getRIndex(sortedNames, insertBefore); found {
					if idx < minIndex {
						sortedNames = append(sortedNames[:idx], sortedNames[idx+1:]...)
						sortMiddleware(middlewaresMap[insertBefore])
						return
					}
				}
			}

			if minIndex >= 0 {
				sortedNames = append(sortedNames[:minIndex+1], append([]string{m.Name}, sortedNames[minIndex+1:]...)...)
			} else if _, has := getRIndex(sortedNames, m.Name); !has {
				// if current callback haven't been sorted, append it to last
				sortedNames = append(sortedNames, m.Name)
			}
		}
	}

	for _, middleware := range stack.middlewares {
		sortMiddleware(middleware)
	}

	for _, name := range sortedNames {
		sortedMiddlewares = append(sortedMiddlewares, middlewaresMap[name])
	}

	return sortedMiddlewares, nil
}

//按insertAfter, intertBefore中的顺序，输出所有的中间件
func (stack *MiddlewareStack) String() string {
	var (
		sortedNames            []string
		sortedMiddlewares, err = stack.sortMiddlewares()
	)

	if err != nil {
		fmt.Println(err)
	}

	for _, middleware := range sortedMiddlewares {
		sortedNames = append(sortedNames, middleware.Name)
	}

	return fmt.Sprintf("MiddlewareStack: %v", strings.Join(sortedNames, ", "))
}

// Apply apply middlewares to handler
// 将Handler按倒序传入给所有的中间件，然后将包装后的handler往前传入,  如果某一个没有返回，则中间件调用传入的handler参数
/*

func TestApply(t *testing.T) {
	stack := &MiddlewareStack{}
	stack.Use(Middleware{
		Name: "B",
		Handler: func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println("B")
				defer handler.ServeHTTP(w, r)
			})
		},
		InsertAfter: []string{"A"},
		InsertBefore:[]string{"C"},
		Requires: []string{"A", "C"},
	})
	stack.Use(Middleware{
		Name: "A",
		Handler: func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println("A")
				defer handler.ServeHTTP(w, r)
			})
		},
	})
	stack.Use(Middleware{
		Name: "C",
		Handler: func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println("C")
				defer handler.ServeHTTP(w, r)
			})
		},
		InsertAfter: []string{"B"},
	})
	log.Println(stack.String())
	result := stack.Apply(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Println("hanlder function")
	}))

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	result.ServeHTTP(res, req)
}

*/
func (stack *MiddlewareStack) Apply(handler http.Handler) http.Handler {
	var (
		compiledHandler        http.Handler
		sortedMiddlewares, err = stack.sortMiddlewares()
	)

	if err != nil {
		fmt.Println(err)
	}

	for idx := len(sortedMiddlewares) - 1; idx >= 0; idx-- {

		middleware := sortedMiddlewares[idx]

		if compiledHandler == nil {
			compiledHandler = middleware.Handler(handler)
		} else {
			compiledHandler = middleware.Handler(compiledHandler)
		}
	}

	return compiledHandler
}
