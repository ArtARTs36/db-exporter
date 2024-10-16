package {{ _file.Package.Name }}

import (
{% for pkg in _file.Imports.List() %}	"{{ pkg }}"{% if loop.last == false %}
{% endif %}{% endfor %}
)