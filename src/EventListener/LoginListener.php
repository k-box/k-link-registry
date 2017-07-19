<?php

namespace App\EventListener;

use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\Routing\Generator\UrlGeneratorInterface;

class LoginListener implements EventSubscriberInterface
{
    private $router;

    public function __construct(UrlGeneratorInterface $router)
    {
        $this->router = $router;
    }

    public static function getSubscribedEvents()
    {
        return [
            FOSUserEvents::SECURITY_IMPLICIT_LOGIN => 'onImplicitLogin',
        ];
    }

    public function onImplicitLogin(UserEvent $userEvent)
    {
        $user = $userEvent->getAuthenticationToken()->getUser();
        $request = $event->getRequest();

        if ($this->securityContext->isGranted('IS_AUTHENTICATED_FULLY')) {
        }
    }
}
