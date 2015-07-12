<?php
/* Prototype  : int sleep  ( int $seconds  )
 * Description: Delays the program execution for the given number of seconds . 
 * Source code: ext/standard/basic_functions.c
 */
 set_time_limit(20);
 
echo "*** Testing sleep() : error conditions ***\n";

echo "\n-- Testing sleep() function with zero arguments --\n";
var_dump( sleep() );

echo "\n-- Testing sleep() function with more than expected no. of arguments --\n";
$seconds = 10;
$extra_arg = 10;
var_dump( sleep($seconds, $extra_arg) );

echo "\n-- Testing sleep() function with negative interval --\n";
$seconds = -10;
var_dump( sleep($seconds) );

?>
===DONE===
