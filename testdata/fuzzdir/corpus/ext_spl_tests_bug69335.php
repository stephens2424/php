<?php
$a = array(1=>1, 3=>3, 5=>5, 7=>7);
$a = new ArrayObject($a);

foreach ($a as $k => $v) {
	var_dump("$k => $v");
	if ($k == 3) {
		$a['a'] = "?";
	}
}
?>
