<?php
$dftz021 = date_default_timezone_get(); //UTC

$dtms021 = new DateTime(); 

date_timestamp_set($dtms021, 123456789, 'error');
?>
