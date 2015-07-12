<?php
class Foo {

  function __construct() {
  }

  function test() {
    return $this->str;
  }
}

$foo = new Foo();
$server = new SoapServer(null,array('uri'=>"http://testuri.org"));
$server->setObject($foo);
var_dump($server->getfunctions());
echo "ok\n";
?>
