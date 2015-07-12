<?php
libxml_use_internal_errors(true);
$doc = new DomDocument();
$doc->load(__DIR__ . DIRECTORY_SEPARATOR . 'bug69753.xml');
$error = libxml_get_last_error();
var_dump($error->file);
?>
