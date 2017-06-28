package dash

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gobuffalo/plush"
	"github.com/spf13/viper"
	"html/template"
)

var api_proto, api_uri, api_port, api_user, api_password, web_port string
var tmplCache map[string]string

const useCache = false

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := plush.NewContext()
	ctx.Set("api_proto", api_proto)
	ctx.Set("api_uri", api_uri)
	ctx.Set("api_port", api_port)
	ctx.Set("content", "food01/content.html")

	ctx.Set("partial", partial(ctx))

	s, err := plush.Render(loadTemplate("main.html"), ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, s)

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
	api_proto, api_uri, api_port, api_user, api_password, web_port = readConfig()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	fmt.Println("Server Up")
	http.ListenAndServe(":"+web_port, nil)
	fmt.Println("Server Down")
}

func readConfig() (api_proto, api_uri, api_port, api_user, api_password, web_port string) {
	if _, err := os.Stat("./config.yml"); err != nil {
		fmt.Println("Error: config.yml file does not exist")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	var au, ap bool
	api_user, au = os.LookupEnv("API_USER")
	api_password, ap = os.LookupEnv("API_PASSWORD")
	api_proto = os.Getenv("API_PROTO")
	api_uri = os.Getenv("API_URI")
	api_port = os.Getenv("API_PORT")
	web_port = os.Getenv("WEB_PORT")

	if api_proto == "" {
		api_proto = fmt.Sprint(viper.Get("api_proto"))
	}
	if api_uri == "" {
		api_uri = fmt.Sprint(viper.Get("api_uri"))
	}
	if api_port == "" {
		api_port = fmt.Sprint(viper.Get("api_port"))
	}
	if !au {
		api_user = fmt.Sprint(viper.Get("api_user"))
	}
	if !ap {
		api_password = fmt.Sprint(viper.Get("api_password"))
	}

	if web_port == "" {
		web_port = "1313"
	}
	if api_proto == "" {
		api_proto = "http"
	}

	return

}
