import store from "@/store";
import router from "@/router";
import axios from "axios";

function apiURL(fragment) {
  return `${store.state.baseURL}/api/2.0${fragment}`;
}

function saveToken(token) {
  localStorage.setItem("auth", token);
  store.commit("setJWT", token);
}

function saveUser(user) {
  store.commit("setUser", user);
}

function getToken() {
  return localStorage.getItem("auth");
}

function loggedIn() {
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL("/auth/session"), {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      })
      .then(response => {
        if (response.status === 200) {
          saveToken(response.data.token);
          saveUser({
            id: response.data.id,
            role: response.data.role
          });
          resolve();
        } else {
          reject(new Error(response.data));
        }
      })
      .catch(e => {
        reject(new Error("Could not finish the request"));
      });
  });
}

function login(email, password) {
  let data = {
    email: email,
    password: password
  };
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL("/auth/session"), data)
      .then(response => {
        if (response.status === 200) {
          saveToken(response.data.token);
          resolve();
        } else {
          reject(response.data);
        }
      })
      .catch(e => {
        reject(new Error("Could not finish the request: " + e));
      });
  });
}

// logout deletes the session and navigates to the login page
function logout() {
  localStorage.clear();
  router.push({
    path: "/login"
  });
}

export default {
  loggedIn,
  login,
  logout
};