<?php

use Illuminate\Database\Migrations\Migration;

class {{ migration.Meta["Name"] }} extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        {% for query in migration.UpQueries %}\Illuminate\Support\Facades\DB::unprepared(<<<SQL
{{ query }}
SQL);{% if loop.last == false %}
        {% endif %}{% endfor %}
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        {% for query in migration.DownQueries %}\Illuminate\Support\Facades\DB::unprepared(<<<SQL
{{ query }}
SQL);{% if loop.last == false %}
        {% endif %}{% endfor %}
    }
}
