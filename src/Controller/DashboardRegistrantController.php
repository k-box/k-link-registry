<?php

namespace App\Controller;

use App\Form\Type\RegistrantEditorType;
use App\Form\Type\RegistrantType;
use App\Service\EmailVerificationService;
use App\Service\PasswordChangeVerificationService;
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
class DashboardRegistrantController extends Controller
{
    /**
     * Dashboard add registrant controller action method.
     *
     * This method is called upon a browser POST request with the objective of adding a new registrant
     * to the system. A valid & active user account with administrative privileges is required to use
     * this function.
     *
     * @Route("/dashboard/registrant/add", name="dashboard_add_registrant")
     *
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to populate, update & create registrants
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardAddRegistrantAction(Request $request, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user || !in_array('ROLE_ADMIN', $user->getRoles(), true)) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = [];
        $form = $this->createForm(RegistrantEditorType::class);
        $form->handleRequest($request);
        //if ($form->isSubmitted() && $form->isValid()) {
        $parameters = [];
        if ($form['roles']->getData()) {
            $parameters['roles'] = $form['roles']->getData();
        }
        if (strlen($form['status']->getData()) > 0) {
            $parameters['status'] = $form['status']->getData();
        }
        if ($form['name']->getData()) {
            $parameters['name'] = $form['name']->getData();
        }
        if ($form['email']->getData()) {
            $parameters['email'] = $form['email']->getData();
        }
        if ($form['password']->getData()) {
            $parameters['password'] = $form['password']->getData();
        }
        if (count($parameters) > 0) {
            $registrant = $registrationService->updateAccountInformation($parameters, $messages);
        } else {
            $registrant = $registrationService->emptyRegistrant();
        }
        //} else {
        //    $registrant = $registrationService->emptyRegistrant();
        //}

        return $this->render('default/dashboard_add_registrant.html.twig', ['form' => $form->createView(), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Dashboard add registrant controller action method.
     *
     * This method is called upon a browser POST request with the objective of adding a new registrant
     * to the system. A valid & active user account with administrative privileges is required to use
     * this function.
     *
     * @Route("/dashboard/registrant/edit/{id}", name="dashboard_edit_registrant")
     *
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to populate, update & create registrants
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardEditRegistrantAction(int $id, Request $request, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user || !in_array('ROLE_ADMIN', $user->getRoles(), true)) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = [];
        $form = $this->createForm(RegistrantEditorType::class, $registrationService->getRegistrant($id));
        $form->handleRequest($request);
        //if ($form->isSubmitted() && $form->isValid()) {
        $parameters = ['registrant_id' => $id];
        if ($form['roles']->getData()) {
            $parameters['roles'] = $form['roles']->getData();
        }
        if (strlen($form['status']->getData()) > 0) {
            $parameters['status'] = $form['status']->getData();
        }
        if ($form['name']->getData()) {
            $parameters['name'] = $form['name']->getData();
        }
        if ($form['email']->getData()) {
            $parameters['email'] = $form['email']->getData();
        }
        if ($form['password']->getData()) {
            $parameters['password'] = $form['password']->getData();
        }
        if (count($parameters) > 1) {
            $registrant = $registrationService->updateAccountInformation($parameters, $messages);
        } else {
            $registrant = $registrationService->emptyRegistrant();
        }
        //} else {
        //    $registrant = $registrationService->emptyRegistrant();
        //}

        return $this->render('default/dashboard_add_registrant.html.twig', ['form' => $form->createView(), 'email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Registration controller action method.
     *
     * This method is called upon a browser request with the objective of rendering the registration page Twig
     * template. The session token is validated to determine if the user is to be redirected to the dashboard,
     * prompted with a form to register or process form data provided by said form. The K-Link Registry
     * registration requires administrator approval for new users. A noonce needs to be introduced to avoid
     * having this form used as a user query tool.
     *
     * @Route("/register", name="register")
     *
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to generate a registration form & process it
     *
     * @return Response returns a response object containing the rendered registration page
     */
    public function registerAction(Request $request, RegistrationService $registrationService)
    {
        if ($this->get('security.authorization_checker')->isGranted('IS_AUTHENTICATED_REMEMBERED')) {
            return $this->redirect($this->generateUrl('dashboard', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        $messages = [];
        $form = $this->createForm(RegistrantType::class);
        $form->handleRequest($request);
        if ($form->isSubmitted() && $form->isValid()) {
            $messages = $registrationService->register($form['name']->getData(), $form['email']->getData(), $this->getParameter('mailer_sender_address'), $this->generateUrl('verify_email', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->render('default/register.html.twig', ['message_type_class' => array_shift($messages), 'messages' => $messages, 'form' => $form->createView(), 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Recovery controller action method.
     *
     * This method is called upon a browser request with the objective of rendering the password recover page
     * Twig template. The session token is validated to determine if the user is to be redirected to the
     * dashboard, prompted with a form to recover their password or process form data provided by said form. A
     * token is sent if a valid e-mail address is supplied with the confirmation to reset the password in
     * a password reset dialog provided in the link. This tool will respond with the same message regardless
     * of whether the token was sent or not to avoid its' exploitation as a user query tool.
     *
     * @Route("/recover", name="recover")
     *
     * @param Request                           $request         the controller browser request object
     * @param PasswordChangeVerificationService $passwordService the password change service instance used to issue a verification token to an email address on file
     *
     * @return Response returns a response object containing the rendered recovery page
     */
    public function recoverAction(Request $request, PasswordChangeVerificationService $passwordService)
    {
        if ($this->get('security.authorization_checker')->isGranted('IS_AUTHENTICATED_REMEMBERED')) {
            return $this->redirect($this->generateUrl('dashboard', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        if ($request->get('email')) { // TODO CSRF & rate-limit token
            $passwordService->issueToken($request->get('email'), $this->getParameter('mailer_sender_address'), $this->generateUrl('password', [], UrlGeneratorInterface::ABSOLUTE_URL));
            $messages = ['info', 'A password reset link was sent to the email address on file.'];
        } else {
            $messages = [''];
        }

        return $this->render('default/recover.html.twig', ['message_type_class' => array_shift($messages), 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }

    /**
     * Password controller action method.
     *
     * This method is called upon a browser request with the objective of rendering the password entry page
     * Twig template. The session token is validated to determine if the user is logged in or not. Should a
     * valid user token be provided, a confirmation e-mail is required to verify the password change. This
     * could be modified to include the original password instead. Only newly registered users are able to
     * enter a new password otherwise. Intermediary passwords are not stored & the change is transactional.
     *
     * @Route("/password", name="password")
     *
     * @param Request                           $request         the controller browser request object
     * @param PasswordChangeVerificationService $passwordService the password change service instance used to validate a token issued by the recovery method
     *
     * @return Response returns a response object containing the rendered password setting page
     */
    public function passwordAction(Request $request, PasswordChangeVerificationService $passwordService)
    {
        // TODO might want to work with post tokens instead for increased security
        $user = $this->getUser(); // TODO force logout
        // TODO password & verify set via POST for non-null password requests
        if ($request->get('e') && $request->get('t')) { // TODO CSRF & rate-limit token
            $messages = ['info', 'Please enter a new password.'];
            if ($request->get('password') && $request->get('verify')) {
                if ($passwordService->validatePassword($request->get('password'), $request->get('verify')) && $passwordService->setPassword($request->get('e'), $request->get('t'), $request->get('password'))) {
                    return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
                }
                $messages = ['error', 'That didn\'t work, please try again later.'];
            }

            return $this->render('default/password.html.twig', ['email' => $request->get('e'), 'token' => $request->get('t'), 'message_type_class' => array_shift($messages), 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
        } elseif ($user) { // TODO CSRF & rate-limit token
            $passwordService->issueToken($user->getEmail(), $this->getParameter('mailer_sender_address'), $this->generateUrl('password', [], UrlGeneratorInterface::ABSOLUTE_URL));

            return $this->redirect($this->generateUrl('dashboard', ['c' => 'info', 'm' => 'A password reset link was sent to the email address on file.'], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
    }

    /**
     * Email verification controller action method.
     *
     * This method is called upon a browser request with the objective of rendering the email verification
     * page Twig template. A token is verified & sent by the controller to the email address provided to
     * validate access to the account provided. If the request is initiated through the registration process,
     * the user will then be redirected to set their password as well. All other requests will result in
     * redirection. It may be desirable to do additional black list & email validation prior to sending an
     * email as this could easily be abused from both resource, reputation & security perspectives.
     *
     * @Route("/verify", name="verify_email")
     *
     * @param Request                  $request             the controller browser request object
     * @param EmailVerificationService $emailService        the email change service instance used to validate a token issued by the registration service or issue one if the user is logged in
     * @param RegistrationService      $registrationService the registration service instance used to validate registrations
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function verifyEmailAction(Request $request, EmailVerificationService $emailService, RegistrationService $registrationService)
    {
        // TODO might want to work with post tokens instead for increased security
        if ($request->get('e') && $request->get('t')) { // TODO CSRF & rate-limit token
            if ($emailService->verifyToken($request->get('e'), $request->get('t'))) {
                if ($registrationService->hasEmptyPassword($request->get('e'))) {
                    return $this->redirect($this->generateUrl('password', ['e' => $request->get('e'), 't' => $request->get('t')], UrlGeneratorInterface::ABSOLUTE_URL));
                }

                return $this->render('default/index.html.twig', ['email' => $request->get('e'), 'message_type_class' => 'info', 'messages' => ['Email verified successfully. Please proceed to login.'], 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
            }
        } elseif (($user = $this->getUser()) && $request->get('e')) { // TODO CSRF & rate-limit token
            $messages = $emailService->issueToken($user, $request->get('e'));

            return $this->redirect($this->generateUrl('dashboard', ['c' => array_shift($messages), 'm' => implode(',', $messages)], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
    }

    /**
     * De-authentication controller action method.
     *
     * This method is called upon a browser request to de-auhtenticate a user token. It is handled by the
     * internal Symfony authentication handler & additional logic is added only to redirect the user back
     * to the default page.
     *
     * @Route("/exit", name="logout")
     *
     * @param Request $request the controller browser request object
     *
     * @return Response returns a response object containing the login redirect
     */
    public function logoutAction(Request $request)
    {
        return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
    }

    /**
     * Registrant delete controller action method.
     *
     * This method is called upon a browser request to delete a registrant from the system. The ROLE_ADMIN
     * is required to perform this action.
     *
     * @Route("/dashboard/registrant/delete/{id}", name="delete_user")
     *
     * @param int                 $id                  the registrant ID to delete
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to delete registrants
     *
     * @return Response returns a response object containing the login redirect
     */
    public function deleteUserAction(int $id, Request $request, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (in_array('ROLE_ADMIN', $user->getRoles(), true) && $registrationService->deleteUser($id)) {
            return $this->redirect($this->generateUrl('dashboard_registrants', ['c' => 'info', 'm' => 'User with ID "'.$id.'" deleted successfully.'], UrlGeneratorInterface::ABSOLUTE_URL));
        }

        return $this->redirect($this->generateUrl('dashboard_registrants', ['c' => 'error', 'm' => 'Permissions: you need them for this operation!'], UrlGeneratorInterface::ABSOLUTE_URL));
    }

    /**
     * Dashboard registrant list controller action method.
     *
     * This method is called upon a browser request with the objective of listing all the registrants
     * in the system. Only a user with administrative privileges can use this tool.
     *
     * @Route("/dashboard/registrants", name="dashboard_registrants")
     *
     * @param Request             $request             the controller browser request object
     * @param RegistrationService $registrationService the registration service instance used to populate registrant data
     *
     * @return Response returns a response object containing the redirected page or login template
     */
    public function dashboardRegistrantsAction(Request $request, RegistrationService $registrationService)
    {
        $user = $this->getUser();
        if (!$user || !in_array('ROLE_ADMIN', $user->getRoles(), true)) {
            return $this->redirect($this->generateUrl('index', [], UrlGeneratorInterface::ABSOLUTE_URL));
        }
        if ($request->get('c') && $request->get('m')) {
            $messages = [$request->get('m') => $request->get('c')];
        } else {
            $messages = [];
        }

        return $this->render('default/dashboard_registrants.html.twig', ['email' => method_exists($user, 'getEmail') ? $user->getEmail() : false, 'name' => method_exists($user, 'getName') ? $user->getName() : false, 'registrants' => $registrationService->getAllUsers(), 'messages' => $messages, 'base_dir' => realpath($this->container->getParameter('kernel.root_dir').'/../web')]);
    }
}
