<?php

$method = new ReflectionMethod('SplFileObject', 'fputcsv');
$params = $method->getParameters(); 
var_dump($params);

?>
===DONE===
