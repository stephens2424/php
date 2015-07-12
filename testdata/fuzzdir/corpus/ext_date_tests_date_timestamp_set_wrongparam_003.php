<?php
$dftz021 = date_default_timezone_get(); //UTC

$dtms021 = new DateTime(); 

$wrong_parameter = array();

date_timestamp_set($dtms021, $wrong_parameter);
?>
