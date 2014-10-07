package siesta

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var services map[string]*Service = make(map[string]*Service)

type Service struct {
	baseURI string

	pre  []contextHandler
	post []contextHandler

	handlers map[*regexp.Regexp]contextHandler
}

func NewService(baseURI string) *Service {
	if services[baseURI] != nil {
		panic("service already registered")
	}

	return &Service{
		baseURI:  strings.TrimRight(baseURI, "/"),
		handlers: make(map[*regexp.Regexp]contextHandler),
	}
}

func addToChain(f interface{}, chain []contextHandler) []contextHandler {
	m := toContextHandler(f)

	if m == nil {
		panic(errors.New("unsupported middleware type"))
	}

	return append(chain, m)
}

func (s *Service) AddPre(f interface{}) {
	s.pre = addToChain(f, s.pre)
}

func (s *Service) AddPost(f interface{}) {
	s.post = addToChain(f, s.post)
}

// Service satisfies the http.Handler interface.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ServeHTTPInContext(NewSiestaContext(), w, r)
}

func (s *Service) ServeHTTPInContext(c Context, w http.ResponseWriter, r *http.Request) {
	quit := false
	for _, m := range s.pre {
		if quit {
			return
		}

		m(c, w, r, func() {
			quit = true
		})
	}

	if !quit {

		var handler contextHandler

		for re, h := range s.handlers {
			req := r.Method + " " + r.URL.Path

			if matches := re.FindStringSubmatch(req); len(matches) > 0 {
				r.ParseForm()
				for i, match := range matches {
					if i > 0 {
						param := re.SubexpNames()[i]
						r.Form.Set(param, match)
					}
				}

				handler = h
				break
			}
		}

		if handler == nil {
			http.NotFoundHandler().ServeHTTP(w, r)
		} else {
			handler(c, w, r, func() {
				quit = true
			})
		}
	}

	for _, m := range s.post {
		if quit {
			return
		}

		m(c, w, r, func() {
			quit = true
		})
	}
}

func (s *Service) Route(verb, pattern, usage string, f interface{}) {
	handler := toContextHandler(f)

	expr := strings.TrimRight(strings.TrimLeft(pattern, "/"), "/")
	expr = strings.Replace(expr, "<", "(?P<", -1)
	expr = strings.Replace(expr, ">", ">[\\d\\w\\-\\_]+)", -1)

	end := "?$"
	if len(expr) == 0 {
		end = "/?$"
	}

	if len(expr) > 0 {
		expr += "/"
		if s.baseURI != "/" {
			expr = "/" + expr
		}
	}

	expr = "^" + verb + " " + s.baseURI + expr + end
	re := regexp.MustCompile(expr)

	if _, ok := s.handlers[re]; ok {
		panic("already registered handler for " + verb + " " + pattern)
	} else {
		log.Println("Handling", expr)
	}

	s.handlers[re] = handler
}

func (s *Service) Register() {
	http.Handle(s.baseURI, s)
	http.Handle(s.baseURI+"/", s)
}
