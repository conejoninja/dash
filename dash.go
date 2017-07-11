package dash

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"html/template"
	"log"
	"strings"

	"github.com/gobuffalo/plush"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

type dashConfig struct {
	apiProto, apiURI, apiPort, apiUser, apiPassword string
	webPort, webUser, webPassword, webProto, webURI string
	wsServer, wsPort                                   string
	webAuth                                             bool
}

var cfg dashConfig

var tmplCache map[string]string

const useCache = false

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dash := ps.ByName("id")
	period := ps.ByName("period")
	ctx := plush.NewContext()
	ctx.Set("abspath", cfg.webURI)
	ctx.Set("api_proto", cfg.apiProto)
	ctx.Set("api_uri", cfg.apiURI)
	ctx.Set("api_port", cfg.apiPort)
	ctx.Set("content", dash+"/content.html")
	ctx.Set("script", dash+"/script.js")
	ctx.Set("title", "Food-01")
	ctx.Set("period", period)

	ctx.Set("partial", partial(ctx))
	ctx.Set("url", url(ctx))

	s, err := plush.Render(loadTemplate("main.html"), ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, s)

}

func handlerDash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := plush.NewContext()
	ctx.Set("abspath", cfg.webURI)
	ctx.Set("content", "dash/content.html")
	ctx.Set("script", "dash/script.js")
	ctx.Set("ws_server", cfg.wsServer)
	ctx.Set("ws_port", cfg.wsPort)
	ctx.Set("title", "Dash")

	ctx.Set("partial", partial(ctx))

	s, err := plush.Render(loadTemplate("main.html"), ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, s)

}

func handlerAjaxGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := ps.ByName("request")
	resp, err := http.Get(apiPath() + req)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(body))
}

func handlerAjaxPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	device := ps.ByName("device")
	f := ps.ByName("f")

	r.ParseForm()
	resp, err := http.PostForm(apiPath()+"/call/"+device+"/"+f, r.Form)
	if err != nil {
		fmt.Println(err)
	}
	if resp == nil {
		fmt.Fprint(w, "{\"type\":\"error\",\"message\":\"Request failed\"}")
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	fmt.Fprint(w, string(body))
}

func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		user, password, hasAuth := r.BasicAuth()
		if !cfg.webAuth || (hasAuth && user == cfg.webUser && password == cfg.webPassword) {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func partial(ctx *plush.Context) func(string) (template.HTML, error) {
	return func(name string) (template.HTML, error) {
		t, err := plush.Render(loadTemplate(name), ctx)
		return template.HTML(t), err
	}
}

func url(ctx *plush.Context) func(string) (template.HTML, error) {
	return func(str string) (template.HTML, error) {
		str = strings.Trim(str, " ")
		return template.HTML(cfg.webURI + str), nil
	}
}

func loadTemplate(name string) string {
	if content, ok := tmplCache[name]; ok && useCache {
		return content
	}
	content, err := ioutil.ReadFile("./views/" + name)
	fmt.Println("loading:", "./views/"+name)
	if err != nil {
		fmt.Println(err, name)
		return ""
	}
	tmplCache[name] = string(content)
	return string(content)
}

func apiPath() string {
	if (cfg.apiProto == "http" && cfg.apiPort == "80") || (cfg.apiProto == "https" && cfg.apiPort == "443") {
		return cfg.apiProto + "://" + cfg.apiURI
	}
	return cfg.apiProto + "://" + cfg.apiURI + ":" + cfg.apiPort
}

func Start() {
	tmplCache = make(map[string]string)
	cfg = readConfig()

	router := httprouter.New()
	router.GET("/dash/:id", auth(handler))
	router.GET("/dash/:id/:period", auth(handler))
	router.GET("/", auth(handlerDash))
	router.GET("/ajax/*request", auth(handlerAjaxGet))
	router.POST("/ajax/:device/:f", auth(handlerAjaxPost))
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	fmt.Println("Server Up")
	log.Fatal(http.ListenAndServe(":"+cfg.webPort, router))
	fmt.Println("Server Down")

}

func readConfig() (cfg dashConfig) {
	if _, err := os.Stat("./config.yml"); err != nil {
		fmt.Println("Error: config.yml file does not exist")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	var au, ap, wu, wp bool
	cfg.apiUser, au = os.LookupEnv("API_USER")
	cfg.apiPassword, ap = os.LookupEnv("API_PASSWORD")
	cfg.apiProto = os.Getenv("API_PROTO")
	cfg.apiURI = os.Getenv("API_URI")
	cfg.apiPort = os.Getenv("API_PORT")

	webAuthStr := os.Getenv("WEB_AUTH")
	cfg.webUser, wu = os.LookupEnv("WEB_USER")
	cfg.webPassword, wp = os.LookupEnv("WEB_PASSWORD")
	cfg.webPort = os.Getenv("WEB_PORT")
	cfg.webProto = os.Getenv("WEB_PROTO")
	cfg.webURI = os.Getenv("WEB_URI")

	cfg.wsServer = os.Getenv("WS_SERVER")
	cfg.wsPort = os.Getenv("WS_PORT")

	if cfg.apiProto == "" {
		cfg.apiProto = fmt.Sprint(viper.Get("api_proto"))
	}
	if cfg.apiURI == "" {
		cfg.apiURI = fmt.Sprint(viper.Get("api_uri"))
	}
	if cfg.apiPort == "" {
		cfg.apiPort = fmt.Sprint(viper.Get("api_port"))
	}
	if !au {
		cfg.apiUser = fmt.Sprint(viper.Get("api_user"))
	}
	if !ap {
		cfg.apiPassword = fmt.Sprint(viper.Get("api_password"))
	}
	if cfg.apiPort == "" {
		cfg.apiPort = "8080"
	}
	if cfg.apiProto == "" {
		cfg.apiProto = "http"
	}

	if cfg.webPort == "" {
		cfg.webPort = fmt.Sprint(viper.Get("web_port"))
	}
	if webAuthStr == "" {
		webAuthStr = fmt.Sprint(viper.Get("web_auth"))
	}
	if cfg.webProto == "" {
		cfg.webProto = fmt.Sprint(viper.Get("web_proto"))
	}
	if cfg.webURI == "" {
		cfg.webURI = fmt.Sprint(viper.Get("web_uri"))
	}
	if !wu {
		cfg.webUser = fmt.Sprint(viper.Get("web_user"))
	}
	if !wp {
		cfg.webPassword = fmt.Sprint(viper.Get("web_password"))
	}
	if cfg.webPort == "" {
		cfg.webPort = "1313"
	}

	if webAuthStr == "1" || webAuthStr == "true" {
		cfg.webAuth = true
	}
	if cfg.wsServer == "" {
		cfg.wsServer = fmt.Sprint(viper.Get("ws_server"))
	}
	if cfg.wsPort == "" {
		cfg.wsPort = fmt.Sprint(viper.Get("ws_port"))
	}
	if cfg.wsPort == "" {
		cfg.wsPort = "8055"
	}

	return

}
