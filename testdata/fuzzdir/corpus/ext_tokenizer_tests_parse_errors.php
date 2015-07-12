<?php

function test_parse_error($code) {
    try {
        var_dump(token_get_all($code, TOKEN_PARSE));
    } catch (ParseError $e) {
        echo $e->getMessage(), "\n";
    }

    foreach (token_get_all($code) as $token) {
        if (is_array($token)) {
            echo token_name($token[0]), " ($token[1])\n";
        } else {
            echo "$token\n";
        }
    }
    echo "\n";
}

test_parse_error('<?php var_dump(078);');
test_parse_error('<?php var_dump("\u{xyz}");');
test_parse_error('<?php var_dump("\u{ffffff}");');
test_parse_error('<?php var_dump(078 + 078);');

?>
