<?php
$start = new DateTime('2000-01-01 00:00:00', new DateTimeZone('Europe/Berlin'));
$end   = new DateTime('2000-01-31 00:00:00', new DateTimeZone('UTC'));
$interval = new DateInterval('P1D');
$period   = new DatePeriod($start, $interval, $end);

var_dump($period->getStartDate()->format('Y-m-d H:i:s'));
var_dump($period->getStartDate()->getTimeZone()->getName());

var_dump($period->getEndDate()->format('Y-m-d H:i:s'));
var_dump($period->getEndDate()->getTimeZone()->getName());

var_dump($period->getDateInterval()->format('%R%y-%m-%d-%h-%i-%s'));
?>
