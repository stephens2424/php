<?php
// Could also be set with a cookie like "PHPSESSID=; path=/"
session_id('');

// Will still start the session and return true
var_dump(session_start());

// Returns an empty string
var_dump(session_id());
?>
