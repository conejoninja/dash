package dash

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gobuffalo/plush"
	"github.com/spf13/viper"
	"html/template"
	"github.com/julienschmidt/httprouter"
	"log"
)


type dashConfig struct {
	api_proto, api_uri, api_port, api_user, api_password, web_port, web_user, web_password, ws_server, ws_port string
	web_auth bool
}
var cfg dashConfig

var tmplCache map[string]string

const useCache = false

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dash := ps.ByName("id")
	ctx := plush.NewContext()
	ctx.Set("abspath", "http://localhost:"+cfg.web_port)
	ctx.Set("api_proto", cfg.api_proto)
	ctx.Set("api_uri", cfg.api_uri)
	ctx.Set("api_port", cfg.api_port)
	ctx.Set("content", dash + "/content.html")
	ctx.Set("script", dash + "/script.js")
	ctx.Set("title", "Food-01")


	ctx.Set("partial", partial(ctx))

	s, err := plush.Render(loadTemplate("main.html"), ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, s)

}

func handlerDash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := plush.NewContext()
	ctx.Set("abspath", "http://localhost:"+cfg.web_port)
	ctx.Set("content", "dash/content.html")
	ctx.Set("script", "dash/script.js")
	ctx.Set("ws_server", cfg.ws_server)
	ctx.Set("ws_port", cfg.ws_port)
	ctx.Set("title", "Dash")

	ctx.Set("partial", partial(ctx))

	s, err := plush.Render(loadTemplate("main.html"), ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, s)

}


func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		user, password, hasAuth := r.BasicAuth()
		if !cfg.web_auth || (hasAuth && user == cfg.web_user && password == cfg.web_password) {
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

func loadTemplate(name string) string {
	if content, ok := tmplCache[name]; ok && useCache {
		return content
	}
	content, err := ioutil.ReadFile("views/" + name)
	if err != nil {
		fmt.Println(err, name)
		return ""
	}
	tmplCache[name] = string(content)
	return string(content)
}

func Start() {
	tmplCache = make(map[string]string)
	cfg = readConfig()


	router := httprouter.New()
	router.GET("/dash/:id", auth(handler))
	router.GET("/", auth(handlerDash))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	
	fmt.Println("Server Up")
	log.Fatal(http.ListenAndServe(":"+cfg.web_port, router))
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
	cfg.api_user, au = os.LookupEnv("API_USER")
	cfg.api_password, ap = os.LookupEnv("API_PASSWORD")
	cfg.api_proto = os.Getenv("API_PROTO")
	cfg.api_uri = os.Getenv("API_URI")
	cfg.api_port = os.Getenv("API_PORT")

	web_auth_str := os.Getenv("WEB_AUTH")
	cfg.web_user, wu = os.LookupEnv("WEB_USER")
	cfg.web_password, wp = os.LookupEnv("WEB_PASSWORD")
	cfg.web_port = os.Getenv("WEB_PORT")

	cfg.ws_server = os.Getenv("WS_SERVER")
	cfg.ws_port = os.Getenv("WS_PORT")


	if cfg.api_proto == "" {
		cfg.api_proto = fmt.Sprint(viper.Get("api_proto"))
	}
	if cfg.api_uri == "" {
		cfg.api_uri = fmt.Sprint(viper.Get("api_uri"))
	}
	if cfg.api_port == "" {
		cfg.api_port = fmt.Sprint(viper.Get("api_port"))
	}
	if !au {
		cfg.api_user = fmt.Sprint(viper.Get("api_user"))
	}
	if !ap {
		cfg.api_password = fmt.Sprint(viper.Get("api_password"))
	}
	if cfg.api_port == "" {
		cfg.api_port = "8080"
	}
	if cfg.api_proto == "" {
		cfg.api_proto = "http"
	}



	if cfg.web_port == "" {
		cfg.web_port = fmt.Sprint(viper.Get("web_port"))
	}
	if web_auth_str == "" {
		web_auth_str = fmt.Sprint(viper.Get("web_auth"))
	}
	if !wu {
		cfg.web_user = fmt.Sprint(viper.Get("web_user"))
	}
	if !wp {
		cfg.web_password = fmt.Sprint(viper.Get("web_password"))
	}
	if cfg.web_port == "" {
		cfg.web_port = "1313"
	}


	if web_auth_str == "1" || web_auth_str == "true" {
		cfg.web_auth = true
	}
	if cfg.ws_server == "" {
		cfg.ws_server = fmt.Sprint(viper.Get("ws_server"))
	}
	if cfg.ws_port == "" {
		cfg.ws_port = fmt.Sprint(viper.Get("ws_port"))
	}
	if cfg.ws_port == "" {
		cfg.ws_port = "8055"
	}

	return

}
