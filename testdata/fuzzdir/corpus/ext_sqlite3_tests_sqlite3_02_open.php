<?php

try {
  $db = new SQLite3();
} catch (TypeError $e) {
  var_dump($e->getMessage());
}

?>
