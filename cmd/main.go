package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/deckarep/golang-set"
  "io/ioutil"
  "strings"
  "regexp"
  "encoding/json"
)

type Comparison struct {
	JaccardDistance float64
	SiteA string
  SiteB string
}

// JaccardSimilarity, as known as the Jaccard Index, compares the similarity of sample sets.
// This doesn't measure similarity between texts, but if regarding a text as bag-of-word,
// it can apply.
func JaccardSimilarity(s1, s2 string, f func(string) mapset.Set) float64 {
	if s1 == s2 {
		return 1.0
	}
	if f == nil {
		f = convertStringToSet
	}
	s1set := f(s1)
	s2set := f(s2)
	s1ands2 := s1set.Intersect(s2set).Cardinality()
	s1ors2 := s1set.Union(s2set).Cardinality()
	return float64(s1ands2) / float64(s1ors2)
}

func convertStringToSet(s string) mapset.Set {
	set := mapset.NewSet()
	for _, token := range strings.Fields(s) {
		set.Add(token)
	}
	return set
}

func distance(w http.ResponseWriter, r *http.Request) {

  //urls := strings.Split(r.FormValue("urls"),"\n")
  if err := r.ParseForm(); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  urls := r.Form["urls[]"]

  if len(urls) < 1 {
    http.Error(w, "no urls sent", http.StatusInternalServerError)
  } else if len(urls) < 2 {
    urls = append(urls, "http://www.safer-shopping.de")
  }

  client := &http.Client{}
  sites := make(map[string]string)
  rex, _ := regexp.Compile("[ \t\r\n]+")
  for idx := range urls {
    //solves NXDomain Error
    url := strings.Trim(urls[idx], "\r\n ")
    if url == "" {
      continue
    }
    //TODO: Cache Requests
    req, _ := http.NewRequest("GET", url, nil)
    res, err := client.Do(req)
    //resp, err := http.Get(urls[idx])

    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      //TODO: Cache error message
      continue
    }

    body, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close()

    bodyS := string(body)

    //fmt.Fprintf(w, "%d\t%s\t%d\t%d\n", idx, url, len(bodyS), len(sites))
    sites[url] = rex.ReplaceAllString(bodyS, " ")
  }

  var simi []Comparison
  for idxA := range sites {
    for idxB := range sites {

      if idxA == idxB {
        continue
      }
      //TODO: idxB and idxA already in simi

      dist := JaccardSimilarity(sites[idxA], sites[idxB], convertStringToSet)
      simi = append(simi, Comparison{JaccardDistance: dist, SiteA: idxA, SiteB: idxB})
      //fmt.Fprintf(w, "%f\t%s\t%s\n", dist, idxA, idxB)
    }
  }

  b, err := json.Marshal(simi)
  if err != nil {
	  fmt.Fprintf(w, "execution failed: %s", err)
  }
  fmt.Fprint(w, string(b))
}

func handler(w http.ResponseWriter, r *http.Request) {
  t, _:= template.New("t_start").ParseFiles("start.tmpl")
  err := t.Execute(w, nil)
  if err != nil {
	  fmt.Fprintf(w, "execution failed: %s", err)
  }
}

func main() {
  http.HandleFunc("/distance", distance )
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8081", nil)
}
