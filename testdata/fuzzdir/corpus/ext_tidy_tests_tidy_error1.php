<?php
$buffer = '<html></html>';
$config = array('bogus' => 'willnotwork');

$tidy = new tidy();
var_dump($tidy->parseString($buffer, $config));
?>
