<?php

namespace App\Form\Type;

use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\EmailType;
use Symfony\Component\Form\Extension\Core\Type\SubmitType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolverInterface;

class RegistrantType extends AbstractType
{
    //private $entityManager;

    public function __construct()//EntityManagerInterface $entityManagerInterface)
    {
        //$this->entityManager = $entityManagerInterface;
    }

    public function buildForm(FormBuilderInterface $builder, array $options)
    {
        $builder->add('name', TextType::class, ['label' => false, 'attr' => ['placeholder' => 'Name']])
            ->add('email', EmailType::class, ['label' => false, 'attr' => ['placeholder' => 'Email']])
            ->add('Proceed', SubmitType::class);
    }

    public static function getRoleName(string $role)
    {
        switch ($role) {
            case 'ROLE_USER':
                return 'User';
            case 'ROLE_ADMIN':
                return 'Administrator';
            default:
                return 'Guest';
        }
    }

    public static function getRolesNames(array $roles)
    {
        $r = [];
        foreach ($roles as $role) {
            $r[] = self::getRoleName($role);
        }

        return $r;
    }

    public function getName()
    {
        return 'registrant';
    }

    public function setDefaultOptions(OptionsResolverInterface $resolver)
    {
        $resolver->setDefaults(['data_class' => 'App\Entity\Registrant']);
    }
}
