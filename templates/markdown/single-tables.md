# Tables

{% for table in schema.Tables %}
- [{{ table.Name.Value }}](#{{ table.Name.Replace('_', '') }})
{% endfor %}

{% if diagram.Valid() %}![](./{{ diagram.FileName }}){% endif%}

{% for table in schema.Tables %}
## {{ table.Name.Value }}

| Column                         | Type                                                                                                                                                                                | Nullable              | Unique | Comment                       |
|--------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|---------------------------------------|---------------------------------------|
 {% for col in table.Columns %} | **{{ col.Name.String }}** {% if col.IsPrimaryKey() %} 🔑 {% endif %} {% if col.HasForeignKey() %} → {{ col.ForeignKey.Table.Value }}.{{ col.ForeignKey.Column.Value }} {% endif %} | {{ col.Type.String }} | {{ col.Nullable ? 'true' : 'false' }} | {{ col.IsUniqueKey() ? 'true' : 'false' }} | {{ col.Comment.Value }}  |
{% endfor %}
{% endfor %}
