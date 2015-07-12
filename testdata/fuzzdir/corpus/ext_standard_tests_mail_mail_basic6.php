<?php
/* Prototype  : int mail(string to, string subject, string message [, string additional_headers [, string additional_parameters]])
 * Description: Send an email message with invalid addtional_headers
 * Source code: ext/standard/mail.c
 * Alias to functions:
 */

echo "*** Testing mail() : basic functionality ***\n";


// Valid header
$to = 'user@example.com';
$subject = 'Test Subject';
$message = 'A Message';
$additional_headers = "HEAD1: a\r\nHEAD2: b\r\n";
$outFile = "mailBasic.out";
@unlink($outFile);

echo "-- Valid Header --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo file_get_contents($outFile);
unlink($outFile);

// Valid header
$additional_headers = "HEAD1: a\nHEAD2: b\n";
@unlink($outFile);

echo "-- Valid Header --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Valid header
// \r is accepted as valid. This may be changed to invalid.
$additional_headers = "HEAD1: a\rHEAD2: b\r";
@unlink($outFile);

echo "-- Valid Header --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

//===============================================================================
// Invalid header
$additional_headers = "\nHEAD1: a\nHEAD2: b\n";
@unlink($outFile);

echo "-- Invalid Header - preceeding newline--\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "\rHEAD1: a\nHEAD2: b\r";
@unlink($outFile);

echo "-- Invalid Header - preceeding newline--\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "\r\nHEAD1: a\r\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - preceeding newline--\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "\r\n\r\nHEAD1: a\r\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - preceeding newline--\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "\n\nHEAD1: a\r\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - preceeding newline--\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "\r\rHEAD1: a\r\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - preceeding newline--\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "HEAD1: a\r\n\r\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - multiple newlines in the middle --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "HEAD1: a\r\n\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - multiple newlines in the middle --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "HEAD1: a\n\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - multiple newlines in the middle --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "HEAD1: a\r\rHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - multiple newlines in the middle --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "HEAD1: a\n\rHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - multiple newlines in the middle --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
$additional_headers = "HEAD1: a\n\r\nHEAD2: b\r\n";
@unlink($outFile);

echo "-- Invalid Header - multiple newlines in the middle --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
// Invalid, but PHP_FUNCTION(mail) trims newlines
$additional_headers = "HEAD1: a\r\nHEAD2: b\r\n\n";
@unlink($outFile);

echo "-- Invalid Header - trailing newlines --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
// Invalid, but PHP_FUNCTION(mail) trims newlines
$additional_headers = "HEAD1: a\r\nHEAD2: b\n\n";
@unlink($outFile);

echo "-- Invalid Header - trailing newlines --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
// Invalid, but PHP_FUNCTION(mail) trims newlines
$additional_headers = "HEAD1: a\r\nHEAD2: b\n";
@unlink($outFile);

echo "-- Invalid Header - trailing newlines --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

// Invalid header
// Invalid, but PHP_FUNCTION(mail) trims newlines
$additional_headers = "HEAD1: a\r\nHEAD2: b\r";
@unlink($outFile);

echo "-- Invalid Header - trailing newlines --\n";
// Calling mail() with all additional headers
var_dump( mail($to, $subject, $message, $additional_headers) );
echo @file_get_contents($outFile);
@unlink($outFile);

?>
===DONE===
