<?php
try {
    $comment = new DOMComment("comment1", "comment2");
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
?>
