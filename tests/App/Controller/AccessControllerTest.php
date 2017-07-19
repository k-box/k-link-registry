<?php

namespace Tests\App\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;

class AccessControllerTest extends WebTestCase
{
    public function testIndex()
    {
        $client = static::createClient();

        $client->request('GET', getenv('KREGISTRY_BASE_URL_PATH').'api/1.0/application.authenticate');

        $this->assertEquals(200, $client->getResponse()->getStatusCode());
        $this->assertJsonStringEqualsJsonString('{"id":false,"error":{"code":-32700,"message":"Invalid JSON object."}}', $client->getResponse()->getContent());

        $client->request('POST', getenv('KREGISTRY_BASE_URL_PATH').'api/1.0/application.authenticate');

        $this->assertEquals(200, $client->getResponse()->getStatusCode());
        $this->assertJsonStringEqualsJsonString('{"id":false,"error":{"code":-32700,"message":"Invalid JSON object."}}', $client->getResponse()->getContent());
    }
}
