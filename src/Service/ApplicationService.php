<?php

namespace App\Service;

use App\Entity\Application;
use App\Entity\Registrant;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Security\Core\User\UserInterface;

/**
 * Application service class.
 *
 * This class provides methods that manage various service endpoints & aspects of the
 * application ecosystem.
 */
class ApplicationService
{
    /**
     * @var EntityManagerInterface the entity manager interface associated with the application service
     */
    private $entityManager;

    /**
     * Application service instance constructor.
     *
     * The application service constructor is used to instantiate a service object connected to a manager
     * interface that is used to control the database associated with applications.
     *
     * @param EntityManagerInterface $entityManager the ORM database query backend interface
     *
     * @return ApplicationService returns the ApplicationService instance
     */
    public function __construct(EntityManagerInterface $entityManager)
    {
        $this->entityManager = $entityManager;
    }

    /**
     * Application deletion service method.
     *
     * The application deletion service method is used to delete an application with the given input ID.
     *
     * @param int $id the application ID to locate & delete
     *
     * @return bool returns success or failure depending on whether the delete operation completed successfully or not
     */
    public function deleteApplication(int $id)
    {
        return 1 === $this->entityManager->createQuery('DELETE FROM App\Entity\Application A WHERE A.application_id=:id')->setParameter('id', $id)->getResult();
    }

    public function emptyApplication()
    {
        return new Application();
    }

    /**
     * Application deletion service method.
     *
     * The application deletion service method is used to delete an application with the given input ID &
     * registrant ID associated as the owner of said application.
     *
     * @param int $owner the registrant ID of the application owner
     * @param int $id    the application ID to locate & delete
     *
     * @return bool returns success or failure depending on whether the delete operation completed successfully or not
     */
    public function deleteApplicationByOwner(int $owner, int $id)
    {
        return 1 === $this->entityManager->createQuery('DELETE FROM App\Entity\Application A WHERE A.id=:id AND A.registrant_id=:owner')->setParameter('id', $id)->setParameter('owner', $owner)->getResult();
    }

    public function updateApplication(array $params, array &$messages)
    {
        try {
            if (!array_key_exists('application_id', $params)) {
                $app = new Application();
                if (!array_key_exists('registrant_id', $params)) {
                    $messages['No owner specified!'] = 'error';

                    return new Application();
                }
                $app->setRegistrant($this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R WHERE R.registrant_id=:id')->setParameter('id', $params['registrant_id'])->getSingleResult());
                $app->setAuthToken(hash('sha512', bin2hex(openssl_random_pseudo_bytes(32)))); // TODO decide on an auth token format, this doesn't work
            } else {
                $app = $this->entityManager->createQuery('SELECT A FROM App\Entity\Application A WHERE A.application_id=:id')->setParameter('id', $params['application_id'])->getSingleResult();
            }
            if (array_key_exists('name', $params) && ($params['name'] !== $app->getName())) {
                if (!empty($app->getName())) {
                    $messages['Name successfully changed from "'.$app->getName().'" to "'.$params['name'].'".'] = 'info'; // TODO HTML/JS injections... need to clean up
                }
                $app->setName($params['name']);
            }
            if (array_key_exists('domain', $params) && ($params['domain'] !== $app->getDomain())) {
                if (!empty($app->getDomain())) {
                    $messages['Domain successfully changed from "'.$app->getDomain().'" to "'.$params['domain'].'".'] = 'info'; // TODO HTML/JS injections... need to clean up
                }
                $app->setDomain($params['domain']);
            }
            if (array_key_exists('registrant_id', $params) && ($params['registrant_id'] !== $app->getRegistrantID())) {
                if (!empty($app->getRegistrantID())) {
                    $messages['Owner successfully changed.'] = 'info'; // TODO HTML/JS injections... need to clean up
                }
                $app->setRegistrant($this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R WHERE R.registrant_id=:id')->setParameter('id', $params['registrant_id'])->getSingleResult());
            }
            if (array_key_exists('status', $params) && ($params['status'] !== $app->getStatus())) {
                $app->setStatus($params['status']);
            }
            if (array_key_exists('permissions', $params)) { // TODO multiple permissions or single vertical?
                $app->setPermissions(is_array($params['permissions']) ? $params['permissions'] : [$params['permissions']]);
            }
            $this->entityManager->persist($app);
            $this->entityManager->flush();

            return $app;
        } catch (\Exception $exception) {
            if (!$app || !$app->getRegistrant()) {
                $messages['No registrant specified!'] = 'error';
            } else {
                $messages['Duplicate application detected!'] = 'error';
            }

            return $app ? $app : new Application();
        }
        $messages['Permission denied!'] = 'error';

        return $app ? $app : new Application();
    }

    /**
     * Application listing service method.
     *
     * The application listing service method is used to list all applications belonging to a specific user ID.
     *
     * @param UserInterface $user the user making the service request
     *
     * @return array $result returns success or failure depending on whether the delete operation completed successfully or not
     */
    public function getApplications(UserInterface $user = null)
    {
        $result = [];
        if ((null === $user) || (!method_exists($user, 'getID'))) {
            $apps = $this->entityManager->createQuery('SELECT A FROM App\Entity\Application A')->getResult();
        } else {
            $apps = $this->entityManager->createQuery('SELECT A FROM App\Entity\Application A WHERE A.registrant_id=:id')->setParameter('id', $user->getID())->getResult();
        }
        foreach ($apps as $app) {
            $result[] = ['registrant_name' => method_exists($user, 'getName') ? $user->getName() : $user->getUsername(), 'registrant_id' => method_exists($user, 'getID') ? $user->getID() : $user->getUsername(), 'auth_token' => $app->getAuthToken(), 'application_id' => $app->getID(), 'name' => $app->getName(), 'status' => $app->getStatus(), 'status_name' => $app->getStatusName(), 'app_domain' => $app->getDomain(), 'permissions' => $app->getPermissions()];
        }

        return $result;
    }

    /**
     * Application permissions service method.
     *
     * This method returns a list of permissions available in HTML format.
     *
     * @return string $result returns a string containing all the permissions as options for an HTML select form (should be moved to Twig)
     */
    public function getPermissionsHTML(array $name = [])
    {
        $options = '';
        foreach ($this->entityManager->createQuery('SELECT P FROM App\Entity\Permission P')->getResult() as $p) {
            $options .= '<option '.(in_array($p->getName(), $name, true) ? 'selected="selected" ' : '').'value="'.$p->getName().'">'.$p->getName().'</option>';
        }

        return $options;
    }

    public function getApplication(int $user, int $id)
    {
        if (0 === $user) {
            return $this->entityManager->createQuery('SELECT A FROM App\Entity\Application A WHERE A.application_id=:app')->setParameter('app', $id)->getOneOrNullResult();
        }

        return $this->entityManager->createQuery('SELECT A FROM App\Entity\Application A WHERE A.application_id=:app AND A.registrant_id=:id')->setParameter('app', $id)->setParameter('id', $user)->getOneOrNullResult();
    }

    /**
     * Application registration service method.
     *
     * This method allows a registrant to register a new application.
     *
     * @param Registrant $user        The user registering the application
     * @param string     $name        The name of the application
     * @param string     $domain      The unique URL hosting the application
     * @param array      $permissions The set of permissions assigned to the application
     * @param string     $status      The status code of the application
     *
     * @return array $result returns an array with the message class & an array of messages collected during the operation
     */
    public function register(Registrant $user, string $name, string $domain, array $permissions, boolean $status)
    {
        $rollback = false;
        try {
            $app = new Application();
            $app->setName($name);
            $app->setAppDomain($domain);
            $app->setStatus($status);
            $app->setRegistrantID($user->getID());
            $app->setAuthToken(hash('sha512', bin2hex(openssl_random_pseudo_bytes(32)))); // TODO decide on an auth token format, this doesn't work
            $app->setPermissions($permissions);
            $this->entityManager->persist($app);
            $this->entityManager->flush();
            $result = ['info', 'Application <b>'.$name.'</b> added successfully.'];
        } catch (\Exception $exception) {
            $rollback = true;
            $result = ['error', 'The application domain you specified is already taken!'];
        }
        if ($rollback) {
            try {
                $this->entityManager->remove($app);
                $this->entityManager->flush();
            } catch (\Exception $exception) {
            }
        }

        return $result;
    }
}
