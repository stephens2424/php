<?php
class SplMinHeap1 extends SplMinHeap {
  public function compare($a, $b) {
    return -parent::notexist($a, $b);
  }
}
$h = new SplMinHeap1();
$h->insert(1);
$h->insert(6);
?>
===DONE===
