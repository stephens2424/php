<?php
try {
    $fragment = new DOMDocumentFragment("root");
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
?>
