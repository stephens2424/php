<?php

$values = Array(
'http://example.com/index.html',	
'http://www.example.com/index.php',	
'http://www.example/img/test.png',	
'http://www.example/img/dir/',	
'http://www.example/img/dir',	
'http://www.thelongestdomainnameintheworldandthensomeandthensomemoreandmore.com/',
'http://toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong.com',
'http://eauBcFReEmjLcoZwI0RuONNnwU4H9r151juCaqTI5VeIP5jcYIqhx1lh5vV00l2rTs6y7hOp7rYw42QZiq6VIzjcYrRm8gFRMk9U9Wi1grL8Mr5kLVloYLthHgyA94QK3SaXCATklxgo6XvcbXIqAGG7U0KxTr8hJJU1p2ZQ2mXHmp4DhYP8N9SRuEKzaCPcSIcW7uj21jZqBigsLsNAXEzU8SPXZjmVQVtwQATPWeWyGW4GuJhjP4Q8o0.com',
'http://kDTvHt1PPDgX5EiP2MwiXjcoWNOhhTuOVAUWJ3TmpBYCC9QoJV114LMYrV3Zl58.kDTvHt1PPDgX5EiP2MwiXjcoWNOhhTuOVAUWJ3TmpBYCC9QoJV114LMYrV3Zl58.kDTvHt1PPDgX5EiP2MwiXjcoWNOhhTuOVAUWJ3TmpBYCC9QoJV114LMYrV3Zl58.CQ1oT5Uq3jJt6Uhy3VH9u3Gi5YhfZCvZVKgLlaXNFhVKB1zJxvunR7SJa.com.',
'http://kDTvHt1PPDgX5EiP2MwiXjcoWNOhhTuOVAUWJ3TmpBYCC9QoJV114LMYrV3Zl58R.example.com',
'http://[2001:0db8:0000:85a3:0000:0000:ac1f:8001]',
'http://[2001:db8:0:85a3:0:0:ac1f:8001]:123/me.html',
'http://[2001:db8:0:85a3::ac1f:8001]/',
'http://[::1]',
'http://cont-ains.h-yph-en-s.com',
'http://..com',
'http://a.-bc.com',
'http://ab.cd-.com',
'http://-.abc.com',
'http://abc.-.abc.com',
'http://underscore_.example.com',
'http//www.example/wrong/url/',	
'http:/www.example',	
'file:///tmp/test.c',	
'ftp://ftp.example.com/tmp/',	
'/tmp/test.c',	
'/',	
'http://',	
'http:/',	
'http:',	
'http',	
'',	
-1,	
array(),	
'mailto:foo@bar.com',
'news:news.php.net',
'file://foo/bar',
"http://\r\n/bar",
"http://example.com:qq",
"http://example.com:-2",
"http://example.com:65536",
"http://example.com:65537",
);
foreach ($values as $value) {
	var_dump(filter_var($value, FILTER_VALIDATE_URL));
}


var_dump(filter_var("qwe", FILTER_VALIDATE_URL, FILTER_FLAG_SCHEME_REQUIRED));
var_dump(filter_var("http://qwe", FILTER_VALIDATE_URL, FILTER_FLAG_SCHEME_REQUIRED));
var_dump(filter_var("http://", FILTER_VALIDATE_URL, FILTER_FLAG_HOST_REQUIRED));
var_dump(filter_var("/tmp/test", FILTER_VALIDATE_URL, FILTER_FLAG_HOST_REQUIRED));
var_dump(filter_var("http://www.example.com", FILTER_VALIDATE_URL, FILTER_FLAG_HOST_REQUIRED));
var_dump(filter_var("http://www.example.com", FILTER_VALIDATE_URL, FILTER_FLAG_PATH_REQUIRED));
var_dump(filter_var("http://www.example.com/path/at/the/server/", FILTER_VALIDATE_URL, FILTER_FLAG_PATH_REQUIRED));
var_dump(filter_var("http://www.example.com/index.html", FILTER_VALIDATE_URL, FILTER_FLAG_QUERY_REQUIRED));
var_dump(filter_var("http://www.example.com/index.php?a=b&c=d", FILTER_VALIDATE_URL, FILTER_FLAG_QUERY_REQUIRED));

echo "Done\n";
?>
