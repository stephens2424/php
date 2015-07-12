<?php
echo ini_get('filter.default') . "\n";
echo ini_get('filter.default_flags') . "\n";
var_dump(FILTER_FLAG_STRIP_LOW == 4);
echo addcslashes($_GET['a'],"\0") . "\n";
?>
