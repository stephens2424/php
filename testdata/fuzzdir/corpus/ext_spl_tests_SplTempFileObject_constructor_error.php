<?php
try {
    new SplTempFileObject('invalid');
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
?>
