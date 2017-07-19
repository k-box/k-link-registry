<?php

namespace App\Entity;

use Doctrine\ORM\Mapping as ORM;

/**
 * PasswordChangeVerification entity class.
 *
 * This class outlines the password change verification model & contains accessor & modifier methods used
 * to acquire & persist information respectively.
 *
 * @ORM\Entity
 */
class PasswordChangeVerification
{
    /**
     * @var int the registrant ID associated with the password change verification object
     *
     * @ORM\Column(type="bigint")
     * @ORM\Id
     * @ORM\GeneratedValue(strategy="NONE")
     * @ORM\ManyToOne(targetEntity="Registrant", inversedBy="registrant_id")
     */
    private $registrant_id;

    /**
     * @var string the security token associated with the password change verification object
     *
     * @ORM\Column(type="string", length=255)
     */
    private $token;

    /**
     * @var int the timestamp issued during the initiation of the password change verification object
     *
     * @ORM\Column(name="`timestamp`", type="integer")
     */
    private $timestamp;

    /**
     * PasswordChangeVerification class token accessor method.
     *
     * This method enables the ability to read the token from the password change verification object.
     *
     * @return string $token returns the token associated with the process
     */
    public function getToken()
    {
        return $this->token;
    }

    /**
     * EmailVerification class registrant ID accessor method.
     *
     * This method enables the ability to read the registrant ID from the password change verification object.
     *
     * @return string $registrant_id returns the registrant ID associated with the process
     */
    public function getRegistrantID()
    {
        return $this->registrant_id;
    }

    /**
     * PasswordChangeVerification class token modifier method.
     *
     * This method enables the ability to change the token issued for the password change verification process.
     *
     * @param string $token the new token to assign to the password change verification process
     */
    public function setToken(string $token)
    {
        $this->token = $token;
    }

    /**
     * PasswordChangeVerification class registrant ID modifier method.
     *
     * This method enables the ability to change the registrant ID issued for the password change verification process.
     *
     * @param int $id the new registrant ID to assign to the password change verification process
     */
    public function setRegistrant(int $id)
    {
        $this->registrant_id = $id;
    }

    /**
     * PasswordChangeVerification class timestamp modifier method.
     *
     * This method enables the ability to change the timestamp issued for the password change verification process.
     *
     * @param int $timestamp the new timestamp to assign to the password change verification process
     */
    public function setTimestamp(int $timestamp)
    {
        $this->timestamp = $timestamp;
    }
}
