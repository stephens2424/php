<?php
$ret = filter_input(INPUT_GET, 'a', FILTER_VALIDATE_INT);
var_dump($ret);

$ret = filter_input(INPUT_GET, 'a', FILTER_VALIDATE_INT, array('flags'=>FILTER_FLAG_ALLOW_OCTAL));
var_dump($ret);

$ret = filter_input(INPUT_GET, 'ar', FILTER_VALIDATE_INT, array('flags'=>FILTER_REQUIRE_ARRAY));
var_dump($ret);

$ret = filter_input(INPUT_GET, 'ar', FILTER_VALIDATE_INT, array('flags'=>FILTER_FLAG_ALLOW_OCTAL|FILTER_REQUIRE_ARRAY));
var_dump($ret);

?>
