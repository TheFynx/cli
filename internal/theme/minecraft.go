package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// minecraftVariants are the three installed omarchy slugs and the icon
// themes they inherit from, matching the original builder's fallbacks.
var minecraftVariants = []struct{ slug, inherits string }{
	{"minecraft-deepstone", "Yaru-sage"},
	{"minecraft-nether", "Yaru-red"},
	{"minecraft-end", "Yaru-purple"},
}

// minecraftShapes are native 32x32 pixel icons; {a} is the accent color
// and {c} the cyan. Shape-rendering keeps edges crisp at any scale.
var minecraftShapes = []struct{ key, shape string }{
	{"grass", `<polygon points="2,11 16,3 30,11 16,19" fill="{a}"/><polygon points="2,11 16,19 16,30 2,22" fill="#70513a"/><polygon points="16,19 30,11 30,22 16,30" fill="#493626"/><path d="M2 11v5l4 2v-3l4 2v3l6 3v-4" fill="#6ca64c"/>`},
	{"diamond", `<path d="M12 3h8v4h5v6h4v7h-5v5h-6v4h-7v-4H6v-6H3v-7h5V7h4z" fill="{c}"/><path d="M12 7h8v4h-8v8H8v-8h4z" fill="#c5fff0"/><path d="M18 22h7v-7h-4v4h-3z" fill="#328f95"/>`},
	{"pickaxe", `<path d="M7 4h16v4h5v5h-5v-3H12v4H7z" fill="{c}"/><path d="M15 10h5v6h-5v5h-5v6H5v-5h5v-6h5z" fill="#a17a46"/>`},
	{"sword", `<path d="M24 2h6v6h-4v4h-4v4h-4v4h-5v-5h4v-4h4V7h3z" fill="{c}"/><path d="M9 13h4v4h4v4h-4v-4H9z" fill="#e5c777"/><path d="M8 20h5v5H8v5H3v-5h5z" fill="#a17a46"/>`},
	{"chest", `<path d="M3 8h26v20H3z" fill="#986431"/><path d="M3 8h26v7H3z" fill="#bb8b45"/><path d="M3 15h26M5 8v20M27 8v20" stroke="#432d20" stroke-width="2"/><path d="M14 13h5v8h-5z" fill="#ebc76c"/>`},
	{"crafting", `<path d="M3 6h26v23H3z" fill="#785137"/><path d="M5 7h22v13H5z" fill="#c2996a"/><path d="M12 7v13M20 7v13M5 13h22" stroke="#624129" stroke-width="2"/><path d="M7 23h18v3H7z" fill="#30271e"/>`},
	{"compass", `<path d="M8 3h16v4h5v17h-5v5H8v-5H3V8h5z" fill="#9da6a2"/><path d="M9 7h14v4h3v12h-4v3H10v-3H6V11h3z" fill="#24312c"/><path d="M18 6l-7 13h7z" fill="#e87b6b"/><path d="M11 19h7l-4 7z" fill="#e2e7d7"/>`},
	{"torch", `<path d="M13 12h7v18h-7z" fill="#987045"/><path d="M9 4h15v13H9z" fill="#ed9c43"/><path d="M12 1h9v12h-9z" fill="#f8d56c"/><path d="M15 4h3v8h-3z" fill="#fff1b4"/>`},
	{"heart", `<path d="M2 9h3V6h7v3h8V6h7v3h3v10h-4v4h-4v4h-4v3h-4v-3h-4v-4H6v-4H2z" fill="#dc6963"/><path d="M6 10h5v4H6z" fill="#ffb0a0"/>`},
	{"folder", `<path d="M2 7h11l3 4h14v17H2z" fill="{a}"/><path d="M2 13h28v15H2z" fill="#647d4c"/><path d="M3 15h26v3H3z" fill="{a}"/><path d="M6 20h5v4H6zM19 20h5v4h-5z" fill="#263a24"/>`},
	{"creeper", `<path d="M3 3h26v26H3z" fill="{a}"/><path d="M7 8h6v6H7zM19 8h6v6h-6zM13 14h6v10h-6zM10 18h3v9h-3zM19 18h3v9h-3z" fill="#19311f"/>`},
	{"food", `<path d="M5 7h18v4h5v12h-5v5H8v-4H3V12h2z" fill="#aa6b38"/><path d="M8 11h14v4H8zM8 19h17v3H8z" fill="#deb172"/>`},
}

// minecraftMappings map icon-theme file names to generated assets.
var minecraftMappings = map[string]map[string]string{
	"places": {
		"folder": "folder", "folder-documents": "crafting", "folder-download": "chest",
		"folder-pictures": "grass", "folder-videos": "compass", "folder-music": "torch",
		"user-home": "grass", "user-trash": "chest",
	},
	"apps": {
		"org.gnome.Nautilus": "chest", "system-file-manager": "chest",
		"utilities-terminal": "crafting", "org.gnome.Terminal": "crafting",
		"com.mitchellh.ghostty": "crafting", "code": "pickaxe",
		"visual-studio-code": "pickaxe", "brave-browser": "compass",
		"brave": "compass", "org.blender.Blender": "diamond",
	},
	"status": {
		"dialog-information": "diamond", "dialog-warning": "torch", "dialog-error": "creeper",
	},
}

// minecraftQTColors positions match QPalette color roles, in order.
var minecraftQTDisabled = map[int]bool{0: true, 6: true, 8: true, 19: true}

func buildMinecraft(home string) error {
	for _, v := range minecraftVariants {
		if e := buildMinecraftVariant(home, v.slug, v.inherits); e != nil {
			return e
		}
	}
	// Regenerate the menu icon from the finished grass icon, as before.
	for _, v := range minecraftVariants {
		p := filepath.Join(home, ".config/omarchy/themes", v.slug, "assets")
		if e := copyFile(filepath.Join(p, "grass.svg"), filepath.Join(p, "menu.svg")); e != nil {
			return e
		}
	}
	return nil
}

func buildMinecraftVariant(home, slug, inherits string) error {
	c, e := load(home, slug)
	if e != nil {
		return e
	}
	a, cy := c.Accent, c.Cyan
	bg, fg := c.Background, c.Foreground
	panel, dark := c.LighterBackground, c.DarkBackground
	last := slug[strings.LastIndex(slug, "-")+1:]
	name := "Minecraft-" + title(last)
	themeDir := filepath.Join(home, ".config/omarchy/themes", slug)
	iconroot := filepath.Join(home, ".local/share/icons", name)
	if e := os.MkdirAll(filepath.Join(home, ".local/share/fonts"), 0o755); e != nil {
		return e
	}
	for _, sub := range []string{"places", "apps", "status"} {
		if e := os.MkdirAll(filepath.Join(iconroot, "scalable", sub), 0o755); e != nil {
			return e
		}
	}

	fmt.Println(slug + ":")
	// SVG assets.
	assets := map[string]string{}
	for _, s := range minecraftShapes {
		body := strings.ReplaceAll(s.shape, "{a}", a)
		body = strings.ReplaceAll(body, "{c}", cy)
		assets[s.key] = `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" shape-rendering="crispEdges">` + body + `</svg>`
	}
	for _, s := range minecraftShapes {
		if e := write(filepath.Join(themeDir, "assets", s.key+".svg"), assets[s.key]); e != nil {
			return e
		}
	}
	// Icon theme index and hard-linked applications of the assets.
	index := "[Icon Theme]\nName=" + name + "\nComment=Minecraft pixel accents\nInherits=" + inherits +
		"\nDirectories=scalable/places,scalable/apps,scalable/status\n"
	for _, sub := range []struct{ dir, ctx string }{
		{"places", "Places"}, {"apps", "Applications"}, {"status", "Status"},
	} {
		index += fmt.Sprintf("\n[scalable/%s]\nSize=32\nMinSize=16\nMaxSize=256\nType=Scalable\nContext=%s\n", sub.dir, sub.ctx)
	}
	if e := write(filepath.Join(iconroot, "index.theme"), index); e != nil {
		return e
	}
	for _, ctx := range []string{"places", "apps", "status"} {
		for file, asset := range minecraftMappings[ctx] {
			if e := write(filepath.Join(iconroot, "scalable", ctx, file+".svg"), assets[asset]); e != nil {
				return e
			}
		}
	}
	if e := write(filepath.Join(themeDir, "icons.theme"), name+"\n"); e != nil {
		return e
	}
	// GTK 3/4 and libadwaita named roles, with square beveled controls.
	defines := []string{
		"accent_color", a, "accent_bg_color", a, "accent_fg_color", bg,
		"window_bg_color", bg, "window_fg_color", fg,
		"view_bg_color", dark, "view_fg_color", fg,
		"headerbar_bg_color", panel, "headerbar_fg_color", fg,
		"card_bg_color", panel, "card_fg_color", fg,
		"dialog_bg_color", bg, "dialog_fg_color", fg,
		"popover_bg_color", bg, "popover_fg_color", fg,
		"sidebar_bg_color", dark, "sidebar_fg_color", fg,
		"theme_bg_color", bg, "theme_fg_color", fg,
		"theme_base_color", dark, "theme_text_color", fg,
		"theme_selected_bg_color", a, "theme_selected_fg_color", bg,
		"warning_color", c.Yellow, "error_color", c.Red, "success_color", c.Green,
	}
	lines := make([]string, 0, len(defines)/2)
	for i := 0; i < len(defines); i += 2 {
		lines = append(lines, fmt.Sprintf("@define-color %s %s;", defines[i], defines[i+1]))
	}
	css := strings.Join(lines, "\n") + "\n" + minecraftCSSBody(a, cy, bg, fg, panel, dark, c.Muted, c.Selection, c.BrightForeground)
	if e := write(filepath.Join(themeDir, "gtk.css"), css); e != nil {
		return e
	}
	// Qt color roles and stylesheet. Values become #ffRRGGBB (ARGB).
	roles := []string{fg, panel, c.LightForeground, panel, dark, c.Muted, fg, c.BrightForeground, fg, dark, bg, "#000000", a, bg, cy, c.Magenta, panel, "#000000", c.Yellow, fg, c.DarkForeground, a}
	role := func(v string) string { return "#ff" + v[1:] }
	for i := range roles {
		roles[i] = role(roles[i])
	}
	disabled := append([]string(nil), roles...)
	for i := range minecraftQTDisabled {
		disabled[i] = role(c.DarkForeground)
	}
	qtc := "[ColorScheme]\n"
	for _, v := range []struct {
		k string
		v []string
	}{{"active_colors", roles}, {"inactive_colors", roles}, {"disabled_colors", disabled}} {
		qtc += v.k + "=" + strings.Join(v.v, ", ") + "\n"
	}
	if e := write(filepath.Join(themeDir, "qt.colors"), qtc); e != nil {
		return e
	}
	qss := minecraftQSSTemplate
	for _, kv := range c.pairs() {
		qss = strings.ReplaceAll(qss, "@"+kv.key+"@", kv.value)
	}
	if e := write(filepath.Join(themeDir, "qt.qss"), qss); e != nil {
		return e
	}
	// Center every row on one shared canvas; the Creeper mask is symmetric.
	face := []string{
		"▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄",
		"████████████████",
		"██    ████    ██",
		"██    ████    ██",
		"██████    ██████",
		"████        ████",
		"████        ████",
		"████  ████  ████",
		"▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀",
	}
	rows := append(face, "", "M I N E C R A F T",
		strings.ToUpper(last), "", "MINE · CRAFT · EXPLORE · REST")
	width := 0
	for _, r := range rows {
		if n := utf8.RuneCountInString(r); n > width {
			width = n
		}
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if r == "" {
			out = append(out, "")
			continue
		}
		pad := width - utf8.RuneCountInString(r)
		left := pad / 2
		// CPython's str.center keeps odd leftovers on the left; the
		// screensaver is diffed against those generated canvases.
		if pad%2 == 1 && width%2 == 1 {
			left++
		}
		out = append(out, strings.TrimRight(strings.Repeat(" ", left)+r+strings.Repeat(" ", pad-left), " "))
	}
	if e := write(filepath.Join(themeDir, "screensaver.txt"), strings.Join(out, "\n")+"\n"); e != nil {
		return e
	}
	fmt.Println("  icons, GTK, Qt, screensaver assets built")
	return nil
}

func minecraftCSSBody(a, cy, bg, fg, panel, dark, muted, selection, bright string) string {
	return "window, dialog, .background { background-color: " + bg + "; color: " + fg + "; }\n" +
		"headerbar, .titlebar { background: " + panel + "; color: " + fg + "; border-bottom: 2px solid " + a + "; border-radius: 0; }\n" +
		"button { border-radius: 0; background-image: none; background-color: " + panel + "; color: " + fg + "; border: 2px solid " + muted + "; box-shadow: inset 2px 2px alpha(" + fg + ",0.12), inset -2px -2px alpha(#000000,0.35); }\n" +
		"button:hover { border-color: " + a + "; background-color: " + selection + "; }\n" +
		"button:checked, button.suggested-action { background: " + a + "; color: " + bg + "; border-color: " + cy + "; }\n" +
		"entry, textview, .view { background: " + dark + "; color: " + fg + "; border-radius: 0; }\n" +
		"entry:focus { border-color: " + a + "; box-shadow: 0 0 0 1px " + a + "; }\n" +
		"selection, row:selected { background-color: " + selection + "; color: " + bright + "; }\n" +
		"progressbar progress, scale highlight { background: " + a + "; border-radius: 0; }\n" +
		"scale slider { background: " + cy + "; border-radius: 0; }\n" +
		"scrollbar slider { background: " + muted + "; border-radius: 0; min-width: 8px; min-height: 8px; }\n" +
		"popover, menu { background: " + bg + "; color: " + fg + "; border: 2px solid " + a + "; border-radius: 0; }\n"
}

const minecraftQSSTemplate = "QPushButton, QToolButton { border: 2px solid @muted@; border-radius: 0; " +
	"border-top-color: @dark_foreground@; border-left-color: @dark_foreground@; " +
	"border-bottom-color: @darker_background@; border-right-color: @darker_background@; " +
	"padding: 5px; background: @lighter_background@; color: @foreground@; }\n" +
	"QPushButton:hover, QToolButton:hover { border-color: @accent@; }\n" +
	"QLineEdit, QComboBox { border: 2px solid @muted@; border-radius: 0; background: @dark_background@; padding: 4px; }\n" +
	"QProgressBar::chunk { background: @accent@; }\n" +
	"QLineEdit:focus { border-color: @accent@; }\n" +
	"QMenu { background: @background@; color: @foreground@; border: 2px solid @accent@; }\n" +
	"QMenu::item:selected { background: @selection@; }\n" +
	"QHeaderView::section { background: @lighter_background@; color: @foreground@; border: 1px solid @muted@; padding: 5px; }\n" +
	"QScrollBar::handle { background: @muted@; border-radius: 0; }\n"
