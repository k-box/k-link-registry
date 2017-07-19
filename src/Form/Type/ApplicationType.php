<?php

namespace App\Form\Type;

use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\SelectType;
use Symfony\Component\Form\Extension\Core\Type\SubmitType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolverInterface;

class ApplicationType extends AbstractType
{
    //private $entityManager;

    public function __construct()//EntityManagerInterface $entityManagerInterface)
    {
        //$this->entityManager = $entityManagerInterface;
    }

    public function buildForm(FormBuilderInterface $builder, array $options)
    {
        $builder->add('name', TextType::class, ['label' => false, 'attr' => ['placeholder' => 'Name']])
            ->add('app_domain', TextType::class, ['label' => false, 'attr' => ['placeholder' => 'App Domain']])
            ->add('auth_token', TextType::class, ['label' => false, 'attr' => ['disabled' => 'disabled', 'placeholder' => 'Auth Token']])
            //->add('permissions', SelectType::class, ['label' => false, 'attr' => ['disabled' => 'disabled', 'placeholder' => 'Auth Token']])
            //->add('status', SelectType::class, ['label' => false, 'attr' => ['placeholder' => 'Status']])
            ->add('Save', SubmitType::class);
    }

    public function getName()
    {
        return 'application';
    }

    public function setDefaultOptions(OptionsResolverInterface $resolver)
    {
        $resolver->setDefaults(['data_class' => 'App\Entity\Application']);
    }
}
