<?php

$here = <<<EOD
This is a big long $string
EOD;

$now = <<<'EOT'
This is a big long string with no variables
EOT;

$now = <<<"EOR"
This is another heredoc
EOR;
