<?php
$serverCode = <<<'CODE'
    $serverUri = "ssl://127.0.0.1:64321";
    $serverFlags = STREAM_SERVER_BIND | STREAM_SERVER_LISTEN;
    $serverCtx = stream_context_create(['ssl' => [
        'local_cert' => __DIR__ . '/san-cert.pem',
    ]]);

    $server = stream_socket_server($serverUri, $errno, $errstr, $serverFlags, $serverCtx);
    phpt_notify();

    stream_socket_accept($server, 30);
    stream_socket_accept($server, 30);
    stream_socket_accept($server, 30);
    stream_socket_accept($server, 30);
CODE;

$clientCode = <<<'CODE'
    $serverUri = "ssl://127.0.0.1:64321";
    $clientFlags = STREAM_CLIENT_CONNECT;

    phpt_wait();

    $ctx = stream_context_create(['ssl' => ['verify_peer'=> false, 'peer_fingerprint' => true]]);
    $sock = stream_socket_client($serverUri, $errno, $errstr, 30, $clientFlags, $ctx);
    var_dump($sock);

    $ctx = stream_context_create(['ssl' => ['verify_peer'=> false, 'peer_fingerprint' => null]]);
    $sock = stream_socket_client($serverUri, $errno, $errstr, 30, $clientFlags, $ctx);
    var_dump($sock);

    $ctx = stream_context_create(['ssl' => ['verify_peer'=> false, 'peer_fingerprint' => []]]);
    $sock = stream_socket_client($serverUri, $errno, $errstr, 30, $clientFlags, $ctx);
    var_dump($sock);

    $ctx = stream_context_create(['ssl' => ['verify_peer'=> false, 'peer_fingerprint' => ['foo']]]);
    $sock = stream_socket_client($serverUri, $errno, $errstr, 30, $clientFlags, $ctx);
    var_dump($sock);
CODE;

include 'ServerClientTestCase.inc';
ServerClientTestCase::getInstance()->run($clientCode, $serverCode);
