<?php

var_dump(xmlrpc_get_type(new stdclass));
var_dump(xmlrpc_get_type(array()));
$var = array(1,2,3);
var_dump(xmlrpc_get_type($var));
$var = array("test"=>1,2,3);
var_dump(xmlrpc_get_type($var));
$var = array("test"=>1,"test2"=>2);
var_dump(xmlrpc_get_type($var));

echo "Done\n";
?>
