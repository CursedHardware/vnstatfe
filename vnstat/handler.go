package vnstat

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"path"
	"slices"
	"strconv"
	"strings"
)

//go:embed assets
var assets embed.FS

//goland:noinspection SpellCheckingInspection
var views = map[string]string{
	"5min":         "--fiveminutes",
	"5min-graph":   "--fivegraph",
	"hourly":       "--hours",
	"hourly-graph": "--hoursgraph",
	"daily":        "--days",
	"monthly":      "--months",
	"yearly":       "--years",
	"top":          "--top",
	"summary":      "--summary",
	"hsummary":     "--hsummary",
	"vsummary":     "--vsummary",
}

type Handler struct {
	Template  *template.Template
	Arguments []string
}

func NewHandler(arguments []string) (http.Handler, error) {
	var err error
	handler := new(Handler)
	handler.Arguments = arguments
	handler.Template, err = template.New("vnstat").
		Funcs(template.FuncMap{"vnstati": handler.makeURL}).
		ParseFS(assets, "assets/*")
	return handler, err
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	if r.URL.Path == "/" {
		header.Set("Content-Type", "text/html")
		_ = h.Template.ExecuteTemplate(w, "index.gohtml", nil)
	} else if pathname, ok := strings.CutPrefix(r.URL.Path, "/vnstati/"); ok {
		header.Set("Content-Type", "image/png")
		header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		header.Set("Pragma", "no-cache")
		var view, scale string
		var scaling float64
		if view, scale, ok = strings.Cut(pathname, "@"); ok {
			scaling, _ = strconv.ParseFloat(strings.TrimSuffix(scale, "x"), 10)
		} else {
			view = pathname
		}
		scaling = min(max(1, scaling), 5)
		h.ServeImage(w, view, scaling)
	} else {
		http.ServeFileFS(w, r, assets, path.Join("assets", r.URL.Path))
	}
}

//goland:noinspection SpellCheckingInspection
func (h *Handler) makeURL(view string, scale int) template.URL {
	if scale != 1 {
		view = fmt.Sprintf("%s@%dx", view, scale)
	}
	return template.URL(path.Join("vnstati", view))
}

func (h *Handler) ServeImage(w http.ResponseWriter, view string, scale float64) {
	cmd := exec.Command("vnstati")
	cmd.Args = slices.Concat(h.Arguments, []string{"--large", "--output", "-"})
	if v, ok := views[view]; ok {
		cmd.Args = append(cmd.Args, v)
	}
	if scale != 1 {
		cmd.Args = append(cmd.Args, "--scale", strconv.Itoa(int(scale*100)))
	}
	log.Println(cmd)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		_, _ = w.Write(output)
	}
}
