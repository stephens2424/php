<?php
$buffer = '<img src="file.png" /><php>';

$tidy = tidy_parse_string($buffer);
var_dump(tidy_error_count($tidy));
?>
