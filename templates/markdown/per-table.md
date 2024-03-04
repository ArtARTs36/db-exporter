# {{ table.Name.Value }}

| Column                         | Type                                                                                                                                                                               | Nullable              | Comment                               |
|--------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|---------------------------------------|
 {% for col in table.Columns %} | **{{ col.Name.String }}** {% if col.IsPrimaryKey() %} 🔑 {% endif %} {% if col.HasForeignKey() %} → {{ col.ForeignKey.Table.Value }}.{{ col.ForeignKey.Column.Value }} {% endif %} | {{ col.Type.String }} | {{ col.Nullable ? 'true' : 'false' }} | {{ col.Comment.Value }}  |
{% endfor %}
