<?php

$tz = new DateTimeZone('Europe/London');
$dt = new DateTimeImmutable('2014-09-20', $tz);

echo $tz->getOffset($dt);
echo $tz->getOffset(1);
