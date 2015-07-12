<?php
if (false) {
	    $willNeverBeDefined = true;
}
$result = compact('willNeverBeDefined');
var_dump($result, empty($result), $result === array(), empty($willNeverBeDefined));
