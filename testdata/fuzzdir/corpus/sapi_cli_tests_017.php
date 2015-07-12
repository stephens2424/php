<?php
$php = getenv('TEST_PHP_EXECUTABLE');

$codes = array();

$codes[1] = <<<EOT
echo 'Hello world';
exit
EOT;

$codes[] = <<<EOT
echo 'multine
single
quote';
exit
EOT;

$codes[] = <<<EOT
echo <<<HEREDOC
Here
comes
the
doc
HEREDOC;
EOT;

$codes[] = <<<EOT
if (0) {
    echo "I'm not there";
}
echo "Done";
EOT;

$codes[] = <<<EOT
function a_function_with_some_name() {
    echo "I was called!";
}
a_function_w	);
EOT;

foreach ($codes as $key => $code) {
	echo "\n--------------\nSnippet no. $key:\n--------------\n";
	$code = escapeshellarg($code);
	echo `echo $code | "$php" -a`, "\n";
}

echo "\nDone\n";
?>
