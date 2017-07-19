<?php

namespace App\Controller;

use App\Form\Type\RegistrantType;
use App\Service\ApplicationService;
use App\Service\RegistrationService;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Method;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Route;
use Symfony\Bundle\FrameworkBundle\Controller\Controller;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Generator\UrlGeneratorInterface;
use Symfony\Component\Security\Http\Authentication\AuthenticationUtils;

/**
 * Dashboard controller class (will be divided into sub-controllers).
 *
 * This class provides methods that render various components of the user interface, processing
 * input data & displaying components with forms. This controller is not RPC compliant & uses
 * a Twig back-end. All URL rendering is done centrally using an overriding class for the URL
 * generator. See the "App\Routing" namespace for more details on routing configuration.
 */
class DashboardController extends Controller
{
    /**
     * Index page controller action method.
     *
     * This method is called upon a browser request with the objective of rendering the index page Twig
     * template. The session token is validated to determine if the user is to be redirected to the dashboard
     * or prompted to authenticate. The K-Link Registry is only accessible to authenticated users.
     *
     * @Route("/", name="index")
     *
     * @param Request             $request   the controller browser request object
     * @param AuthenticationUtils $authUtils the authentication utilites instance used to validate user access
     *
     * @return Response returns a response object containing the rendered index page
     */
    public function indexAction(Request $request, AuthenticationUtils $authUtils)
    {
        if ($this->get('security.authorization_checker')->isGranted('IS_AUTHENTICATED_REMEMBERED')) {
            return $this->redirect($this->generateUrl('dashboard', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = $authUtils->getLastAuthenticationError() ? ['That didn\'t work, please try again.'] : false;

        return $this->render('default/index.html.twig', ['email' => (($messages) ? $authUtils->getLastUsername() : ''), 'message_type_class' => (($messages) ? 'error' : ''), 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Dashboard controller action method.
     *
     * This method is called upon a browser POST request with the objective of responding to various
     * dashboard API calls from the dashboard Twig template. A valid & active user account is required
     * to use the dashboard. User roles also determine which interfaces are rendered. Currently only
     * two roles are configured with the unauthenticated user being referred to as a guest. The
     * ROLE_ADMIN represents the administration role & is used to manage all aspects of the registry.
     * The ROLE_USER represents a registrant & a K-Box institution with the ability to create & manage
     * applications. Forms for system settings, permissions & registrants are available to administrators
     * only while forms for user settings & applications are availble to both.
     *
     * @Route("/dashboard", name="dashboard_form")
     * @Method({"POST"})
     *
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to populate, update & create registrants
     * @param ApplicationService  $applicationService  the application service instance used to populate, update & create the applications
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardFormAction(Request $request, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = $parameters = [];
        if (in_array('ROLE_ADMIN', $user->getRoles(), true)) {
            if ($request->get('registrant_id')) {
                $parameters['registrant_id'] = $request->get('registrant_id');
                $parameters['roles'] = $request->get('registrant_roles');
                $parameters['status'] = $request->get('registrant_status');
            }
        } else {
            $parameters['registrant_id'] = $user->getID();
        }
        if ($request->get('registrant_name')) {
            $parameters['name'] = $request->get('registrant_name');
        }
        if ($request->get('registrant_email')) {
            $parameters['email'] = $request->get('registrant_email');
        }
        if ($request->get('registrant_password')) {
            $parameters['password'] = $request->get('registrant_password');
        }
        if (count($parameters) > 0) {
            $messages = $registrationService->updateAccountInformation($parameters);
        }

        return $this->render('default/dashboard.html.twig', ['role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Dashboard controller action method.
     *
     * This method is called upon a browser request with the objective of rendering the dashboard Twig
     * template. A valid & active user account is required to access this template & various roles
     * enable use of the following interfaces: permission editor (ROLE_ADMIN), permissions viewer
     * (ROLE_ADMIN), registrant editor (ROLE_ADMIN), registrants viewer (ROLE_ADMIN), server settings
     * (ROLE_ADMIN), account settings (all users), application editor (all) & applications viewer (all).
     * Note: each tool will be moved to a controller method of its own.
     *
     * @Route("/dashboard", name="dashboard")
     *
     * @param Request             $request             the controller browser request object
     * @param ApplicationService  $applicationService  the application service instance used to populate the applications
     * @param RegistrationService $registrationService the registration service instance used to populate the registrants
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardAction(Request $request, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->render('default/dashboard.html.twig', ['role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => [], 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }
}
