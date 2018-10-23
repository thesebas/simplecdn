<?php

require_once 'vendor/autoload.php';

$client = new Predis\Client([
    'scheme' => 'tcp',
    'host' => '127.0.0.1',
    'port' => 6379,
]);

$response = $client->mget([
    'jquery',
    'requirejs',
    'ga.js',
]);

\var_dump($response);

foreach ($response as $file ){
    echo file_get_contents($file);
}