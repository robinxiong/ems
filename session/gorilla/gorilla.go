package gorilla

import (
	"context"
	"ems/session"
	"net/http"

	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"ems/core/utils"
	"fmt"
	"encoding/json"
)

var reader utils.ContextKey = "gorilla_reader"

//New initialize session manager for gorilla
func New(sessionName string, store sessions.Store) *Gorilla {
	return &Gorilla{SessionName: sessionName, Store: store}
}

type Gorilla struct {
	SessionName string
	Store       sessions.Store
}


func (gorilla Gorilla) getSession(req *http.Request) (*sessions.Session, error) {
	if r, ok := req.Context().Value(reader).(*http.Request); ok {
		return gorilla.Store.Get(r, gorilla.SessionName)
	}

	/*	gorilla.Store调用Get方法, 它将调用Cookiestore的Get方法
		CookieStore的Get方法将查找内存对像gorilla/context中保存的以req为键的Registry(每一次请求完成时gorilla/context会清空当前的req所对应的Registry)

		type Registry struct {
			request  *http.Request
			sessions map[string]sessionInfo
		}

		Registry在Get(cookieStore, gorilla.SessionName)
		如果的registery的session中没有保存到name所对应的session, 则调用store创建一个session, 它将调用CookieStore的New方法，从request的cookie中解析数据
		并且保存进registry.sessions中

		保存session到cookie, 则是调用regsitry.Get返回的session的Save, 它会调用CookieStore在保存到cookie中


	 */
	return gorilla.Store.Get(req, gorilla.SessionName)
}

func (gorilla Gorilla) saveSession(w http.ResponseWriter, req *http.Request) {
	if session, err := gorilla.getSession(req); err == nil {
		if err := session.Save(req, w); err != nil {
			fmt.Printf("No error should happen when saving session data, but got %v", err)
		}
	}
}

//ManagerInterface method
func (gorilla Gorilla) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	defer gorilla.saveSession(w, req)

	session, err := gorilla.getSession(req)
	if err != nil {
		return err
	}

	if str, ok := value.(string); ok {
		session.Values[key] = str
	} else {
		result, _ := json.Marshal(value)
		session.Values[key] = string(result)
	}

	return nil
}


// Pop value from session data
func (gorilla Gorilla) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	defer gorilla.saveSession(w, req)

	if session, err := gorilla.getSession(req); err == nil {
		if value, ok := session.Values[key]; ok {
			delete(session.Values, key)
			return fmt.Sprint(value)
		}
	}
	return ""
}

func (gorilla Gorilla) Get(req *http.Request, key string) string {
	if session, err := gorilla.getSession(req); err == nil {
		if value, ok := session.Values[key]; ok {
			return fmt.Sprint(value)
		}
	}
	return ""
}
func (gorilla Gorilla) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := gorilla.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return gorilla.Add(w, req, "_flashes", messages)
}

func (gorilla Gorilla) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	gorilla.PopLoad(w, req, "_flashes", &messages)
	return messages
}

func (gorilla Gorilla) Load(req *http.Request, key string, result interface{}) error {
	value := gorilla.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// PopLoad pop value from session data and unmarshal it into result
func (gorilla Gorilla) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := gorilla.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// Middleware returns a new session manager middleware instance
func (gorilla Gorilla) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer gorillaContext.Clear(req) //清空gorillaContext中保存的此次请求的上下文数据，否则导致内存泄漏
		ctx := context.WithValue(req.Context(), reader, req)  //可以忽略此步骤
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
