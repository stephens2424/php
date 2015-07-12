<?php
$file_path = dirname(__FILE__);
@mkdir("$file_path/realpath_basic/home/test", 0777, true);
@symlink("$file_path/realpath_basic/home", "$file_path/realpath_basic/link1");
@symlink("$file_path/realpath_basic/link1", "$file_path/realpath_basic/link2");
echo "1. " . realpath("$file_path/realpath_basic/link2") . "\n";
echo "2. " . realpath("$file_path/realpath_basic/link2/test") . "\n";
?>
