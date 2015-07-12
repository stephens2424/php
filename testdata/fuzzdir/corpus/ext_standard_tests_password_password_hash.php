<?php
//-=-=-=-

var_dump(strlen(password_hash("foo", PASSWORD_BCRYPT)));

$hash = password_hash("foo", PASSWORD_BCRYPT);

var_dump($hash === crypt("foo", $hash));

echo "OK!";
?>
