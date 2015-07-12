<?php
session_start();
define ("user", "foo");
var_dump(session_regenerate_id());
?>
