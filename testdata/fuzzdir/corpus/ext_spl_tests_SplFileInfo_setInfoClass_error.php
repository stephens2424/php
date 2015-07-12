<?php

$info = new SplFileInfo(__FILE__);

try {
    $info->setInfoClass('stdClass');
} catch (UnexpectedValueException $e) {
    echo $e->getMessage(), "\n";
}

?>
