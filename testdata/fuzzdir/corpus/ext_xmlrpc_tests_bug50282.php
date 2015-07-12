<?php

class One { var $x = 10; }

$o = new One();
var_dump($o);
var_dump(xmlrpc_encode_request('test', $o));
var_dump($o);

?>
