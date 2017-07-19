<?php

namespace App\Entity;

use App\Form\Type\RegistrantType;
use Doctrine\ORM\Mapping as ORM;
use Symfony\Bridge\Doctrine\Validator\Constraints\UniqueEntity;
use Symfony\Component\Security\Core\User\AdvancedUserInterface;
use Symfony\Component\Validator\Constraints as Assert;

/**
 * Registrant entity class.
 *
 * This class outlines the registrant model & contains accessor & modifier methods used
 * to acquire & persist information respectively. It implements all required Symfony user
 * methods in order to comply with internal authentication processes.
 *
 * @ORM\Table
 * @ORM\Entity
 * @UniqueEntity(fields="email", message="Email already in use")
 */
class Registrant implements AdvancedUserInterface, \Serializable
{
    /**
     * @var string the password of the registrant object
     *
     * @ORM\Column(name="`password`", type="string", length=64, nullable=true)
     */
    protected $password;

    /**
     * @var string the salt associated with the password of the registrant object
     *
     * @ORM\Column(type="string", length=64, nullable=true)
     */
    protected $salt;

    /**
     * @var string the email address associated with the registrant object
     *
     * @ORM\Column(type="string", length=150, unique=true)
     * @Assert\NotBlank()
     * @Assert\Email()
     */
    protected $email;

    /**
     * @var int the unique registrant ID associated with the registrant object (this will actually be encoded as a string or require the use of "bmath" should numbers start becoming large enough)
     *
     * @ORM\Column(type="bigint")
     * @ORM\Id
     * @ORM\GeneratedValue(strategy="AUTO")
     */
    private $registrant_id;

    /**
     * @var array the applications owned by the registrant object
     *
     * @ORM\OneToMany(targetEntity="Application", mappedBy="registrant_id")
     */
    private $applications;

    /**
     * @var string the name or the registrant
     *
     * @ORM\Column(name="`name`", type="string", length=255)
     */
    private $name;

    /**
     * @var string the role slug name given to the registrant
     *
     * @ORM\Column(type="string", length=255)
     */
    private $role;

    /**
     * @var bool the status of the registrant
     *
     * @ORM\Column(name="`status`", type="boolean")
     */
    private $status;

    /**
     * @var int the last timestamp the registrant tried to login
     *
     * @ORM\Column(type="integer", nullable=true)
     */
    private $last_login;

    /**
     * Registrant class constructor.
     *
     * This constructor creates a new registrant object with a default user role & a disabled status.
     *
     * @return Registrant $this returns a new registrant object
     */
    public function __construct()
    {
        $this->role = 'ROLE_USER';
        $this->status = 0;
    }

    /**
     * Registrant class username accessor method.
     *
     * This method enables the ability to read the username from the registrant object.
     *
     * @return string $email returns the username of the registrant object (which is actually their email address)
     */
    public function getUsername()
    {
        return $this->email;
    }

    /**
     * Registrant class password salt accessor method.
     *
     * This method enables the ability to read the password salt from the registrant object.
     *
     * @return string $salt returns the password salt of the registrant object (which is always null with bcrypt)
     */
    public function getSalt()
    {
        return null;
    }

    /**
     * Registrant class password accessor method.
     *
     * This method enables the ability to read the password from the registrant object.
     *
     * @return string $password returns the password of the registrant object
     */
    public function getPassword()
    {
        return $this->password;
    }

    /**
     * Registrant class roles accessor method.
     *
     * This method enables the ability to read the user roles from the registrant object.
     *
     * @return array $role returns the user roles of the registrant object
     */
    public function getRoles()
    {
        return [$this->role];
    }

    public function getRole()
    {
        return $this->role;
    }

    public function setRoles(array $roles)
    {
        $this->role = $roles[0];
    }

    /**
     * Registrant class role names accessor method.
     *
     * This method enables the ability to read the user role names from the registrant object.
     *
     * @return array $role returns the user role names of the registrant object
     */
    public function getRolesNames()
    {
        return [RegistrantType::getRoleName($this->role)];
    }

    /**
     * Registrant class secure credentials disposal method.
     *
     * This method enables the ability to clear sensitive security information from the registrant object.
     */
    public function eraseCredentials()
    {
    }

    /**
     * Registrant class account expiration method.
     *
     * This method enables the ability to determine if a registrant account has expired or not. Due to the
     * simplicity of the K-Link Registry, this is only determined by user status.
     *
     * @return bool $status returns true if & only if the user account status is in good standing
     */
    public function isAccountNonExpired()
    {
        return $this->status;
    }

    /**
     * Registrant class account locked method.
     *
     * This method enables the ability to determine if a registrant account is locked or not. Due to the
     * simplicity of the K-Link Registry, this is only determined by user status.
     *
     * @return bool $status returns true if & only if the user account status is in good standing
     */
    public function isAccountNonLocked()
    {
        return $this->status;
    }

    /**
     * Registrant class account enabled method.
     *
     * This method enables the ability to determine if a registrant account is enabled or not. Due to the
     * simplicity of the K-Link Registry, this is only determined by user status.
     *
     * @return bool $status returns true if & only if the user account status is in good standing
     */
    public function isEnabled()
    {
        return $this->status;
    }

    /**
     * Registrant class account credentials expiration method.
     *
     * This method enables the ability to determine if a registrant account has expired credentialsor not.
     * Due to the simplicity of the K-Link Registry, this is only determined by user status.
     *
     * @return bool $status returns true if & only if the user account status is in good standing
     */
    public function isCredentialsNonExpired()
    {
        return $this->status;
    }

    /**
     * Registrant class serialization method.
     *
     * This method enables the PHP serialization of a registrant object.
     *
     * @return string $this returns a serialized string containing the persisted fields of a registrant account
     */
    public function serialize()
    {
        return serialize([$this->registrant_id, $this->email, $this->password, $this->salt]);
    }

    /**
     * Registrant class de-serialization method.
     *
     * This method enables the PHP de-serialization of a registrant object.
     *
     * @param string $serialized the input serialized registrant object string
     *
     * @return Registrant $this returns a registrant object obtained from a de-serialized string containing the persisted fields of a registrant account
     */
    public function unserialize($serialized)
    {
        return list($this->registrant_id, $this->email, $this->password, $this->salt) = unserialize($serialized);
    }

    /**
     * Registrant class name accessor method.
     *
     * This method enables the ability to read the name from the registrant object.
     *
     * @return string $name returns the name of the registrant
     */
    public function getName()
    {
        return $this->name;
    }

    /**
     * Registrant class email accessor method.
     *
     * This method enables the ability to read the email from the registrant object.
     *
     * @return string $email returns the email of the registrant
     */
    public function getEmail()
    {
        return $this->email;
    }

    /**
     * Registrant class name modifier method.
     *
     * This method enables the ability to change the name of the registrant object.
     *
     * @param string $name the new name to assign the registrant
     */
    public function setName(string $name)
    {
        $this->name = $name;
    }

    /**
     * Registrant class email modifier method.
     *
     * This method enables the ability to change the email of the registrant object.
     *
     * @param string $email the new email to assign the registrant
     */
    public function setEmail(string $email)
    {
        $this->email = $email;
    }

    /**
     * Registrant class role modifier method.
     *
     * This method enables the ability to change the role of the registrant object.
     *
     * @param string $role the new role to assign the registrant
     */
    public function setRole(string $role)
    {
        $this->role = $role; // TODO switch with robust verification
    }

    /**
     * Registrant class password modifier method.
     *
     * This method enables the ability to change the password of the registrant object.
     *
     * @param string $password the new password to assign the registrant
     */
    public function setPassword(string $password)
    {
        $this->password = $password;
    }

    /**
     * Registrant class status modifier method.
     *
     * This method enables the ability to change the status of the registrant object.
     *
     * @param int $status the new status to assign the registrant
     */
    public function setStatus(int $status)
    {
        $this->status = $status;
    }

    /**
     * Registrant class ID accessor method.
     *
     * This method enables the ability to read the ID from the registrant object.
     *
     * @return string $name returns the ID of the registrant
     */
    public function getID()
    {
        return $this->registrant_id;
    }

    /**
     * Registrant class ID accessor method.
     *
     * This method enables the ability to read the ID from the registrant object.
     *
     * @return string $name returns the ID of the registrant
     */
    public function getRegistrantID()
    {
        return $this->registrant_id;
    }

    /**
     * Registrant class status accessor method.
     *
     * This method enables the ability to read the status from the registrant object.
     *
     * @return int $status returns the status of the registrant object
     */
    public function getStatus()
    {
        return $this->status;
    }

    /**
     * Registrant class human readable status accessor method.
     *
     * This method enables the ability to read the status name from the registrant object.
     *
     * @return string $status returns the status name of the registrant object
     */
    public function getStatusName()
    {
        return $this->status ? 'Enabled' : 'Disabled';
    }
}
