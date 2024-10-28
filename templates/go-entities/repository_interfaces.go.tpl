{% for repo in repositories %}{% set entityCall = repo.Entity.Call(_file.Package) %}type {{ repo.Interface.Name }} interface {
	Get(ctx context.Context, filter *{{ repo.Filters.Get.Name }}) (*{{ entityCall }}, error)
	List(ctx context.Context, filter *{{ repo.Filters.List.Name }}) ([]*{{ entityCall }}, error)
	Create(ctx context.Context, {{ repo.Entity.AsVarName }} *{{ entityCall }}) (*{{ entityCall }}, error)
	Update(ctx context.Context, {{ repo.Entity.AsVarName }} *{{ entityCall }}) (*{{ entityCall }}, error)
	Delete(ctx context.Context, filter *{{ repo.Filters.Delete.Name }}) (count int64, err error)
}{% if loop.last == false %}
{% endif %}{% endfor %}