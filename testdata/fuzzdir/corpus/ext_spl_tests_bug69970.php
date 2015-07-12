<?php

$count = 10;

class RecursiveArrayIteratorIterator extends RecursiveIteratorIterator {
  function rewind() {
    echo "dummy\n";
  }
  function endChildren() {
	  global $count;
	  echo $this->getDepth();
	  if (--$count > 0) {
		  // Trigger use-after-free
		  parent::rewind();
	  }
  }
}
$arr = array("a", array("ba", array("bba", "bbb")));
$obj = new RecursiveArrayIterator($arr);
$rit = new RecursiveArrayIteratorIterator($obj);

foreach ($rit as $k => $v) {
  echo ($rit->getDepth()) . "$k=>$v\n";
}
?>
