<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * @property string $entity_type
 * @property string $entity_id
 * @property string $name
 */
class Entity extends Model
{
    // @todo entities_pk: Eloquent native not supported multiple primary keys
    protected $table = 'entities';
}
