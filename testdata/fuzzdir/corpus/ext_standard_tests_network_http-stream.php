<?php
require __DIR__.'/../http/server.inc';

$pid = http_server("tcp://127.0.0.1:12342", [__DIR__."/news.rss"]);

$d = new DomDocument;
$e = $d->load("http://127.0.0.1:12342/news.rss");
echo "ALIVE\n";
http_server_kill($pid);
?>
