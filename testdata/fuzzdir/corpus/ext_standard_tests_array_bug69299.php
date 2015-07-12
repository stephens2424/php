<?php
$toFilter = array('foo' => 'bar', 'fiz' => 'buz');
$filtered = array_filter($toFilter, function ($value, $key) {
	if ($value === 'buz'
		|| $key === 'foo'
	) {
		return false;
	}
	return true;
}, ARRAY_FILTER_USE_BOTH);
var_dump($filtered);
?>
