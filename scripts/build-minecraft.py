from pathlib import Path
import tomllib,shutil,ctypes as C
from PIL import Image,ImageDraw
home=Path.home(); root=home/'.config/omarchy/themes'; fonts=home/'.local/share/fonts'
fonts.mkdir(parents=True,exist_ok=True)
# Native vector icons: stay crisp at every desktop scale.
shapes={
'grass':'<polygon points="2,11 16,3 30,11 16,19" fill="{a}"/><polygon points="2,11 16,19 16,30 2,22" fill="#70513a"/><polygon points="16,19 30,11 30,22 16,30" fill="#493626"/><path d="M2 11v5l4 2v-3l4 2v3l6 3v-4" fill="#6ca64c"/>',
'diamond':'<path d="M12 3h8v4h5v6h4v7h-5v5h-6v4h-7v-4H6v-6H3v-7h5V7h4z" fill="{c}"/><path d="M12 7h8v4h-8v8H8v-8h4z" fill="#c5fff0"/><path d="M18 22h7v-7h-4v4h-3z" fill="#328f95"/>',
'pickaxe':'<path d="M7 4h16v4h5v5h-5v-3H12v4H7z" fill="{c}"/><path d="M15 10h5v6h-5v5h-5v6H5v-5h5v-6h5z" fill="#a17a46"/>',
'sword':'<path d="M24 2h6v6h-4v4h-4v4h-4v4h-5v-5h4v-4h4V7h3z" fill="{c}"/><path d="M9 13h4v4h4v4h-4v-4H9z" fill="#e5c777"/><path d="M8 20h5v5H8v5H3v-5h5z" fill="#a17a46"/>',
'chest':'<path d="M3 8h26v20H3z" fill="#986431"/><path d="M3 8h26v7H3z" fill="#bb8b45"/><path d="M3 15h26M5 8v20M27 8v20" stroke="#432d20" stroke-width="2"/><path d="M14 13h5v8h-5z" fill="#ebc76c"/>',
'crafting':'<path d="M3 6h26v23H3z" fill="#785137"/><path d="M5 7h22v13H5z" fill="#c2996a"/><path d="M12 7v13M20 7v13M5 13h22" stroke="#624129" stroke-width="2"/><path d="M7 23h18v3H7z" fill="#30271e"/>',
'compass':'<path d="M8 3h16v4h5v17h-5v5H8v-5H3V8h5z" fill="#9da6a2"/><path d="M9 7h14v4h3v12h-4v3H10v-3H6V11h3z" fill="#24312c"/><path d="M18 6l-7 13h7z" fill="#e87b6b"/><path d="M11 19h7l-4 7z" fill="#e2e7d7"/>',
'torch':'<path d="M13 12h7v18h-7z" fill="#987045"/><path d="M9 4h15v13H9z" fill="#ed9c43"/><path d="M12 1h9v12h-9z" fill="#f8d56c"/><path d="M15 4h3v8h-3z" fill="#fff1b4"/>',
'heart':'<path d="M2 9h3V6h7v3h8V6h7v3h3v10h-4v4h-4v4h-4v3h-4v-3h-4v-4H6v-4H2z" fill="#dc6963"/><path d="M6 10h5v4H6z" fill="#ffb0a0"/>',
'folder':'<path d="M2 7h11l3 4h14v17H2z" fill="{a}"/><path d="M2 13h28v15H2z" fill="#647d4c"/><path d="M3 15h26v3H3z" fill="{a}"/><path d="M6 20h5v4H6zM19 20h5v4h-5z" fill="#263a24"/>',
'creeper':'<path d="M3 3h26v26H3z" fill="{a}"/><path d="M7 8h6v6H7zM19 8h6v6h-6zM13 14h6v10h-6zM10 18h3v9h-3zM19 18h3v9h-3z" fill="#19311f"/>',
'food':'<path d="M5 7h18v4h5v12h-5v5H8v-4H3V12h2z" fill="#aa6b38"/><path d="M8 11h14v4H8zM8 19h17v3H8z" fill="#deb172"/>',
}
for slug in ['minecraft-deepstone','minecraft-nether','minecraft-end']:
 p=root/slug;c=tomllib.loads((p/'colors.toml').read_text());a=c['accent'];cy=c['cyan'];bg=c['background'];fg=c['foreground'];panel=c['lighter_background'];name='Minecraft-'+slug.split('-')[-1].title()
 iconroot=home/'.local/share/icons'/name
 for sub in ['scalable/places','scalable/apps','scalable/status']:(iconroot/sub).mkdir(parents=True,exist_ok=True)
 index='[Icon Theme]\nName='+name+'\nComment=Minecraft pixel accents\nInherits='+('Yaru-purple' if 'end' in slug else 'Yaru-red' if 'nether' in slug else 'Yaru-sage')+'\nDirectories=scalable/places,scalable/apps,scalable/status\n'
 for sub,ctx in [('places','Places'),('apps','Applications'),('status','Status')]:index+=f'\n[scalable/{sub}]\nSize=32\nMinSize=16\nMaxSize=256\nType=Scalable\nContext={ctx}\n'
 (iconroot/'index.theme').write_text(index)
 (p/'assets').mkdir(exist_ok=True)
 for key,shape in shapes.items():
  svg='<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" shape-rendering="crispEdges">'+shape.replace('{a}',a).replace('{c}',cy)+'</svg>'
  (p/'assets'/f'{key}.svg').write_text(svg)
 mappings={'places':{'folder':'folder','folder-documents':'crafting','folder-download':'chest','folder-pictures':'grass','folder-videos':'compass','folder-music':'torch','user-home':'grass','user-trash':'chest'},'apps':{'org.gnome.Nautilus':'chest','system-file-manager':'chest','utilities-terminal':'crafting','org.gnome.Terminal':'crafting','com.mitchellh.ghostty':'crafting','code':'pickaxe','visual-studio-code':'pickaxe','brave-browser':'compass','brave':'compass','org.blender.Blender':'diamond'},'status':{'dialog-information':'diamond','dialog-warning':'torch','dialog-error':'creeper'}}
 for context,entries in mappings.items():
  for filename,asset in entries.items():shutil.copy2(p/'assets'/f'{asset}.svg',iconroot/'scalable'/context/(filename+'.svg'))
 (p/'icons.theme').write_text(name+'\n')
 # GTK 3/4 and libadwaita named roles, and square beveled controls.
 css='\n'.join(f'@define-color {role} {value};' for role,value in {'accent_color':a,'accent_bg_color':a,'accent_fg_color':bg,'window_bg_color':bg,'window_fg_color':fg,'view_bg_color':c['dark_background'],'view_fg_color':fg,'headerbar_bg_color':panel,'headerbar_fg_color':fg,'card_bg_color':panel,'card_fg_color':fg,'dialog_bg_color':bg,'dialog_fg_color':fg,'popover_bg_color':bg,'popover_fg_color':fg,'sidebar_bg_color':c['dark_background'],'sidebar_fg_color':fg,'theme_bg_color':bg,'theme_fg_color':fg,'theme_base_color':c['dark_background'],'theme_text_color':fg,'theme_selected_bg_color':a,'theme_selected_fg_color':bg,'warning_color':c['yellow'],'error_color':c['red'],'success_color':c['green']}.items())
 css+=f'''\nwindow, dialog, .background {{ background-color: {bg}; color: {fg}; }}
headerbar, .titlebar {{ background: {panel}; color: {fg}; border-bottom: 2px solid {a}; border-radius: 0; }}
button {{ border-radius: 0; background-image: none; background-color: {panel}; color: {fg}; border: 2px solid {c['muted']}; box-shadow: inset 2px 2px alpha({fg},0.12), inset -2px -2px alpha(#000000,0.35); }}
button:hover {{ border-color: {a}; background-color: {c['selection']}; }}
button:checked, button.suggested-action {{ background: {a}; color: {bg}; border-color: {cy}; }}
entry, textview, .view {{ background: {c['dark_background']}; color: {fg}; border-radius: 0; }}
entry:focus {{ border-color: {a}; box-shadow: 0 0 0 1px {a}; }}
selection, row:selected {{ background-color: {c['selection']}; color: {c['bright_foreground']}; }}
progressbar progress, scale highlight {{ background: {a}; border-radius: 0; }}
scale slider {{ background: {cy}; border-radius: 0; }}
scrollbar slider {{ background: {c['muted']}; border-radius: 0; min-width: 8px; min-height: 8px; }}
popover, menu {{ background: {bg}; color: {fg}; border: 2px solid {a}; border-radius: 0; }}
'''
 (p/'gtk.css').write_text(css)
 roles=[fg,panel,c['light_foreground'],panel,c['dark_background'],c['muted'],fg,c['bright_foreground'],fg,c['dark_background'],bg,'#000000',a,bg,cy,c['magenta'],panel,'#000000',c['yellow'],fg,c['dark_foreground'],a]
 disabled=roles.copy()
 for i in [0,6,8,19]:disabled[i]=c['dark_foreground']
 (p/'qt.colors').write_text('[ColorScheme]\n'+''.join(k+'='+', '.join('#ff'+x[1:] for x in v)+'\n' for k,v in [('active_colors',roles),('inactive_colors',roles),('disabled_colors',disabled)]))
 qss = 'QPushButton, QToolButton { border: 2px solid @muted@; border-radius: 0; border-top-color: @dark_foreground@; border-left-color: @dark_foreground@; border-bottom-color: @darker_background@; border-right-color: @darker_background@; padding: 5px; background: @lighter_background@; color: @foreground@; }\nQPushButton:hover, QToolButton:hover { border-color: @accent@; }\nQLineEdit, QComboBox { border: 2px solid @muted@; border-radius: 0; background: @dark_background@; padding: 4px; }\nQProgressBar::chunk { background: @accent@; }\nQLineEdit:focus { border-color: @accent@; }\nQMenu { background: @background@; color: @foreground@; border: 2px solid @accent@; }\nQMenu::item:selected { background: @selection@; }\nQHeaderView::section { background: @lighter_background@; color: @foreground@; border: 1px solid @muted@; padding: 5px; }\nQScrollBar::handle { background: @muted@; border-radius: 0; }\n'
 for key,value in c.items():
  if isinstance(value,str):qss=qss.replace('@'+key+'@',value)
 (p/'qt.qss').write_text(qss)
 # Center every row within one shared canvas; the Creeper mask is symmetric.
 face = [
  '▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄',
  '████████████████',
  '██    ████    ██',
  '██    ████    ██',
  '██████    ██████',
  '████        ████',
  '████        ████',
  '████  ████  ████',
  '▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀',
 ]
 lines = face + ['', 'M I N E C R A F T', slug.split('-')[-1].upper(), '', 'MINE · CRAFT · EXPLORE · REST']
 width = max(map(len, lines))
 (p/'screensaver.txt').write_text('\n'.join(line.center(width).rstrip() if line else '' for line in lines)+'\n')
 print(slug,'icons, GTK, Qt, screensaver assets built')

for slug in ['minecraft-deepstone','minecraft-nether','minecraft-end']:
 shutil.copy2(root/slug/'assets/grass.svg',root/slug/'assets/menu.svg')
