<?php
$x = unserialize('O:9:"exception":1:{s:16:"'."\0".'Exception'."\0".'trace";s:4:"ryat";}');
echo $x;
$x =  unserialize('O:4:"test":1:{s:27:"__PHP_Incomplete_Class_Name";R:1;}');
$x->test();

?>
