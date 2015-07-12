<?php
$a = new ArrayIterator(array('test1', 'test2', 'test3'));
$i = new RegexIterator($a, '/^(test)(\d+)/', RegexIterator::REPLACE);
$r = '$2:$1';
$i->replacement =& $r;
var_dump(iterator_to_array($i));
?>
