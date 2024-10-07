<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * @property int $id
 * @property string $name
 */
class User extends Model
{
    public $incrementing = true;

    protected $primaryKey = 'id';
    protected $keyType = 'int';

    protected $table = 'users';
}
