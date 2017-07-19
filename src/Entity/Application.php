<?php

namespace App\Entity;

use Doctrine\ORM\Mapping as ORM;
use Symfony\Bridge\Doctrine\Validator\Constraints\UniqueEntity;
use Symfony\Component\Validator\Constraints as Assert;

/**
 * Application entity class.
 *
 * This class outlines the application model & contains accessor & modifier methods used
 * to acquire & persist information respectively.
 *
 * @ORM\Table
 * @ORM\Entity
 * @UniqueEntity(fields="app_domain", message="Domain already in use")
 */
class Application
{
    /**
     * @var int the application ID of the application object
     *
     * @ORM\Column(type="integer")
     * @ORM\Id
     * @ORM\GeneratedValue(strategy="AUTO")
     */
    private $application_id;

    /**
     * @var int the owning registrant ID of the application object
     *
     * @ORM\ManyToOne(targetEntity="Registrant", inversedBy="applications")
     * @ORM\JoinColumn(name="registrant_id", referencedColumnName="registrant_id")
     */
    private $registrant_id;

    /**
     * @var string the name of the application object
     *
     * @ORM\Column(name="`name`", type="string", length=255)
     * @Assert\NotBlank()
     */
    private $name;

    /**
     * @var string the unique domain of the application object
     *
     * @ORM\Column(type="string", length=150, unique=true)
     * @Assert\NotBlank()
     */
    private $app_domain;

    /**
     * @var string the unique & automatically generated authentication token of the application object
     *
     * @ORM\Column(type="string", length=255)
     */
    private $auth_token;

    /**
     * @var array the set of permissions the application object is given access to
     *
     * @ORM\Column(type="simple_array", nullable=true)
     * @ORM\ManyToMany(targetEntity="Permission", mappedBy="name")
     * @ORM\JoinColumn(name="name", referencedColumnName="name")
     */
    private $permissions;

    /**
     * @var bool the status of the application object (may be converted to an integer for a wider range of available options)
     *
     * @ORM\Column(name="`status`", type="boolean")
     */
    private $status;

    /**
     * Application class constructor.
     *
     * This constructor creates a new application object with an empty permissions set & a disabled status.
     *
     * @return Application $this returns a new application object
     */
    public function __construct()
    {
        $this->permissions = [];
        $this->status = 0;
    }

    /**
     * Application class name accessor method.
     *
     * This method enables the ability to read the name from the application object.
     *
     * @return string $name returns the name of the application object
     */
    public function getName()
    {
        return $this->name;
    }

    /**
     * Application class domain accessor method.
     *
     * This method enables the ability to read the domain from the application object.
     *
     * @return string $app_domain returns the domain of the application object
     */
    public function getDomain()
    {
        return $this->app_domain;
    }

    /**
     * Application class ID accessor method.
     *
     * This method enables the ability to read the ID from the application object.
     *
     * @return int $id returns the ID of the application object
     */
    public function getID()
    {
        return $this->application_id;
    }

    /**
     * Application class status accessor method.
     *
     * This method enables the ability to read the status from the application object.
     *
     * @return int $status returns the status of the application object
     */
    public function getStatus()
    {
        return $this->status;
    }

    /**
     * Application class human readable status accessor method.
     *
     * This method enables the ability to read the status name from the application object.
     *
     * @return string $status returns the status name of the application object
     */
    public function getStatusName()
    {
        return $this->status ? 'Enabled' : 'Disabled';
    }

    /**
     * Application class authentication token accessor method.
     *
     * This method enables the ability to read the authentication token from the application object.
     *
     * @return string $auth_token returns the authentication token of the application object
     */
    public function getAuthToken()
    {
        return $this->auth_token;
    }

    /**
     * Application class owning registrant ID accessor method.
     *
     * This method enables the ability to read the owning registrant ID from the application object.
     *
     * @return string $registrant_id returns the owning registrant ID of the application object
     */
    public function getRegistrantID() // TODO Doctrine seems to want the actual object & flares when it's just an id...
    {
        return null === $this->registrant_id ? null : $this->registrant_id->getID();
    }

    /**
     * Application class name modifier method.
     *
     * This method enables the ability to change the name of the application object.
     *
     * @param string $name the new name to assign the application
     */
    public function setName(string $name)
    {
        $this->name = $name;
    }

    /**
     * Application class domain modifier method.
     *
     * This method enables the ability to change the domain of the application object.
     *
     * @param string $domain the new domain to assign the application
     */
    public function setDomain(string $domain)
    {
        $this->app_domain = $domain;
    }

    /**
     * Application class status modifier method.
     *
     * This method enables the ability to change the status of the application object.
     *
     * @param int $status the new status to assign the application
     */
    public function setStatus(int $status)
    {
        $this->status = $status;
    }

    /**
     * Application class authentication token modifier method.
     *
     * This method enables the ability to change the authentication token of the application object.
     *
     * @param string $authToken the new authentication token to assign the application
     */
    public function setAuthToken(string $authToken)
    {
        $this->auth_token = $authToken;
    }

    /**
     * Application class registrant modifier method.
     *
     * This method enables the ability to change the owning registrant of the application object.
     *
     * @param Registrant $registrant the new registrant to assign the application
     */
    public function setRegistrant(Registrant $registrant) // TODO Doctrine seems to want the actual object & flares when it's just an id...
    {
        $this->registrant_id = $registrant;
    }

    /**
     * Application class permissions modifier method.
     *
     * This method enables the ability to change the permissions of the application object.
     *
     * @param array $ids the new permissions to assign the application
     */
    public function setPermissions(array $ids)
    {
        $this->permissions = $ids; // TODO integration with Permission class/object & sanitization + validation to string
    }

    /**
     * Application class permissions accessor method.
     *
     * This method enables the ability to read the permissions of the application object.
     *
     * @return array $permissions returns the set of permissions available to the application object
     */
    public function getPermissions()
    {
        return $this->permissions;
    }

    /**
     * Application class object to JSON format accessor method.
     *
     * This method enables the ability to represent the application object as a JSON string.
     *
     * @return string $this returns the JSON formatted end user accessible components of the application as a JSON array encoded as a string
     */
    public function toJSON()
    {
        return ['name' => $this->name, 'app_url' => $this->app_domain, 'app_id' => $this->application_id, 'permissions' => $this->permissions, 'email' => $this->registrant_id->getEmail()];
    }

    /**
     * Application class permission checker method.
     *
     * This method enables the ability to determine if an application object has all the permissions requested.
     *
     * @param array $permissions the set of permissions to check
     *
     * @return bool $result returns true if the application object has all the permissions specified or false otherwise
     */
    public function hasPermissions(array $permissions)
    {
        $result = true;
        foreach ($permissions as $p) {
            if (!in_array($p, $this->permissions, true)) {
                $result = false;
                break;
            }
        }

        return $result;
    }
}
