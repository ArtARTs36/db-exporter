<?php

namespace {{ namespace }};

use Illuminate\Database\Eloquent\Model;
{% for model in schema.Models %}
/**
{% for property in model.Properties %} * @property {{ property.Type }} ${{ property.Name }}{% if loop.last == false %}
{% endif %}{% endfor %}
 */
class {{ model.Name }} extends Model
{
{% if model.PrimaryKey.Exists %}{% if model.PrimaryKey.IsMultiple %}    // @todo {{ model.PrimaryKey.Name }}: Eloquent native not supported multiple primary keys{% else %}    public $incrementing = {{ bool_string(model.PrimaryKey.Incrementing) }};

    protected $primaryKey = '{{ model.PrimaryKey.Column }}';
    protected $keyType = '{{ model.PrimaryKey.Type }}';
{% endif %}{% endif %}
    protected $table = '{{ model.Table }}';{% if model.Dates|length > 0 %}
    protected $dates = [
{% for name in model.Dates %}        '{{ name }}',{% endfor %}
    ];{% endif %}
}{% if loop.last == false %}
{% endif %}{% endfor %}
