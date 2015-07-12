<?php
try {
    $attr = new DOMAttr();
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
?>
