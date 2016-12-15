package darwini_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/tv42/darwini"
)

var (
	ErrNotFound = errors.New("not found")
)

type gistStore struct {
	gists map[uint64]*gist
}

func (s *gistStore) list(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(s.gists)
}

func (s *gistStore) listPublic(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(s.gists)
}

func (s *gistStore) listStarred(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(s.gists)
}

func (s *gistStore) create(w http.ResponseWriter, req *http.Request) {
	i := uint64(len(s.gists)) + 1
	g := &gist{ID: uint64(i), store: s}
	s.gists[i] = g
	json.NewEncoder(w).Encode(i)
}

func (s *gistStore) get(seg string) (*gist, error) {
	id, err := strconv.ParseUint(seg, 10, 64)
	if err != nil {
		return nil, err
	}
	g := s.gists[id]
	if g == nil {
		return nil, ErrNotFound
	}
	return g, nil
}

type gist struct {
	ID    uint64
	Text  string
	store *gistStore
	Star  bool
}

func (g *gist) get(w http.ResponseWriter, req *http.Request)   {}
func (g *gist) patch(w http.ResponseWriter, req *http.Request) {}
func (g *gist) del(w http.ResponseWriter, req *http.Request)   {}
func (g *gist) isStar(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(g.Star)
}
func (g *gist) star(w http.ResponseWriter, req *http.Request) {
	g.Star = true
}
func (g *gist) unstar(w http.ResponseWriter, req *http.Request) {
	g.Star = false
}
func (g *gist) forks(w http.ResponseWriter, req *http.Request) {
}

func removeSlash(w http.ResponseWriter, req *http.Request) {}

/*
GET /gists
GET /gists/public
GET /gists/starred
GET /gists/:id
POST /gists
PATCH /gists/:id
PUT /gists/:id/star
DELETE /gists/:id/star
GET /gists/:id/star
POST /gists/:id/forks
DELETE /gists/:id
*/

func Example() {
	gists := &gistStore{
		gists: map[uint64]*gist{},
	}
	m := darwini.Map{
		"gists": darwini.Dir{
			Parent: darwini.Method{
				GET:  gists.list,
				POST: gists.create,
			},
			Child: darwini.Var{
				Child: func(seg string) http.Handler {
					// Mixing dynamic and static segments is just bad,
					// so we won't bother to assist in that. Write code.
					switch seg {
					case "public":
						return http.HandlerFunc(gists.listPublic)
					case "starred":
						return http.HandlerFunc(gists.listStarred)
					}
					g, err := gists.get(seg)
					if err != nil {
						return nil // not found
					}
					return darwini.Dir{
						Parent: darwini.Method{
							GET:    g.get,
							PATCH:  g.patch,
							DELETE: g.del,
						},
						Child: darwini.Map{
							"": http.HandlerFunc(removeSlash),
							"star": darwini.Method{
								GET:    g.isStar,
								PUT:    g.star,
								DELETE: g.unstar,
							},
							"forks": darwini.Method{
								GET: g.forks,
							},
						},
					}
				},
			},
		},
	}
	s := httptest.NewServer(m)
	defer s.Close()

	resp, err := http.Post(s.URL+"/gists", "text/plain", strings.NewReader("hello, world"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return
	}
	var id uint64
	if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("created gist", id)

	resp, err = http.Get(s.URL + "/gists/" + strconv.FormatUint(id, 10) + "/star")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return
	}
	var starred bool
	if err := json.NewDecoder(resp.Body).Decode(&starred); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("gist %d is starred: %v\n", id, starred)

	// Output:
	// created gist 1
	// gist 1 is starred: false
}
