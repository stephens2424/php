<?php

echo "INPUT: \n";
echo file_get_contents("php://input") . "\n";
echo "\n\n-----------\n\n";

function test() {
  return "Hello World";
}

$server = new soapserver(null,array('uri'=>"http://testuri.org"));
$server->addfunction("test");
$server->handle();
echo "ok\n";
?>
