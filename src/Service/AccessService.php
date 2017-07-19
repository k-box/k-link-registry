<?php

namespace App\Service;

use App\Entity\Application;
use App\Exception\ApplicationNotFoundException;
use Doctrine\ORM\EntityManagerInterface;

/**
 * Access validation service class.
 *
 * This class provides methods that manage various service endpoints & aspects of the
 * application access subsystem.
 */
class AccessService
{
    /**
     * @var EntityManagerInterface the entity manager interface associated with the access validation service
     */
    private $entityManager;

    /**
     * Access validation service instance constructor.
     *
     * The access validation service constructor is used to instantiate a service object connected to a manager
     * interface that is used to control the database associated with application access validation processes.
     *
     * @param EntityManagerInterface $entityManager the ORM database query backend interface
     *
     * @return AccessService returns the AccessService instance
     */
    public function __construct(EntityManagerInterface $entityManager)
    {
        $this->entityManager = $entityManager;
    }

    /**
     * Access validation service method.
     *
     * The access validation service method is used to determine if an application identified by an authentication
     * token & domain provided has the permissions requested.
     *
     * @param string $appURL      The unique URL identifying the application to access
     * @param int[]  $permissions An array of permissions to request
     * @param string $authToken   The authorization token used to access the K-Registry
     *
     * @return bool returns true if & only if the application has all the permissions specified & is identified by the provided domain & security token
     */
    public function checkAccess(string $appURL, array $permissions, string $authToken)
    {
        try {
            $app = $this->entityManager->createQuery('SELECT A FROM App\Entity\Application A WHERE A.app_domain=:domain AND A.auth_token=:token')->setParameter('domain', $appURL)->setParameter('token', $authToken)->getSingleResult(); // TODO check permissions in SQL directly
            if (!$app) {
                //throw new ApplicationNotFoundException(['No application found identified by URL "'.$appURL.'" or authentication token incorrect.']);
                return false;
            }
        } catch (\Exception $exception) {
            return false;
        }

        return $permissions === [] ? $app->toJSON() : ($app->hasPermissions($permissions) ? $app->toJSON() : false);
    }

    /**
     * Access validation service IP address matcher method.
     *
     * The access validation IP address matcher service is used to determine if a provided IP address falls within a
     * range provided. The method provided is static because it does not require a context.
     *
     * @param string $ip     the IPv4 address to check
     * @param string $ranges IP/CIDR ranges to match to
     *
     * @return bool returns true if the supplied IP is in any of the ranges provided
     */
    public static function matchIP(string $ip, array $ranges)
    {
        foreach ($ranges as $range) {
            list($subnet, $mask) = explode('/', $range);
            if (((ip2long($ip) & ($mask = ~((1 << (32 - $mask)) - 1))) === (ip2long($subnet) & $mask))) {
                return true;
            }
        }

        return false;
    }
}
