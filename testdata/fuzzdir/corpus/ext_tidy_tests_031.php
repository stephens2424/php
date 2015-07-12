<?php
$buffer = '<html></html>';
$config = array('doctype' => 'php');

$tidy = tidy_parse_string($buffer, $config);
var_dump(tidy_config_count($tidy));
?>
