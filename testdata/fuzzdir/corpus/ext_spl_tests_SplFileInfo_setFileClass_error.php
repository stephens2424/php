<?php

$info = new SplFileInfo(__FILE__);

try {
    $info->setFileClass('stdClass');
} catch (UnexpectedValueException $e) {
    echo $e->getMessage(), "\n";
}

?>
