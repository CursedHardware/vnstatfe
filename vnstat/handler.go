package vnstat

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
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
		Funcs(template.FuncMap{"vnstati": handler.vnstati}).
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
	} else if r.URL.Path == "/vnstat.json" {
		header.Set("Content-Type", "text/json")
		query := r.URL.Query()
		callback := query.Get("callback")
		begin, _ := time.Parse(time.DateOnly, query.Get("begin"))
		end, _ := time.Parse(time.DateOnly, query.Get("end"))
		h.ServeJSON(w, callback, begin, end)
	} else if pathname, ok := strings.CutPrefix(r.URL.Path, "/vnstati/"); ok {
		header.Set("Content-Type", "image/png")
		header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		header.Set("Pragma", "no-cache")
		var view, scale string
		var scaling int
		if view, scale, ok = strings.Cut(pathname, "@"); ok {
			scaling, _ = strconv.Atoi(scale)
		} else {
			view = pathname
		}
		h.ServeImage(w, view, min(max(100, scaling), 500))
	} else {
		http.ServeFileFS(w, r, assets, path.Join("assets", r.URL.Path))
	}
}

//goland:noinspection SpellCheckingInspection
func (h *Handler) vnstati(view string) template.HTML {
	img := &Image{Program: "vnstati", View: view}
	return img.HTML(2, 3, 4, 5)
}

func (h *Handler) ServeJSON(w http.ResponseWriter, callback string, begin, end time.Time) {
	cmd := exec.Command("vnstat", "--json")
	if !begin.IsZero() {
		cmd.Args = append(cmd.Args, "--begin", begin.Format(time.DateOnly))
	}
	if !end.IsZero() {
		cmd.Args = append(cmd.Args, "--end", end.Format(time.DateOnly))
	}
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var buf bytes.Buffer
	if len(callback) > 0 {
		buf.WriteString(callback)
		buf.WriteByte('(')
		buf.Write(output)
		buf.WriteByte(')')
		buf.WriteByte(';')
	} else {
		buf.Write(output)
	}
	_, _ = buf.WriteTo(w)
}

func (h *Handler) ServeImage(w http.ResponseWriter, view string, scale int) {
	cmd := exec.Command("vnstati")
	cmd.Stdout = w
	cmd.Args = slices.Concat(h.Arguments, []string{
		"--large",
		"--transparent",
		"--noedge",
		"--scale", strconv.Itoa(scale),
		"--output", "-",
	})
	if v, ok := views[view]; ok {
		cmd.Args = append(cmd.Args, v)
	}
	log.Println(cmd)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}
}
