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
