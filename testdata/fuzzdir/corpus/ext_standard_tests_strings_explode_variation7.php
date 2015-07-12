<?php

/* Prototype  : array explode  ( string $delimiter  , string $string  [, int $limit  ] )
 * Description: Split a string by string.
 * Source code: ext/standard/string.c
*/

var_dump(count(explode('|', implode(range(1,65),'|'), -1)));

?>
