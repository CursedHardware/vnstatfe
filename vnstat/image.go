package vnstat

import (
	"fmt"
	"html/template"
	"path"
	"sort"
	"strings"
)

type Image struct {
	Program string
	View    string
}

func (i *Image) Alt() string {
	//goland:noinspection SpellCheckingInspection
	switch strings.ToLower(i.View) {
	case "5min":
		return "Output traffic statistics with a 5 minute resolution for the last hours."
	case "5min-graph":
		return "Output traffic statistics with a 5 minute resolution for the last 48 hours using a bar graph."
	case "hourly":
		return "Output traffic statistics on a hourly basis."
	case "hourly-graph":
		return "Output traffic statistics on a hourly basis for the last 24 hours using a bar graph."
	case "daily":
		return "Output traffic statistics on a daily basis for the last days."
	case "monthly":
		return "Output traffic statistics on a monthly basis for the last months."
	case "yearly":
		return "Output traffic statistics on a yearly basis for the last years."
	case "top":
		return "Output all time top traffic days."
	case "summary":
		return "Output traffic statistics summary."
	case "hsummary":
		return "Output traffic summary including hourly data using a horizontal layout."
	case "vsummary":
		return "Output traffic summary including hourly data using a vertical layout."
	}
	return i.View
}

func (i *Image) HTML(scales ...int) template.HTML {
	sort.Ints(scales)
	view := path.Join(i.Program, i.View)
	var srcset []string
	for _, scale := range scales {
		srcset = append(srcset, fmt.Sprintf("%s@%d %dx", view, scale*100, scale))
	}
	return template.HTML(fmt.Sprintf(
		"<img src=%q srcset=%q alt=%q/>",
		view,
		strings.Join(srcset, ","),
		i.Alt(),
	))
}
