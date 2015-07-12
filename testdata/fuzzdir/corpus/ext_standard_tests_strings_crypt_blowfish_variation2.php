<?php
$crypt = crypt(b'U*U', b'$2a$CCCCCCCCCCCCCCCCCCCCC.E5YPO9kmyuRGyh0XouQYb4YMJKvyOeW');
if ($crypt==='*0') {
    echo "OK\n";
} else {
    echo "Not OK\n";
}
?>
