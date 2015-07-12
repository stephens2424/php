<?php

$data = pack('VV', 1, 2);

$result = unpack('Va/X' ,$data);
var_dump($result);

$result = unpack('Va/X4' ,$data);
var_dump($result);

$result = unpack('V1a/X4/V1b/V1c/X4/V1d', $data);
var_dump($result);

?>
===DONE===
