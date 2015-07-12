<?php
date_default_timezone_set('America/Buenos_Aires');

$date = new DateTime('2009-09-28 09:45:31.918312');

var_dump($date->format(DateTime::RFC3339_EXTENDED));
var_dump($date->format('u'));
var_dump($date->format('v'));
