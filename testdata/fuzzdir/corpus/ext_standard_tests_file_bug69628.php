<?php

$file_path = dirname(__FILE__);

// temp dirname used here
$dirname = "$file_path/bug69628";

// temp dir created
mkdir($dirname);

// temp files created
file_put_contents("$dirname/image.jPg", '');
file_put_contents("$dirname/image.gIf", '');
file_put_contents("$dirname/image.png", '');

sort_var_dump(glob("$dirname/*.{[jJ][pP][gG],[gG][iI][fF]}", GLOB_BRACE));

function sort_var_dump($results) {
   sort($results);
   var_dump($results);
}

?>
