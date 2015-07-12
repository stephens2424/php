<?php

$array = array(array(7,8,9),1,2,3,array(4,5,6));
$arrayIterator = new ArrayIterator($array);
try {
    $test = new IteratorIterator($arrayIterator);

    $test = new IteratorIterator($arrayIterator, 1);
    $test = new IteratorIterator($arrayIterator, 1, 1);
    $test = new IteratorIterator($arrayIterator, 1, 1, 1);
    $test = new IteratorIterator($arrayIterator, 1, 1, 1, 1);
} catch (TypeError $e){
  echo $e->getMessage() . "\n";
}

?>
===DONE===
