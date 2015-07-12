<?php
// We currently don't have any deprecated functions :/
$rc = new ReflectionFunction('var_dump');
var_dump($rc->isDeprecated());
