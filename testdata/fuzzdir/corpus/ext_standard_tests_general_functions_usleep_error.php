<?php
/* Prototype  : void usleep  ( int $micro_seconds  )
 * Description: Delays program execution for the given number of micro seconds. 
 * Source code: ext/standard/basic_functions.c
 */
 
set_time_limit(20); 

echo "*** Testing usleep() : error conditions ***\n";

echo "\n-- Testing usleep() function with zero arguments --\n";
var_dump( usleep() );

echo "\n-- Testing usleep() function with more than expected no. of arguments --\n";
$seconds = 10;
$extra_arg = 10;
var_dump( usleep($seconds, $extra_arg) );

echo "\n-- Testing usleep() function with negative interval --\n";
$seconds = -10;
var_dump( usleep($seconds) );

?>
===DONE===
