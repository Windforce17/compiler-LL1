package main

import (
	"./analyzer"
	"./chart"
	"./first_set"
	"./follow_set"
	"./rule"
	"./util/feedback"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "html")
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		f, err := os.Open("./static/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		io.Copy(w, f)
	})
	mux.HandleFunc("/api/solve", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json
		fb := feedback.NewFeedBack(w)
		err := r.ParseForm()
		if err != nil {
			fb.Code(http.StatusInternalServerError).Msg(err.Error()).Response()
			return
		}
		raw := r.FormValue("grammar")
		if len(raw) <= 0 {
			fb.Code(http.StatusBadRequest).Msg("grammar is empty").Response()
			return
		}
		grammar := strings.Split(raw, "\n")
		rules := rule.NewRules()
		for i := range grammar {
			rules.AddRules(grammar[i])
		}
		firstSet := first_set.GetFirstFrom(rules)
		start := grammar[0][0]
		followSet := follow_set.GetFollowFrom(rules, start, firstSet)
		ch := chart.GetChartFrom(firstSet, followSet, rules)
		data := map[string]interface{}{
			"table":  ch.CoverToTable(),
			"first":  firstSet.Strings(),
			"follow": followSet.Strings(),
		}
		fb.Data(data).Response()
	})
	mux.HandleFunc("/api/input", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json
		fb := feedback.NewFeedBack(w)
		err := r.ParseForm()
		if err != nil {
			fb.Code(http.StatusInternalServerError).Msg(err.Error()).Response()
			return
		}
		raw := r.FormValue("grammar")
		if len(raw) <= 0 {
			fb.Code(http.StatusBadRequest).Msg("grammar is empty").Response()
			return
		}
		input := r.FormValue("input")
		grammar := strings.Split(raw, "\n")
		rules := rule.NewRules()
		for i := range grammar {
			rules.AddRules(grammar[i])
		}
		firstSet := first_set.GetFirstFrom(rules)
		start := grammar[0][0]
		followSet := follow_set.GetFollowFrom(rules, start, firstSet)
		ch := chart.GetChartFrom(firstSet, followSet, rules)
		step, err := analyzer.Analyze(ch, start, input)
		if err != nil {
			fb.Code(http.StatusBadRequest).Msg(err.Error()).Response()
			return
		}
		data := map[string]interface{}{
			"step": step,
		}
		fb.Data(data).Response()
	})
	http.ListenAndServe(":8081", mux)
}
