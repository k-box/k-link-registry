<?php

namespace App\Entity;

use Doctrine\ORM\Mapping as ORM;

/**
 * Permission entity class.
 *
 * This class outlines the permission model & contains accessor & modifier methods used
 * to acquire & persist information respectively.
 *
 * @ORM\Entity
 */
class Permission
{
    /**
     * @var string the variable storing the unique permission name
     *
     * @ORM\Column(name="`name`", type="string", length=150, unique=true)
     * @ORM\Id
     * @ORM\GeneratedValue(strategy="NONE")
     * @ORM\ManyToMany(targetEntity="Application", mappedBy="permissions")
     * @ORM\JoinColumn(name="permissions", referencedColumnName="permissions")
     */
    private $name;

    /**
     * Permission class name accessor method.
     *
     * This method enables the ability to read the permission slug name from the permission object.
     *
     * @return string $name returns the name of the permission object
     */
    public function getName()
    {
        return $this->name;
    }

    /**
     * Permission class name modifier method.
     *
     * This method enables the ability to change the permission slug name of the permission object.
     *
     * @param string $name the new name to assign the permission
     */
    public function setName(string $name)
    {
        $this->name = $name;
    }
}
