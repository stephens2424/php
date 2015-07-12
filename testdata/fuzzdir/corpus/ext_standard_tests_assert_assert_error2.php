<?php
function f1($script, $line, $message, $user_message) 
{
	echo "f1 called\n";
}

//bail out on error
var_dump($rao = assert_options(ASSERT_BAIL, 1));
var_dump($r2 = assert("0 != 0"));
echo "If this is printed BAIL hasn't worked";
