<?php

namespace App\Service;

use App\Entity\PasswordChangeVerification;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Security\Core\Encoder\UserPasswordEncoderInterface;

/**
 * Password change verification service class.
 *
 * This class provides methods that manage various service endpoints & aspects of the
 * password modification/verification subsystem.
 */
class PasswordChangeVerificationService
{
    /**
     * @var EntityManagerInterface the entity manager interface associated with the password change service
     */
    private $entityManager;
    /**
     * @var UserPasswordEncoderInterface the user password encoder interface associated with the password change service
     */
    private $encoder;
    /**
     * @var Swift_Mailer the mailer service associated with the password change service used to send out emails
     */
    private $mailer;
    /**
     * @var Twig_Environment the Twig parser environment associated with the password change service used to render email templates
     */
    private $twig;

    /**
     * Password change service instance constructor.
     *
     * The password change service constructor is used to instantiate a service object connected to a manager
     * interface that is used to control the database associated with password change processes.
     *
     * @param EntityManagerInterface       $entityManager the ORM database query backend interface
     * @param UserPasswordEncoderInterface $encoder       the user password encoder interface
     * @param Swift_Mailer                 $mailer        the mailer service
     * @param Twig_Environment             $twig          the Twig rendering service
     *
     * @return PasswordChangeVerificationService returns the PasswordChangeVerificationService instance
     */
    public function __construct(EntityManagerInterface $entityManager, UserPasswordEncoderInterface $encoder, \Swift_Mailer $mailer, \Twig_Environment $twig)
    {
        $this->entityManager = $entityManager;
        $this->encoder = $encoder;
        $this->mailer = $mailer;
        $this->twig = $twig;
    }

    /**
     * Password change service token issuing method.
     *
     * The password change service token issuing method is used to issue a verification token for the provided email
     * address.
     *
     * @param string $email   The email address of the user requesting a password reset
     * @param string $from    The email address from which to send the reset request
     * @param string $baseURL The URL on the site used to reset the password
     *
     * @return bool $result returns as successful response if the token was issued & sent successfully
     */
    public function issueToken(string $email, string $from, string $baseURL)
    {
        try {
            $user = $this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R WHERE R.email=:email')->setParameter('email', $email)->getSingleResult();
            if ($user) {
                $this->entityManager->createQuery('DELETE FROM App\Entity\PasswordChangeVerification P WHERE P.registrant_id=:id')->setParameter('id', $user->getID())->getSingleScalarResult();
                $token = hash('sha512', bin2hex(openssl_random_pseudo_bytes(32))); // TODO separate service... sha512 might be overkill
                $password = new PasswordChangeVerification();
                $password->setRegistrant($user->getID());
                $password->setToken($token);
                $password->setTimestamp(time());
                $this->entityManager->persist($password);
                $this->entityManager->flush();

                return $this->mailer->send((new \Swift_Message('K-Link Registry Reset Password Request'))->setFrom($from)->setTo($email)->setBody($this->twig->render('email/reset_password.txt.twig', ['name' => $user->getName(), 'url' => $baseURL.'?t='.$token.'&e='.$email]))); // TODO send as an HTML with a form to POST
            }
        } catch (\Exception $exception) {
        }

        return false;
    }

    /**
     * Password change service password validation method.
     *
     * The password change service password validation method is used to ensure all new passwords pass certain security
     * & validation standards.
     *
     * @param string $password The new password to be set
     * @param string $verify   The password again
     *
     * @return bool $result returns true if & only if the password validates security standards successfully
     */
    public function validatePassword(string $password, string $verify)
    {
        // TODO other fancy checks
        return ($password === $verify) && (strlen($password) > 6); // TODO turn this into a parameter... or use a RegExp
    }

    /**
     * Password change service password assignment method.
     *
     * The password change service password assignment method is used to apply the password to the account upon token
     * verification & validation.
     *
     * @param string $email    the email address associated with the account password reset request
     * @param string $token    the security token sent to the stored user email in order to confirm the password reset request
     * @param string $password the new password to be assigned to the user identified by the given email address
     *
     * @return bool $result returns whether or not the new password assignment & storage was persisted to the database
     */
    public function setPassword(string $email, string $token, string $password)
    {
        try {
            $user = $this->entityManager->createQuery('SELECT R FROM App\Entity\Registrant R INNER JOIN App\Entity\PasswordChangeVerification P WHERE R.registrant_id=P.registrant_id AND R.email=:email AND P.token=:token AND P.timestamp+:expiration>=:time')->setParameter('email', $email)->setParameter('token', $token)->setParameter('expiration', getenv('TOKEN_EXPIRATION_SECONDS'))->setParameter('time', time())->getSingleResult();
            if ($user) {
                $user->setPassword($this->encoder->encodePassword($user, $password));
                $this->entityManager->persist($user);
                $this->entityManager->flush();

                return 1 === (int) ($this->entityManager->createQuery('DELETE FROM App\Entity\PasswordChangeVerification P WHERE P.registrant_id=:id')->setParameter('id', $user->getID())->getSingleScalarResult());
            }
        } catch (\Exception $exception) {
        }

        return false;
    }

    /**
     * Password change service token verification method.
     *
     * The password change service token verification method is used to verify a token issued to validate a password change
     * request. Currently not in use & pending removal.
     *
     * @param string $email The email address associated with the account password reset request
     * @param string $token The security token sent to the stored user email in order to confirm the password reset request
     *
     * @return bool $result returns true if the token is valid for the given email address
     */
    public function verifyToken(string $email, string $token)
    {
        // TODO

        return false;
    }
}
