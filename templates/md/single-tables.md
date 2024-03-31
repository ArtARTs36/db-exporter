# Tables

{% for table in schema.Tables.List() %}
- [{{ table.Name.Value }}](#{{ table.Name.Replace('_', '') }})
{% endfor %}

{% if diagram.Valid() %}![](./{{ diagram.FileName }}){% endif%}

{% for table in schema.Tables.List() %}
## {{ table.Name.Value }}

| Column                         | Type                                                                                                                                                                                | Nullable              | Unique | Comment                                                        |
|--------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|---------------------------------------|----------------------------------------------------------------|
 {% for col in table.Columns %} | **{{ col.Name.String }}** {% if col.IsPrimaryKey() %} 🔑 {% endif %} {% if col.HasForeignKey() %} → {{ col.ForeignKey.ForeignTable.Value }}.{{ col.ForeignKey.ForeignColumn.Value }} {% endif %} | {{ col.Type.String }} | {{ col.Nullable ? 'true' : 'false' }} | {{ col.IsUniqueOrPrimaryKey() ? 'true' : 'false' }} | {{ col.Comment.Value }}  |
{% endfor %}
{% if table.HasUniqueKeys() %}
{% if table.HasSingleUniqueKey() %}
Unique: {% if table.GetFirstUniqueKey().ColumnsNames.Once() %} {{ table.GetFirstUniqueKey().ColumnsNames.Join(", ") }}{% else %}({% for name in table.GetFirstUniqueKey().ColumnsNames.List() %}{{ name }}{% if loop.last == false %}, {% endif %}{% endfor %}) {% endif %}
{% else %}
Unique:
{% for key in table.UniqueKeys %}
- {% if key.ColumnsNames.Once() %} {{ key.ColumnsNames.Join(", ") }}{% else %}({% for name in key.ColumnsNames.List() %}{{ name }}{% if loop.last == false %}, {% endif %}{% endfor %}) {% endif %}
{% endfor %}
{% endif %}
{% endif %}
{% endfor %}
