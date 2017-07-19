<?php

namespace App\Controller;

use App\Service\AccessService;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Method;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Route;
use Swagger\Annotations as SWG;
use Symfony\Bundle\FrameworkBundle\Controller\Controller;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

/**
 * Access controller class.
 *
 * This class provides methods that determine access of certain applications via a dynamically
 * defined permission backend while verifying authenticity of requesting clients. This client
 * is RPC-compliant.
 *
 * @SWG\Swagger(
 *     schemes={"http","https"},
 *     host="kregistry.local",
 *     basePath="/",
 *     @SWG\Info(
 *         version="1.0",
 *         title="K-Registry API",
 *         description="The K-Registry API definition",
 *         termsOfService="https://klink.asia/terms/",
 *         @SWG\Contact(
 *             email="api@klink.asia"
 *         )
 *     )
 * )
 */
class AccessController extends Controller
{
    /**
     * Application access controller action method.
     *
     * This method is called upon a POST RPC browser request with the objective of identifying application
     * permissions from the K-Search registry client. Checks are done to ensure that only certain trusted
     * clients may make such requests & these are configured in the environment variable identified by
     * "KSEARCH_CORE_IP_ADDRESS". See ".env.dist" format for details.
     *
     * @Route("/api/1.0/application.authenticate", name="application_access")
     *
     * @param Request       $request       the controller RPC request object (POST only)
     * @param AccessService $accessService the access service instance used to filter requests
     *
     * @return Response returns a Response object containing the JSON encoded application if access is granted or false
     */
    public function applicationAccessAction(Request $request, AccessService $accessService)
    {
        try {
            $r = json_decode($request->getContent());
        } catch (\Exception $exception) {
            $r = false;
        }
        if ($r && ($request->isFromTrustedProxy() || $accessService->matchIP($request->getClientIp(), explode(',', getenv('KSEARCH_CORE_IP_ADDRESS')))) && property_exists($r, 'id') && ctype_alnum($r->id) && property_exists($r, 'params') && property_exists($r->params, 'app_url') && property_exists($r->params, 'app_secret') && ($app = $accessService->checkAccess($r->params->app_url, (property_exists($r->params, 'permissions') && is_array($r->params->permissions)) ? $r->params->permissions : [], $r->params->app_secret))) {
            $response = new Response(json_encode(['id' => $r->id, 'result' => $app]));
        } else {
            if (!$r) {
                $c = -32700;
                $m = 'Invalid JSON object.';
            } elseif (isset($app)) {
                $c = -32000;
                $m = 'Permission denied.';
            } else {
                $c = -32602;
                $m = 'Invalid request.';
            }
            $response = new Response(json_encode(['id' => $r && property_exists($r, 'id') ? $r->id : false, 'error' => ['code' => $c, 'message' => $m]]));
        }
        $response->setStatusCode(200);
        $response->headers->set('Content-Type', 'application/json');

        return $response;
    }

    /**
     * Application access controller action method.
     *
     * This method is called upon a POST RPC browser request with the objective of identifying application
     * permissions from the K-Search registry client. Checks are done to ensure that only certain trusted
     * clients may make such requests & these are configured in the environment variable identified by
     * "KSEARCH_CORE_IP_ADDRESS". See ".env.dist" format for details.
     *
     * @Route("/application/access", name="application_access_compatibility")
     * @Method({"POST"})
     *
     * @param Request       $request       the controller RPC request object (POST only)
     * @param AccessService $accessService the access service instance used to filter requests
     *
     * @return Response returns a Response object containing the JSON encoded application if access is granted or false
     */
    public function applicationAccessActionCompatibility(Request $request, AccessService $accessService)
    {
        try {
            $r = json_decode($request->getContent());
        } catch (\Exception $exception) {
            $r = false;
        }
        if ($r && ($request->isFromTrustedProxy() || $accessService->matchIP($request->getClientIp(), explode(',', getenv('KSEARCH_CORE_IP_ADDRESS')))) && property_exists($r, 'app_url') && property_exists($r, 'permissions') && property_exists($r, 'auth_token') && ($app = $accessService->checkAccess($r->app_url, $r->permissions, $r->auth_token))) {
            $response = new Response(json_encode($app));
            $response->setStatusCode(200);
        } else {
            if (!$r) {
                $c = -32700;
                $m = 'Invalid JSON object.';
            } elseif (isset($app)) {
                $c = -32000;
                $m = 'Permission denied.';
            } else {
                $c = -32602;
                $m = 'Invalid request.';
            }
            $response = new Response(json_encode(['id' => property_exists($r, 'id') ? $r->id : false, 'error' => ['code' => $c, 'message' => $m]]));
        }
        $response->headers->set('Content-Type', 'application/json');

        return $response;
    }

    /**
     * Application access controller action method.
     *
     * This method is a catch-all used to trap incorrect requests made to the API endpoint it represents.
     *
     * @Route("/application/access", name="application_access_bad_method_compatibility")
     *
     * @param Request $request the controller RPC request object
     *
     * @return Response returns a Response object with a JSON encoded false value
     */
    public function applicationAccessBadMethodActionCompatibility(Request $request)
    {
        $response = new Response(json_encode(false));
        $response->setStatusCode(400);
        $response->headers->set('Content-Type', 'application/json');

        return $response;
    }
}
