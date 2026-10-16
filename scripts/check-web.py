#!/usr/bin/env python3
"""Validate packaged local links, images, scripts, CSS, and HTML IDs."""
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import urlsplit, unquote
root = Path('dist/browser').resolve()
errors = []
class Page(HTMLParser):
    def __init__(self, file):
        super().__init__(); self.file=file; self.ids=set()
    def handle_starttag(self, tag, attrs):
        attrs=dict(attrs)
        ident=attrs.get('id')
        if ident:
            if ident in self.ids: errors.append(f'{self.file.name}: duplicate ID {ident}')
            self.ids.add(ident)
        for key in ('src','href'):
            value=attrs.get(key,''); parsed=urlsplit(value)
            if not value or parsed.scheme or parsed.netloc or not parsed.path or parsed.path.startswith('/api/'): continue
            dest=(self.file.parent/unquote(parsed.path)).resolve()
            if not dest.is_relative_to(root) or not dest.exists(): errors.append(f'{self.file.name}: missing {value}')
for file in root.glob('*.html'):
    Page(file).feed(file.read_text())
for required in ('assets/wasm/cinebite.wasm','assets/wasm/wasm_exec.js','sw.js','manifest.json'):
    if not (root/required).is_file(): errors.append(f'missing {required}')
if errors: raise SystemExit('\n'.join(errors))
print('Packaged HTML references, unique IDs, and required runtime assets passed.')
