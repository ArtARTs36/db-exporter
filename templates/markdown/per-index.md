# Tables

{% for table in tables %}
- [{{ table.Name.Value }}]({{ table.FileName }})
  {% endfor %}

{% if diagram.Valid() %}![](./{{ diagram.FileName }}){% endif%}
