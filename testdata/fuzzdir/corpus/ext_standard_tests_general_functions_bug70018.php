<?php
$output = array();
exec('/bin/echo -n -e "abc\f\n \n"',$output);
var_dump($output);
?>
