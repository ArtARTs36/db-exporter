{% for repo in repositories %}type {{ repo.Interface.Name }} interface {
	Get(ctx context.Context, filter *{{ repo.Filters.Get.Name }}) (*{{ repo.EntityCall }}, error)
	List(ctx context.Context, filter *{{ repo.Filters.List.Name }}) ([]*{{ repo.EntityCall }}, error)
	Create(ctx context.Context, {{ repo.Entity.AsVarName }} *{{ repo.EntityCall }}) (*{{ repo.EntityCall }}, error)
	Update(ctx context.Context, {{ repo.Entity.AsVarName }} *{{ repo.EntityCall }}) (*{{ repo.EntityCall }}, error)
	Delete(ctx context.Context, filter *{{ repo.Filters.Delete.Name }}) (count int64, err error)
}{% endfor %}
