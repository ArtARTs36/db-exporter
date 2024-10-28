package {{ _file.Package.Name }}{% for importPkg in stdImports %}{% set _ = _file.Imports.AddStd(importPkg) %}{% endfor %}{% for importPkg in sharedImports %}{% set _ = _file.Imports.AddShared(importPkg) %}{% endfor %}{% if _file.Imports.Valid() %}

import (
{% for group in _file.Imports.Sorted() %}{% for pkg in group %}	"{{ pkg }}"{% if loop.last == false %}
{% endif %}{% endfor %}{% if loop.last == false %}

{% endif %}{% endfor %}
){% endif %}