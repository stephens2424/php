<?php
$buffer = '<img src="file.png" /><php>';
$config = array(
  'accessibility-check' => 1);

$tidy = tidy_parse_string($buffer, $config);
$tidy->diagnose();
var_dump(tidy_access_count($tidy));
?>
