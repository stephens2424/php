<?php 
$a = tidy_parse_string('<HTML></HTML>');
var_dump(tidy_diagnose($a));
echo str_replace("\r", "", tidy_get_error_buffer($a));

$html = <<< HTML
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2//EN">
<html>
<head><title>foo</title></head>
<body><p>hello</p></body>
</html>
HTML;
$a = tidy_parse_string($html);
var_dump(tidy_diagnose($a));
echo tidy_get_error_buffer($a);
?>
