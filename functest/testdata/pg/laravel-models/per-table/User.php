<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * @property int $id
 * @property string $name
 * @property float $balance
 * @property float $prev_balance
 * @property \Illuminate\Support\Carbon $created_at
 * @property  $current_mood
 * @property \Illuminate\Support\Carbon $updated_at
 */
class User extends Model
{
    public $incrementing = false;

    protected $primaryKey = 'id';
    protected $keyType = 'int';

    protected $table = 'users';
}
