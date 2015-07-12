<?php

try {
    token_get_all('<?php invalid code;', TOKEN_PARSE);
} catch (ParseError $e) {
    echo $e->getMessage(), PHP_EOL;
}

echo "Done";

?>
