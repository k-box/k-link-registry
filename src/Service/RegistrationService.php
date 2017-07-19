<?php

namespace App\Service;

use App\Entity\EmailVerification;
use App\Entity\PasswordChangeVerification;
use App\Entity\Registrant;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Security\Core\Encoder\UserPasswordEncoderInterface;

/**
 * Registration service class.
 *
 * This class provides methods that manage various service endpoints & aspects of the
 * registration subsystem.
 */
class RegistrationService
{
    /**
     * @var EntityManagerInterface the entity manager interface associated with the registration service
     */
    private $entityManager;
    /**
     * @var UserPasswordEncoderInterface the user password encoder interface associated with the registration service
     */
    private $encoder;
    /**
     * @var Swift_Mailer the mailer service associated with the registration service used to send out emails
     */
    private $mailer;
    /**
     * @var Twig_Environment the Twig parser environment associated with the registration service used to render email templates
     */
    private $twig;

    /**
     * Registration service instance constructor.
     *
     * The registration service constructor is used to instantiate a service object connected to a manager
     * interface that is used to control the database associated with registration processes.
     *
     * @param EntityManagerInterface       $entityManager the ORM database query backend interface
     * @param UserPasswordEncoderInterface $encoder       the user password encoder interface
     * @param Swift_Mailer                 $mailer        the mailer system service
     * @param Twig_Environment             $twig          the Twig renderer service
     *
     * @return RegistrationService $this returns the RegistrationService instance
     */
    public function __construct(EntityManagerInterface $entityManager, UserPasswordEncoderInterface $encoder, \Swift_Mailer $mailer, \Twig_Environment $twig)
    {
        $this->entityManager = $entityManager;
        $this->encoder = $encoder;
        $this->mailer = $mailer;
        $this->twig = $twig;
    }

    /**
     * Registration service account information updater method.
     *
     * The registration service account information updater method is used to update user account information given a set
     * of parameters. It will create a new registrant if no "registrant_id" is defined.
     *
     * @param array $params the parameters to update the account in the specified registrant ID with
     *
     * @return array $messages returns an array containing the message type & messages encountered while adding a registrant
     */
    public function updateAccountInformation(array $params, array &$messages)
    {
        try {
            if (!array_key_exists('registrant_id', $params)) {
                $user = new Registrant();
            } else {
                $user = $this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R WHERE R.registrant_id=:id')->setParameter('id', $params['registrant_id'])->getSingleResult();
            }
            // TODO validate parameters
            if (array_key_exists('name', $params) && ($params['name'] !== $user->getName())) {
                if (strlen($params['name']) >= 2) {
                    if (!empty($user->getName())) {
                        $messages['Name successfully changed from "'.$user->getName().'" to "'.$params['name'].'".'] = 'info'; // TODO HTML/JS injections... need to clean up
                    }
                    $user->setName($params['name']);
                } else {
                    $messages['You must certainly have a name my good user.'] = 'error';
                }
            }
            if (array_key_exists('email', $params) && ($params['email'] !== $user->getEmail())) {
                $messages['A confirmation link has been sent to the new email address "'.$params['email'].'".'] = 'info'; // TODO HTML/JS injections... need to clean up
                if (!empty($user->getName())) {
                    $messages['After confirmation the email will change.'] = 'info';
                }
                $user->setEmail($params['email']);
            }
            if (array_key_exists('status', $params) && ($params['status'] !== $user->getStatus())) {
                $user->setStatus($params['status']);
            }
            if (array_key_exists('password', $params)) {
                if (strlen($params['password']) >= 8) {
                    $user->setPassword($this->encoder->encodePassword($user, $params['password']));
                } else {
                    $messages['A password of at least 8 characters is required.'] = 'error';
                }
            }
            if (array_key_exists('roles', $params)) { // TODO multiple roles?
                $user->setRole(is_array($params['roles']) ? $params['roles'][0] : $params['roles']);
            }
            $this->entityManager->persist($user);
            $this->entityManager->flush();
            $messages['Registrant modified successfully.'] = 'info';

            return $user;
        } catch (\Exception $exception) {
            $messages = [];
        }
        $messages['Duplicate entry!'] = 'error';

        return $user ? $user : new Registrant();
    }

    /**
     * Registration service account empty password check method.
     *
     * The registration service account empty password check method is used to determine if a user has an initialized account
     * through the registration process. It is the only process by which a null password can be assigned to an account.
     *
     * @param string $email The email address of the user to query for an empty password field
     *
     * @return bool $result returns true if the password on the account is null or false otherwise
     */
    public function hasEmptyPassword(string $email)
    {
        return 1 === (int) ($this->entityManager->createQuery('SELECT COUNT(R.registrant_id) FROM App\Entity\Registrant R WHERE R.email=:email AND R.password IS NULL')->setParameter('email', $email)->getSingleScalarResult());
    }

    /**
     * Registration service user listing method.
     *
     * The registration service user listing method is used to create a list of all registrants in the system with the exclusion
     * of the global administrator account.
     *
     * @return array $result returns a list of all registrants in the system
     */
    public function getAllUsers()
    {
        $users = [];
        foreach ($this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R')->getResult() as $user) {
            $users[] = $user;
        }

        return $users;
    }

    public function getRegistrant(int $id)
    {
        return $this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R WHERE R.registrant_id=:id')->setParameter('id', $id)->getOneOrNullResult();
    }

    /**
     * Registration service registration method.
     *
     * The registration service registration method is used to register a user & initialize them into the system. It does so in a
     * transactional manner to increase performance & avoid duplicate registrant information. If this succeeds, an email is sent
     * to the provided email address to verify it. All registrants are by default disabled until approved by an administrator & are
     * assigned a null password so it is theoretically impossible to use until the verification process is complete.
     *
     * @param string $name    The name of the user
     * @param string $email   The email used to register as a user
     * @param string $from    The administrator email used as the "from" address for outgoing emails
     * @param string $baseURL The verification URL used to validate email tokens
     *
     * @return array $result returns an array containing the message type & messages encountered while during registration
     */
    public function register(string $name, string $email, string $from, string $baseURL)
    {
        $rollback = $registrant = $request = $password = false; // TODO use a transaction frame instead
        try {
            $registrant = new Registrant();
            $registrant->setName($name);
            $registrant->setEmail($email);
            $registrant->setStatus(false);
            $this->entityManager->persist($registrant);
            $this->entityManager->flush();
            $request = new EmailVerification();
            $request->setEmail($email);
            $request->setTimestamp(time());
            $request->setRegistrant($registrant->getID());
            $token = hash('sha512', bin2hex(openssl_random_pseudo_bytes(32))); // TODO separate service
            $password = new PasswordChangeVerification();
            $password->setRegistrant($registrant->getID());
            $password->setToken($token);
            $password->setTimestamp($request->getTimestamp());
            $this->mailer->send((new \Swift_Message('K-Link Registry Verification Email'))->setFrom($from)->setTo($email)->setBody($this->twig->render('email/verification.txt.twig', ['name' => $name, 'url' => $baseURL.'?t='.$token.'&e='.$email]))); // TODO send as an HTML with a form to POST
            $request->setToken($token);
            $this->entityManager->persist($request);
            $this->entityManager->persist($password);
            $this->entityManager->flush();
            $result = ['info', 'Please check your inbox.'];
        } catch (\Exception $exception) {
            $rollback = true;
            $result = ['error', 'The email provided cannot be used, please try another one!'];
        }
        if ($rollback) {
            try {
                if ($registrant) {
                    $this->entityManager->remove($registrant);
                }
                if ($request) {
                    $this->entityManager->remove($request);
                }
                if ($password) {
                    $this->entityManager->remove($password);
                }
                $this->entityManager->flush();
            } catch (\Exception $exception) {
            }
        }

        return $result;
    }

    public function deleteUser(int $user)
    {
        $this->entityManager->createQuery('DELETE FROM App\Entity\Application A WHERE A.registrant_id=:id')->setParameter('id', $user)->getResult();

        return 1 === $this->entityManager->createQuery('DELETE FROM App\Entity\Registrant R WHERE R.registrant_id=:id')->setParameter('id', $user)->getResult();
    }

    public function emptyRegistrant()
    {
        return new Registrant();
    }
}
