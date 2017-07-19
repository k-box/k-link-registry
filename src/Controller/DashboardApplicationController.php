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

/**
 * User interface controller class (will be divided into sub-controllers).
 *
 * This class provides methods that render various components of the user interface, processing
 * input data & displaying components with forms. This controller is not RPC compliant & uses
 * a Twig back-end. All URL rendering is done centrally using an overriding class for the URL
 * generator. See the "App\Routing" namespace for more details on routing configuration.
 */
class DashboardApplicationController extends Controller
{
    /**
     * Application delete controller action method.
     *
     * This method is called upon a browser request to delete an application from the system. Only
     * users with the ROLE_ADMIN are able to delete any application, others can only delete their
     * own.
     *
     * @Route("/dashboard/application/delete/{id}", name="dashboard_delete_application")
     *
     * @param int                 $id                  the application ID to delete
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to delete registrants
     *
     * @return Response returns a response object containing the login redirect
     */
    public function dashboardDeleteApplicationAction(int $id, Request $request, ApplicationService $applicationService)
    {
        // TODO render template, don't redirect
        $user = $this->getUser();
        if ((in_array('ROLE_ADMIN', $user->getRoles(), true) && $applicationService->deleteApplication($id)) || ($user && $applicationService->deleteApplicationByOwner($user->getID(), $id))) {
            return $this->redirect($this->generateUrl('dashboard_applications', ['c' => 'info', 'm' => 'Application with ID '.$id.' deleted successfully.'], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->redirect($this->generateUrl('dashboard_applications', ['c' => 'error', 'm' => 'Permissions: you need them for this operation.'], UrlGeneratorInterface::ABSOLUTE_URL));
    }

    /**
     * Application edit controller action method.
     *
     * This method is called upon a browser request to delete an application from the system. Only
     * users with the ROLE_ADMIN are able to delete any application, others can only delete their
     * own.
     *
     * @Route("/dashboard/application/edit/{id}", name="dashboard_edit_application_form")
     * @Method("POST")
     *
     * @param int                 $id                  the registrant ID to delete
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to delete registrants
     *
     * @return Response returns a response object containing the login redirect
     */
    public function dashboardEditApplicationFormAction(int $id, Request $request, ApplicationService $applicationService, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = [];
        $application = ['application_id' => $id];
        if ($request->get('name')) {
            $application['name'] = $request->get('name');
        }
        if ($request->get('domain')) {
            $application['domain'] = $request->get('domain');
        }
        if ($request->get('permissions')) {
            $application['permissions'] = $request->get('permissions');
        }
        if (strlen($request->get('status')) > 0) {
            $application['status'] = $request->get('status');
        }
        if ($request->get('registrant_id') && in_array('ROLE_ADMIN', $user->getRoles(), true)) {
            $application['registrant_id'] = $request->get('registrant_id');
        }
        if (count($application) > 0) {
            if (!array_key_exists('registrant_id', $application)) {
                $application['registrant_id'] = method_exists($user, 'getID') ? $user->getID() : 0;
            }
            $app = $applicationService->updateApplication($application, $messages);
            $messages['Application updated successfully.'] = 'info';
        } else {
            $app = $applicationService->emptyApplication();
            $messages['Permissions: you need them for this operation.'] = 'error';
        }

        return $this->render('default/dashboard_add_application.html.twig', ['users' => in_array('ROLE_ADMIN', $user->getRoles(), true) ? $registrationService->getAllUsers() : [], 'application' => $app, 'permissions_html' => $applicationService->getPermissionsHTML($app->getPermissions()), 'role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Application edit controller action method.
     *
     * This method is called upon a browser request to delete an application from the system. Only
     * users with the ROLE_ADMIN are able to delete any application, others can only delete their
     * own.
     *
     * @Route("/dashboard/application/edit/{id}", name="dashboard_edit_application")
     *
     * @param int                 $id                  the registrant ID to delete
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to delete registrants
     *
     * @return Response returns a response object containing the login redirect
     */
    public function dashboardEditApplicationAction(int $id, Request $request, ApplicationService $applicationService, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $app = $applicationService->getApplication(method_exists($user, 'getID') ? $user->getID() : 0, $id);

        return $this->render('default/dashboard_add_application.html.twig', ['users' => in_array('ROLE_ADMIN', $user->getRoles(), true) ? $registrationService->getAllUsers() : [], 'application' => $app, 'permissions_html' => $applicationService->getPermissionsHTML($app->getPermissions()), 'role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => [], 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Dashboard add application controller action method.
     *
     * This method is called upon a browser POST request with the objective of adding a new application
     * to the system. A valid & active user account is required to add an application.
     *
     * @Route("/dashboard/application/add", name="dashboard_add_application")
     * @Method({"POST"})
     *
     * @param Request            $request            the controller browser request object
     * @param ApplicationService $applicationService the application service instance used to populate, update & create the applications
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardAddApplicationAction(Request $request, ApplicationService $applicationService, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = $application = [];
        if ($request->get('application_id')) {
            $application['application_id'] = $request->get('application_id');
        }
        if ($request->get('name')) {
            $application['name'] = $request->get('name');
        }
        if ($request->get('domain')) {
            $application['domain'] = $request->get('domain');
        }
        if ($request->get('permissions')) {
            $application['permissions'] = $request->get('permissions');
        }
        if (strlen($request->get('status')) > 0) {
            $application['status'] = $request->get('status');
        }
        if ($request->get('registrant_id') && in_array('ROLE_ADMIN', $user->getRoles(), true)) {
            $application['registrant_id'] = $request->get('registrant_id');
        }
        if (count($application) > 0) {
            if (!array_key_exists('registrant_id', $application)) {
                $application['registrant_id'] = method_exists($user, 'getID') ? $user->getID() : 0;
            }
            $app = $applicationService->updateApplication($application, $messages);
        } else {
            $app = $applicationService->emptyApplication();
        }

        return $this->render('default/dashboard_add_application.html.twig', ['users' => in_array('ROLE_ADMIN', $user->getRoles(), true) ? $registrationService->getAllUsers() : [], 'application' => $app, 'permissions_html' => $applicationService->getPermissionsHTML($app->getPermissions()), 'role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Dashboard add application controller action method.
     *
     * This method is called upon a browser POST request with the objective of adding a new application
     * to the system. A valid & active user account is required to add an application.
     *
     * @Route("/dashboard/application/add", name="dashboard_add_application_form")
     *
     * @param Request            $request            the controller browser request object
     * @param ApplicationService $applicationService the application service instance used to populate, update & create the applications
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardAddApplicationFormAction(Request $request, ApplicationService $applicationService, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->render('default/dashboard_add_application.html.twig', ['users' => in_array('ROLE_ADMIN', $user->getRoles(), true) ? $registrationService->getAllUsers() : [], 'application' => $applicationService->emptyApplication(), 'permissions_html' => $applicationService->getPermissionsHTML(), 'role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => [], 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Dashboard application list controller action method.
     *
     * This method is called upon a browser request with the objective of listing all the applications
     * owned by a specific user. A user with administrative privileges can see all the applications.
     *
     * @Route("/dashboard/applications", name="dashboard_applications")
     *
     * @param Request            $request            the controller browser request object
     * @param ApplicationService $applicationService the application service instance used to populate the applications
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardApplicationsAction(Request $request, ApplicationService $applicationService)
    {
        $user = $this->getUser();
        if (!$user) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        if ($request->get('c') && $request->get('m')) {
            $messages = [$request->get('m') => $request->get('c')];
        } else {
            $messages = [];
        }

        return $this->render('default/dashboard_applications.html.twig', ['apps' => $applicationService->getApplications($user), 'role' => implode(', ', RegistrantType::getRolesNames($user->getRoles())), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'username' => $user->getUsername(), 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }
}
