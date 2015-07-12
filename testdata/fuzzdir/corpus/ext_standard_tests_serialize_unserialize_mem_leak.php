<?php
function foo() {
	gc_collect_cycles();
}

$str = 'a:1:{i:0;r:1;}';
foo(unserialize($str));
echo "okey";
?>
