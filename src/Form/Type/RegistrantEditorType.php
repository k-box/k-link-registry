<?php

namespace App\Form\Type;

use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\ChoiceType;
use Symfony\Component\Form\Extension\Core\Type\EmailType;
use Symfony\Component\Form\Extension\Core\Type\HiddenType;
use Symfony\Component\Form\Extension\Core\Type\PasswordType;
use Symfony\Component\Form\Extension\Core\Type\SubmitType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolverInterface;

class RegistrantEditorType extends AbstractType
{
    //private $entityManager;

    public function __construct()//EntityManagerInterface $entityManagerInterface)
    {
        //$this->entityManager = $entityManagerInterface;
    }

    public function buildForm(FormBuilderInterface $builder, array $options)
    {
        $builder->add('name', TextType::class, ['label' => 'Name', 'label_attr' => ['class' => 'underlined field'], 'attr' => ['placeholder' => 'Name']])
            ->add('email', EmailType::class, ['label' => 'Email', 'label_attr' => ['class' => 'underlined field'], 'attr' => ['placeholder' => 'Email']])
            ->add('password', PasswordType::class, ['label' => 'Password', 'label_attr' => ['class' => 'underlined field'], 'empty_data' => '', 'required' => false, 'attr' => ['placeholder' => 'Password']])
            ->add('roles', ChoiceType::class, ['label' => 'Role(s)', 'label_attr' => ['class' => 'underlined field'], 'multiple' => true, 'choices' => self::getRoleNames(), 'attr' => ['placeholder' => 'Role(s)']])
            ->add('status', ChoiceType::class, ['label' => 'Status', 'label_attr' => ['class' => 'underlined field'], 'choices' => self::getStatusNames(), 'attr' => ['placeholder' => 'Status']])
            ->add('registrant_id', HiddenType::class)
            ->add('Save', SubmitType::class);
    }

    public static function getStatusNames()
    {
        return ['Disabled' => 0, 'Enabled' => 1];
    }

    public static function getRoleNames()
    {
        return ['User' => 'ROLE_USER', 'Administrator' => 'ROLE_ADMIN'];
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
