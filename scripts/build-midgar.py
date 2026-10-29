"""Build native Midgar assets from colors.toml; no provider jobs or downloads."""
from pathlib import Path
import tomllib, shutil, ctypes as C
from PIL import Image, ImageDraw
home=Path.home();p=home/'.config/omarchy/themes/midgar-mako';c=tomllib.loads((p/'colors.toml').read_text());a=c['accent'];cy=c['cyan'];bg=c['background'];fg=c['foreground'];panel=c['lighter_background'];name='Midgar-Mako'
assets=p/'assets';assets.mkdir(exist_ok=True)
def svg(body):return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">'+body+'</svg>\n'
for kind,key in [('magic','green'),('support','blue'),('command','yellow'),('summon','red'),('independent','magenta')]:
 color=c[key];body=f'<defs><radialGradient id="g" cx="30%" cy="25%"><stop stop-color="#f0faf7"/><stop offset=".25" stop-color="{color}"/><stop offset="1" stop-color="{c["dark_background"]}"/></radialGradient></defs><circle cx="32" cy="32" r="25" fill="url(#g)" stroke="{color}" stroke-width="2"/><ellipse cx="24" cy="21" rx="8" ry="5" fill="#fff" opacity=".45"/>'
 (assets/(kind+'.svg')).write_text(svg(body))
# Approved v3 meteor; keep the authored SVG rather than redrawing its paths.
meteor=(Path(__file__).parent/'assets/midgar-meteor.svg').read_text().replace('#63e6be',a)
(assets/'menu.svg').write_text(meteor)
folder=svg(f'<path d="M5 15h20l6 7h28v33H5z" fill="{panel}" stroke="{a}" stroke-width="2"/><path d="M5 28h54" stroke="{cy}"/><circle cx="32" cy="41" r="9" fill="{a}" opacity=".7"/><circle cx="29" cy="38" r="3" fill="{fg}"/>')
sword=svg(f'<path d="M38 3l14 9-24 32-13-9z" fill="{fg}" stroke="{cy}" stroke-width="2"/><path d="M11 34l19 13M21 41L9 58" stroke="{a}" stroke-width="5"/><circle cx="33" cy="17" r="3" fill="{a}"/><circle cx="29" cy="24" r="3" fill="{c["magenta"]}"/>')
(assets/'sword.svg').write_text(sword)
iconroot=home/'.local/share/icons'/name
for sub in ['scalable/places','scalable/apps','scalable/status']:(iconroot/sub).mkdir(parents=True,exist_ok=True)
(iconroot/'index.theme').write_text('[Icon Theme]\nName=Midgar Mako\nInherits=Yaru-sage,Adwaita,hicolor\nDirectories=scalable/places,scalable/apps,scalable/status\n\n'+''.join(f'[{sub}]\nSize=64\nType=Scalable\nMinSize=16\nMaxSize=512\nContext={ctx}\n\n' for sub,ctx in [('scalable/places','Places'),('scalable/apps','Applications'),('scalable/status','Status')]))
for key in ['folder','folder-documents','folder-download','folder-pictures','folder-videos','folder-music','user-home']:(iconroot/'scalable/places'/(key+'.svg')).write_text(folder)
for key,body in [('brave-browser',meteor),('brave',meteor),('code',sword),('com.mitchellh.ghostty',(assets/'magic.svg').read_text()),('utilities-terminal',(assets/'magic.svg').read_text()),('org.gnome.Nautilus',folder),('system-file-manager',folder),('org.blender.Blender',(assets/'command.svg').read_text())]:(iconroot/'scalable/apps'/(key+'.svg')).write_text(body)
(p/'icons.theme').write_text(name+'\n')
roles={'accent_color':a,'accent_bg_color':a,'accent_fg_color':bg,'window_bg_color':bg,'window_fg_color':fg,'view_bg_color':c['dark_background'],'view_fg_color':fg,'headerbar_bg_color':panel,'headerbar_fg_color':fg,'card_bg_color':panel,'card_fg_color':fg,'dialog_bg_color':bg,'dialog_fg_color':fg,'popover_bg_color':panel,'popover_fg_color':fg,'sidebar_bg_color':c['dark_background'],'sidebar_fg_color':fg,'theme_bg_color':bg,'theme_fg_color':fg,'theme_base_color':c['dark_background'],'theme_text_color':fg,'theme_selected_bg_color':c['selection'],'theme_selected_fg_color':fg,'warning_color':c['yellow'],'error_color':c['red'],'success_color':a}
css=''.join(f'@define-color {k} {v};\n' for k,v in roles.items())
css+=f'''window, dialog, .background {{ background-color: {bg}; color: {fg}; }}
headerbar, .titlebar {{ background: {panel}; color: {fg}; border-bottom: 1px solid {a}; }}
button {{ border-radius: 4px; background-image: none; background-color: {panel}; color: {fg}; border: 1px solid {c['muted']}; }}
button:hover {{ border-color: {cy}; background-color: {c['selection']}; }}
button:checked, button.suggested-action {{ background: {c['selection']}; color: {fg}; border-color: {a}; }}
entry, textview, .view {{ background: {c['dark_background']}; color: {fg}; border-radius: 4px; }}
entry:focus {{ border-color: {a}; box-shadow: 0 0 0 1px {a}; }}
selection, row:selected {{ background-color: {c['selection']}; color: {c['bright_foreground']}; }}
progressbar progress, scale highlight {{ background: {a}; border-radius: 4px; }}
scale slider {{ background: {cy}; }}
scrollbar slider {{ background: {c['muted']}; min-width: 8px; min-height: 8px; }}
popover, menu {{ background: {panel}; color: {fg}; border: 1px solid {a}; border-radius: 4px; }}
'''
(p/'gtk.css').write_text(css)
qtroles=[fg,panel,panel,c['dark_background'],c['muted'],fg,fg,fg,fg,bg,bg,c['dark_background'],c['selection'],c['bright_foreground'],cy,c['magenta'],bg,bg,c['yellow'],fg,c['muted'],c['dark_background']]
(p/'qt.colors').write_text('[ColorScheme]\n'+''.join(k+'='+','.join('#ff'+x[1:] for x in qtroles)+'\n' for k in ['active_colors','inactive_colors','disabled_colors']))
(p/'qt.qss').write_text(f'''QPushButton, QToolButton {{ border: 1px solid {c['muted']}; border-radius: 4px; padding: 5px; background: {panel}; color: {fg}; }}
QPushButton:hover, QToolButton:hover {{ border-color: {cy}; }}
QLineEdit, QComboBox {{ border: 1px solid {c['muted']}; border-radius: 4px; background: {c['dark_background']}; padding: 4px; }}
QLineEdit:focus {{ border-color: {a}; }}
QProgressBar::chunk {{ background: {a}; }}
QMenu {{ background: {panel}; color: {fg}; border: 1px solid {a}; }}
QMenu::item:selected {{ background: {c['selection']}; }}
QHeaderView::section {{ background: {panel}; color: {fg}; border: 1px solid {c['muted']}; padding: 5px; }}
''')
print('Built Midgar materia, meteor menu, icons, GTK and Qt assets.')
