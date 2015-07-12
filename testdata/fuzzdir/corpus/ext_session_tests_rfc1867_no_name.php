<?php
session_start();
var_dump(session_id());
var_dump($_FILES);
var_dump($_SESSION["upload_progress_" . basename(__FILE__)]);
session_destroy();
?>
