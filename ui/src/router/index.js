import Vue from "vue";
import Router from "vue-router";
import Auth from "@/views/Auth";
import Login from "@/views/Login";
import Signup from "@/views/Signup";
import ConfirmEmail from "@/views/ConfirmEmail";
import ConfirmPassword from "@/views/ConfirmPassword";
import ResetPassword from "@/views/ResetPassword";
import Layout from "@/views/Layout";
import Applications from "@/views/Applications";
import Application from "@/views/Application";
import Klinks from "@/views/Klinks";
import Klink from "@/views/Klink";
import Registrants from "@/views/Registrants";
import Registrant from "@/views/Registrant";
import Permissions from "@/views/Permissions";
import auth from "@/utils/auth";

Vue.use(Router);

// Routes that require authentication will have the property
// `meta.requiresAuth=true`. Routes that require admin permission will have
// the property `meta.requiresAdmin`.

const router = new Router({
  base: document.querySelector('meta[name="base"]').getAttribute("content"),
  mode: "history", // fancy url rewriting
  routes: [
    // redirects
    {
      path: "/login",
      redirect: {
        name: "Log in"
      }
    },
    {
      path: "/signup",
      redirect: {
        name: "Sign up"
      }
    },

    {
      path: "/",
      name: "Layout",
      component: Layout,
      redirect: {
        name: "Applications"
      },
      meta: {
        requiresAuth: true
      },
      children: [{
          path: "applications",
          name: "Applications",
          component: Applications,
        },
        {
          path: "applications/:id",
          name: "Application",
          component: Application
        },
        {
          path: "klinks",
          name: "Klinks",
          component: Klinks,
        },
        {
          path: "klinks/:id",
          name: "Klink",
          component: Klink
        },
        {
          path: "registrants",
          name: "Registrants",
          component: Registrants,
        },
        {
          path: "permissions",
          name: "Permissions",
          component: Permissions,
        },
        {
          path: "registrants/:id",
          name: "Registrant",
          component: Registrant
        }
      ]
    },
    {
      path: "/auth",
      component: Auth,
      redirect: {
        name: "Log in"
      },
      meta: {
        requiresAnon: true
      },
      children: [{
          path: "log-in",
          name: "Log in",
          component: Login
        },
        {
          path: "sign-up",
          name: "Sign up",
          component: Signup
        },
        {
          path: "reset-password",
          name: "Reset password",
          component: ResetPassword,
        },
        {
          path: "confirm-email/:token",
          name: "Confirm Email",
          component: ConfirmEmail
        },
        {
          path: "reset-password/:token",
          name: "Confirm Password",
          component: ConfirmPassword
        }
      ]
    }
  ]
});

router.beforeEach((to, from, next) => {
  document.title = to.name;

  if (to.matched.some(record => record.meta.requiresAuth)) {
    // this route requires auth, check if logged in
    // if not, redirect to login page.
    auth
      .loggedIn()
      .then(() => {
        next();
      })
      .catch(e => {
        next({
          name: "Log in",
          query: {
            redirect: to.fullPath
          }
        });
      });

    return;
  }

  if (to.matched.some(record => record.meta.requiresAnon)) {
    // this route requires unauthenticated users,
    // if the user is already logged in, they will be redirected
    // to the main application page
    auth.
    loggedIn()
      .then(() => {
        next({
          name: "Applications",
        });
      }).catch(e => {
        next();
      });
  }

  next();
});

export default router;