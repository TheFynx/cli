package theme

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// meteorSVG is the authored approved meteor glyph; the accent color is
// substituted at build time.
//
//go:embed assets/midgar-meteor.svg
var meteorSVG string

// buildMidgar regenerates the midgar-mako theme assets.
func buildMidgar(home string) error {
	slug := "midgar-mako"
	name := "Midgar-Mako"
	c, e := load(home, slug)
	if e != nil {
		return e
	}
	a, cy := c.Accent, c.Cyan
	bg, fg := c.Background, c.Foreground
	panel, dark := c.LighterBackground, c.DarkBackground
	themeDir := filepath.Join(home, ".config/omarchy/themes", slug)
	assets := filepath.Join(themeDir, "assets")

	fmt.Println(slug + ":")
	if e := os.MkdirAll(assets, 0o755); e != nil {
		return e
	}
	meteor := strings.ReplaceAll(meteorSVG, "#63e6be", a)
	// Materialia orbs, one per FF7 materia class color; menu gets the
	// meteor glyph and sword stands alone. folder.svg is deliberately
	// not written to assets (matching the original builder).
	assetsMap := map[string]string{"menu": meteor}
	// Materia orbs, one per FF7 materia class color.
	for _, m := range []struct{ kind, color string }{
		{"magic", "green"}, {"support", "blue"}, {"command", "yellow"},
		{"summon", "red"}, {"independent", "magenta"},
	} {
		color := mapColor(c, m.color)
		body := `<defs><radialGradient id="g" cx="30%" cy="25%"><stop stop-color="#f0faf7"/><stop offset=".25" stop-color="` +
			color + `"/><stop offset="1" stop-color="` + dark + `"/></radialGradient></defs>` +
			`<circle cx="32" cy="32" r="25" fill="url(#g)" stroke="` + color + `" stroke-width="2"/>` +
			`<ellipse cx="24" cy="21" rx="8" ry="5" fill="#fff" opacity=".45"/>`
		assetsMap[m.kind] = svg(body)
	}
	assetsMap["folder"] = svg(`<path d="M5 15h20l6 7h28v33H5z" fill="` + panel + `" stroke="` + a + `" stroke-width="2"/>` +
		`<path d="M5 28h54" stroke="` + cy + `"/>` +
		`<circle cx="32" cy="41" r="9" fill="` + a + `" opacity=".7"/><circle cx="29" cy="38" r="3" fill="` + fg + `"/>`)
	assetsMap["sword"] = svg(`<path d="M38 3l14 9-24 32-13-9z" fill="` + fg + `" stroke="` + cy + `" stroke-width="2"/>` +
		`<path d="M11 34l19 13M21 41L9 58" stroke="` + a + `" stroke-width="5"/>` +
		`<circle cx="33" cy="17" r="3" fill="` + a + `"/><circle cx="29" cy="24" r="3" fill="` + c.Magenta + `"/>`)
	for key, body := range assetsMap {
		if key == "folder" {
			continue // used for icon mapping only, never written
		}
		if e := write(filepath.Join(assets, key+".svg"), body); e != nil {
			return e
		}
	}
	// Icon theme: all three context dirs exist even when status has no
	// mappings (matching the original builder's mkdirs).
	iconroot := filepath.Join(home, ".local/share/icons", name)
	for _, sub := range []string{"places", "apps", "status"} {
		if e := os.MkdirAll(filepath.Join(iconroot, "scalable", sub), 0o755); e != nil {
			return e
		}
	}
	index := "[Icon Theme]\nName=Midgar Mako\nInherits=Yaru-sage,Adwaita,hicolor\nDirectories=scalable/places,scalable/apps,scalable/status\n\n"
	for _, sub := range []struct{ dir, ctx string }{
		{"scalable/places", "Places"}, {"scalable/apps", "Applications"}, {"scalable/status", "Status"},
	} {
		index += "[" + sub.dir + "]\nSize=64\nType=Scalable\nMinSize=16\nMaxSize=512\nContext=" + sub.ctx + "\n\n"
	}
	if e := write(filepath.Join(iconroot, "index.theme"), index); e != nil {
		return e
	}
	for _, key := range []string{"folder", "folder-documents", "folder-download", "folder-pictures", "folder-videos", "folder-music", "user-home"} {
		if e := write(filepath.Join(iconroot, "scalable/places", key+".svg"), assetsMap["folder"]); e != nil {
			return e
		}
	}
	apps := map[string]string{
		"brave-browser": meteor, "brave": meteor, "code": assetsMap["sword"],
		"com.mitchellh.ghostty": assetsMap["magic"], "utilities-terminal": assetsMap["magic"],
		"org.gnome.Nautilus": assetsMap["folder"], "system-file-manager": assetsMap["folder"],
		"org.blender.Blender": assetsMap["command"],
	}
	for file, body := range apps {
		if e := write(filepath.Join(iconroot, "scalable/apps", file+".svg"), body); e != nil {
			return e
		}
	}
	if e := write(filepath.Join(themeDir, "icons.theme"), name+"\n"); e != nil {
		return e
	}
	// GTK.
	defines := []string{
		"accent_color", a, "accent_bg_color", a, "accent_fg_color", bg,
		"window_bg_color", bg, "window_fg_color", fg,
		"view_bg_color", dark, "view_fg_color", fg,
		"headerbar_bg_color", panel, "headerbar_fg_color", fg,
		"card_bg_color", panel, "card_fg_color", fg,
		"dialog_bg_color", bg, "dialog_fg_color", fg,
		"popover_bg_color", panel, "popover_fg_color", fg,
		"sidebar_bg_color", dark, "sidebar_fg_color", fg,
		"theme_bg_color", bg, "theme_fg_color", fg,
		"theme_base_color", dark, "theme_text_color", fg,
		"theme_selected_bg_color", c.Selection, "theme_selected_fg_color", fg,
		"warning_color", c.Yellow, "error_color", c.Red, "success_color", a,
	}
	lines := make([]string, 0, len(defines)/2)
	for i := 0; i < len(defines); i += 2 {
		lines = append(lines, fmt.Sprintf("@define-color %s %s;\n", defines[i], defines[i+1]))
	}
	css := strings.Join(lines, "") // each line carries its own newline
	css += midgarCSSBody(a, cy, bg, fg, panel, dark, c.Muted, c.Selection, c.BrightForeground)
	if e := write(filepath.Join(themeDir, "gtk.css"), css); e != nil {
		return e
	}
	// Qt color roles (same list three times, as before, in #ffARGB form)
	// and stylesheet.
	roles := []string{fg, panel, panel, dark, c.Muted, fg, fg, fg, fg, bg, bg, dark,
		c.Selection, c.BrightForeground, cy, c.Magenta, bg, bg, c.Yellow, fg, c.Muted, dark}
	for i := range roles {
		roles[i] = "#ff" + roles[i][1:]
	}
	qtc := "[ColorScheme]\n"
	for _, k := range []string{"active_colors", "inactive_colors", "disabled_colors"} {
		qtc += k + "=" + strings.Join(roles, ",") + "\n"
	}
	if e := write(filepath.Join(themeDir, "qt.colors"), qtc); e != nil {
		return e
	}
	qss := "QPushButton, QToolButton { border: 1px solid " + c.Muted + "; border-radius: 4px; padding: 5px; background: " + panel + "; color: " + fg + "; }\n" +
		"QPushButton:hover, QToolButton:hover { border-color: " + cy + "; }\n" +
		"QLineEdit, QComboBox { border: 1px solid " + c.Muted + "; border-radius: 4px; background: " + dark + "; padding: 4px; }\n" +
		"QLineEdit:focus { border-color: " + a + "; }\n" +
		"QProgressBar::chunk { background: " + a + "; }\n" +
		"QMenu { background: " + panel + "; color: " + fg + "; border: 1px solid " + a + "; }\n" +
		"QMenu::item:selected { background: " + c.Selection + "; }\n" +
		"QHeaderView::section { background: " + panel + "; color: " + fg + "; border: 1px solid " + c.Muted + "; padding: 5px; }\n"
	if e := write(filepath.Join(themeDir, "qt.qss"), qss); e != nil {
		return e
	}
	fmt.Println("Built Midgar materia, meteor menu, icons, GTK and Qt assets.")
	return nil
}

func mapColor(c Colors, key string) string {
	switch key {
	case "green":
		return c.Green
	case "blue":
		return c.Blue
	case "yellow":
		return c.Yellow
	case "red":
		return c.Red
	case "magenta":
		return c.Magenta
	}
	return c.Accent
}

func svg(body string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">` + body + "</svg>\n"
}

func midgarCSSBody(a, cy, bg, fg, panel, dark, muted, selection, bright string) string {
	return "window, dialog, .background { background-color: " + bg + "; color: " + fg + "; }\n" +
		"headerbar, .titlebar { background: " + panel + "; color: " + fg + "; border-bottom: 1px solid " + a + "; }\n" +
		"button { border-radius: 4px; background-image: none; background-color: " + panel + "; color: " + fg + "; border: 1px solid " + muted + "; }\n" +
		"button:hover { border-color: " + cy + "; background-color: " + selection + "; }\n" +
		"button:checked, button.suggested-action { background: " + selection + "; color: " + fg + "; border-color: " + a + "; }\n" +
		"entry, textview, .view { background: " + dark + "; color: " + fg + "; border-radius: 4px; }\n" +
		"entry:focus { border-color: " + a + "; box-shadow: 0 0 0 1px " + a + "; }\n" +
		"selection, row:selected { background-color: " + selection + "; color: " + bright + "; }\n" +
		"progressbar progress, scale highlight { background: " + a + "; border-radius: 4px; }\n" +
		"scale slider { background: " + cy + "; }\n" +
		"scrollbar slider { background: " + muted + "; min-width: 8px; min-height: 8px; }\n" +
		"popover, menu { background: " + panel + "; color: " + fg + "; border: 1px solid " + a + "; border-radius: 4px; }\n"
}
