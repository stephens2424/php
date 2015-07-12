<?php
var_dump(intdiv(3, 2));
var_dump(intdiv(-3, 2));
var_dump(intdiv(3, -2));
var_dump(intdiv(-3, -2));
var_dump(intdiv(PHP_INT_MAX, PHP_INT_MAX));
var_dump(intdiv(PHP_INT_MIN, PHP_INT_MIN));
try {
  var_dump(intdiv(PHP_INT_MIN, -1));
} catch (Throwable $e) {
  echo "Exception: " . $e->getMessage() . "\n";
}
try {
  var_dump(intdiv(1, 0));
} catch (Throwable $e) {
  echo "Exception: " . $e->getMessage() . "\n";
}

?>
