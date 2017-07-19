<?php

namespace App\Service;

use App\Entity\Registrant;
use Doctrine\ORM\EntityManagerInterface;

/**
 * Email verification service class.
 *
 * This class provides methods that manage various service endpoints & aspects of the
 * email verification subsystem.
 */
class EmailVerificationService
{
    /**
     * @var EntityManagerInterface the entity manager interface associated with the email verification service
     */
    private $entityManager;

    /**
     * Email verification service instance constructor.
     *
     * The email verification service constructor is used to instantiate a service object connected to a manager
     * interface that is used to control the database associated with email verification processes.
     *
     * @param EntityManagerInterface $entityManager the ORM database query backend interface
     *
     * @return EmailVerificationService returns the EmailVerificationService instance
     */
    public function __construct(EntityManagerInterface $entityManager)
    {
        $this->entityManager = $entityManager;
    }

    /**
     * Email verification service instance token issuer method.
     *
     * The email verification service token issuer method is used to generate a token for the email verification
     * process. Currently not in use & pending removal.
     *
     * @param Registrant $user  the user to be issued a token
     * @param string     $email the email address associated with the account to be verified
     *
     * @return array $result returns an array containing the message type & a set of messages generated during the issuing process
     */
    public function issueToken(Registrant $user, string $email)
    {
        // TODO
        return ['info', 'A password reset link was sent to the email address on file.'];
    }

    /**
     * Email verification service instance token verification method.
     *
     * The email verification service token verification method is used to validate a token for the email verification
     * process.
     *
     * @param string $email the email address associated with the account to be verified
     * @param string $token the security token sent to the stored user email in order to confirm the email address
     *
     * @return bool $result indicates whether the token & email address provided are valid or not
     */
    public function verifyToken(string $email, string $token)
    {
        if (($id = $this->entityManager->createQuery('SELECT E.registrant_id FROM App\Entity\EmailVerification E WHERE E.email=:email AND E.token=:token AND E.timestamp+:expiration>=:time')->setParameter('email', $email)->setParameter('token', $token)->setParameter('expiration', getenv('TOKEN_EXPIRATION_SECONDS'))->setParameter('time', time())->getOneOrNullResult()) && (1 === $this->entityManager->createQuery('DELETE FROM App\Entity\EmailVerification E WHERE E.email=:email AND E.token=:token AND E.timestamp+:expiration>=:time')->setParameter('email', $email)->setParameter('token', $token)->setParameter('expiration', getenv('TOKEN_EXPIRATION_SECONDS'))->setParameter('time', time())->getResult())) {
            $this->entityManager->createQuery('UPDATE App\Entity\Registrant R SET R.email=:email WHERE R.registrant_id=:id')->setParameter('email', $email)->setParameter('id', $id)->getResult();

            return true;
        }

        return false;
    }
}
