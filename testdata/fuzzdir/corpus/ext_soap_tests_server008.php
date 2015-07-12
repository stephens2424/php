<?php
class Foo {

  function __construct() {
  }

  function test() {
    return $this->str;
  }
}

$server = new soapserver(null,array('uri'=>"http://testuri.org"));
$server->setclass("Foo");
var_dump($server->getfunctions());
echo "ok\n";
?>
