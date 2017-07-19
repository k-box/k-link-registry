<?php

namespace App\Entity;

use Doctrine\ORM\Mapping as ORM;

/**
 * EmailVerification entity class.
 *
 * This class outlines the email verification model & contains accessor & modifier methods used
 * to acquire & persist information respectively.
 *
 * @ORM\Entity
 */
class EmailVerification
{
    /**
     * @var string the email associated with the email verification object
     *
     * @ORM\Column(type="string", length=150)
     * @ORM\Id
     * @ORM\GeneratedValue(strategy="NONE")
     * @ORM\ManyToOne(targetEntity="Registrant", inversedBy="email")
     */
    private $email;

    /**
     * @var int the registrant ID associated with the email verification object
     *
     * @ORM\Column(type="bigint")
     * @ORM\ManyToOne(targetEntity="Registrant", inversedBy="registrant_id")
     */
    private $registrant_id;

    /**
     * @var string the security token assigned to the email verification process
     *
     * @ORM\Column(type="string", length=255)
     */
    private $token;

    /**
     * @var int the timestamp issued upon the start of the email verification process
     *
     * @ORM\Column(name="`timestamp`", type="integer")
     */
    private $timestamp;

    /**
     * EmailVerification class timestamp modifier method.
     *
     * This method enables the ability to change the timestamp issued for the email verification process.
     *
     * @param int $timestamp the new timestamp to assign to the email verification process
     */
    public function setTimestamp(int $timestamp)
    {
        $this->timestamp = $timestamp;
    }

    /**
     * EmailVerification class registrant ID modifier method.
     *
     * This method enables the ability to change the registrant associated in the email verification process.
     *
     * @param int $registrant the new registrant to assign to the email verification process
     */
    public function setRegistrant(int $registrant)
    {
        $this->registrant_id = $registrant;
    }

    /**
     * EmailVerification class token modifier method.
     *
     * This method enables the ability to change the token associated in the email verification process.
     *
     * @param string $token the new token to assign to the email verification process
     */
    public function setToken(string $token)
    {
        $this->token = $token;
    }

    /**
     * EmailVerification class email modifier method.
     *
     * This method enables the ability to change the email associated in the email verification process.
     *
     * @param string $email the new email to assign to the email verification process
     */
    public function setEmail(string $email)
    {
        $this->email = $email;
    }

    /**
     * EmailVerification class email accessor method.
     *
     * This method enables the ability to read the email from the email verification object.
     *
     * @return string $email returns the email associated with the process
     */
    public function getEmail()
    {
        return $this->email;
    }

    /**
     * EmailVerification class token accessor method.
     *
     * This method enables the ability to read the token from the email verification object.
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
     * This method enables the ability to read the registrant ID from the email verification object.
     *
     * @return int $registrant_id returns the registrant ID associated with the process
     */
    public function getRegistrant()
    {
        return $this->registrant_id;
    }

    /**
     * EmailVerification class timestamp accessor method.
     *
     * This method enables the ability to read the timestamp from the email verification object.
     *
     * @return int $timestamp returns the timestamp associated with the process
     */
    public function getTimestamp()
    {
        return $this->timestamp;
    }
}
