<?php 
echo filter_input(INPUT_POST, 'foo', FILTER_SANITIZE_STRIPPED);
?>
